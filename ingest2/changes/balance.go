package changes

import (
	"github.com/SafeRE-IT/horizon/db2/history2"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

type balanceStorage interface {
	InsertBalance(rawID xdr.BalanceId, balance history2.Balance) error
	MustBalance(rawID xdr.BalanceId) history2.Balance
}

type balanceHandler struct {
	storage        balanceStorage
	accountStorage accountStorage
}

func newBalanceHandler(accountStorage accountStorage, storage balanceStorage) *balanceHandler {
	return &balanceHandler{
		storage:        storage,
		accountStorage: accountStorage,
	}
}

// Created - stores new instance of balance
func (c *balanceHandler) Created(lc ledgerChange) error {
	balance := lc.LedgerChange.MustCreated().Data.MustBalance()
	account := c.accountStorage.MustAccount(balance.AccountId)
	newBalance := history2.NewBalance(uint64(balance.SequentialId), account.ID, balance.BalanceId.AsString(),
		string(balance.Asset))
	err := c.storage.InsertBalance(balance.BalanceId, newBalance)
	if err != nil {
		return errors.Wrap(err, "failed to insert balance", logan.F{
			"balance_address": newBalance.Address,
		})
	}
	return nil
}
