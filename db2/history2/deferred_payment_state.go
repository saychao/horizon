package history2

// ReviewableRequestState - represents state of reviewable request
type DeferredPaymentState int

const (
	// DeferredPaymentStateOpen - deferred payment is active and can be used
	DeferredPaymentStateOpen DeferredPaymentState = iota + 1
	// DeferredPaymentStateClosed - all amount from request was spent
	DeferredPaymentStateClosed
)

var deferredPaymentStateStr = map[DeferredPaymentState]string{
	DeferredPaymentStateOpen:   "open",
	DeferredPaymentStateClosed: "closed",
}

//String - converts int enum value to string
func (s DeferredPaymentState) String() string {
	return deferredPaymentStateStr[s]
}
