package changes

import (
	history "github.com/saychao/horizon/db2/history2"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

type accountSpecificRuleStorage interface {
	//Inserts rule into DB
	Insert(rule history.AccountSpecificRule) error

	//Removes rule from DB
	Remove(ruleID uint64) error
}

type accountSpecificRuleHandler struct {
	storage accountSpecificRuleStorage
}

func (h accountSpecificRuleHandler) Removed(change ledgerChange) error {
	op := change.Operation

	switch op.Body.Type {
	case xdr.OperationTypeManageAccountSpecificRule:
		body := op.Body.MustManageAccountSpecificRuleOp()
		switch body.Data.Action {
		case xdr.ManageAccountSpecificRuleActionRemove:
			id := uint64(op.Body.MustManageAccountSpecificRuleOp().Data.MustRemoveData().RuleId)
			err := h.storage.Remove(id)
			if err != nil {
				return errors.Wrap(err, "failed to remove account specific rule", logan.F{
					"rule_id": id,
				})
			}
		default:
			return errors.From(errors.New("Unexpected manage account specific rule action"), logan.F{
				"action": body.Data.Action,
			})
		}

	case xdr.OperationTypeManageSale:
	case xdr.OperationTypeCheckSaleState:
	default:
		return errors.From(errors.New("Unexpected operation type"), logan.F{
			"operation_type": op.Body.Type,
		})
	}

	return nil
}

func (h accountSpecificRuleHandler) Created(change ledgerChange) error {
	op := change.Operation

	switch op.Body.Type {
	case xdr.OperationTypeCreateSaleRequest:
		rawRule := change.LedgerChange.Created.Data.MustAccountSpecificRule()
		rule := history.NewAccountSpecificRule(rawRule)
		err := h.storage.Insert(rule)
		if err != nil {
			return errors.Wrap(err, "failed to insert account specific rule", logan.F{
				"rule_id": rule.ID,
			})
		}
	case xdr.OperationTypeReviewRequest:
		reviewRequestOp := op.Body.MustReviewRequestOp()
		switch reviewRequestOp.Action {
		case xdr.ReviewRequestOpActionApprove:
			rawRule := change.LedgerChange.Created.Data.MustAccountSpecificRule()
			rule := history.NewAccountSpecificRule(rawRule)
			err := h.storage.Insert(rule)
			if err != nil {
				return errors.Wrap(err, "failed to insert account specific rule", logan.F{
					"rule_id": rule.ID,
				})
			}
		default:
			return errors.From(errors.New("Unexpected review request action"), logan.F{
				"action": reviewRequestOp.Action,
			})
		}
	case xdr.OperationTypeManageAccountSpecificRule:
		body := op.Body.MustManageAccountSpecificRuleOp()
		switch body.Data.Action {
		case xdr.ManageAccountSpecificRuleActionCreate:
			rawRule := change.LedgerChange.Created.Data.MustAccountSpecificRule()
			rule := history.NewAccountSpecificRule(rawRule)
			err := h.storage.Insert(rule)
			if err != nil {
				return errors.Wrap(err, "failed to insert account specific rule", logan.F{
					"rule_id": rule.ID,
				})
			}
		default:
			return errors.From(errors.New("Unexpected manage account specific rule action"), logan.F{
				"action": body.Data.Action,
			})
		}
	default:
		return errors.From(errors.New("Unexpected operation type"), logan.F{
			"operation_type": op.Body.Type,
		})
	}

	return nil
}

func newAccountSpecificRuleHandler(storage accountSpecificRuleStorage) *accountSpecificRuleHandler {
	return &accountSpecificRuleHandler{
		storage: storage,
	}
}
