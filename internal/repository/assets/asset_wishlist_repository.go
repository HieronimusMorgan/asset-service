package assets

import (
	response "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	"asset-service/internal/utils"
	"database/sql"
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type AssetWishlistRepository interface {
	AddAssetWishlist(asset *assets.Asset) error
	GetAssetWishlistByID(clientID string, assetID uint) (*response.AssetResponse, error)
	GetAssetWishlistList(clientID string) ([]response.AssetResponse, error)
	UpdateAssetWishlist(assetID uint, data map[string]interface{}) error
	DeleteAssetWishlist(clientID, name string, assetID uint) error
	GetAssetWishlistByCategory(clientID string, categoryID uint) ([]response.AssetResponse, error)
	GetAssetWishlistByStatus(clientID string, statusID uint) ([]response.AssetResponse, error)
	GetAssetWishlistByCategoryAndStatus(clientID string, categoryID, statusID uint) ([]response.AssetResponse, error)
}

type assetWishlistRepository struct {
	db    gorm.DB
	audit AssetAuditLogRepository
	asset AssetRepository
}

func NewAssetWishlistRepository(db gorm.DB, audit AssetAuditLogRepository, asset AssetRepository) AssetWishlistRepository {
	return assetWishlistRepository{db: db, audit: audit, asset: asset}
}

func (r assetWishlistRepository) AddAssetWishlist(asset *assets.Asset) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Insert the assets
	if err := tx.Table(utils.TableAssetName).Create(&asset).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create assets: %w", err)
	}
	log.Printf("Asset created: %v", asset)

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
func (r assetWishlistRepository) GetAssetWishlistByID(clientID string, assetID uint) (*response.AssetResponse, error) {
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
        WHERE asset.user_client_id = ? AND asset.asset_id = ? AND asset.deleted_at IS NULL AND asset.is_wishlist = true
        ORDER BY asset.asset_id ASC;
    `

	row := r.db.Raw(selectQuery, clientID, assetID).Row()

	var asset response.AssetResponse
	var category response.AssetCategoryResponse
	var status response.AssetStatusResponse

	// Handling NULL values from SQL
	var serialNumber sql.NullString
	var barcode sql.NullString
	var notes sql.NullString
	var purchaseDate sql.NullTime
	var expiryDate sql.NullTime
	var warrantyExpiryDate sql.NullTime

	err := row.Scan(
		&asset.AssetID,
		&asset.UserClientID,
		&serialNumber,
		&asset.Name,
		&asset.Description,
		&barcode,
		&purchaseDate,
		&expiryDate,
		&warrantyExpiryDate,
		&asset.Price,
		&asset.Stock,
		&notes,
		&category.AssetCategoryID,
		&category.CategoryName,
		&category.Description,
		&status.AssetStatusID,
		&status.StatusName,
		&status.Description,
	)

	if err != nil {
		log.Error().Uint("assetID", assetID).Str("clientID", clientID).Err(err).Msg("❌ Failed to scan asset wishlist")
		return nil, err
	}

	// Convert NULL SQL values to Go `nil`
	if serialNumber.Valid {
		asset.SerialNumber = &serialNumber.String
	}
	if barcode.Valid {
		asset.Barcode = &barcode.String
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

	asset.Category = category
	asset.Status = status

	log.Info().Uint("assetID", assetID).Str("clientID", clientID).Msg("✅ Successfully retrieved asset wishlist")
	return &asset, nil
}
func (r assetWishlistRepository) GetAssetWishlistList(clientID string) ([]response.AssetResponse, error) {
	selectQuery := `
		SELECT 
			asset.asset_id,
			asset.user_client_id,
			asset.name,
			asset.description,
			asset.price,
			category.asset_category_id,
			category.category_name,
			category.description AS category_description,
			status.asset_status_id,
			status.status_name,
			status.description AS status_description,
			asset.purchase_date
		FROM "asset-service"."asset" asset
		INNER JOIN "asset-service"."asset_category" category ON asset.category_id = category.asset_category_id
		INNER JOIN "asset-service"."asset_status" status ON asset.status_id = status.asset_status_id
		WHERE asset.user_client_id = ? AND asset.deleted_at IS NULL AND asset.is_wishlist = true
		ORDER BY asset.asset_id ASC;
	`

	rows, err := r.db.Raw(selectQuery, clientID).Rows()
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("❌ Failed to fetch asset wishlist list")
		return nil, err
	} // Ensures rows are closed after execution

	var assetsResult []response.AssetResponse
	for rows.Next() {
		var assetResult response.AssetResponse
		var categoryResult response.AssetCategoryResponse
		var statusResult response.AssetStatusResponse

		// Handling NULL values from SQL
		var description sql.NullString
		var purchaseDate sql.NullTime
		var price sql.NullFloat64

		err := rows.Scan(
			&assetResult.AssetID,
			&assetResult.UserClientID,
			&assetResult.Name,
			&description, // Handling NULL description
			&price,       // Handling NULL price
			&categoryResult.AssetCategoryID,
			&categoryResult.CategoryName,
			&categoryResult.Description,
			&statusResult.AssetStatusID,
			&statusResult.StatusName,
			&statusResult.Description,
			&purchaseDate, // Handling NULL purchaseDate
		)

		if err != nil {
			log.Error().Str("clientID", clientID).Err(err).Msg("❌ Failed to scan asset wishlist row")
			return nil, err
		}

		// Convert NULL SQL values to Go `nil`
		if description.Valid {
			assetResult.Description = description.String
		}
		if price.Valid {
			assetResult.Price = price.Float64
		}
		if purchaseDate.Valid {
			assetResult.PurchaseDate = (*response.DateOnly)(&purchaseDate.Time)
		}

		// Assign category and status details
		assetResult.Category = categoryResult
		assetResult.Status = statusResult

		// Append to result slice
		assetsResult = append(assetsResult, assetResult)
	}

	log.Info().Str("clientID", clientID).Int("assets_count", len(assetsResult)).Msg("✅ Successfully fetched asset wishlist list")
	return assetsResult, nil
}

func (r assetWishlistRepository) UpdateAssetWishlist(assetID uint, data map[string]interface{}) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Table(utils.TableAssetName).
		Where("asset_id = ?", assetID).
		Updates(data).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update assets: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r assetWishlistRepository) DeleteAssetWishlist(clientID, name string, assetID uint) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Table(utils.TableAssetName).Where("asset_id = ? AND user_client_id = ?", assetID, clientID).
		First(assets.Asset{}).Error; err != nil {
		return fmt.Errorf("failed to find asset: %w", err)
	}

	if err := tx.Table(utils.TableAssetName).
		Where("asset_id = ?", assetID).
		Updates(map[string]interface{}{
			"deleted_at": gorm.Expr("NOW()"),
			"deleted_by": name,
		}).Delete(&assets.Asset{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete assets: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r assetWishlistRepository) GetAssetWishlistByCategory(clientID string, categoryID uint) ([]response.AssetResponse, error) {
	selectQuery := `
		SELECT 
			asset.asset_id,
			asset.user_client_id,
			asset.name,
			asset.description,
			asset.price,
			category.asset_category_id,
			category.category_name,
			category.description AS category_description,
			status.asset_status_id,
			status.status_name,
			status.description AS status_description,
			asset.purchase_date,
			asset.price
		FROM "asset-service"."asset" asset
		INNER JOIN "asset-service"."asset_category" category ON asset.category_id = category.asset_category_id
		INNER JOIN "asset-service"."asset_status" status ON asset.status_id = status.asset_status_id
		WHERE asset.user_client_id = ? AND asset.category_id = ? AND asset.deleted_at IS NULL AND asset.is_wishlist = true
		ORDER BY asset.asset_id ASC;
	`
	rows, err := r.db.Raw(selectQuery, clientID, categoryID).Rows()
	if err != nil {
		return nil, err
	}

	var assetsResult []response.AssetResponse
	for rows.Next() {
		var assetResult response.AssetResponse
		var categoryResult response.AssetCategoryResponse
		var statusResult response.AssetStatusResponse

		err := rows.Scan(
			&assetResult.AssetID,
			&assetResult.UserClientID,
			assetResult.Name,
			assetResult.Description,
			assetResult.Price,
			categoryResult.AssetCategoryID,
			categoryResult.CategoryName,
			categoryResult.Description,
			statusResult.AssetStatusID,
			statusResult.StatusName,
			statusResult.Description,
			assetResult.PurchaseDate,
			assetResult.Price,
		)

		if err != nil {
			return nil, err
		}

		assetResult.Category = categoryResult
		assetResult.Status = statusResult

		assetsResult = append(assetsResult, assetResult)
	}

	return assetsResult, nil
}

func (r assetWishlistRepository) GetAssetWishlistByStatus(clientID string, statusID uint) ([]response.AssetResponse, error) {
	selectQuery := `
		SELECT 
			asset.asset_id,
			asset.user_client_id,
			asset.name,
			asset.description,
			asset.price,
			category.asset_category_id,
			category.category_name,
			category.description AS category_description,
			status.asset_status_id,
			status.status_name,
			status.description AS status_description,
			asset.purchase_date
		FROM "asset-service"."asset" asset
		INNER JOIN "asset-service"."asset_category" category ON asset.category_id = category.asset_category_id
		INNER JOIN "asset-service"."asset_status" status ON asset.status_id = status.asset_status_id
		WHERE asset.user_client_id = ? AND asset.status_id = ? AND asset.deleted_at IS NULL AND asset.is_wishlist = true
		ORDER BY asset.asset_id ASC;
	`
	rows, err := r.db.Raw(selectQuery, clientID, statusID).Rows()
	if err != nil {
		return nil, err
	}

	var assetsResult []response.AssetResponse
	for rows.Next() {
		var assetResult response.AssetResponse
		var categoryResult response.AssetCategoryResponse
		var statusResult response.AssetStatusResponse

		// Handling NULL values from SQL
		var description sql.NullString
		var purchaseDate sql.NullTime
		var price sql.NullFloat64

		err := rows.Scan(
			&assetResult.AssetID,
			&assetResult.UserClientID,
			&assetResult.Name,
			&description, // Handling NULL description
			&price,       // Handling NULL price
			&categoryResult.AssetCategoryID,
			&categoryResult.CategoryName,
			&categoryResult.Description,
			&statusResult.AssetStatusID,
			&statusResult.StatusName,
			&statusResult.Description,
			&purchaseDate, // Handling NULL purchaseDate
		)

		if err != nil {
			log.Error().Str("clientID", clientID).Err(err).Msg("❌ Failed to scan asset wishlist row")
			return nil, err
		}

		// Convert NULL SQL values to Go `nil`
		if description.Valid {
			assetResult.Description = description.String
		}
		if price.Valid {
			assetResult.Price = price.Float64
		}
		if purchaseDate.Valid {
			assetResult.PurchaseDate = (*response.DateOnly)(&purchaseDate.Time)
		}

		// Assign category and status details
		assetResult.Category = categoryResult
		assetResult.Status = statusResult
		assetsResult = append(assetsResult, assetResult)
	}

	return assetsResult, nil
}

func (r assetWishlistRepository) GetAssetWishlistByCategoryAndStatus(clientID string, categoryID, statusID uint) ([]response.AssetResponse, error) {
	selectQuery := `
		SELECT 
			asset.asset_id,
			asset.user_client_id,
			asset.name,
			asset.description,
			asset.price,
			category.asset_category_id,
			category.category_name,
			category.description AS category_description,
			status.asset_status_id,
			status.status_name,
			status.description AS status_description,
			asset.purchase_date
		FROM "asset-service"."asset" asset
		INNER JOIN "asset-service"."asset_category" category ON asset.category_id = category.asset_category_id
		INNER JOIN "asset-service"."asset_status" status ON asset.status_id = status.asset_status_id
		WHERE asset.user_client_id = ? AND asset.category_id = ? AND asset.status_id = ? AND asset.deleted_at IS NULL AND asset.is_wishlist = true
		ORDER BY asset.asset_id ASC;
	`
	rows, err := r.db.Raw(selectQuery, clientID, categoryID, statusID).Rows()
	if err != nil {
		return nil, err
	}

	var assetsResult []response.AssetResponse
	for rows.Next() {
		var assetResult response.AssetResponse
		var categoryResult response.AssetCategoryResponse
		var statusResult response.AssetStatusResponse

		// Handling NULL values from SQL
		var description sql.NullString
		var purchaseDate sql.NullTime
		var price sql.NullFloat64

		err := rows.Scan(
			&assetResult.AssetID,
			&assetResult.UserClientID,
			&assetResult.Name,
			&description, // Handling NULL description
			&price,       // Handling NULL price
			&categoryResult.AssetCategoryID,
			&categoryResult.CategoryName,
			&categoryResult.Description,
			&statusResult.AssetStatusID,
			&statusResult.StatusName,
			&statusResult.Description,
			&purchaseDate, // Handling NULL purchaseDate
		)

		if err != nil {
			log.Error().Str("clientID", clientID).Err(err).Msg("❌ Failed to scan asset wishlist row")
			return nil, err
		}

		// Convert NULL SQL values to Go `nil`
		if description.Valid {
			assetResult.Description = description.String
		}
		if price.Valid {
			assetResult.Price = price.Float64
		}
		if purchaseDate.Valid {
			assetResult.PurchaseDate = (*response.DateOnly)(&purchaseDate.Time)
		}

		// Assign category and status details
		assetResult.Category = categoryResult
		assetResult.Status = statusResult

		assetsResult = append(assetsResult, assetResult)
	}

	return assetsResult, nil
}
