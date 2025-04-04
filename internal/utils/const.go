package utils

import (
	"time"
)

const (
	User      = "user"
	PinVerify = "pin_verify"
)
const ClientID = "client_id"

const (
	Authorization = "Authorization"
)

const (
	TableAssetAuditLogName              = "asset_audit_log"
	TableAssetCategoryName              = "asset_category"
	TableAssetMaintenanceRecordName     = "asset_maintenance_record"
	TableAssetMaintenanceName           = "asset_maintenance"
	TableAssetMaintenanceTypeName       = "asset_maintenance_type"
	TableAssetName                      = "asset"
	TableAssetStatusName                = "asset_status"
	TableAssetImageName                 = "asset_image"
	TableAssetStockName                 = "asset_stock"
	TableAssetStockHistoryName          = "asset_stock_history"
	TableAssetGroupName                 = "asset_group"
	TableAssetGroupPermissionName       = "asset_group_permission"
	TableAssetGroupMemberName           = "asset_group_member"
	TableAssetGroupMemberPermissionName = "asset_group_member_permission"
	TableAssetGroupAssetName            = "asset_group_asset"
)

const (
	NatsAssetImageDelete = "asset.image.delete"
	NatsAssetImageUsage  = "asset.image.usage"
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
