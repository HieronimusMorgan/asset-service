package utils

import "time"

const Token = "token"
const User = "user"
const UserSession = "user_session"
const ClientID = "client_id"
const UserID = "user_id"
const RoleID = "role_id"

const (
	Admin         = "Admin"
	SuperAdmin    = "Super Admin"
	NormalUser    = "User"
	Anonymous     = "Anonymous"
	Authorization = "Authorization"
)

const (
	TableAssetAuditLogName          = "my-home.asset_audit_log"
	TableAssetCategoryName          = "my-home.asset_category"
	TableAssetMaintenanceRecordName = "my-home.asset_maintenance_record"
	TableAssetMaintenanceName       = "my-home.asset_maintenance"
	TableAssetMaintenanceTypeName   = "my-home.asset_maintenance_type"
	TableAssetName                  = "my-home.asset"
	TableAssetStatus                = "my-home.asset_status"
)

func ParseOptionalDate(str *string) (*time.Time, error) {
	if str == nil {
		return nil, nil
	}
	parsedDate, err := time.Parse("2006-01-02", *str)
	if err != nil {
		return nil, err
	}
	return &parsedDate, nil
}
