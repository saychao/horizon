package operations

type ReviewForfeitRequest struct {
	Base
	RequestID uint64 `json:"request_id"`
	Balance   string `json:"balance"`
	Accept    bool   `json:"accept"`
	Amount    string `json:"amount"`
	Asset     string `json:"asset"`
}
