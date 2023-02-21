package horizon

import (
	"github.com/saychao/horizon/db2/core"
	"github.com/saychao/horizon/render/hal"
	"github.com/saychao/horizon/render/problem"
	"github.com/saychao/horizon/resource"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type SaleAnteAction struct {
	Action
	q                    core.SaleAnteQI
	SaleID               string
	ParticipantBalanceID string
	Records              []core.SaleAnte
	Page                 hal.Page
}

func (action *SaleAnteAction) JSON() {
	action.Do(
		action.loadParams,
		action.checkAllowed,
		action.loadRecords,
		action.loadPage,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

func (action *SaleAnteAction) loadParams() {
	action.SaleID = action.GetString("sale_id")
	action.ParticipantBalanceID = action.GetString("participant_balance_id")
	action.Page.Filters = map[string]string{
		"sale_id":                action.SaleID,
		"participant_balance_id": action.ParticipantBalanceID,
	}
}

func (action *SaleAnteAction) checkAllowed() {
	if action.ParticipantBalanceID == "" {
		action.IsAllowed("")
		return
	}

	participantBalance, err := action.CoreQ().Balances().ByID(action.ParticipantBalanceID)
	if err != nil {
		action.Log.WithError(err).Error("failed to get sale ante participant balance from core DB")
		action.Err = &problem.ServerError
		return
	}

	if participantBalance == nil {
		action.SetInvalidField("participant_balance_id", errors.New("sale ante participant balance does not exist in core DB"))
		return
	}

	action.IsAllowed(participantBalance.AccountID)
}

func (action *SaleAnteAction) loadRecords() {
	action.q = action.CoreQ().SaleAntes()

	if action.SaleID != "" {
		action.q = action.q.ForSale(action.SaleID)
	}

	if action.ParticipantBalanceID != "" {
		action.q = action.q.ForBalance(action.ParticipantBalanceID)
	}

	if action.Err != nil {
		return
	}

	var err error

	action.Records, err = action.q.Select()
	if err != nil {
		action.Log.WithError(err).Error("failed to get sale antes from core DB")
		action.Err = &problem.ServerError
		return
	}
}

func (action *SaleAnteAction) loadPage() {
	for _, saleAnte := range action.Records {
		participantBalance, err := action.CoreQ().Balances().ByID(saleAnte.ParticipantBalanceID)
		if err != nil || participantBalance == nil {
			action.Log.WithError(err).WithFields(logan.F{
				"sale_id":                saleAnte.SaleID,
				"participant_balance_id": saleAnte.ParticipantBalanceID,
			}).Error("failed to get participant balance for sale ante from core DB")
			action.Err = &problem.ServerError
			return
		}
		var res resource.SaleAnte
		res.Populate(&saleAnte, participantBalance.Asset)
		action.Page.Add(&res)
	}

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.PopulateLinks()
}
