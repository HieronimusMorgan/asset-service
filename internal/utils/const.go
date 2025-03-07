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
	TableAssetAuditLogName          = "asset-service.asset_audit_log"
	TableAssetCategoryName          = "asset-service.asset_category"
	TableAssetMaintenanceRecordName = "asset-service.asset_maintenance_record"
	TableAssetMaintenanceName       = "asset-service.asset_maintenance"
	TableAssetMaintenanceTypeName   = "asset-service.asset_maintenance_type"
	TableAssetName                  = "asset-service.asset"
	TableAssetStatusName            = "asset-service.asset_status"
	TableAssetImageName             = "asset-service.asset_image"
	TableAssetStockName             = "asset-service.asset_stock"
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
