package storage

import (
	"github.com/SafeRE-IT/horizon/db2/core2"
	"github.com/SafeRE-IT/horizon/db2/history2"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

type accountStorage interface {
	// MustAccount - gets account by rawAddress
	MustAccount(rawAddress xdr.AccountId) history2.Account
}

// Balance is helper struct to aid in loading balances and storing balances
// it uses lazy change, which stores all selected or inserted balances locally.
type Balance struct {
	balances map[xdr.BalanceId]*history2.Balance

	balanceQ       history2.BalancesQ
	coreBalances   core2.BalancesQ
	historyRepo    *pgdb.DB
	accountStorage accountStorage
}

// NewBalance - creates new instance of Balance
func NewBalance(repo *pgdb.DB, coreRepo *pgdb.DB, accountStorage accountStorage) *Balance {
	return &Balance{
		historyRepo:    repo,
		balanceQ:       history2.NewBalancesQ(repo),
		balances:       make(map[xdr.BalanceId]*history2.Balance),
		coreBalances:   core2.NewBalancesQ(coreRepo),
		accountStorage: accountStorage,
	}
}

// MustBalance - gets balance by rawAddress from local change or storage
// In case of error panics
// If balance does not exists in horizon db, tries to get it from core and inserts into horizon
func (s *Balance) MustBalance(rawAddress xdr.BalanceId) history2.Balance {
	balance, err := s.getBalance(rawAddress)
	if err != nil {
		panic(errors.Wrap(err, "failed to get balance"))
	}
	return balance
}

func (s *Balance) getBalance(rawAddress xdr.BalanceId) (history2.Balance, error) {
	balance, ok := s.balances[rawAddress]
	if ok {
		return *balance, nil
	}

	address := rawAddress.AsString()
	balance, err := s.balanceQ.GetByAddress(address)
	if err != nil {
		return history2.Balance{}, errors.Wrap(err, "failed to load balance by address", logan.F{
			"address": address,
		})
	}

	// it seems that we have ingest partial history
	if balance == nil {
		balance = new(history2.Balance)
		*balance, err = s.getBalanceFromCore(rawAddress, address)
		if err != nil {
			return history2.Balance{}, errors.Wrap(err, "failed to load balance from core db",
				logan.F{"address": address})
		}
	}

	s.balances[rawAddress] = balance
	return *balance, nil
}

func (s *Balance) getBalanceFromCore(rawAddress xdr.BalanceId, address string) (history2.Balance, error) {
	coreBalance, err := s.coreBalances.GetByAddress(address)
	if err != nil {
		return history2.Balance{}, errors.Wrap(err, "failed to load balance by address from core", logan.F{
			"address": address,
		})
	}

	if coreBalance == nil {
		return history2.Balance{}, errors.From(errors.New("balance not found in core db"),
			logan.F{"balance_address": address})
	}

	var rawAccountAddress xdr.AccountId
	err = rawAccountAddress.SetAddress(coreBalance.AccountAddress)
	if err != nil {
		return history2.Balance{}, errors.Wrap(err, "failed to set address for account_id",
			logan.F{"str_address": coreBalance.AccountAddress})
	}

	account := s.accountStorage.MustAccount(rawAccountAddress)
	balance := history2.NewBalance(coreBalance.BalanceSeqID, account.ID, coreBalance.BalanceAddress, coreBalance.AssetCode)
	err = s.InsertBalance(rawAddress, balance)
	if err != nil {
		return history2.Balance{}, errors.Wrap(err, "failed to insert balance into db after fatched it from core")
	}

	return balance, nil
}

// InsertBalance - adds balance to local cache and inserts it into storage
func (s *Balance) InsertBalance(rawBalanceID xdr.BalanceId, balance history2.Balance) error {
	// it's ok if the balance already exists in the map.
	// Such case could occur during roll back of transaction and retry to process same ledger
	s.balances[rawBalanceID] = &balance
	err := s.historyRepo.ExecRaw("INSERT INTO balances (id, account_id, address, asset_code) VALUES($1, $2, $3, $4)",
		balance.ID, balance.AccountID, balance.Address, balance.AssetCode)
	if err != nil {
		return errors.Wrap(err, "failed to insert new balance", logan.F{
			"address": balance.Address,
			"id":      balance.ID,
		})
	}

	return nil
}

// MustBalanceID - returns balanceID. Panics if fails to find one.
func (s *Balance) MustBalanceID(rawBalanceID xdr.BalanceId) uint64 {
	return s.MustBalance(rawBalanceID).ID
}
