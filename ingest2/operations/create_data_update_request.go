package operations

import (
	"github.com/saychao/horizon/db2/history2"
	"github.com/saychao/horizon/ingest2/internal"
	"gitlab.com/tokend/go/xdr"
)

type createDataUpdateRequestHandler struct {
	effectsProvider
}

func (h *createDataUpdateRequestHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	oper := op.Body.MustCreateDataUpdateRequestOp()
	response := opRes.MustCreateDataUpdateRequestResult().Success

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateDataUpdateRequest,
		CreateDataUpdateRequest: &history2.CreateDataUpdateRequest{
			ID:             uint64(oper.DataUpdateRequest.Id),
			Value:          internal.MarshalCustomDetails(oper.DataUpdateRequest.Value),
			RequestID:      uint64(response.RequestId),
			CreatorDetails: internal.MarshalCustomDetails(oper.DataUpdateRequest.CreatorDetails),
		},
	}, nil
}

func (h *createDataUpdateRequestHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	participants := []history2.ParticipantEffect{h.Participant(sourceAccountID)}

	return participants, nil
}
