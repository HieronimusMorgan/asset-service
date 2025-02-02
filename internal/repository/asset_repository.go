package repository

import (
	"asset-service/internal/dto/out"
	"asset-service/internal/models/assets"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
)

type AssetRepository struct {
	DB       *gorm.DB
	logAudit *AssetAuditLogRepository
}

const tableAssetName = "my-home.asset"

func NewAssetRepository(db *gorm.DB) *AssetRepository {
	return &AssetRepository{DB: db, logAudit: NewAssetAuditLogRepository(db)}
}

func (r AssetRepository) AddAsset(asset *assets.Asset, maintenance *assets.AssetMaintenance) (*out.AssetResponse, error) {
	if asset == nil {
		return nil, errors.New("assets cannot be nil")
	}

	// Validate CategoryID and StatusID
	if asset.CategoryID == 0 || asset.StatusID == 0 {
		return nil, errors.New("category_id and status_id cannot be null or zero")
	}

	// Verify the existence of CategoryID and StatusID in a single query
	var exists int
	err := r.DB.Raw(
		`SELECT 
			(SELECT COUNT(*) FROM "my-home"."asset_category" WHERE asset_category_id = ?) + 
			(SELECT COUNT(*) FROM "my-home"."asset_status" WHERE asset_status_id = ?) AS exists`,
		asset.CategoryID, asset.StatusID).Scan(&exists).Error
	if err != nil {
		return nil, fmt.Errorf("failed to verify category_id or status_id: %w", err)
	}
	if exists != 2 {
		return nil, errors.New("invalid category_id or status_id")
	}

	// Start a transaction for atomic operations
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Insert the assets
	if err := tx.Table(tableAssetName).Create(&asset).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create assets: %w", err)
	}
	log.Printf("Asset created: %v", asset)
	// Set asset_id for maintenance and insert the maintenance record
	maintenance.AssetID = int(asset.AssetID)
	if err := tx.Table("my-home.asset_maintenance").Create(maintenance).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create maintenance record: %w", err)
	}

	// Retrieve the assets with maintenance information
	var result out.AssetResponse
	selectAssetQuery := `
		SELECT 
			a.asset_id,
			a.user_client_id,
			a.name,
			a.description,
			c.category_name,
			s.status_name,
			a.purchase_date,
			a.value,
			m.maintenance_date,
			m.maintenance_cost
		FROM "my-home"."asset" a
		INNER JOIN "my-home"."asset_category" c ON a.category_id = c.asset_category_id
		INNER JOIN "my-home"."asset_status" s ON a.status_id = s.asset_status_id
		LEFT JOIN "my-home"."asset_maintenance" m ON a.asset_id = m.asset_id
		WHERE a.asset_id = ?
	`
	if err := tx.Raw(selectAssetQuery, asset.AssetID).Scan(&result).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to retrieve assets after creation: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	err = r.logAudit.AfterCreateAsset(asset)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r AssetRepository) GetAssetByNameAndClientID(name string, clientID string) (*assets.Asset, error) {
	var asset assets.Asset
	err := r.DB.Table(tableAssetName).Where("name LIKE ? AND user_client_id = ?", name, clientID).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r AssetRepository) GetListAsset(clientID string) ([]out.AssetResponse, error) {
	selectQuery := `
		SELECT 
			asset.asset_id,
			asset.name,
			asset.description,
			category.category_name,
			status.status_name,
			asset.purchase_date,
			asset.value,
			maintenance.maintenance_date,
			maintenance.maintenance_cost
		FROM "my-home"."asset" asset
		INNER JOIN "my-home"."asset_category" category ON asset.category_id = category.asset_category_id
		INNER JOIN "my-home"."asset_status" status ON asset.status_id = status.asset_status_id
		LEFT JOIN "my-home"."asset_maintenance" maintenance ON asset.asset_id = maintenance.asset_id
		WHERE asset.user_client_id = ? AND asset.deleted_at IS NULL
		ORDER BY asset.name DESC
	`
	var result []out.AssetResponse
	err := r.DB.Raw(selectQuery, clientID).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r AssetRepository) GetAssetByID(clientID string, id uint) (*out.AssetResponse, error) {
	selectQuery := `
		SELECT 
			asset.asset_id,
			asset.name,
			asset.description,
			category.category_name,
			status.status_name,
			asset.purchase_date,
			asset.value,
			maintenance.maintenance_date,
			maintenance.maintenance_cost
		FROM "my-home"."asset" asset
		INNER JOIN "my-home"."asset_category" category ON asset.category_id = category.asset_category_id
		INNER JOIN "my-home"."asset_status" status ON asset.status_id = status.asset_status_id
		LEFT JOIN "my-home"."asset_maintenance" maintenance ON asset.asset_id = maintenance.asset_id
		WHERE asset.user_client_id = ? AND asset.deleted_at IS NULL AND asset.asset_id = ?
		ORDER BY asset.name DESC
	`
	var result *out.AssetResponse
	err := r.DB.Raw(selectQuery, clientID, int(id)).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r AssetRepository) UpdateAsset(asset *assets.Asset, clientID string) (*out.AssetResponse, error) {
	if asset == nil {
		return nil, errors.New("assets cannot be nil")
	}

	check, err := r.GetAssetByID(clientID, asset.AssetID)
	if err != nil {
		return nil, err
	}

	if check == nil {
		return nil, errors.New("asset not found")
	}

	// Retrieve the old asset
	var assetOld assets.Asset
	err = r.DB.Table(tableAssetName).Where("asset_id = ?", asset.AssetID).First(&assetOld).Error

	// Start a transaction
	tx := r.DB.Begin()
	defer tx.Rollback()

	// Update asset fields (only changed fields)
	if err := tx.Table(tableAssetName).
		Where("asset_id = ?", asset.AssetID).
		Updates(map[string]interface{}{
			"description":   asset.Description,
			"purchase_date": asset.PurchaseDate,
			"value":         asset.Value,
			"expiry_date":   asset.ExpiryDate,
			"updated_by":    asset.UpdatedBy,
		}).Error; err != nil {
		return nil, fmt.Errorf("failed to update asset: %w", err)
	}

	// Retrieve the assets with maintenance information
	var result out.AssetResponse
	selectAssetQuery := `
		SELECT 
			a.asset_id,
			a.user_client_id,
			a.name,
			a.description,
			c.category_name,
			s.status_name,
			a.purchase_date,
			a.value,
			m.maintenance_date,
			m.maintenance_cost,
			m.maintenance_details
		FROM "my-home"."asset" a
		INNER JOIN "my-home"."asset_category" c ON a.category_id = c.asset_category_id
		INNER JOIN "my-home"."asset_status" s ON a.status_id = s.asset_status_id
		LEFT JOIN "my-home"."asset_maintenance" m ON a.asset_id = m.asset_id
		WHERE a.asset_id = ?
	`
	if err := tx.Raw(selectAssetQuery, asset.AssetID).Scan(&result).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to retrieve assets after creation: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Log audit
	err = r.logAudit.AfterUpdateAsset(assetOld, asset)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r AssetRepository) UpdateAssetStatus(assetID uint, statusID uint, clientID string, fullName string) error {
	// Start a transaction
	tx := r.DB.Begin()
	defer tx.Rollback()

	var assetOld assets.Asset
	err := r.DB.Table(tableAssetName).Where("asset_id = ?", assetID).First(&assetOld).Error
	if err != nil {
		return fmt.Errorf("failed to find asset: %w", err)
	}

	// Verify the existence of the asset
	var asset assets.Asset
	if err := tx.Table(tableAssetName).Where("asset_id = ? AND user_client_id = ?", assetID, clientID).
		First(&asset).Error; err != nil {
		return fmt.Errorf("failed to find asset: %w", err)
	}

	// Verify the existence of the status
	var status assets.AssetStatus
	if err := tx.Table("my-home.asset_status").Where("asset_status_id = ?", statusID).
		First(&status).Error; err != nil {
		return fmt.Errorf("failed to find status: %w", err)
	}

	// Update the asset status and updated by
	if err := tx.Table(tableAssetName).Model(&asset).
		Where("asset_id = ? AND user_client_id = ?", assetID, clientID).
		Updates(map[string]interface{}{
			"status_id":  statusID,
			"updated_by": fullName,
		}).Error; err != nil {
		return fmt.Errorf("failed to update asset status: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Log audit
	err = r.logAudit.AfterUpdateAsset(assetOld, &asset)
	if err != nil {
		return err
	}

	return nil
}

func (r AssetRepository) UpdateAssetCategory(assetID uint, categoryID uint, clientID string, fullName string) error {
	// Start a transaction
	tx := r.DB.Begin()
	defer tx.Rollback()

	var assetOld assets.Asset
	err := r.DB.Table(tableAssetName).Where("asset_id = ?", assetID).First(&assetOld).Error
	if err != nil {
		return fmt.Errorf("failed to find asset: %w", err)
	}

	// Verify the existence of the asset
	var asset assets.Asset
	if err := tx.Table(tableAssetName).Where("asset_id = ? AND user_client_id = ?", assetID, clientID).
		First(&asset).Error; err != nil {
		return fmt.Errorf("failed to find asset: %w", err)
	}

	// Verify the existence of the category
	var category assets.AssetCategory
	if err := tx.Table("my-home.asset_category").Where("asset_category_id = ?", categoryID).
		First(&category).Error; err != nil {
		return fmt.Errorf("failed to find category: %w", err)
	}

	// Update the asset category and updated by
	if err := tx.Table(tableAssetName).Model(&asset).
		Where("asset_id = ? AND user_client_id = ?", assetID, clientID).
		Updates(map[string]interface{}{
			"category_id": categoryID,
			"updated_by":  fullName,
		}).Error; err != nil {
		return fmt.Errorf("failed to update asset category: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Log audit
	err = r.logAudit.AfterUpdateAsset(assetOld, &asset)
	if err != nil {
		return err
	}

	return nil
}
