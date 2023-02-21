package storage

import (
	"github.com/SafeRE-IT/horizon/db2/core2"
	"github.com/SafeRE-IT/horizon/db2/history2"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

// Account is helper struct to aid in loading accounts and storing accounts
// it uses lazy change, which stores all selected or inserted accounts locally.
type Account struct {
	accounts map[[32]byte]*history2.Account

	accountQ     history2.AccountsQ
	coreAccounts core2.AccountsQ
	repo         *pgdb.DB
}

// NewAccount - creates new instance of Account
func NewAccount(repo *pgdb.DB, coreRepo *pgdb.DB) *Account {
	return &Account{
		repo:         repo,
		accountQ:     history2.NewAccountsQ(repo),
		accounts:     make(map[[32]byte]*history2.Account),
		coreAccounts: core2.NewAccountsQ(coreRepo),
	}
}

// MustAccount - gets account by rawAddress from local change or storage
// In case of error panics
// If account does not exists in horizon db, tries to get it from core and inserts into horizon
func (a *Account) MustAccount(rawAddress xdr.AccountId) history2.Account {
	account, err := a.getAccount(rawAddress)
	if err != nil {
		panic(errors.Wrap(err, "failed to get account"))
	}
	return account
}

func (a *Account) getAccount(rawAddress xdr.AccountId) (history2.Account, error) {
	account, ok := a.accounts[*rawAddress.Ed25519]
	if ok {
		return *account, nil
	}

	address := rawAddress.Address()
	account, err := a.accountQ.ByAddress(address)
	if err != nil {
		return history2.Account{}, errors.Wrap(err, "failed to load account by address", logan.F{
			"address": address,
		})
	}

	// it seems that we have ingest partial history
	if account == nil {
		account = new(history2.Account)
		*account, err = a.getAccountFromCore(rawAddress, address)
		if err != nil {
			return history2.Account{}, errors.Wrap(err, "failed to load account from core")
		}
	}

	a.accounts[*rawAddress.Ed25519] = account
	return *account, nil
}

func (a *Account) getAccountFromCore(rawAddress xdr.AccountId, address string) (history2.Account, error) {
	coreAccount, err := a.coreAccounts.GetByAddress(address)
	if err != nil {
		return history2.Account{}, errors.Wrap(err, "failed to load account by address from core", logan.F{
			"address": address,
		})
	}

	if coreAccount == nil {
		return history2.Account{}, errors.From(errors.New("account is not found in core db"), logan.F{"address": address})
	}

	account := history2.NewAccount(coreAccount.SequenceID, coreAccount.Address)
	err = a.InsertAccount(rawAddress, account)
	if err != nil {
		return history2.Account{}, errors.Wrap(err, "failed to insert account into db after fatched it from core")
	}

	return account, nil
}

// InsertAccount - adds account to local cache and inserts it into storage
func (a *Account) InsertAccount(rawAccountID xdr.AccountId, account history2.Account) error {
	// it's ok if the account already exists in the map.
	// Such case could occur during roll back of transaction and retry to process same ledger
	a.accounts[*rawAccountID.Ed25519] = &account
	err := a.repo.ExecRaw("INSERT INTO accounts (id, address) VALUES($1, $2)", account.ID, account.Address)
	if err != nil {
		return errors.Wrap(err, "failed to insert new account", logan.F{
			"address": account.Address,
			"id":      account.ID,
		})
	}

	return nil
}

func (a *Account) SetKYCRecoveryStatus(address string, status int) error {
	err := a.repo.ExecRaw("UPDATE accounts set kyc_recovery_status = ? where address = ?", status, address)
	if err != nil {
		return errors.Wrap(err, "failed to update account state", logan.F{
			"address":             address,
			"kyc_recovery_status": status,
		})
	}
	return nil
}

// MustAccountID returns int value which corresponds to xdr.AccountId
func (a *Account) MustAccountID(raw xdr.AccountId) uint64 {
	return a.MustAccount(raw).ID
}
