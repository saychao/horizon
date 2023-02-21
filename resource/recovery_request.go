package resource

import (
	"fmt"
	"strconv"
	"time"

	"github.com/saychao/horizon/db2/history"
)

type RecoveryRequest struct {
	PT         string    `json:"paging_token"`
	OldAccount string    `json:"old_account"`
	NewAccount string    `json:"new_account"`
	Accepted   *bool     `json:"accepted"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Populate fills out the resource's fields
func (request *RecoveryRequest) Populate(row *history.RecoveryRequest) error {
	request.PT = strconv.FormatUint(row.ID, 10)
	request.OldAccount = row.OldAccount
	request.NewAccount = row.NewAccount
	request.Accepted = row.Accepted
	request.CreatedAt = row.CreatedAt
	request.UpdatedAt = row.UpdatedAt
	return nil
}

// PagingToken implementation for hal.Pageable
func (request *RecoveryRequest) PagingToken() string {
	return fmt.Sprintf("%s", request.PT)
}
