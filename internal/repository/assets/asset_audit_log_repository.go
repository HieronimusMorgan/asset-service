package assets

import (
	"asset-service/internal/models/assets"
	"asset-service/internal/utils"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type AssetAuditLogRepository interface {
	AfterCreateAsset(asset *assets.Asset) error
	AfterUpdateAsset(old assets.Asset, asset *assets.Asset) error
	AfterDeleteAsset(asset *assets.Asset) error
	AfterCreateAssetMaintenance(assetMaintenance *assets.AssetMaintenance) error
	AfterUpdateAssetMaintenance(old assets.AssetMaintenance, assetMaintenance *assets.AssetMaintenance) error
	AfterDeleteAssetMaintenance(assetMaintenance *assets.AssetMaintenance) error
	AfterCreateAssetCategory(assetCategory *assets.AssetCategory) error
	AfterUpdateAssetCategory(old *assets.AssetCategory, assetCategory *assets.AssetCategory) error
	AfterDeleteAssetCategory(assetCategory *assets.AssetCategory) error
	AfterCreateAssetStatus(assetStatus *assets.AssetStatus) error
	AfterUpdateAssetStatus(old assets.AssetStatus, assetStatus *assets.AssetStatus) error
	AfterDeleteAssetStatus(assetStatus *assets.AssetStatus) error
	AfterCreateAssetMaintenanceRecord(assetMaintenanceRecord *assets.AssetMaintenanceRecord) error
	AfterUpdateAssetMaintenanceRecord(old assets.AssetMaintenanceRecord, assetMaintenanceRecord *assets.AssetMaintenanceRecord) error
	AfterDeleteAssetMaintenanceRecord(assetMaintenanceRecord *assets.AssetMaintenanceRecord) error
}

type assetAuditLogRepository struct {
	db gorm.DB
}

func NewAssetAuditLogRepository(db gorm.DB) AssetAuditLogRepository {
	return &assetAuditLogRepository{db: db}
}

func (a assetAuditLogRepository) AfterCreateAsset(asset *assets.Asset) (err error) {
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

	if err := a.db.Table(utils.TableAssetAuditLogName).Create(&log).Error; err != nil {
		return err
	}
	return nil
}

func (a assetAuditLogRepository) AfterUpdateAsset(old assets.Asset, asset *assets.Asset) error {
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

	if err := a.db.Table(utils.TableAssetAuditLogName).Create(&log).Error; err != nil {
		return err
	}
	return nil
}

func (a assetAuditLogRepository) AfterDeleteAsset(asset *assets.Asset) error {
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

	if err := a.db.Table(utils.TableAssetAuditLogName).Create(&log).Error; err != nil {
		return err
	}
	return nil
}

func (a assetAuditLogRepository) AfterCreateAssetMaintenance(assetMaintenance *assets.AssetMaintenance) error {
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

	if err := a.db.Table(utils.TableAssetAuditLogName).Create(&log).Error; err != nil {
		return err
	}
	return nil
}

func (a assetAuditLogRepository) AfterUpdateAssetMaintenance(old assets.AssetMaintenance, assetMaintenance *assets.AssetMaintenance) error {
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

	if err := a.db.Table(utils.TableAssetAuditLogName).Create(&log).Error; err != nil {
		return err
	}
	return nil
}

func (a assetAuditLogRepository) AfterDeleteAssetMaintenance(assetMaintenance *assets.AssetMaintenance) error {
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

	if err := a.db.Table(utils.TableAssetAuditLogName).Create(&log).Error; err != nil {
		return err
	}
	return nil
}

func (a assetAuditLogRepository) AfterCreateAssetCategory(assetCategory *assets.AssetCategory) error {
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

	if err := a.db.Table(utils.TableAssetAuditLogName).Create(&log).Error; err != nil {
		return err
	}
	return nil
}

func (a assetAuditLogRepository) AfterUpdateAssetCategory(old *assets.AssetCategory, assetCategory *assets.AssetCategory) error {
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

	if err := a.db.Table(utils.TableAssetAuditLogName).Create(&log).Error; err != nil {
		return err
	}
	return nil
}

func (a assetAuditLogRepository) AfterDeleteAssetCategory(assetCategory *assets.AssetCategory) error {
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

	if err := a.db.Table(utils.TableAssetAuditLogName).Create(&log).Error; err != nil {
		return err
	}
	return nil
}

func (a assetAuditLogRepository) AfterCreateAssetStatus(assetStatus *assets.AssetStatus) error {
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

	if err := a.db.Table(utils.TableAssetAuditLogName).Create(&log).Error; err != nil {
		return err
	}

	return nil
}

func (a assetAuditLogRepository) AfterUpdateAssetStatus(old assets.AssetStatus, assetStatus *assets.AssetStatus) error {
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

	if err := a.db.Table(utils.TableAssetAuditLogName).Create(&log).Error; err != nil {
		return err
	}

	return nil
}

func (a assetAuditLogRepository) AfterDeleteAssetStatus(assetStatus *assets.AssetStatus) error {
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

	if err := a.db.Table(utils.TableAssetAuditLogName).Create(&log).Error; err != nil {
		return err
	}

	return nil
}

func (a assetAuditLogRepository) AfterCreateAssetMaintenanceRecord(assetMaintenanceRecord *assets.AssetMaintenanceRecord) error {
	newDataBytes, err := json.Marshal(assetMaintenanceRecord)
	if err != nil {
		return err
	}
	newData := string(newDataBytes)

	log := assets.AssetAuditLog{
		TableName:   "asset_maintenance_record",
		Action:      "CREATE",
		NewData:     &newData,
		PerformedAt: time.Now(),
		PerformedBy: &assetMaintenanceRecord.CreatedBy,
	}

	if err := a.db.Table(utils.TableAssetAuditLogName).Create(&log).Error; err != nil {
		return err
	}

	return nil
}

func (a assetAuditLogRepository) AfterUpdateAssetMaintenanceRecord(old assets.AssetMaintenanceRecord, assetMaintenanceRecord *assets.AssetMaintenanceRecord) error {
	oldDataBytes, err := json.Marshal(old)
	if err != nil {
		return err
	}
	oldData := string(oldDataBytes)

	newDataBytes, err := json.Marshal(assetMaintenanceRecord)
	if err != nil {
		return err
	}
	newData := string(newDataBytes)

	log := assets.AssetAuditLog{
		TableName:   "asset_maintenance_record",
		Action:      "UPDATE",
		OldData:     &oldData,
		NewData:     &newData,
		PerformedAt: time.Now(),
		PerformedBy: &assetMaintenanceRecord.UpdatedBy,
	}

	if err := a.db.Table(utils.TableAssetAuditLogName).Create(&log).Error; err != nil {
		return err
	}

	return nil
}

func (a assetAuditLogRepository) AfterDeleteAssetMaintenanceRecord(assetMaintenanceRecord *assets.AssetMaintenanceRecord) error {
	oldDataBytes, err := json.Marshal(assetMaintenanceRecord)
	if err != nil {
		return err
	}
	oldData := string(oldDataBytes)

	log := assets.AssetAuditLog{
		TableName:   "asset_maintenance_record",
		Action:      "DELETE",
		OldData:     &oldData,
		PerformedAt: time.Now(),
		PerformedBy: assetMaintenanceRecord.DeletedBy,
	}

	if err := a.db.Table(utils.TableAssetAuditLogName).Create(&log).Error; err != nil {
		return err
	}

	return nil
}
