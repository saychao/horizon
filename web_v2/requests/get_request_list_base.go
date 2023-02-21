package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
)

type GetRequestListBaseFilters struct {
	ID                  *uint64 `filter:"id"`
	Requestor           *string `filter:"requestor"`
	Reviewer            *string `filter:"reviewer"`
	State               *uint64 `filter:"state"`
	Type                *uint64 `filter:"type"`
	PendingTasks        *uint64 `filter:"pending_tasks"`
	PendingTasksAnyOf   *uint64 `filter:"pending_tasks_any_of"`
	PendingTasksNotSet  *uint64 `filter:"pending_tasks_not_set"`
	MissingPendingTasks *uint64 `filter:"missing_pending_tasks"`
	CreatedBefore       *int64  `filter:"created_before"`
	CreatedAfter        *int64  `filter:"created_after"`
	AllTasks            *uint64 `filter:"all_tasks"`
	AllTasksAnyOf       *uint64 `filter:"all_tasks_any_of"`
	AllTasksNotSet      *uint64 `filter:"all_tasks_not_set"`
}

type GetRequestsBase struct {
	*base
	Filters    GetRequestListBaseFilters
	PageParams pgdb.CursorPageParams
	Includes   struct {
		RequestDetails bool `include:"request_details"`
	}
}

func NewGetRequestsBase(
	r *http.Request,
	filters map[string]struct{},
	includes map[string]struct{},
) (GetRequestsBase, error) {

	// merge filters
	mergedFilters := map[string]struct{}{}
	for k := range filters {
		mergedFilters[k] = struct{}{}
	}
	for k := range filterTypeRequestListAll {
		mergedFilters[k] = struct{}{}
	}
	// merge includes
	mergedIncludes := map[string]struct{}{}
	for k := range includes {
		mergedIncludes[k] = struct{}{}
	}
	for k := range includeTypeReviewableRequestListAll {
		mergedIncludes[k] = struct{}{}
	}

	base, err := newBase(r, baseOpts{
		supportedFilters:  mergedFilters,
		supportedIncludes: mergedIncludes,
	})
	request := GetRequestsBase{base: base}
	if err != nil {
		return request, err
	}

	return request, nil
}
func PopulateRequest(requestsBase *GetRequestsBase) error {
	var err error

	err = SetDefaultCursorPageParams(&requestsBase.PageParams)
	if err != nil {
		return err
	}

	ID, err := requestsBase.base.getUint64ID()
	requestsBase.Filters.ID = &ID
	if err != nil {
		return err
	}

	return nil
}
