package resource

import (
	"github.com/SafeRE-IT/horizon/db2/history"
)

func (this *HistoryAccount) Populate(row history.Account) {
	this.AccountID = row.Address
}
