package ingestion

import (
	sql2 "database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/saychao/horizon/db2"
	"gitlab.com/distributed_lab/kit/pgdb"

	sq "github.com/Masterminds/squirrel"
	"github.com/guregu/null"

	"github.com/saychao/horizon/db2/core"
	"github.com/saychao/horizon/db2/history"
	"github.com/saychao/horizon/db2/sqx"
	"github.com/saychao/horizon/ingest/participants"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

// Clear removes data from the ledger
func (ingest *Ingestion) Clear(start int64, end int64) error {
	clear := db2.DeleteRange
	err := clear(ingest.DB, start, end, "history_operation_participants", "history_operation_id")
	if err != nil {
		return errors.Wrap(err, "failed to clear history_operations_participants table")
	}
	err = clear(ingest.DB, start, end, "history_operations", "id")
	if err != nil {
		return errors.Wrap(err, "failed to clear history_operations table")
	}
	err = clear(ingest.DB, start, end, "history_transaction_participants", "history_transaction_id")
	if err != nil {
		return errors.Wrap(err, "failed to clear history_transactions_participants table")
	}
	err = clear(ingest.DB, start, end, "history_transactions", "id")
	if err != nil {
		return errors.Wrap(err, "failed to clear history_transactions table")
	}
	err = clear(ingest.DB, start, end, "history_ledgers", "id")
	if err != nil {
		return errors.Wrap(err, "failed to clear history_ledgers table")
	}
	err = clear(ingest.DB, start, end, "history_ledger_changes", "tx_id")
	if err != nil {
		return errors.Wrap(err, "failed to clear history_ledger_changes table")
	}
	err = clear(ingest.DB, start, end, "history_contracts", "id")
	if err != nil {
		return errors.Wrap(err, "failed to clear history_contract table")
	}
	err = clear(ingest.DB, start, end, "history_contracts_details", "id")
	if err != nil {
		return errors.Wrap(err, "failed to clear history_contracts_details table")
	}
	err = clear(ingest.DB, start, end, "history_contracts_disputes", "id")
	if err != nil {
		return errors.Wrap(err, "failed to clear history_contracts_disputes table")
	}

	return nil
}

// Ledger adds a ledger to the current ingestion
func (ingest *Ingestion) Ledger(
	id int64,
	header *core.LedgerHeader,
	txs int,
	ops int,
) error {

	sql := ingest.ledgers.Values(
		id,
		header.Sequence,
		header.LedgerHash,
		null.NewString(header.PrevHash, header.Sequence > 1),
		0, // TODO remove
		0, // TODO remove
		header.Data.BaseFee,
		header.Data.BaseReserve,
		header.Data.MaxTxSetSize,
		time.Unix(header.CloseTime, 0).UTC(),
		txs,
		ops,
	)

	err := ingest.DB.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to execute sql query")
	}

	return nil
}

func (ingest *Ingestion) LedgerChanges(
	txID, opID int64,
	orderNumber, effect int,
	entryType xdr.LedgerEntryType,
	payload interface{},
) error {
	xdrPayload, err := xdr.MarshalBase64(payload)
	if err != nil {
		return errors.Wrap(err, "failed to marshal payload")
	}

	sql := ingest.ledger_changes.Values(txID, opID, orderNumber, effect, entryType, xdrPayload)

	err = ingest.DB.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to execute sql query")
	}

	return nil
}

func (ingest *Ingestion) Contracts(contract history.Contract) error {
	if (int32(contract.State) & int32(xdr.ContractStateDisputing)) != 0 {
		return errors.New("Unexpected contract state. Expected not disputing")
	}

	sql := ingest.contracts.Values(
		contract.ID,
		contract.Contractor,
		contract.Customer,
		contract.Escrow,
		contract.StartTime,
		contract.EndTime,
		contract.InitialDetails,
		contract.CustomerDetails,
		contract.Invoices,
		contract.State,
	)

	err := ingest.DB.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to execute sql query")
	}

	return nil
}

func (ingest *Ingestion) ContractDetails(contractDetails history.ContractDetails) error {
	sql := ingest.contractsDetails.Values(
		contractDetails.ContractID,
		contractDetails.Details,
		contractDetails.Author,
		contractDetails.CreatedAt,
	)

	err := ingest.DB.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to execute sql query")
	}

	return nil
}

func (ingest *Ingestion) ContractDispute(contractDetails history.ContractDispute) error {
	sql := ingest.contractsDisputes.Values(
		contractDetails.ContractID,
		contractDetails.Reason,
		contractDetails.Author,
		contractDetails.CreatedAt,
	)

	err := ingest.DB.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to execute sql query")
	}

	return nil
}

// Operation ingests the provided operation data into a new row in the
// `history_operations` table
func (ingest *Ingestion) Operation(
	id int64,
	txid int64,
	order int32,
	source xdr.AccountId,
	typ xdr.OperationType,
	details map[string]interface{},
	ledgerCloseTime int64,
	identifier uint64,
	state history.OperationState,
) error {
	djson, err := json.Marshal(details)
	if err != nil {
		return errors.Wrap(err, "failed to marshal details")
	}

	sql := ingest.operations.Values(id,
		txid, order, source.Address(),
		typ, djson, time.Unix(ledgerCloseTime, 0).UTC(),
		identifier, state)
	err = ingest.DB.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to execute sql query")
	}

	return nil
}

// OperationParticipants ingests the provided accounts `aids` as participants of
// operation with id `op`, creating a new row in the
// `history_operation_participants` table.
func (ingest *Ingestion) OperationParticipants(op int64, opParticipants []participants.Participant) error {
	sql := ingest.operation_participants
	for _, opParticipant := range opParticipants {
		var djson *[]byte
		var err error
		if opParticipant.Details != nil {
			marshalledDetails, err := json.Marshal(opParticipant.Details)
			djson = &marshalledDetails
			if err != nil {
				return errors.Wrap(err, "failed to marshal operation participant details", logan.F{
					"operation id": op,
				})
			}
		}
		haid, err := ingest.getParticipantID(opParticipant.AccountID)
		if err != nil {
			return errors.Wrap(err, "failed to get operation participant", map[string]interface{}{
				"operation id": op,
			})
		}
		sql = sql.Values(op, haid, opParticipant.BalanceID.AsString(), djson)
	}

	err := ingest.DB.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to execute sql query")
	}

	return nil
}

func (ingest *Ingestion) TransactWithFunction(fn pgdb.TransactionFunc) error {
	return ingest.DB.Transaction(fn)
}

// Start makes the ingestion ready, initializing the insert builders and tx
func (ingest *Ingestion) Start() (err error) {
	ingest.createInsertBuilders()

	return nil
}

// Transaction ingests the provided transaction data into a new row in the
// `history_transactions` table
func (ingest *Ingestion) Transaction(
	ledger *core.LedgerHeader,
	id int64,
	tx *core.Transaction,
	fee *core.TransactionFee,
) error {

	sql := ingest.transactions.Values(
		id,
		tx.TransactionHash,
		tx.LedgerSequence,
		tx.Index,
		tx.SourceAddress(),
		tx.Salt(),
		tx.Fee(),
		len(tx.Envelope.Tx.Operations),
		tx.EnvelopeXDR(),
		tx.ResultXDR(),
		tx.ResultMetaXDR(),
		fee.ChangesXDR(),
		sqx.StringArray(tx.Base64Signatures()),
		ingest.formatTimeBounds(tx.Envelope.Tx.TimeBounds),
		tx.MemoType(),
		tx.Memo(),
		time.Unix(ledger.CloseTime, 0).UTC(),
	)

	err := ingest.DB.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to add new row into history_transactions table")
	}

	return nil
}

// TransactionParticipants ingests the provided account ids as participants of
// transaction with id `tx`, creating a new row in the
// `history_transaction_participants` table.
func (ingest *Ingestion) TransactionParticipants(tx int64, aids []xdr.AccountId) error {
	sql := ingest.transaction_participants

	for _, aid := range aids {
		haid, err := ingest.getParticipantID(aid)
		if err != nil {
			return errors.Wrap(err, "failed to get participant ID")
		}
		sql = sql.Values(tx, haid)
	}

	err := ingest.DB.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failes to exeute query on history_transactions_participants")
	}

	return nil
}

func (ingest *Ingestion) createInsertBuilders() {
	ingest.ledgers = sq.Insert("history_ledgers").Columns(
		"id",
		"sequence",
		"ledger_hash",
		"previous_ledger_hash",
		"total_coins",
		"fee_pool",
		"base_fee",
		"base_reserve",
		"max_tx_set_size",
		"closed_at",
		"transaction_count",
		"operation_count",
	)

	ingest.transactions = sq.Insert("history_transactions").Columns(
		"id",
		"transaction_hash",
		"ledger_sequence",
		"application_order",
		"account",
		"salt",
		"fee_paid",
		"operation_count",
		"tx_envelope",
		"tx_result",
		"tx_meta",
		"tx_fee_meta",
		"signatures",
		"time_bounds",
		"memo_type",
		"memo",
		"ledger_close_time",
	)

	ingest.transaction_participants = sq.Insert("history_transaction_participants").Columns(
		"history_transaction_id",
		"history_account_id",
	)

	ingest.operations = sq.Insert("history_operations").Columns(
		"id",
		"transaction_id",
		"application_order",
		"source_account",
		"type",
		"details",
		"ledger_close_time",
		"identifier",
		"state",
	)

	ingest.operation_participants = sq.Insert("history_operation_participants").Columns(
		"history_operation_id",
		"history_account_id",
		"balance_id",
		"effects",
	)

	ingest.trades = sq.Insert("history_trades").Columns(
		"order_book_id",
		"base_asset",
		"quote_asset",
		"base_amount",
		"quote_amount",
		"price",
		"created_at",
	)

	ingest.priceHistory = sq.Insert("history_price").Columns("base_asset", "quote_asset", "timestamp", "price")

	ingest.ledger_changes = sq.Insert("history_ledger_changes").Columns(
		"tx_id",
		"op_id",
		"order_number",
		"effect",
		"entry_type",
		"payload",
	)

	ingest.contracts = sq.Insert("history_contracts").Columns("id",
		"contractor",
		"customer",
		"escrow",
		"start_time",
		"end_time",
		"initial_details",
		"customer_details",
		"invoices",
		"state",
	)

	ingest.contractsDetails = sq.Insert("history_contracts_details").Columns(
		"contract_id",
		"details",
		"author",
		"created_at",
	)

	ingest.contractsDisputes = sq.Insert("history_contracts_disputes").Columns(
		"contract_id",
		"reason",
		"author",
		"created_at",
	)
}

func (ingest *Ingestion) TryIngestAccount(aid string) (result int64, err error) {
	q := history.Q{DB: ingest.DB}
	var existing history.Account
	err = q.AccountByAddress(&existing, aid)

	if err != nil && err != sql2.ErrNoRows {
		return 0, errors.Wrap(err, "failed to get account from DB")
	}

	// already imported, return the found value
	if err != sql2.ErrNoRows {
		result = existing.ID
		return result, nil
	}

	coreQ := core.Q{DB: ingest.CoreDB}
	account, err := coreQ.Accounts().ByAddress(aid)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get core balance")
	}

	err = ingest.DB.GetRaw(
		&result,
		`INSERT INTO history_accounts (address, account_type) VALUES (?, ?) RETURNING id`,
		aid, account.RoleID,
	)
	if err != nil {
		return 0, errors.Wrap(err, "failed to execute insert query on history_accounts table")
	}

	return result, nil
}

func (ingest *Ingestion) getParticipantID(
	aid xdr.AccountId,
) (id int64, err error) {
	id, err = ingest.TryIngestAccount(aid.Address())
	if err != nil {
		return 0, errors.Wrap(err, "failed to ingest account")
	}
	return id, err
}

func (ingest *Ingestion) formatTimeBounds(bounds xdr.TimeBounds) interface{} {
	return sq.Expr("?::int8range", fmt.Sprintf("[%d,%d]", bounds.MinTime, bounds.MaxTime))
}
