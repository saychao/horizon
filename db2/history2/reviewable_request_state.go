package history2

// ReviewableRequestState - represents state of reviewable request
type ReviewableRequestState int

const (
	// ReviewableRequestStatePending - request was just created or updated
	ReviewableRequestStatePending ReviewableRequestState = iota + 1
	// ReviewableRequestStateCanceled - was canceled by requestor
	ReviewableRequestStateCanceled
	// ReviewableRequestStateApproved - was approved by reviewer
	ReviewableRequestStateApproved
	// ReviewableRequestStateRejected - was rejected by reviewer, but still can be updated
	ReviewableRequestStateRejected
	// ReviewableRequestStatePermanentlyRejected - was rejected by reviewer, can't be updated
	ReviewableRequestStatePermanentlyRejected
)

var reviewableRequestStateStr = map[ReviewableRequestState]string{
	ReviewableRequestStatePending:             "pending",
	ReviewableRequestStateCanceled:            "canceled",
	ReviewableRequestStateApproved:            "approved",
	ReviewableRequestStateRejected:            "rejected",
	ReviewableRequestStatePermanentlyRejected: "permanently_rejected",
}

//String - converts int enum value to string
func (s ReviewableRequestState) String() string {
	return reviewableRequestStateStr[s]
}
