package horizon

import (
	"fmt"
	"github.com/SafeRE-IT/horizon/db2"
	"time"

	"github.com/SafeRE-IT/horizon/db2/history"
	"github.com/SafeRE-IT/horizon/render/hal"
	"github.com/SafeRE-IT/horizon/render/problem"
	"github.com/SafeRE-IT/horizon/resource"
)

type Sort int64

const (
	SortTypeDefaultPage Sort = iota
	SortTypeMostFounded
	SortTypeByEndTime
	SortTypeByPopularity
	SortTypeStartTime
)

// SaleIndexAction renders slice of reviewable requests
type SaleIndexAction struct {
	Action
	Owner        string
	BaseAsset    string
	OpenOnly     bool
	Upcoming     bool
	Voting       bool
	Promotions   bool
	SortType     *int64
	Name         string
	Records      []history.Sale
	PagingParams db2.PageQueryV2
	Page         hal.Page
}

// JSON is a method for actions.JSON
func (action *SaleIndexAction) JSON() {
	action.Do(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.loadRecord,
		action.loadPage,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

func (action *SaleIndexAction) loadParams() {
	action.PagingParams = action.GetPageQueryV2()
	action.Owner = action.GetString("owner")
	action.BaseAsset = action.GetString("base_asset")
	action.Name = action.GetString("name")

	// TODO: refactoring required: switch to state
	action.OpenOnly = action.GetBool("open_only")
	action.Upcoming = action.GetBool("upcoming")
	action.Voting = action.GetBool("voting")
	action.Promotions = action.GetBool("promotions")

	action.SortType = action.GetOptionalInt64("sort_by")
	action.Page.Filters = map[string]string{
		"owner":      action.Owner,
		"base_asset": action.BaseAsset,
		"name":       action.Name,
		"open_only":  action.GetString("open_only"),
		"upcoming":   action.GetString("upcoming"),
		"voting":     action.GetString("voting"),
		"promotions": action.GetString("promotions"),
		"sort_by":    action.GetString("sort_by"),
	}
}

func (action *SaleIndexAction) loadRecord() {
	q := action.HistoryQ().Sales().PageV2(action.PagingParams)

	if action.Owner != "" {
		q = q.ForOwner(action.Owner)
	}

	if action.BaseAsset != "" {
		q = q.ForBaseAsset(action.BaseAsset)
	}

	if action.Name != "" {
		q = q.ForName(action.Name)
	}

	if action.OpenOnly {
		q = q.Open(time.Now().UTC())
	}

	if action.Upcoming {
		q = q.Upcoming(time.Now().UTC())
	}

	if action.Voting {
		q = q.Voting()
	}

	if action.Promotions {
		q = q.Promotions()
	}

	sortBy := SortTypeDefaultPage
	if action.SortType != nil {
		sortBy = Sort(*action.SortType)
	}

	switch sortBy {
	case SortTypeDefaultPage:
		// FIXME tmp:
		q = q.OrderById("asc")
	case SortTypeStartTime:
		q = q.OrderByStartTime()
	case SortTypeByEndTime:
		q = q.OrderByEndTime()
	case SortTypeByPopularity:
		values, err := action.CoreQ().OrderBook().InvestorsCount()
		if err != nil {
			action.Log.WithError(err).Error("Unable to load investors count")
			action.Err = &problem.ServerError
			return
		}
		q = q.OrderByPopularity(values)
	default:
		action.SetInvalidField("sort_by", fmt.Errorf("invalid value %d", sortBy))
		return
	}

	converter, err := action.CreateConverter()
	if err != nil {
		action.Log.WithError(err).Error("Failed to init converter")
		action.Err = &problem.ServerError
		return
	}

	action.Records, err = selectSalesWithCurrentCap(q, converter)
	if err != nil {
		action.Log.WithError(err).Error("failed to load sales")
		action.Err = &problem.ServerError
		return
	}
}

func (action *SaleIndexAction) loadPage() {
	for i := range action.Records {
		var res resource.Sale
		res.Populate(&action.Records[i])
		err := populateSaleWithStats(action.Records[i].ID, &res, action.CoreQ())
		if err != nil {
			action.Log.WithError(err).
				WithField("sale_id", action.Records[i].ID).
				Error("failed to populate stat for sale")
			action.Err = &problem.ServerError
			return
		}
		action.Page.Add(&res)
	}

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Page = action.PagingParams.Page
	action.Page.Limit = action.PagingParams.Limit
	action.Page.PopulateLinksV2()
}
