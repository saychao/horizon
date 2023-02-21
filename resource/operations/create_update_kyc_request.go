package operations

type CreateChangeRoleRequest struct {
	Base
	RequestID          uint64                 `json:"request_id"`
	DestinationAccount string                 `json:"destination_account"`
	AccountRoleToSet   int32                  `json:"account_role_to_set"`
	KYCData            map[string]interface{} `json:"kyc_data"`
	AllTasks           *uint32                `json:"all_tasks"`
}
