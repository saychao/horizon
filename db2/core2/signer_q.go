package core2

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

//SignerQ - helper struct to load signers from db
type SignerQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

//NewSignerQ - returns new instance of SignerQ with empty filter
func NewSignerQ(repo *pgdb.DB) SignerQ {
	return SignerQ{
		repo: repo,
		selector: sq.Select("signers.account_id",
			"signers.public_key",
			"signers.weight",
			"signers.role_id",
			"signers.identity",
			"signers.details",
		).From("signers"),
	}
}

func (q SignerQ) Count(address string) (int64, error) {
	q.selector = sq.Select("COUNT(*)").From("signers").Where("signers.account_id = ?", address)
	var result int64
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		return 0, errors.Wrap(err, "failed to load signers count")
	}

	return result, nil
}

//FilterByPublicKey - return new instance of SignerQ with filter by public key
func (q SignerQ) FilterByPublicKey(publicKey string) SignerQ {
	q.selector = q.selector.Where("signers.public_key = ?", publicKey)
	return q
}

//FilterByAccountAddress - return new instance of SignerQ with filter by account address
func (q SignerQ) FilterByAccountAddress(address string) SignerQ {
	q.selector = q.selector.Where("signers.account_id = ?", address)
	return q
}

//Select - selects slice of signers from db using specified filters. Returns nil, nil if none exists
func (q SignerQ) Select() ([]Signer, error) {
	var result []Signer
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to select signers")
	}

	return result, nil
}

// Page - returns Q with specified limit and offset params
func (q SignerQ) Page(params pgdb.OffsetPageParams) SignerQ {
	q.selector = params.ApplyTo(q.selector)
	return q
}
