package horizon

import (
	"net/http"
	"time"

	"github.com/SafeRE-IT/horizon/cache"

	"github.com/SafeRE-IT/horizon/corer"
	"gitlab.com/tokend/go/resources"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"

	"github.com/SafeRE-IT/horizon/db2/core2"
	hdoorman "github.com/SafeRE-IT/horizon/doorman"
	"github.com/SafeRE-IT/horizon/log"
	"github.com/SafeRE-IT/horizon/web_v2"
	"github.com/SafeRE-IT/horizon/web_v2/ctx"
	"github.com/SafeRE-IT/horizon/web_v2/handlers"
	v2middleware "github.com/SafeRE-IT/horizon/web_v2/middleware"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/go/doorman"
	"gitlab.com/tokend/go/xdr"
)

const (
	kvKeyLicenseAdminSignerRole = "license_admin_signer_role"
	kvKeyKycRecoverySignerRole  = "kyc_recovery_signer_role"
	kvKeyKycRecoveryEnabled     = "kyc_recovery_enabled"
)

type WebV2 struct {
	mux     *chi.Mux
	metrics *web_v2.WebMetrics
}

func initWebV2(app *App) {
	mux := chi.NewMux()
	metrics := web_v2.NewWebMetrics()

	app.webV2 = &WebV2{
		mux:     mux,
		metrics: metrics,
	}
}

func initWebV2Middleware(app *App) {
	m := app.webV2.mux

	logger := &log.DefaultLogger.Entry

	drm, err := createDoorman(app)
	if err != nil {
		panic(errors.Wrap(err, "failed to init doorman, stopping"))
	}

	cacher := cache.NewMiddlewareCache(app.config.CacheSize, app.config.CachePeriod)

	m.Use(
		middleware.StripSlashes,
		middleware.SetHeader(upstreamHeader, app.config.Hostname),
		middleware.RequestID,
		ape.LoganMiddleware(logger, time.Second, ape.LoggerSetter(ctx.SetLog),
			ape.RequestIDProvider(middleware.GetReqID)),
		ape.RecoverMiddleware(logger),
		ape.CtxMiddleWare(
			// log will be set by logger setter on handler call
			ctx.SetCoreRepo(app.CoreRepoLogged(nil)),
			// log will be set by logger setter on handler call
			ctx.SetHistoryRepo(app.HistoryRepoLogged(nil)),
			ctx.SetDoorman(drm),
			ctx.SetCoreInfo(func() corer.Info {
				// required to make sure that we return most recent core info
				if app.CoreInfo == nil {
					panic("unexpected state: core info is nil, while we are preparing context for the request")
				}

				return *app.CoreInfo
			}),
			ctx.SetSubmitter(app.submitterV2),
			ctx.SetConfig(&app.config),
		),
		v2middleware.WebMetrics(app),
		cacher.Middleware,
	)

	if app.config.CORSEnabled {
		// TODO: chi doesn't provide an analogue, should write own implementation?
		//r.Use(middleware.AutomaticOptions)

		c := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedHeaders:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "DELETE", "PATCH", "PUT"},
			AllowCredentials: true,
		})

		m.Use(c.Handler)
	}

	signatureValidator := &SignatureValidator{app.config.SkipCheck}

	m.Use(signatureValidator.MiddlewareV2)
	m.NotFound(func(w http.ResponseWriter, r *http.Request) {
		err := problems.NotFound()
		err.Detail = "Unknown path"
		err.Meta = &map[string]interface{}{
			"url": r.URL.String(),
		}
		ape.RenderErr(w, err)
	})
}

func createDoorman(app *App) (doorman.Doorman, error) {
	// TODO: Add constraint for KYC recovery signer when clients will be fully ready
	lazyLicenseSignerConstrain := newLazySignerConstrain(
		core2.NewKeyValueQ(app.CoreRepoLogged(nil)),
		kvKeyLicenseAdminSignerRole,
		func(_ *core2.KeyValueQ) bool {
			return true
		},
	)

	constrains := []doorman.SignerOfExt{
		lazyLicenseSignerConstrain,
	}

	signersProvider := hdoorman.NewSignersQ(core2.NewSignerQ(app.CoreRepoLogged(nil)))
	return doorman.NewWithOpts(app.Conf().SkipCheck, signersProvider, doorman.SignerOfOpts{Constraints: constrains}), nil
}

type lazyRestrictedRoleConstrain struct {
	kvQ   *core2.KeyValueQ
	kvKey string

	// check prerequisites
	// should be executed once and save result into `enabled` field on success
	// if error occurs while ensuring - should panic
	mustEnsureEnabled func(app *core2.KeyValueQ) bool
	enabled           *bool

	lazyRole *uint64
}

func newLazySignerConstrain(kvQ *core2.KeyValueQ, kvKey string, shouldCheck func(kvQ *core2.KeyValueQ) bool) *lazyRestrictedRoleConstrain {
	return &lazyRestrictedRoleConstrain{
		kvQ:               kvQ,
		kvKey:             kvKey,
		mustEnsureEnabled: shouldCheck,
	}
}

func (l *lazyRestrictedRoleConstrain) Check(signer resources.Signer) bool {
	if l.enabled == nil {
		enabled := l.mustEnsureEnabled(l.kvQ)
		l.enabled = &enabled
	}

	isEnabled := *l.enabled
	if !isEnabled {
		return true
	}

	if l.lazyRole == nil {
		roleVal, err := getKVValue(l.kvQ, l.kvKey, xdr.KeyValueEntryTypeUint64)
		if err != nil {
			panic(errors.Wrap(err, "failed to get role kv value", logan.F{
				"key":  l.kvKey,
				"type": xdr.KeyValueEntryTypeUint64,
			}))
		}

		lazyRestrictedRole := uint64(*roleVal.Ui64Value)
		l.lazyRole = &lazyRestrictedRole
	}

	return *l.lazyRole != signer.Role
}

func getKVValue(kvQ *core2.KeyValueQ, key string, typ xdr.KeyValueEntryType) (*xdr.KeyValueEntryValue, error) {
	kv, err := kvQ.ByKey(key)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get key value from core")
	}

	if kv == nil {
		return nil, errors.New("kv not found in db")
	}

	if kv.Value.Type != typ {
		return nil, errors.From(errors.New("provided kv type conflicts with one got from db, check terraform"), logan.F{
			"expected_type": typ.String(),
			"got_type":      kv.Value.Type.String(),
		})
	}

	switch kv.Value.Type {
	case xdr.KeyValueEntryTypeUint32:
		if kv.Value.Ui32Value == nil {
			return nil, errors.New("u32 kv value is nil, check terraform")
		}
	case xdr.KeyValueEntryTypeUint64:
		if kv.Value.Ui64Value == nil {
			return nil, errors.New("u64 kv value is nil, check terraform")
		}
	case xdr.KeyValueEntryTypeString:
		if kv.Value.StringValue == nil {
			return nil, errors.New("string kv value is nil, check terraform")
		}
	default:
		return nil, errors.New("invalid kv type") // should never happen unless new kv types arrive
	}

	value := xdr.KeyValueEntryValue(kv.Value)
	return &value, nil
}

func initWebV2Actions(app *App) {
	m := app.webV2.mux

	// DEPRECATED
	m.Get("/v3", handlers.GetRoot)

	m.Get("/v3/info", handlers.GetRoot)

	m.Get("/v3/license", handlers.GetCurrentLicenseInfo)

	m.Post("/v3/transactions", handlers.CreateTransaction)

	m.Get("/v3/accounts", handlers.GetAccountList)
	m.Get("/v3/accounts/{id}", handlers.GetAccount)
	m.Get("/v3/accounts/{id}/signers", handlers.GetAccountSigners)
	m.Get("/v3/accounts/{id}/calculated_fees", handlers.GetCalculatedFees)
	m.Get("/v3/accounts/{id}/sales", handlers.GetSaleListForAccount)
	m.Get("/v3/accounts/{id}/sales/{sale_id}", handlers.GetSaleForAccount)
	m.Get("/v3/accounts/{id}/converted_balances/{asset_code}", handlers.GetConvertedBalances)
	m.Get("/v3/accounts/{id}/balances_statistic/{asset_code}", handlers.GetBalanceStatistic)
	m.Get("/v3/assets/{code}", handlers.GetAsset)
	m.Get("/v3/assets", handlers.GetAssetList)
	m.Get("/v3/balances/{id}", handlers.GetBalance)
	m.Get("/v3/balances", handlers.GetBalanceList)
	m.Get("/v3/fees", handlers.GetFeeList)
	m.Get("/v3/limits", handlers.GetLimitsList)
	m.Get("/v3/history", handlers.GetHistory)
	m.Get("/v3/history/{id}", handlers.GetParticipantEffect)
	m.Get("/v3/matches", handlers.GetMatchList)
	m.Get("/v3/movements", handlers.GetMovements)
	m.Get("/v3/asset_pairs/{id}", handlers.GetAssetPair)
	m.Get("/v3/asset_pairs", handlers.GetAssetPairList)
	m.Get("/v3/offers/{id}", handlers.GetOffer)
	m.Get("/v3/offers", handlers.GetOfferList)
	m.Get("/v3/public_key_entries/{id}", handlers.GetPublicKeyEntry)
	m.Get("/v3/transactions", handlers.GetTransactions)
	m.Get("/v3/transactions/{id}", handlers.GetTransaction)

	// reviewable requests
	m.Get("/v3/requests", handlers.GetRequests)
	m.Get("/v3/requests/{id}", handlers.GetRequests)
	m.Get("/v3/create_asset_requests", handlers.GetCreateAssetRequests)
	m.Get("/v3/create_asset_requests/{id}", handlers.GetCreateAssetRequests)
	m.Get("/v3/create_sale_requests", handlers.GetCreateSaleRequests)
	m.Get("/v3/create_sale_requests/{id}", handlers.GetCreateSaleRequests)
	m.Get("/v3/update_asset_requests", handlers.GetUpdateAssetRequests)
	m.Get("/v3/update_asset_requests/{id}", handlers.GetUpdateAssetRequests)
	m.Get("/v3/create_pre_issuance_requests", handlers.GetCreatePreIssuanceRequests)
	m.Get("/v3/create_pre_issuance_requests/{id}", handlers.GetCreatePreIssuanceRequests)
	m.Get("/v3/create_issuance_requests", handlers.GetCreateIssuanceRequests)
	m.Get("/v3/create_issuance_requests/{id}", handlers.GetCreateIssuanceRequests)
	m.Get("/v3/create_withdraw_requests", handlers.GetCreateWithdrawRequests)
	m.Get("/v3/create_withdraw_requests/{id}", handlers.GetCreateWithdrawRequests)
	m.Get("/v3/update_limits_requests", handlers.GetUpdateLimitsRequests)
	m.Get("/v3/update_limits_requests/{id}", handlers.GetUpdateLimitsRequests)
	m.Get("/v3/create_aml_alert_requests", handlers.GetCreateAmlAlertRequests)
	m.Get("/v3/create_aml_alert_requests/{id}", handlers.GetCreateAmlAlertRequests)
	m.Get("/v3/change_role_requests", handlers.GetChangeRoleRequests)
	m.Get("/v3/change_role_requests/{id}", handlers.GetChangeRoleRequests)
	m.Get("/v3/update_sale_details_requests", handlers.GetUpdateSaleDetailsRequests)
	m.Get("/v3/update_sale_details_requests/{id}", handlers.GetUpdateSaleDetailsRequests)
	m.Get("/v3/create_atomic_swap_ask_requests", handlers.GetCreateAtomicSwapAskRequests)
	m.Get("/v3/create_atomic_swap_ask_requests/{id}", handlers.GetCreateAtomicSwapAskRequests)
	m.Get("/v3/create_atomic_swap_bid_requests", handlers.GetCreateAtomicSwapBidRequests)
	m.Get("/v3/create_atomic_swap_bid_requests/{id}", handlers.GetCreateAtomicSwapBidRequests)

	m.Get("/v3/atomic_swap_asks/{id}", handlers.GetAtomicSwapAsk)
	m.Get("/v3/atomic_swap_asks", handlers.GetAtomicSwapAskList)

	m.Get("/v3/create_poll_requests", handlers.GetCreatePollRequests)
	m.Get("/v3/create_poll_requests/{id}", handlers.GetCreatePollRequests)
	m.Get("/v3/kyc_recovery_requests", handlers.GetKYCRecoveryRequests)
	m.Get("/v3/kyc_recovery_requests/{id}", handlers.GetKYCRecoveryRequests)
	m.Get("/v3/manage_offer_requests", handlers.GetManageOfferRequests)
	m.Get("/v3/manage_offer_requests/{id}", handlers.GetManageOfferRequests)
	m.Get("/v3/create_payment_requests", handlers.GetCreatePaymentRequests)
	m.Get("/v3/create_payment_requests/{id}", handlers.GetCreatePaymentRequests)

	m.Get("/v3/redemption_requests", handlers.GetRedemptionRequests)
	m.Get("/v3/redemption_requests/{id}", handlers.GetRedemptionRequests)

	m.Get("/v3/data_creation_requests", handlers.GetDataCreationRequests)
	m.Get("/v3/data_creation_requests/{id}", handlers.GetDataCreationRequests)

	m.Get("/v3/create_deferred_payments_requests", handlers.GetCreateDeferredPaymentRequests)
	m.Get("/v3/create_deferred_payments_requests/{id}", handlers.GetCreateDeferredPaymentRequests)

	m.Get("/v3/create_close_deferred_payments_requests", handlers.GetCloseDeferredPaymentRequests)
	m.Get("/v3/create_close_deferred_payments_requests/{id}", handlers.GetCloseDeferredPaymentRequests)

	m.Get("/v3/data_update_requests", handlers.GetDataUpdateRequests)
	m.Get("/v3/data_update_requests/{id}", handlers.GetDataUpdateRequests)

	m.Get("/v3/data_remove_requests", handlers.GetDataRemoveRequests)
	m.Get("/v3/data_remove_requests/{id}", handlers.GetDataRemoveRequests)

	m.Get("/v3/key_values", handlers.GetKeyValueList)
	m.Get("/v3/key_values/{key}", handlers.GetKeyValue)

	m.Get("/v3/polls/{id}", handlers.GetPoll)
	m.Get("/v3/polls", handlers.GetPollList)
	m.Get("/v3/polls/{id}/relationships/votes", handlers.GetVoteList)
	m.Get("/v3/polls/{id}/relationships/votes/{voter}", handlers.GetVote)

	m.Get("/v3/votes/{voter}", handlers.GetVoterVotesList)

	m.Get("/v3/sales", handlers.GetSaleList)
	m.Get("/v3/sales/{id}", handlers.GetSale)
	m.Get("/v3/sales/{id}/relationships/whitelist", handlers.GetSaleWhitelist)
	m.Get("/v3/sales/{id}/relationships/participation", handlers.GetSaleParticipations)

	m.Get("/v3/order_book/{id}", handlers.DeprecatedGetOrderBook)
	m.Get("/v3/order_books/{base}:{quote}:{order_book_id}", handlers.GetOrderBook)

	m.Get("/v3/account_roles/{id}", handlers.GetAccountRole)
	m.Get("/v3/account_roles", handlers.GetAccountRoleList)
	m.Get("/v3/account_rules/{id}", handlers.GetAccountRule)
	m.Get("/v3/account_rules", handlers.GetAccountRuleList)

	m.Get("/v3/signer_roles/{id}", handlers.GetSignerRole)
	m.Get("/v3/signer_roles", handlers.GetSignerRoleList)
	m.Get("/v3/signer_rules/{id}", handlers.GetSignerRule)
	m.Get("/v3/signer_rules", handlers.GetSignerRuleList)

	m.Get("/v3/swaps/{id}", handlers.GetSwap)
	m.Get("/v3/swaps", handlers.GetSwapList)

	m.Get("/v3/data/{id}", handlers.GetData)
	m.Get("/v3/data", handlers.GetDataList)

	m.Get("/v3/operations", handlers.GetOperations)
	m.Get("/v3/operations/{id}", handlers.GetOperation)

	m.Get("/v3/deferred_payments", handlers.GetDeferredPaymentList)
	m.Get("/v3/deferred_payments/{id}", handlers.GetDeferredPayment)

	m.Get("/v3/data_owner_update_requests", handlers.GetDataOwnerUpdateRequests)
	m.Get("/v3/data_owner_update_requests/{id}", handlers.GetDataOwnerUpdateRequests)

	cop := app.config.Copus()
	if err := cop.RegisterChi(m); err != nil {
		panic(errors.Wrap(err, "failed to register service"))
	}
}

func init() {
	appInit.Add(
		"web2.init",
		initWebV2,
		"app-context", "core-info", "memory_cache", "ledger-state", "submitter_v2",
	)

	appInit.Add(
		"web2.middleware",
		initWebV2Middleware,
		"web2.init",
		"web.metrics",
	)

	appInit.Add(
		"web2.actions",
		initWebV2Actions,
		"web2.middleware", "web2.init",
	)
}
