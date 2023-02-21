package resource

import (
	"time"

	"github.com/SafeRE-IT/horizon/ledger"
	"gitlab.com/tokend/go/amount"
	"golang.org/x/net/context"
)

// Root is the initial map of links into the api.
type Root struct {
	LedgersState       ledger.SystemState `json:"ledgers_state"`
	NetworkPassphrase  string             `json:"network_passphrase"`
	AdminAccountID     string             `json:"admin_account_id"`
	MasterExchangeName string             `json:"master_exchange_name"`
	TxExpirationPeriod int64              `json:"tx_expiration_period"`
	CurrentTime        int64              `json:"current_time"`
	Precision          int64              `json:"precision"`
	XDRRevision        string             `json:"xdr_revision"`
	HorizonRevision    string             `json:"horizon_revision"`
	MasterAccountID    string             `json:"master_account_id"`
	EnvironmentName    string             `json:"environment_name"`
}

// Populate fills in the details
func (res *Root) PopulateLedgerState(
	ctx context.Context,
	ledgerState ledger.SystemState,
) {
	res.LedgersState = ledgerState
	res.CurrentTime = time.Now().Unix()
	res.Precision = amount.One
}
