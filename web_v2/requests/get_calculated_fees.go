package requests

import (
	"net/http"

	. "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/logan/v3/errors"
	addr "gitlab.com/tokend/go/address"
	amount2 "gitlab.com/tokend/go/amount"
	regources "gitlab.com/tokend/regources/generated"
)

//GetCalculatedFees - represents params to be specified for Get Fees handler
type GetCalculatedFees struct {
	*base
	Address string
	Asset   string
	Subtype int64
	FeeType int32
	Amount  regources.Amount
}

// NewGetCalculatedFees returns the new instance of GetCalculatedFees request
func NewGetCalculatedFees(r *http.Request) (*GetCalculatedFees, error) {
	b, err := newBase(r, baseOpts{})
	if err != nil {
		return nil, err
	}

	return makeCalculatedFees(b)
}

func makeCalculatedFees(b *base) (*GetCalculatedFees, error) {
	address := b.getString("id")
	asset := b.getString("asset")
	amountRaw := b.getString("amount")

	subtype, err := b.getInt64("subtype")
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse subtype")
	}

	feeType, err := b.getInt32("fee_type")
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse fee_type")
	}

	data := struct {
		address addr.Address
		asset   string
		amount  string
		subtype int64
		feeType int32 `json:"fee_type"`
	}{
		address: addr.Address(address),
		asset:   asset,
		amount:  amountRaw,
		subtype: subtype,
		feeType: feeType,
	}

	err = ValidateStruct(&data,
		Field(&data.address, Required, addr.IsAddress),
		Field(&data.asset, Required),
		Field(&data.amount, Required, &isNonNegAmount{}),
	)
	if err != nil {
		return nil, err
	}

	return &GetCalculatedFees{
		base:    b,
		Address: data.address.String(),
		Amount:  regources.Amount(amount2.MustParseU(data.amount)),
		Asset:   data.asset,
		FeeType: data.feeType,
		Subtype: data.subtype,
	}, nil
}

type isNonNegAmount struct{}

func (ia *isNonNegAmount) Validate(value interface{}) error {
	strVal, ok := value.(string)
	if !ok {
		return Errors{"amount": errors.New("amount is not string")}
	}

	a, err := amount2.Parse(strVal)
	if err != nil {
		return Errors{"amount": err}
	}

	if a < 0 {
		return Errors{"amount": errors.New("must be >= 0")}
	}

	return nil
}
