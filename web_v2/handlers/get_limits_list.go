package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"github.com/SafeRE-IT/horizon/db2/core2"
	"github.com/SafeRE-IT/horizon/web_v2/ctx"
	"github.com/SafeRE-IT/horizon/web_v2/requests"
	"github.com/SafeRE-IT/horizon/web_v2/resources"
	regources "gitlab.com/tokend/regources/generated"
)

// GetLimitsList - processes request to get the list of fees
func GetLimitsList(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getLimitsListHandler{
		AccountQ: core2.NewAccountsQ(coreRepo),
		AssetsQ:  core2.NewAssetsQ(coreRepo),
		LimitsQ:  core2.NewLimitsQ(coreRepo),
		Log:      ctx.Log(r),
	}

	request, err := requests.NewGetLimitsList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetLimitsList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get fee list ", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getLimitsListHandler struct {
	LimitsQ  core2.LimitsQ
	AccountQ core2.AccountsQ
	AssetsQ  core2.AssetsQ
	Log      *logan.Entry
}

// GetLimitsList returns the list of fees with related resources
func (h *getLimitsListHandler) GetLimitsList(request *requests.GetLimitsList) (*regources.LimitsListResponse, error) {
	q := h.LimitsQ.Page(request.PageParams)
	if request.Filters.Account != nil {
		q = q.FilterByAccount(*request.Filters.Account)
	}
	if request.Filters.AccountRole != nil {
		q = q.FilterByAccountRole(*request.Filters.AccountRole)
	}

	if request.Filters.Asset != nil {
		q = q.FilterByAsset(*request.Filters.Asset)
	}

	if request.Filters.StatsOpType != nil {
		q = q.FilterByStatsOpType(*request.Filters.StatsOpType)
	}

	limits, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get fee list")
	}

	response := &regources.LimitsListResponse{
		Data:  make([]regources.Limits, 0, len(limits)),
		Links: request.GetOffsetLinks(request.PageParams),
	}

	assets := make([]string, 0, len(limits))
	for i := range limits {
		limit := resources.NewLimits(limits[i])

		limit.Relationships.Asset = resources.NewAssetKey(limits[i].AssetCode).AsRelation()
		if limits[i].AccountId != nil {
			limit.Relationships.Account = resources.NewAccountKey(*limits[i].AccountId).AsRelation()
		}
		if limits[i].AccountType != nil {
			limit.Relationships.AccountRole = resources.NewAccountRoleKey(*limits[i].AccountType).AsRelation()
		}
		assets = append(assets, limits[i].AssetCode)
		response.Data = append(response.Data, limit)
	}

	if request.ShouldInclude(requests.IncludeTypeLimitsListAsset) {
		assetRecords, err := h.AssetsQ.FilterByCodes(assets).Select()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get assets")
		}
		if assetRecords == nil {
			return nil, errors.New("Assets not found")
		}

		for _, v := range assetRecords {
			asset := resources.NewAsset(v)
			response.Included.Add(&asset)
		}
	}

	return response, nil
}
