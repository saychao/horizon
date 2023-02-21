package resource

import (
	"github.com/saychao/horizon/db2/history"
)

func (this *HistoryAccount) Populate(row history.Account) {
	this.AccountID = row.Address
}
