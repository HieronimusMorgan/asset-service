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

func (a AssetAuditLogRepository) AfterCreate(asset *assets.Asset) (err error) {
	newDataBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	newData := string(newDataBytes)

	log := assets.AssetAuditLog{
		AssetID:     asset.AssetID,
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

func (a AssetAuditLogRepository) AfterUpdate(old assets.Asset, asset *assets.Asset) error {
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
		AssetID:     asset.AssetID,
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

func (a AssetAuditLogRepository) AfterDelete(asset *assets.Asset) error {
	oldDataBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	oldData := string(oldDataBytes)

	log := assets.AssetAuditLog{
		AssetID:     asset.AssetID,
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

func (a AssetAuditLogRepository) GetAssetAuditLogByAssetID(assetID uint) ([]assets.AssetAuditLog, error) {
	var logs []assets.AssetAuditLog
	err := a.DB.Table(tableAssetAuditLogName).Where("asset_id = ?", assetID).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (a AssetAuditLogRepository) GetAssetAuditLogByAction(action string) ([]assets.AssetAuditLog, error) {
	var logs []assets.AssetAuditLog
	err := a.DB.Table(tableAssetAuditLogName).Where("action = ?", action).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (a AssetAuditLogRepository) GetAssetAuditLogByPerformedBy(performedBy string) ([]assets.AssetAuditLog, error) {
	var logs []assets.AssetAuditLog
	err := a.DB.Table(tableAssetAuditLogName).Where("performed_by = ?", performedBy).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (a AssetAuditLogRepository) GetAssetAuditLogByPerformedAt(performedAt time.Time) ([]assets.AssetAuditLog, error) {
	var logs []assets.AssetAuditLog
	err := a.DB.Table(tableAssetAuditLogName).Where("performed_at = ?", performedAt).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (a AssetAuditLogRepository) GetAssetAuditLogByAssetIDAndAction(assetID uint, action string) ([]assets.AssetAuditLog, error) {
	var logs []assets.AssetAuditLog
	err := a.DB.Table(tableAssetAuditLogName).Where("asset_id = ? AND action = ?", assetID, action).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (a AssetAuditLogRepository) GetAssetAuditLogByAssetIDAndPerformedBy(assetID uint, performedBy string) ([]assets.AssetAuditLog, error) {
	var logs []assets.AssetAuditLog
	err := a.DB.Table(tableAssetAuditLogName).Where("asset_id = ? AND performed_by = ?", assetID, performedBy).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (a AssetAuditLogRepository) GetAssetAuditLogByAssetIDAndPerformedAt(assetID uint, performedAt time.Time) ([]assets.AssetAuditLog, error) {
	var logs []assets.AssetAuditLog
	err := a.DB.Table(tableAssetAuditLogName).Where("asset_id = ? AND performed_at = ?", assetID, performedAt).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (a AssetAuditLogRepository) GetAssetAuditLogByActionAndPerformedBy(action string, performedBy string) ([]assets.AssetAuditLog, error) {
	var logs []assets.AssetAuditLog
	err := a.DB.Table(tableAssetAuditLogName).Where("action = ? AND performed_by = ?", action, performedBy).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (a AssetAuditLogRepository) GetAssetAuditLogByActionAndPerformedAt(action string, performedAt time.Time) ([]assets.AssetAuditLog, error) {
	var logs []assets.AssetAuditLog
	err := a.DB.Table(tableAssetAuditLogName).Where("action = ? AND performed_at = ?", action, performedAt).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (a AssetAuditLogRepository) GetAssetAuditLogByPerformedByAndPerformedAt(performedBy string, performedAt time.Time) ([]assets.AssetAuditLog, error) {
	var logs []assets.AssetAuditLog
	err := a.DB.Table(tableAssetAuditLogName).Where("performed_by = ? AND performed_at = ?", performedBy, performedAt).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (a AssetAuditLogRepository) GetAssetAuditLogByAssetIDAndActionAndPerformedBy(assetID uint, action string, performedBy string) ([]assets.AssetAuditLog, error) {
	var logs []assets.AssetAuditLog
	err := a.DB.Table(tableAssetAuditLogName).Where("asset_id = ? AND action = ? AND performed_by = ?", assetID, action, performedBy).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (a AssetAuditLogRepository) GetAssetAuditLogByAssetIDAndActionAndPerformedAt(assetID uint, action string, performedAt time.Time) ([]assets.AssetAuditLog, error) {
	var logs []assets.AssetAuditLog
	err := a.DB.Table(tableAssetAuditLogName).Where("asset_id = ? AND action = ? AND performed_at = ?", assetID, action, performedAt).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (a AssetAuditLogRepository) GetAssetAuditLogByAssetIDAndPerformedByAndPerformedAt(assetID uint, performedBy string, performedAt time.Time) ([]assets.AssetAuditLog, error) {
	var logs []assets.AssetAuditLog
	err := a.DB.Table(tableAssetAuditLogName).Where("asset_id = ? AND performed_by = ? AND performed_at = ?", assetID, performedBy, performedAt).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (a AssetAuditLogRepository) GetAssetAuditLogByActionAndPerformedByAndPerformedAt(action string, performedBy string, performedAt time.Time) ([]assets.AssetAuditLog, error) {
	var logs []assets.AssetAuditLog
	err := a.DB.Table(tableAssetAuditLogName).Where("action = ? AND performed_by = ? AND performed_at = ?", action, performedBy, performedAt).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (a AssetAuditLogRepository) GetAssetAuditLogByAssetIDAndActionAndPerformedByAndPerformedAt(assetID uint, action string, performedBy string, performedAt time.Time) ([]assets.AssetAuditLog, error) {
	var logs []assets.AssetAuditLog
	err := a.DB.Table(tableAssetAuditLogName).Where("asset_id = ? AND action = ? AND performed_by = ? AND performed_at = ?", assetID, action, performedBy, performedAt).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}
