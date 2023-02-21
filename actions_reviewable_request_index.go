package horizon

import (
	"github.com/SafeRE-IT/horizon/db2"
	"github.com/SafeRE-IT/horizon/db2/history"
	"github.com/SafeRE-IT/horizon/render/hal"
	"github.com/SafeRE-IT/horizon/render/problem"
	"github.com/SafeRE-IT/horizon/resource/reviewablerequest"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
)

// ReviewableRequestIndexAction renders slice of reviewable requests
type ReviewableRequestIndexAction struct {
	Action
	CustomFilter       func(action *ReviewableRequestIndexAction)
	CustomCheckAllowed func(action *ReviewableRequestIndexAction)
	q                  history.ReviewableRequestQI
	Reviewer           string
	Requestor          string
	State              *int64
	UpdatedAfter       *int64
	HistRecords        []history.ReviewableRequest
	Records            []regources.ReviewableRequest

	RequestTypes []xdr.ReviewableRequestType

	PagingParams  db2.PageQuery
	Page          hal.Page
	DisablePaging bool
}

// JSON is a method for actions.JSON
func (action *ReviewableRequestIndexAction) JSON() {
	action.Do(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.checkAllowed,
		action.loadRecord,
		action.loadPage,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

func (action *ReviewableRequestIndexAction) loadParams() {
	action.DisablePaging = false
	action.PagingParams = action.GetPageQuery()
	if action.PagingParams.Cursor == "" {
		action.DisablePaging = true
	}
	action.Reviewer = action.GetString("reviewer")
	action.Requestor = action.GetString("requestor")
	action.State = action.GetOptionalInt64("state")
	action.UpdatedAfter = action.GetOptionalInt64("updated_after")
	action.Page.Filters = map[string]string{
		"reviewer":      action.Reviewer,
		"requestor":     action.Requestor,
		"state":         action.GetString("state"),
		"updated_after": action.GetString("updated_after"),
	}
}

func (action *ReviewableRequestIndexAction) checkAllowed() {

	if action.CustomCheckAllowed != nil {
		action.CustomCheckAllowed(action)
		return
	}

	ownersOfData := []string{action.App.CoreInfo.AdminAccountID}
	if action.Requestor != "" {
		ownersOfData = append(ownersOfData, action.Requestor)
	}

	if action.Reviewer != "" {
		ownersOfData = append(ownersOfData, action.Reviewer)
	}

	action.IsAllowed(ownersOfData...)
}

func (action *ReviewableRequestIndexAction) loadRecord() {
	if !action.DisablePaging {
		return
	}
	action.q = action.HistoryQ().ReviewableRequests()

	if action.CustomFilter != nil {
		action.CustomFilter(action)
	}

	if action.Err != nil {
		return
	}

	if action.Reviewer != "" {
		action.q = action.q.ForReviewer(action.Reviewer)
	}

	if action.Requestor != "" {
		action.q = action.q.ForRequestor(action.Requestor)
	}

	if action.State != nil {
		action.q = action.q.ForState(*action.State)
	}

	if action.UpdatedAfter != nil {
		action.q = action.q.UpdatedAfter(*action.UpdatedAfter)
	}

	action.q = action.q.ForTypes(action.RequestTypes)

	var err error
	action.HistRecords, err = action.q.Page(action.PagingParams).Select()
	if err != nil {
		action.Log.WithError(err).Error("failed to load reviewable requests")
		action.Err = &problem.ServerError
		return
	}
}

func (action *ReviewableRequestIndexAction) loadPage() {
	for i := range action.HistRecords {
		res, err := reviewablerequest.PopulateReviewableRequest(&action.HistRecords[i])
		if err != nil {
			action.Log.WithError(err).Error("Failed to populate reviewable request")
			action.Err = &problem.ServerError
			return
		}
		action.Page.Add(res)
	}

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}
