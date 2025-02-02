package repository

import (
	"asset-service/internal/models/assets"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

const tableAssetAuditLogName = "my-home.asset_audit_log"

type AssetAuditLogRepository struct {
	DB *gorm.DB
}

func NewAssetAuditLogRepository(db *gorm.DB) *AssetAuditLogRepository {
	return &AssetAuditLogRepository{DB: db}
}

func (a AssetAuditLogRepository) AfterCreateAsset(asset *assets.Asset) (err error) {
	newDataBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	newData := string(newDataBytes)

	log := assets.AssetAuditLog{
		TableName:   "asset",
		Action:      "CREATE",
		NewData:     &newData,
		PerformedAt: time.Now(),
		PerformedBy: &asset.CreatedBy,
	}

	if err := a.DB.Table(tableAssetAuditLogName).Create(&log).Error; err != nil {
		return err
	}
	return nil
}

func (a AssetAuditLogRepository) AfterUpdateAsset(old assets.Asset, asset *assets.Asset) error {
	oldDataBytes, err := json.Marshal(old)
	if err != nil {
		return err
	}
	oldData := string(oldDataBytes)

	newDataBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	newData := string(newDataBytes)

	log := assets.AssetAuditLog{
		TableName:   "asset",
		Action:      "UPDATE",
		OldData:     &oldData,
		NewData:     &newData,
		PerformedAt: time.Now(),
		PerformedBy: &asset.UpdatedBy,
	}

	if err := a.DB.Table(tableAssetAuditLogName).Create(&log).Error; err != nil {
		return err
	}
	return nil
}

func (a AssetAuditLogRepository) AfterDeleteAsset(asset *assets.Asset) error {
	oldDataBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	oldData := string(oldDataBytes)

	log := assets.AssetAuditLog{
		TableName:   "asset",
		Action:      "DELETE",
		OldData:     &oldData,
		PerformedAt: time.Now(),
		PerformedBy: asset.DeletedBy,
	}

	if err := a.DB.Table(tableAssetAuditLogName).Create(&log).Error; err != nil {
		return err
	}
	return nil
}

func (a AssetAuditLogRepository) AfterCreateAssetMaintenance(assetMaintenance *assets.AssetMaintenance) error {
	newDataBytes, err := json.Marshal(assetMaintenance)
	if err != nil {
		return err
	}
	newData := string(newDataBytes)

	log := assets.AssetAuditLog{
		TableName:   "asset_maintenance",
		Action:      "CREATE",
		NewData:     &newData,
		PerformedAt: time.Now(),
		PerformedBy: &assetMaintenance.CreatedBy,
	}

	if err := a.DB.Table(tableAssetAuditLogName).Create(&log).Error; err != nil {
		return err
	}
	return nil
}

func (a AssetAuditLogRepository) AfterUpdateAssetMaintenance(old assets.AssetMaintenance, assetMaintenance *assets.AssetMaintenance) error {
	oldDataBytes, err := json.Marshal(old)
	if err != nil {
		return err
	}
	oldData := string(oldDataBytes)

	newDataBytes, err := json.Marshal(assetMaintenance)
	if err != nil {
		return err
	}
	newData := string(newDataBytes)

	log := assets.AssetAuditLog{
		TableName:   "asset_maintenance",
		Action:      "UPDATE",
		OldData:     &oldData,
		NewData:     &newData,
		PerformedAt: time.Now(),
		PerformedBy: &assetMaintenance.UpdatedBy,
	}

	if err := a.DB.Table(tableAssetAuditLogName).Create(&log).Error; err != nil {
		return err
	}
	return nil
}

func (a AssetAuditLogRepository) AfterDeleteAssetMaintenance(assetMaintenance *assets.AssetMaintenance) error {
	oldDataBytes, err := json.Marshal(assetMaintenance)
	if err != nil {
		return err
	}
	oldData := string(oldDataBytes)

	log := assets.AssetAuditLog{
		TableName:   "asset_maintenance",
		Action:      "DELETE",
		OldData:     &oldData,
		PerformedAt: time.Now(),
		PerformedBy: assetMaintenance.DeletedBy,
	}

	if err := a.DB.Table(tableAssetAuditLogName).Create(&log).Error; err != nil {
		return err
	}
	return nil
}

func (a AssetAuditLogRepository) AfterCreateAssetCategory(assetCategory *assets.AssetCategory) error {
	newDataBytes, err := json.Marshal(assetCategory)
	if err != nil {
		return err
	}
	newData := string(newDataBytes)

	log := assets.AssetAuditLog{
		TableName:   "asset_category",
		Action:      "CREATE",
		NewData:     &newData,
		PerformedAt: time.Now(),
		PerformedBy: &assetCategory.CreatedBy,
	}

	if err := a.DB.Table(tableAssetAuditLogName).Create(&log).Error; err != nil {
		return err
	}
	return nil
}

func (a AssetAuditLogRepository) AfterUpdateAssetCategory(old *assets.AssetCategory, assetCategory *assets.AssetCategory) error {
	oldDataBytes, err := json.Marshal(old)
	if err != nil {
		return err
	}
	oldData := string(oldDataBytes)

	newDataBytes, err := json.Marshal(assetCategory)
	if err != nil {
		return err
	}
	newData := string(newDataBytes)

	log := assets.AssetAuditLog{
		TableName:   "asset_category",
		Action:      "UPDATE",
		OldData:     &oldData,
		NewData:     &newData,
		PerformedAt: time.Now(),
		PerformedBy: &assetCategory.UpdatedBy,
	}

	if err := a.DB.Table(tableAssetAuditLogName).Create(&log).Error; err != nil {
		return err
	}
	return nil
}

func (a AssetAuditLogRepository) AfterDeleteAssetCategory(assetCategory *assets.AssetCategory) error {
	oldDataBytes, err := json.Marshal(assetCategory)
	if err != nil {
		return err
	}
	oldData := string(oldDataBytes)

	log := assets.AssetAuditLog{
		TableName:   "asset_category",
		Action:      "DELETE",
		OldData:     &oldData,
		PerformedAt: time.Now(),
		PerformedBy: assetCategory.DeletedBy,
	}

	if err := a.DB.Table(tableAssetAuditLogName).Create(&log).Error; err != nil {
		return err
	}
	return nil
}

func (a AssetAuditLogRepository) AfterCreateAssetStatus(assetStatus *assets.AssetStatus) error {
	newDataBytes, err := json.Marshal(assetStatus)
	if err != nil {
		return err
	}
	newData := string(newDataBytes)

	log := assets.AssetAuditLog{
		TableName:   "asset_status",
		Action:      "CREATE",
		NewData:     &newData,
		PerformedAt: time.Now(),
		PerformedBy: &assetStatus.CreatedBy,
	}

	if err := a.DB.Table(tableAssetAuditLogName).Create(&log).Error; err != nil {
		return err
	}

	return nil
}

func (a AssetAuditLogRepository) AfterUpdateAssetStatus(old assets.AssetStatus, assetStatus *assets.AssetStatus) error {
	oldDataBytes, err := json.Marshal(old)
	if err != nil {
		return err
	}
	oldData := string(oldDataBytes)

	newDataBytes, err := json.Marshal(assetStatus)
	if err != nil {
		return err
	}
	newData := string(newDataBytes)

	log := assets.AssetAuditLog{
		TableName:   "asset_status",
		Action:      "UPDATE",
		OldData:     &oldData,
		NewData:     &newData,
		PerformedAt: time.Now(),
		PerformedBy: &assetStatus.UpdatedBy,
	}

	if err := a.DB.Table(tableAssetAuditLogName).Create(&log).Error; err != nil {
		return err
	}

	return nil
}

func (a AssetAuditLogRepository) AfterDeleteAssetStatus(assetStatus *assets.AssetStatus) error {
	oldDataBytes, err := json.Marshal(assetStatus)
	if err != nil {
		return err
	}
	oldData := string(oldDataBytes)

	log := assets.AssetAuditLog{
		TableName:   "asset_status",
		Action:      "DELETE",
		OldData:     &oldData,
		PerformedAt: time.Now(),
		PerformedBy: assetStatus.DeletedBy,
	}

	if err := a.DB.Table(tableAssetAuditLogName).Create(&log).Error; err != nil {
		return err
	}

	return nil
}
