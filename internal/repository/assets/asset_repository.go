package assets

import (
	response "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
)

type AssetRepository struct {
	DB    *gorm.DB
	audit *AssetAuditLogRepository
}

const tableAssetName = "my-home.asset"

func NewAssetRepository(db *gorm.DB) *AssetRepository {
	return &AssetRepository{DB: db, audit: NewAssetAuditLogRepository(db)}
}

func (r AssetRepository) AddAsset(asset *assets.Asset) (*response.AssetResponseList, error) {
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

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	selectQuery := `
        SELECT 
            asset.asset_id,
            asset.user_client_id,
            asset.asset_code,
            asset.name,
            asset.description,
            asset.barcode,
            asset.purchase_date,
            asset.expiry_date,
            asset.warranty_expiry_date,
            asset.price,
            asset.stock,
            category.asset_category_id,
            category.category_name,
            category.description AS category_description,
            status.asset_status_id,
            status.status_name,
            status.description AS status_description
        FROM "my-home"."asset" asset
        INNER JOIN "my-home"."asset_category" category ON asset.category_id = category.asset_category_id
        INNER JOIN "my-home"."asset_status" status ON asset.status_id = status.asset_status_id
        WHERE asset.user_client_id = ? AND asset.asset_id = ? AND asset.deleted_at IS NULL AND asset.is_wishlist = false
        ORDER BY asset.name DESC;
    `

	row := r.DB.Raw(selectQuery, asset.UserClientID, asset.AssetID).Row()
	var assetResult response.AssetResponseList
	var categoryResult response.AssetCategoryResponse
	var statusResult response.AssetStatusResponse

	err = row.Scan(
		&assetResult.ID,
		&assetResult.ClientID,
		&assetResult.AssetCode,
		&assetResult.Name,
		&assetResult.Description,
		&assetResult.Barcode,
		&assetResult.PurchaseDate,
		&assetResult.ExpiryDate,
		&assetResult.WarrantyExpiry,
		&assetResult.Price,
		&assetResult.Stock,
		&categoryResult.AssetCategoryID,
		&categoryResult.CategoryName,
		&categoryResult.Description,
		&statusResult.AssetStatusID,
		&statusResult.StatusName,
		&statusResult.Description,
	)

	if err != nil {
		return nil, err
	}

	assetResult.Category = categoryResult
	assetResult.Status = statusResult

	err = r.audit.AfterCreateAsset(asset)
	if err != nil {
		return nil, err
	}

	return &assetResult, nil
}

func (r AssetRepository) GetAssetByNameAndClientID(name string, clientID string) (*assets.Asset, error) {
	var asset assets.Asset
	err := r.DB.Table(tableAssetName).Where("name = ? AND user_client_id = ?", name, clientID).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r AssetRepository) AssetNameExists(name string, clientID string) (bool, error) {
	var count int64
	err := r.DB.Table(tableAssetName).
		Where("name = ? AND user_client_id = ?", name, clientID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (r AssetRepository) GetAsset(assetID uint, clientID string) (*assets.Asset, error) {
	var asset assets.Asset
	err := r.DB.Table(tableAssetName).
		Where("id = ? AND user_client_id = ?", assetID, clientID).First(asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r AssetRepository) GetListAsset(clientID string) ([]response.AssetResponseList, error) {
	selectQuery := `
        SELECT 
            asset.asset_id,
            asset.user_client_id,
            asset.asset_code,
            asset.name,
            asset.description,
            asset.barcode,
            asset.purchase_date,
            asset.expiry_date,
            asset.warranty_expiry_date,
            asset.price,
            asset.stock,
            category.asset_category_id,
            category.category_name,
            category.description AS category_description,
            status.asset_status_id,
            status.status_name,
            status.description AS status_description,
            asset.purchase_date,
            asset.price
        FROM "my-home"."asset" asset
        INNER JOIN "my-home"."asset_category" category ON asset.category_id = category.asset_category_id
        INNER JOIN "my-home"."asset_status" status ON asset.status_id = status.asset_status_id
        WHERE asset.user_client_id = ? AND asset.deleted_at IS NULL AND asset.is_wishlist = false
        ORDER BY asset.name DESC;
    `
	rows, err := r.DB.Raw(selectQuery, clientID).Rows()
	if err != nil {
		return nil, err
	}

	var result []response.AssetResponseList
	for rows.Next() {
		var asset response.AssetResponseList
		var category response.AssetCategoryResponse
		var status response.AssetStatusResponse

		var purchaseDate sql.NullTime
		var expiryDate sql.NullTime
		var warrantyExpiryDate sql.NullTime

		err := rows.Scan(
			&asset.ID,
			&asset.ClientID,
			&asset.AssetCode,
			&asset.Name,
			&asset.Description,
			&asset.Barcode,
			&purchaseDate,
			&expiryDate,
			&warrantyExpiryDate,
			&asset.Price,
			&asset.Stock,
			&category.AssetCategoryID,
			&category.CategoryName,
			&category.Description,
			&status.AssetStatusID,
			&status.StatusName,
			&status.Description,
			&asset.PurchaseDate,
			&asset.Price,
		)

		if err != nil {
			return nil, err
		}

		asset.Category = category
		asset.Status = status

		result = append(result, asset)
	}

	return result, nil
}

func (r AssetRepository) GetAssetByID(clientID string, id uint) (*response.AssetResponseList, error) {
	selectQuery := `
        SELECT 
            asset.asset_id,
            asset.user_client_id,
            asset.asset_code,
            asset.name,
            asset.description,
            asset.barcode,
            asset.purchase_date,
            asset.expiry_date,
            asset.warranty_expiry_date,
            asset.price,
            asset.stock,
            category.asset_category_id,
            category.category_name,
            category.description AS category_description,
            status.asset_status_id,
            status.status_name,
            status.description AS status_description,
            asset.purchase_date,
            asset.price
        FROM "my-home"."asset" asset
        INNER JOIN "my-home"."asset_category" category ON asset.category_id = category.asset_category_id
        INNER JOIN "my-home"."asset_status" status ON asset.status_id = status.asset_status_id
        WHERE asset.user_client_id = ? AND asset.asset_id = ? AND asset.deleted_at IS NULL AND asset.is_wishlist = false
        ORDER BY asset.name DESC;
    `
	row := r.DB.Raw(selectQuery, clientID, id).Row()
	var assetResult response.AssetResponseList
	var categoryResult response.AssetCategoryResponse
	var statusResult response.AssetStatusResponse

	var purchaseDate sql.NullTime
	var expiryDate sql.NullTime
	var warrantyExpiryDate sql.NullTime

	err := row.Scan(
		&assetResult.ID,
		&assetResult.ClientID,
		&assetResult.AssetCode,
		&assetResult.Name,
		&assetResult.Description,
		&assetResult.Barcode,
		&purchaseDate,
		&expiryDate,
		&warrantyExpiryDate,
		&assetResult.Price,
		&assetResult.Stock,
		&categoryResult.AssetCategoryID,
		&categoryResult.CategoryName,
		&categoryResult.Description,
		&statusResult.AssetStatusID,
		&statusResult.StatusName,
		&statusResult.Description,
		&assetResult.PurchaseDate,
		&assetResult.Price,
	)

	if err != nil {
		return nil, err
	}

	assetResult.Category = categoryResult
	assetResult.Status = statusResult

	return &assetResult, nil
}

func (r AssetRepository) UpdateAsset(asset *assets.Asset, clientID string) (*response.AssetResponse, error) {
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
			"price":         asset.Price,
			"expiry_date":   asset.ExpiryDate,
			"updated_by":    asset.UpdatedBy,
		}).Error; err != nil {
		return nil, fmt.Errorf("failed to update asset: %w", err)
	}

	// Retrieve the assets with maintenance information
	var result response.AssetResponse
	selectAssetQuery := `
		SELECT 
			a.asset_id,
			a.user_client_id,
			a.name,
			a.description,
			c.category_name,
			s.status_name,
			a.purchase_date,
			a.price
		FROM "my-home"."asset" a
		INNER JOIN "my-home"."asset_category" c ON a.category_id = c.asset_category_id
		INNER JOIN "my-home"."asset_status" s ON a.status_id = s.asset_status_id
		WHERE a.asset_id = ? AND a.is_wishlist = false AND a.deleted_at IS NULL
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
	err = r.audit.AfterUpdateAsset(assetOld, asset)
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
	err = r.audit.AfterUpdateAsset(assetOld, &asset)
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
	err = r.audit.AfterUpdateAsset(assetOld, &asset)
	if err != nil {
		return err
	}

	return nil
}

func (r AssetRepository) DeleteAsset(id uint, clientID string, fullName string) error {
	// Start a transaction
	tx := r.DB.Begin()
	defer tx.Rollback()

	// Verify the existence of the asset
	var asset assets.Asset
	var assetMaintenance assets.AssetMaintenance
	var assetMaintenanceRecord assets.AssetMaintenanceRecord

	if err := tx.Table(tableAssetName).Where("asset_id = ? AND user_client_id = ?", id, clientID).
		First(&asset).Error; err != nil {
		return fmt.Errorf("failed to find asset: %w", err)
	}

	// Soft delete the asset
	if err := tx.Table(tableAssetName).Model(&asset).
		Where("asset_id = ? AND user_client_id = ?", id, clientID).
		Updates(map[string]interface{}{
			"deleted_by": fullName,
		}).
		Delete(&assets.Asset{}).Error; err != nil {
		return fmt.Errorf("failed to delete asset: %w", err)
	}

	// Soft delete the asset_maintenance if it exists
	if err := tx.Table("my-home.asset_maintenance").Model(&assetMaintenance).
		Where("asset_id = ?", id).
		Updates(map[string]interface{}{
			"deleted_by": fullName,
		}).
		Delete(&assets.AssetMaintenance{}).Error; err != nil {
		return fmt.Errorf("failed to delete asset maintenance: %w", err)
	}

	// soft delete the asset_maintenance_record
	if err := tx.Table("my-home.asset_maintenance_record").Model(&assetMaintenanceRecord).
		Where("asset_id = ?", id).
		Updates(map[string]interface{}{
			"deleted_by": fullName,
		}).
		Delete(&assets.AssetMaintenance{}).Error; err != nil {
		return fmt.Errorf("failed to delete asset maintenance: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	err := r.audit.AfterDeleteAsset(&asset)
	if err != nil {
		return err
	}

	err = r.audit.AfterDeleteAssetMaintenance(&assetMaintenance)
	if err != nil {
		return err
	}

	err = r.audit.AfterDeleteAssetMaintenanceRecord(&assetMaintenanceRecord)
	if err != nil {
		return err
	}

	return nil
}

func (r AssetRepository) GetAssetByIDForMaintenance(id uint, clientID string) (*assets.Asset, error) {
	var asset assets.Asset
	err := r.DB.Table(tableAssetName).Where("asset_id = ? AND user_client_id = ?", id, clientID).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}
