package utils

import (
	"time"
)

const User = "user"
const ClientID = "client_id"

const (
	Authorization = "Authorization"
)

const (
	TableAssetAuditLogName          = "my-home.asset_audit_log"
	TableAssetCategoryName          = "my-home.asset_category"
	TableAssetMaintenanceRecordName = "my-home.asset_maintenance_record"
	TableAssetMaintenanceName       = "my-home.asset_maintenance"
	TableAssetMaintenanceTypeName   = "my-home.asset_maintenance_type"
	TableAssetName                  = "my-home.asset"
	TableAssetStatusName            = "my-home.asset_status"
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

func CalculateNextDueDate(date *time.Time, days *int) (*time.Time, error) {
	if date == nil || days == nil {
		return nil, nil
	}
	nextDueDate := date.AddDate(0, 0, *days)
	return &nextDueDate, nil
}
