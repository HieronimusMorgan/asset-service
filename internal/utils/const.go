package utils

const (
	User          = "user"
	PinVerify     = "pin_verify"
	CredentialKey = "credential_key"
	PageIndex     = "page_index"
	PageSize      = "page_size"
)

const (
	Authorization  = "Authorization"
	XCredentialKey = "X-CREDENTIAL-KEY"
	ClientID       = "client_id"
)

const (
	TableAssetAuditLogName              = "asset_audit_log"
	TableAssetCategoryName              = "asset_category"
	TableAssetMaintenanceRecordName     = "asset_maintenance_record"
	TableAssetMaintenanceName           = "asset_maintenance"
	TableAssetMaintenanceTypeName       = "asset_maintenance_type"
	TableAssetName                      = "asset"
	TableAssetWishlistName              = "asset_wishlist"
	TableAssetStatusName                = "asset_status"
	TableAssetImageName                 = "asset_image"
	TableAssetStockName                 = "asset_stock"
	TableAssetStockHistoryName          = "asset_stock_history"
	TableAssetGroupName                 = "asset_group"
	TableAssetGroupPermissionName       = "asset_group_permission"
	TableAssetGroupMemberName           = "asset_group_member"
	TableAssetGroupMemberPermissionName = "asset_group_member_permission"
	TableAssetGroupAssetName            = "asset_group_asset"
	TableAssetGroupInvitationName       = "asset_group_invitation"

	TableUserSettingName = "user_settings"
)

const (
	InvitationStatusPending  = "pending"
	InvitationStatusAccepted = "accepted"
	InvitationStatusRejected = "rejected"
	InvitationStatusExpired  = "expired"
)

const (
	NatsAssetImageDelete = "asset.image.delete"
	NatsAssetImageUsage  = "asset.image.usage"
)
