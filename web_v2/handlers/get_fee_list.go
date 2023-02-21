package handlers

import (
	"net/http"

	"github.com/SafeRE-IT/horizon/db2/core2"
	"github.com/SafeRE-IT/horizon/web_v2/ctx"
	"github.com/SafeRE-IT/horizon/web_v2/requests"
	"github.com/SafeRE-IT/horizon/web_v2/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"
)

// GetFeeList - processes request to get the list of fees
func GetFeeList(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getFeeListHandler{
		AccountQ: core2.NewAccountsQ(coreRepo),
		AssetsQ:  core2.NewAssetsQ(coreRepo),
		FeesQ:    core2.NewFeesQ(coreRepo),
		Log:      ctx.Log(r),
	}

	request, err := requests.NewGetFeeList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetFeeList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get fee list ", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getFeeListHandler struct {
	FeesQ    core2.FeesQ
	AccountQ core2.AccountsQ
	AssetsQ  core2.AssetsQ
	Log      *logan.Entry
}

// GetFeeList returns the list of fees with related resources
func (h *getFeeListHandler) GetFeeList(request *requests.GetFeeList) (*regources.FeeRecordListResponse, error) {
	q := h.FeesQ.Page(*request.PageParams)
	if request.ShouldFilter(requests.FilterTypeFeeListAccount) {
		q = q.FilterByAddress(request.Filters.Account)
	}
	if request.ShouldFilter(requests.FilterTypeFeeListAccountRole) {
		q = q.FilterByAccountRole(request.Filters.AccountRole)
	}

	if request.ShouldFilter(requests.FilterTypeFeeListAsset) {
		q = q.FilterByAsset(request.Filters.Asset)
	}

	if request.ShouldFilter(requests.FilterTypeFeeListSubtype) {
		q = q.FilterBySubtype(request.Filters.Subtype)
	}

	if request.ShouldFilter(requests.FilterTypeFeeListFeeType) {
		q = q.FilterByType(request.Filters.FeeType)
	}

	if request.ShouldFilter(requests.FilterTypeFeeListLowerBound) {
		q = q.FilterByLowerBound(int64(request.Filters.LowerBound))
	}

	if request.ShouldFilter(requests.FilterTypeFeeListUpperBound) {
		q = q.FilterByUpperBound(int64(request.Filters.UpperBound))
	}

	fees, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get fee list")
	}

	response := &regources.FeeRecordListResponse{
		Data:  make([]regources.FeeRecord, 0, len(fees)),
		Links: request.GetOffsetLinks(*request.PageParams),
	}

	assets := make([]string, 0, len(fees))
	for i := range fees {
		fee := resources.NewFee(fees[i])

		fee.Relationships.Asset = resources.NewAssetKey(fees[i].Asset).AsRelation()
		if fees[i].AccountID != "" {
			fee.Relationships.Account = resources.NewAccountKey(fees[i].AccountID).AsRelation()
		} else if fees[i].AccountRole != core2.FeesEmptyRole {
			fee.Relationships.AccountRole =
				resources.NewAccountRoleKey(fees[i].AccountRole).AsRelation()
		}
		assets = append(assets, fees[i].Asset)
		response.Data = append(response.Data, fee)
	}

	if request.ShouldInclude(requests.IncludeTypeFeeListAsset) {
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
