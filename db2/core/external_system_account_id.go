package core

// ExternalSystemAccountID - represents id generated for user in external system like Bitcoin/Ethereum
type ExternalSystemAccountID struct {
	// Address - account id of the owner of ID
	AccountID string `db:"account_id"`
	// ExternalSystemType - type of the external system Bitcoin/Ethereum
	ExternalSystemType int32 `db:"external_system_type"`
	// Data - ID from external system
	Data string `db:"data"`
	// ExpiresAt - expire time of external system account ID
	ExpiresAt *int64 `db:"pool_entry_expires_at"`
}
