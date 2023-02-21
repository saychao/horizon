package history

type SaleState int

const (
	SaleStateOpen SaleState = 1 << iota
	SaleStateClosed
	SaleStateCanceled
	SaleStatePromotion
	SaleStateVoting
)

var saleStateStr = map[SaleState]string{
	SaleStateOpen:      "open",
	SaleStateClosed:    "closed",
	SaleStateCanceled:  "canceled",
	SaleStatePromotion: "promotion",
	SaleStateVoting:    "voting",
}

func (s SaleState) String() string {
	return saleStateStr[s]
}
