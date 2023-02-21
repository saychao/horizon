package history

type OperationState int32

const (
	OperationStatePending OperationState = 1 + iota
	OperationStateSuccess
	OperationStateRejected
	OperationStateCanceled
	OperationStateFailed
	OperationStatePartiallyMatched
	OperationStateFullyMatched
	OperationStateExternallyFullyMatched
)

var operationStateStr = map[OperationState]string{
	OperationStatePending:                "pending",
	OperationStateSuccess:                "success",
	OperationStateRejected:               "rejected",
	OperationStateCanceled:               "canceled",
	OperationStateFailed:                 "failed",
	OperationStatePartiallyMatched:       "partially matched",
	OperationStateFullyMatched:           "fully matched",
	OperationStateExternallyFullyMatched: "externally_fully_matched",
}

func (s OperationState) String() string {
	return operationStateStr[s]
}
