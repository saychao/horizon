package history

import (
	"time"
)

type RecoveryRequest struct {
	ID         uint64    `db:"id"`
	OldAccount string    `db:"old_account"`
	NewAccount string    `db:"new_account"`
	Accepted   *bool     `db:"accepted"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
