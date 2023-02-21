package handlers

import (
	"net/http"
	"time"

	"gitlab.com/tokend/go/amount"

	"gitlab.com/distributed_lab/ape"

	"github.com/SafeRE-IT/horizon/web_v2/ctx"

	"gitlab.com/tokend/regources/generated"

	"github.com/SafeRE-IT/horizon/ledger"
	"gitlab.com/tokend/go/xdr"
)

// GetRoot - handles requests for `/v3` - returns state of the horizon
func GetRoot(w http.ResponseWriter, r *http.Request) {
	coreInfo := ctx.CoreInfo(r)
	currentState := ledger.CurrentState()
	currentTime := time.Unix(time.Now().Unix(), 0).UTC()
	response := regources.HorizonStateAttributes{
		Core:               stateToLedgerInfo(currentState.Core),
		History:            stateToLedgerInfo(currentState.History),
		HistoryV2:          stateToLedgerInfo(currentState.History2),
		CurrentTime:        currentTime,
		CurrentTimeUnix:    currentTime.Unix(),
		EnvironmentName:    coreInfo.MasterExchangeName,
		MasterAccountId:    coreInfo.AdminAccountID,
		NetworkPassphrase:  coreInfo.NetworkPassphrase,
		Precision:          amount.One,
		TxExpirationPeriod: coreInfo.TxExpirationPeriod,
		XdrRevision:        xdr.Revision,
		CoreVersion:        coreInfo.CoreVersion,
	}

	ape.Render(w, regources.HorizonStateResponse{
		Data: regources.HorizonState{
			Key: regources.Key{
				ID:   currentTime.Format(time.RFC3339),
				Type: regources.HORIZON_STATE,
			},
			Attributes: response,
		},
	})
}

func stateToLedgerInfo(state ledger.State) regources.LedgerInfo {
	return regources.LedgerInfo{
		LastLedgerIncreaseTime: time.Unix(state.LastLedgerIncreaseTime.Unix(), 0).UTC(),
		Latest:                 uint64(state.Latest),
		OldestOnStart:          uint64(state.OldestOnStart),
	}
}
