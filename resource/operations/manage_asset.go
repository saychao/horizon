package operations

type ManageAsset struct {
	Base
	RequestID    uint64 `json:"request_id"`
	Action       int32  `json:"action"`
	ActionString string `json:"action_string"`
}
