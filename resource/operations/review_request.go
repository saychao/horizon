package operations

type ReviewRequest struct {
	Base
	Action            string                 `json:"action"`
	Reason            string                 `json:"reason"`
	RequestHash       string                 `json:"request_hash"`
	RequestID         uint64                 `json:"request_id"`
	RequestType       string                 `json:"request_type"`
	IsFulfilled       bool                   `json:"is_fulfilled"`
	Details           map[string]interface{} `json:"details"`
	AtomicSwapDetails map[string]interface{} `json:"atomic_swap_details"`
}
