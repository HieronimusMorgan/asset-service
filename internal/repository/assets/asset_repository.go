package assets

import (
	response "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	"asset-service/internal/utils"
	"database/sql"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

type AssetRepository interface {
	AddAsset(asset *assets.Asset) error
	GetAssetByNameAndClientID(name string, clientID string) (*assets.Asset, error)
	AssetNameExists(name string, clientID string) (bool, error)
	GetAsset(assetID uint, clientID string) (*assets.Asset, error)
	GetListAssets(clientID string) ([]response.AssetResponse, error)
	GetAssetResponseByID(clientID string, id uint) (*response.AssetResponse, error)
	GetAssetByID(clientID string, id uint) (*assets.Asset, error)
	UpdateAsset(asset *assets.Asset, clientID string) error
	UpdateMaintenanceDateAsset(assetID uint, maintenanceDate *time.Time, clientID string) error
	UpdateAssetStatus(assetID uint, statusID uint, clientID string) (*assets.Asset, error)
	UpdateAssetCategory(assetID uint, categoryID uint, clientID string) (*assets.Asset, error)
	DeleteAsset(id uint, clientID string) error
	GetAssetByIDForMaintenance(id uint, clientID string) (*assets.Asset, error)
	GetAssetByCategoryID(assetCategoryID uint, clientID string) ([]assets.Asset, error)
	GetAssetDeleted() ([]assets.Asset, error)
}

type assetRepository struct {
	db    gorm.DB
	audit AssetAuditLogRepository
}

func NewAssetRepository(db gorm.DB, audit AssetAuditLogRepository) AssetRepository {
	return assetRepository{db: db, audit: audit}
}

func (r assetRepository) AddAsset(asset *assets.Asset) error {
	if asset == nil {
		return errors.New("assets cannot be nil")
	}

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Table(utils.TableAssetName).Create(&asset).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create assets: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r assetRepository) GetAssetByNameAndClientID(name string, clientID string) (*assets.Asset, error) {
	var asset assets.Asset
	err := r.db.Table(utils.TableAssetName).Where("name = ? AND user_client_id = ?", name, clientID).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r assetRepository) AssetNameExists(name string, clientID string) (bool, error) {
	var count int64
	err := r.db.Table(utils.TableAssetName).
		Where("name = ? AND user_client_id = ?", name, clientID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r assetRepository) GetAsset(assetID uint, clientID string) (*assets.Asset, error) {
	var asset assets.Asset
	err := r.db.Table(utils.TableAssetName).
		Where("asset_id = ? AND user_client_id = ?", assetID, clientID).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r assetRepository) GetListAssets(clientID string) ([]response.AssetResponse, error) {
	selectQuery := `
       SELECT 
           asset.asset_id,
           asset.user_client_id,
           asset.serial_number,
           asset.name,
           asset.description,
           asset.barcode,
           asset.purchase_date,
           asset.expiry_date,
           asset.warranty_expiry_date,
           asset.price,
           asset.stock,
           asset.notes,
           category.asset_category_id,
           category.category_name,
           category.description AS category_description,
           status.asset_status_id,
           status.status_name,
           status.description AS status_description
       FROM "asset-service"."asset" asset
       INNER JOIN "asset-service"."asset_category" category ON asset.category_id = category.asset_category_id
       INNER JOIN "asset-service"."asset_status" status ON asset.status_id = status.asset_status_id
       WHERE asset.user_client_id = ? AND asset.deleted_at IS NULL
       ORDER BY asset.created_at ASC;
   `

	rows, err := r.db.Raw(selectQuery, clientID).Rows()
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("❌ Failed to fetch asset list")
		return nil, err
	}

	var assetsList []response.AssetResponse
	for rows.Next() {
		var asset response.AssetResponse
		var category response.AssetCategoryResponse
		var status response.AssetStatusResponse

		// Handling NULL values from SQL
		var serialNumber sql.NullString
		var barcode sql.NullString
		var description sql.NullString
		var purchaseDate sql.NullTime
		var expiryDate sql.NullTime
		var warrantyExpiryDate sql.NullTime
		var price sql.NullFloat64
		var stock sql.NullInt64
		var notes sql.NullString

		err := rows.Scan(
			&asset.AssetID,
			&asset.UserClientID,
			&serialNumber,
			&asset.Name,
			&description,
			&barcode,
			&purchaseDate,
			&expiryDate,
			&warrantyExpiryDate,
			&price,
			&stock,
			&notes,
			&category.AssetCategoryID,
			&category.CategoryName,
			&category.Description,
			&status.AssetStatusID,
			&status.StatusName,
			&status.Description,
		)

		if err != nil {
			log.Error().Str("clientID", clientID).Err(err).Msg("❌ Failed to scan asset row")
			return nil, err
		}

		// Convert NULL SQL values to Go `nil`
		if serialNumber.Valid {
			asset.SerialNumber = &serialNumber.String
		}
		if barcode.Valid {
			asset.Barcode = &barcode.String
		}
		if description.Valid {
			asset.Description = description.String
		}
		if price.Valid {
			asset.Price = price.Float64
		}
		if stock.Valid {
			asset.Stock = int(stock.Int64)
		}
		if notes.Valid {
			asset.Notes = &notes.String
		}
		if purchaseDate.Valid {
			asset.PurchaseDate = (*response.DateOnly)(&purchaseDate.Time)
		}
		if expiryDate.Valid {
			asset.ExpiryDate = (*response.DateOnly)(&expiryDate.Time)
		}
		if warrantyExpiryDate.Valid {
			asset.WarrantyExpiryDate = (*response.DateOnly)(&warrantyExpiryDate.Time)
		}

		// Assign category and status details
		asset.Category = category
		asset.Status = status

		// Append to result slice
		assetsList = append(assetsList, asset)
	}

	log.Info().Str("clientID", clientID).Int("assets_count", len(assetsList)).Msg("✅ Successfully fetched asset list")
	return assetsList, nil
}

func (r assetRepository) GetAssetResponseByID(clientID string, id uint) (*response.AssetResponse, error) {
	selectQuery := `
       SELECT 
           asset.asset_id,
           asset.user_client_id,
           asset.serial_number,
           asset.name,
           asset.description,
           asset.barcode,
           asset.purchase_date,
           asset.expiry_date,
           asset.warranty_expiry_date,
           asset.price,
           asset.stock,
           asset.notes,
           category.asset_category_id,
           category.category_name,
           category.description AS category_description,
           status.asset_status_id,
           status.status_name,
           status.description AS status_description
       FROM "asset-service"."asset" asset
       INNER JOIN "asset-service"."asset_category" category ON asset.category_id = category.asset_category_id
       INNER JOIN "asset-service"."asset_status" status ON asset.status_id = status.asset_status_id
       WHERE asset.user_client_id = ? AND asset.asset_id = ?  AND asset.deleted_at IS NULL
       ORDER BY asset.asset_id ASC;
   `

	rows := r.db.Raw(selectQuery, clientID, id).Row()

	var asset response.AssetResponse
	var category response.AssetCategoryResponse
	var status response.AssetStatusResponse

	// Handling NULL values from SQL
	var serialNumber sql.NullString
	var barcode sql.NullString
	var description sql.NullString
	var purchaseDate sql.NullTime
	var expiryDate sql.NullTime
	var warrantyExpiryDate sql.NullTime
	var price sql.NullFloat64
	var stock sql.NullInt64
	var notes sql.NullString

	err := rows.Scan(
		&asset.AssetID,
		&asset.UserClientID,
		&serialNumber,
		&asset.Name,
		&description,
		&barcode,
		&purchaseDate,
		&expiryDate,
		&warrantyExpiryDate,
		&price,
		&stock,
		&notes,
		&category.AssetCategoryID,
		&category.CategoryName,
		&category.Description,
		&status.AssetStatusID,
		&status.StatusName,
		&status.Description,
	)

	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("❌ Failed to scan asset row")
		return nil, err
	}

	// Convert NULL SQL values to Go `nil`
	if serialNumber.Valid {
		asset.SerialNumber = &serialNumber.String
	}
	if barcode.Valid {
		asset.Barcode = &barcode.String
	}
	if description.Valid {
		asset.Description = description.String
	}
	if price.Valid {
		asset.Price = price.Float64
	}
	if stock.Valid {
		asset.Stock = int(stock.Int64)
	}
	if notes.Valid {
		asset.Notes = &notes.String
	}
	if purchaseDate.Valid {
		asset.PurchaseDate = (*response.DateOnly)(&purchaseDate.Time)
	}
	if expiryDate.Valid {
		asset.ExpiryDate = (*response.DateOnly)(&expiryDate.Time)
	}
	if warrantyExpiryDate.Valid {
		asset.WarrantyExpiryDate = (*response.DateOnly)(&warrantyExpiryDate.Time)
	}

	// Assign category and status details
	asset.Category = category
	asset.Status = status

	return &asset, nil
}

func (r assetRepository) GetAssetByID(clientID string, id uint) (*assets.Asset, error) {
	var asset assets.Asset
	err := r.db.Table(utils.TableAssetName).Where("asset_id = ? AND user_client_id = ?", id, clientID).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r assetRepository) UpdateAsset(asset *assets.Asset, clientID string) error {
	tx := r.db.Begin()
	defer tx.Rollback()

	// Update asset fields (only changed fields)
	if err := tx.Table(utils.TableAssetName).
		Where("asset_id = ? AND user_client_id = ?", asset.AssetID, clientID).
		Updates(asset).Error; err != nil {
		return fmt.Errorf("failed to update asset: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r assetRepository) UpdateMaintenanceDateAsset(assetID uint, maintenanceDate *time.Time, clientID string) error {
	tx := r.db.Begin()
	defer tx.Rollback()

	// Update asset fields (only changed fields)
	if err := tx.Table(utils.TableAssetName).
		Where("user_client_id = ? AND asset_id = ?", clientID, assetID).
		Updates(map[string]interface{}{"maintenance_date": maintenanceDate}).Error; err != nil {
		return fmt.Errorf("failed to update maintenance date: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r assetRepository) UpdateAssetStatus(assetID uint, statusID uint, clientID string) (*assets.Asset, error) {
	// Start a transaction
	tx := r.db.Begin()
	defer tx.Rollback()

	var assetOld assets.Asset
	err := r.db.Table(utils.TableAssetName).Where("asset_id = ?", assetID).First(&assetOld).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find asset: %w", err)
	}

	// Verify the existence of the asset
	var asset assets.Asset
	if err := tx.Table(utils.TableAssetName).Where("asset_id = ? AND user_client_id = ?", assetID, clientID).
		First(&asset).Error; err != nil {
		return nil, fmt.Errorf("failed to find asset: %w", err)
	}

	// Verify the existence of the status
	var status assets.AssetStatus
	if err := tx.Table("asset-service.asset_status").Where("asset_status_id = ?", statusID).
		First(&status).Error; err != nil {
		return nil, fmt.Errorf("failed to find status: %w", err)
	}

	// Update the asset status and updated by
	if err := tx.Table(utils.TableAssetName).Model(&asset).
		Where("asset_id = ? AND user_client_id = ?", assetID, clientID).
		Updates(map[string]interface{}{
			"status_id":  statusID,
			"updated_by": clientID,
		}).Error; err != nil {
		return nil, fmt.Errorf("failed to update asset status: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &asset, nil
}

func (r assetRepository) UpdateAssetCategory(assetID uint, categoryID uint, clientID string) (*assets.Asset, error) {
	// Start a transaction
	tx := r.db.Begin()
	defer tx.Rollback()

	var assetOld assets.Asset
	err := r.db.Table(utils.TableAssetName).Where("asset_id = ?", assetID).First(&assetOld).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find asset: %w", err)
	}

	// Verify the existence of the asset
	var asset assets.Asset
	if err := tx.Table(utils.TableAssetName).Where("asset_id = ? AND user_client_id = ?", assetID, clientID).
		First(&asset).Error; err != nil {
		return nil, fmt.Errorf("failed to find asset: %w", err)
	}

	// Verify the existence of the category
	var category assets.AssetCategory
	if err := tx.Table("asset-service.asset_category").Where("asset_category_id = ?", categoryID).
		First(&category).Error; err != nil {
		return nil, fmt.Errorf("failed to find category: %w", err)
	}

	// Update the asset category and updated by
	if err := tx.Table(utils.TableAssetName).Model(&asset).
		Where("asset_id = ? AND user_client_id = ?", assetID, clientID).
		Updates(map[string]interface{}{
			"category_id": categoryID,
			"updated_by":  clientID,
		}).Error; err != nil {
		return nil, fmt.Errorf("failed to update asset category: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &asset, nil
}
func (r assetRepository) DeleteAsset(id uint, clientID string) error {
	if err := r.db.Table(utils.TableAssetName).Model(assets.Asset{}).
		Where("asset_id = ? AND user_client_id = ?", id, clientID).
		Updates(map[string]interface{}{"deleted_by": clientID, "deleted_at": time.Now()}).
		Delete(&assets.Asset{}).Error; err != nil {
		return fmt.Errorf("failed to delete asset: %w", err)
	}
	return nil
}

func (r assetRepository) GetAssetByIDForMaintenance(id uint, clientID string) (*assets.Asset, error) {
	var asset assets.Asset
	err := r.db.Table(utils.TableAssetName).Where("asset_id = ? AND user_client_id = ?", id, clientID).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r assetRepository) GetAssetByCategoryID(assetCategoryID uint, clientID string) ([]assets.Asset, error) {
	var asset []assets.Asset
	err := r.db.Table(utils.TableAssetName).Where("category_id = ? AND user_client_id = ?", assetCategoryID, clientID).Find(&asset).Error
	if err != nil {
		return nil, err
	}
	return asset, nil
}

func (r assetRepository) GetAssetDeleted() ([]assets.Asset, error) {
	var asset []assets.Asset
	err := r.db.Unscoped().Table(utils.TableAssetName).
		Where("deleted_at IS NOT NULL").
		Find(&asset).Error
	if err != nil {
		return nil, err
	}
	return asset, nil
}
