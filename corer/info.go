package corer

import (
	"fmt"

	"gitlab.com/tokend/go/xdr"
)

// Info -- response for /info request
type Info struct {
	// CoreVersion - version of the core
	CoreVersion string `json:"build"`
	// NetworkPassphrase - passphrase of the network
	NetworkPassphrase string `json:"network"`
	// MasterExchangeName - name of the exchange managed by master key
	MasterExchangeName string `json:"base_exchange_name"`
	// TxExpirationPeriod - max allowed period for tx time bounds max
	TxExpirationPeriod int64 `json:"tx_expiration_period"`
	// WithdrawalDetailsMaxLength - max length of details field for withdrawal operation
	WithdrawalDetailsMaxLength int64 `json:"withdrawal_details_max_length"`
	// Array of the base assets
	BaseAssets []string `json:"base_assets"`

	// MasterAccountID - account ID of master
	AdminAccountID string `json:"admin_account_id"`

	// MasterAccountIDXDR - masterAccountID parsed into xdr.AccountID
	MasterAccountIDXDR xdr.AccountId `json:"-"`
	// CoreURL - url of the core
	CoreURL string `json:"-"`
}

type infoResponse struct {
	Info Info `json:"info"`
}

func (i *Info) validate() error {
	errorProvider := func(name string) error {
		return fmt.Errorf("%s must not be empty. Please check connection with stellar-core", name)
	}
	if i.NetworkPassphrase == "" {
		return errorProvider("NetworkPassphrase")
	}

	if i.AdminAccountID == "" {
		return errorProvider("AdminAccountID")
	}

	err := i.MasterAccountIDXDR.SetAddress(i.AdminAccountID)
	if err != nil {
		return errorProvider("MasterAccountID is invalid")
	}

	if i.TxExpirationPeriod <= 0 {
		return errorProvider("TxExpirationPeriod")
	}

	if i.WithdrawalDetailsMaxLength <= 0 {
		return errorProvider("WithdrawalDetailsMaxLength")
	}
	return nil
}
