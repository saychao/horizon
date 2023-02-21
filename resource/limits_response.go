package resource

type LimitResponse struct {
	Limit      LimitsV2     `json:"limit"`
	Statistics StatisticsV2 `json:"statistics"`
}

type LimitsResponse struct {
	Limits []LimitResponse `json:"limits"`
}
