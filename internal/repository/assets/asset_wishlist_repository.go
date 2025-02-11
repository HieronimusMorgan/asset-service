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

type AssetWishlistRepository interface {
	AddAssetWishlist(asset *assets.Asset) (*response.AssetResponseList, error)
	GetAssetWishlistByID(clientID string, assetID uint) (*response.AssetResponseList, error)
	GetAssetWishlistList(clientID string) ([]response.AssetResponseList, error)
	UpdateAssetWishlist(asset *assets.Asset) (*response.AssetResponseList, error)
	DeleteAssetWishlist(clientID, name string, assetID uint) error
	GetAssetWishlistByCategory(clientID string, categoryID uint) ([]response.AssetResponseList, error)
	GetAssetWishlistByStatus(clientID string, statusID uint) ([]response.AssetResponseList, error)
	GetAssetWishlistByCategoryAndStatus(clientID string, categoryID, statusID uint) ([]response.AssetResponseList, error)
}

type assetWishlistRepository struct {
	db    gorm.DB
	audit AssetAuditLogRepository
	asset AssetRepository
}

const tableAssetWishlistName = "my-home.asset"

func NewAssetWishlistRepository(db gorm.DB, audit AssetAuditLogRepository, asset AssetRepository) AssetWishlistRepository {
	return assetWishlistRepository{db: db, audit: audit, asset: asset}
}

func (r assetWishlistRepository) AddAssetWishlist(asset *assets.Asset) (*response.AssetResponseList, error) {
	if asset == nil {
		return nil, errors.New("assets cannot be nil")
	}

	// Validate CategoryID and StatusID
	if asset.CategoryID == 0 || asset.StatusID == 0 {
		return nil, errors.New("category_id and status_id cannot be null or zero")
	}

	// Verify the existence of CategoryID and StatusID in a single query
	var exists int
	err := r.db.Raw(
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
	tx := r.db.Begin()
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
        FROM "my-home"."asset" asset
        INNER JOIN "my-home"."asset_category" category ON asset.category_id = category.asset_category_id
        INNER JOIN "my-home"."asset_status" status ON asset.status_id = status.asset_status_id
        WHERE asset.user_client_id = ? AND asset.asset_id = ? AND asset.deleted_at IS NULL AND asset.is_wishlist = true
        ORDER BY asset.name DESC;
    `
	row := r.db.Raw(selectQuery, asset.UserClientID, asset.AssetID).Row()
	var assetResult response.AssetResponseList
	var categoryResult response.AssetCategoryResponse
	var statusResult response.AssetStatusResponse

	err = row.Scan(
		&assetResult.ID,
		&assetResult.ClientID,
		&assetResult.Name,
		&assetResult.Description,
		&assetResult.Price,
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

	err = r.audit.AfterCreateAsset(asset)
	if err != nil {
		return nil, err
	}

	return &assetResult, nil
}

func (r assetWishlistRepository) GetAssetWishlistByID(clientID string, assetID uint) (*response.AssetResponseList, error) {
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
		FROM "my-home"."asset" asset
		INNER JOIN "my-home"."asset_category" category ON asset.category_id = category.asset_category_id
		INNER JOIN "my-home"."asset_status" status ON asset.status_id = status.asset_status_id
		WHERE asset.user_client_id = ? AND asset.asset_id = ? AND asset.deleted_at IS NULL AND asset.is_wishlist = true
		ORDER BY asset.name DESC;
	`
	row := r.db.Raw(selectQuery, clientID, assetID).Row()
	var assetResult response.AssetResponseList
	var categoryResult response.AssetCategoryResponse
	var statusResult response.AssetStatusResponse

	err := row.Scan(
		&assetResult.ID,
		&assetResult.ClientID,
		&assetResult.Name,
		&assetResult.Description,
		&assetResult.Price,
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

func (r assetWishlistRepository) GetAssetWishlistList(clientID string) ([]response.AssetResponseList, error) {
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
		FROM "my-home"."asset" asset
		INNER JOIN "my-home"."asset_category" category ON asset.category_id = category.asset_category_id
		INNER JOIN "my-home"."asset_status" status ON asset.status_id = status.asset_status_id
		WHERE asset.user_client_id = ? AND asset.deleted_at IS NULL AND asset.is_wishlist = true
		ORDER BY asset.name DESC;
	`
	rows, err := r.db.Raw(selectQuery, clientID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assetsResult []response.AssetResponseList
	for rows.Next() {
		var assetResult response.AssetResponseList
		var categoryResult response.AssetCategoryResponse
		var statusResult response.AssetStatusResponse

		var description sql.NullString
		var purchaseDate sql.NullTime
		var price sql.NullFloat64

		err := rows.Scan(
			&assetResult.ID,
			&assetResult.ClientID,
			&assetResult.Name,
			&description,
			&price,
			&categoryResult.AssetCategoryID,
			&categoryResult.CategoryName,
			&categoryResult.Description,
			&statusResult.AssetStatusID,
			&statusResult.StatusName,
			&statusResult.Description,
			&purchaseDate,
			&price,
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

func (r assetWishlistRepository) UpdateAssetWishlist(asset *assets.Asset) (*response.AssetResponseList, error) {
	if asset == nil {
		return nil, errors.New("assets cannot be nil")
	}

	// Validate CategoryID and StatusID
	if asset.CategoryID == 0 || asset.StatusID == 0 {
		return nil, errors.New("category_id and status_id cannot be null or zero")
	}

	// Verify the existence of CategoryID and StatusID in a single query
	var exists int
	err := r.db.Raw(
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

	var assetOld assets.Asset
	err = r.db.Table(tableAssetName).Where("asset_id = ?", asset.AssetID).First(&assetOld).Error

	// Start a transaction for atomic operations
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update the assets
	if err := tx.Table(tableAssetName).
		Where("asset_id = ?", asset.AssetID).
		Updates(map[string]interface{}{
			"description":   asset.Description,
			"purchase_date": asset.PurchaseDate,
			"category_id":   asset.CategoryID,
			"status_id":     asset.StatusID,
			"price":         asset.Price,
			"is_wishlist":   asset.IsWishlist,
			"updated_by":    asset.UpdatedBy,
		}).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update assets: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

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
		FROM "my-home"."asset" asset
		INNER JOIN "my-home"."asset_category" category ON asset.category_id = category.asset_category_id
		INNER JOIN "my-home"."asset_status" status ON asset.status_id = status.asset_status_id
		WHERE asset.user_client_id = ? AND asset.asset_id = ? AND asset.deleted_at IS NULL AND asset.is_wishlist = true
		ORDER BY asset.name DESC;
	`
	row := r.db.Raw(selectQuery, asset.UserClientID, asset.AssetID).Row()
	var assetResult response.AssetResponseList
	var categoryResult response.AssetCategoryResponse
	var statusResult response.AssetStatusResponse

	err = row.Scan(
		&assetResult.ID,
		&assetResult.ClientID,
		&assetResult.Name,
		&assetResult.Description,
		&assetResult.Price,
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

	err = r.audit.AfterUpdateAsset(assetOld, asset)
	if err != nil {
		return nil, err
	}

	return &assetResult, nil
}

func (r assetWishlistRepository) DeleteAssetWishlist(clientID, name string, assetID uint) error {

	var asset assets.Asset

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Table(tableAssetName).Where("asset_id = ? AND user_client_id = ?", assetID, clientID).
		First(&asset).Error; err != nil {
		return fmt.Errorf("failed to find asset: %w", err)
	}

	// Soft delete the assets
	if err := tx.Table(tableAssetName).
		Where("asset_id = ?", assetID).
		Updates(map[string]interface{}{
			"deleted_at": gorm.Expr("NOW()"),
			"deleted_by": name,
		}).Delete(&assets.Asset{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete assets: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	err := r.audit.AfterDeleteAsset(&asset)
	if err != nil {
		return err
	}

	return nil
}

func (r assetWishlistRepository) GetAssetWishlistByCategory(clientID string, categoryID uint) ([]response.AssetResponseList, error) {
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
		FROM "my-home"."asset" asset
		INNER JOIN "my-home"."asset_category" category ON asset.category_id = category.asset_category_id
		INNER JOIN "my-home"."asset_status" status ON asset.status_id = status.asset_status_id
		WHERE asset.user_client_id = ? AND asset.category_id = ? AND asset.deleted_at IS NULL AND asset.is_wishlist = true
		ORDER BY asset.name DESC;
	`
	rows, err := r.db.Raw(selectQuery, clientID, categoryID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assetsResult []response.AssetResponseList
	for rows.Next() {
		var assetResult response.AssetResponseList
		var categoryResult response.AssetCategoryResponse
		var statusResult response.AssetStatusResponse

		err := rows.Scan(
			&assetResult.ID,
			&assetResult.ClientID,
			&assetResult.Name,
			&assetResult.Description,
			&assetResult.Price,
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

		assetsResult = append(assetsResult, assetResult)
	}

	return assetsResult, nil
}

func (r assetWishlistRepository) GetAssetWishlistByStatus(clientID string, statusID uint) ([]response.AssetResponseList, error) {
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
		FROM "my-home"."asset" asset
		INNER JOIN "my-home"."asset_category" category ON asset.category_id = category.asset_category_id
		INNER JOIN "my-home"."asset_status" status ON asset.status_id = status.asset_status_id
		WHERE asset.user_client_id = ? AND asset.status_id = ? AND asset.deleted_at IS NULL AND asset.is_wishlist = true
		ORDER BY asset.name DESC;
	`
	rows, err := r.db.Raw(selectQuery, clientID, statusID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assetsResult []response.AssetResponseList
	for rows.Next() {
		var assetResult response.AssetResponseList
		var categoryResult response.AssetCategoryResponse
		var statusResult response.AssetStatusResponse

		err := rows.Scan(
			&assetResult.ID,
			&assetResult.ClientID,
			&assetResult.Name,
			&assetResult.Description,
			&assetResult.Price,
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

		assetsResult = append(assetsResult, assetResult)
	}

	return assetsResult, nil
}

func (r assetWishlistRepository) GetAssetWishlistByCategoryAndStatus(clientID string, categoryID, statusID uint) ([]response.AssetResponseList, error) {
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
		FROM "my-home"."asset" asset
		INNER JOIN "my-home"."asset_category" category ON asset.category_id = category.asset_category_id
		INNER JOIN "my-home"."asset_status" status ON asset.status_id = status.asset_status_id
		WHERE asset.user_client_id = ? AND asset.category_id = ? AND asset.status_id = ? AND asset.deleted_at IS NULL AND asset.is_wishlist = true
		ORDER BY asset.name DESC;
	`
	rows, err := r.db.Raw(selectQuery, clientID, categoryID, statusID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assetsResult []response.AssetResponseList
	for rows.Next() {
		var assetResult response.AssetResponseList
		var categoryResult response.AssetCategoryResponse
		var statusResult response.AssetStatusResponse

		err := rows.Scan(
			&assetResult.ID,
			&assetResult.ClientID,
			&assetResult.Name,
			&assetResult.Description,
			&assetResult.Price,
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

		assetsResult = append(assetsResult, assetResult)
	}

	return assetsResult, nil
}
