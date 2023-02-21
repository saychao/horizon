package ingest

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/saychao/horizon/db2/core"
	"github.com/saychao/horizon/db2/history"
)

// IntegrityCheck validates history database state and tries to fix known issues
func (i *System) IntegrityCheck() error {
	err := i.ensureAccountTypes()
	if err != nil {
		return errors.Wrap(err, "failed to ensure account types")
	}
	return nil
}

func (i *System) ensureAccountTypes() error {
	var accounts []history.Account
	var cursor int64
	var limit uint64 = 100
	coreQ := core.NewQ(i.CoreDB)
	for {
		stmt := sq.
			Select("*").
			From("history_accounts").
			OrderBy("id").
			Where("id > ?", cursor).
			Where("account_type = 0").
			Limit(limit)
		err := i.HorizonDB.Select(&accounts, stmt)
		if err != nil {
			return errors.Wrap(err, "failed to get history accounts")
		}
		// goes through all accounts in history database for which account type
		// is not set
		for _, account := range accounts {
			cursor = account.ID
			// gets core account
			coreAccount, err := coreQ.Accounts().ByAddress(account.Address)
			if err != nil {
				return errors.Wrap(err, "failed to get core account")
			}
			// and updates account type for history record
			stmt := sq.
				Update("history_accounts").
				Set("account_type", coreAccount.RoleID).
				Where("id = ?", account.ID)
			err = i.HorizonDB.Exec(stmt)
			if err != nil {
				return errors.Wrap(err, "failed to update history account")
			}
		}
		if uint64(len(accounts)) < limit {
			// no more accounts to check
			return nil
		}
	}
}
