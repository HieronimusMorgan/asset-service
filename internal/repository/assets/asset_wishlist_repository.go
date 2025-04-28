package assets

import (
	response "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	"asset-service/internal/utils"
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type AssetWishlistRepository interface {
	AssetWishlistNameExists(assetName string, clientID string) (bool, error)
	AddAssetWishlist(asset *assets.AssetWishlist) error
	GetAssetWishlistByID(clientID string, assetWishlistID uint) (*assets.AssetWishlist, error)
	GetAssetWishlistResponseByID(clientID string, assetWishlistID uint) (*response.AssetWishlistResponse, error)
	GetListAssetWishlist(clientID string, size, index int) ([]response.AssetWishlistResponse, error)
	UpdateAssetWishlist(assetWishlist *assets.AssetWishlist) error
	DeleteAssetWishlist(id uint, clientID string) error
	GetListAssetWishlistCount(clientID string) (int64, error)
}

type assetWishlistRepository struct {
	db    gorm.DB
	audit AssetAuditLogRepository
}

func NewAssetWishlistRepository(db gorm.DB, audit AssetAuditLogRepository) AssetWishlistRepository {
	return assetWishlistRepository{db: db, audit: audit}
}

func (r assetWishlistRepository) AssetWishlistNameExists(assetName string, clientID string) (bool, error) {
	var count int64
	err := r.db.Table(utils.TableAssetWishlistName).
		Where("user_client_id = ? AND asset_name = ? AND deleted_at IS NULL", clientID, assetName).
		Count(&count).Error

	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("❌ Failed to check asset wishlist name existence")
		return false, err
	}

	if count > 0 {
		log.Info().Str("clientID", clientID).Msg("✅ Asset wishlist name already exists")
		return true, nil
	}

	log.Info().Str("clientID", clientID).Msg("✅ Asset wishlist name does not exist")
	return false, nil
}

func (r assetWishlistRepository) AddAssetWishlist(asset *assets.AssetWishlist) error {
	tx := r.db.Begin()
	defer tx.Rollback()

	if err := tx.Table(utils.TableAssetWishlistName).Create(&asset).Error; err != nil {
		return fmt.Errorf("failed to create assets: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r assetWishlistRepository) GetAssetWishlistByID(clientID string, assetWishlistID uint) (*assets.AssetWishlist, error) {
	selectQuery := `
		SELECT 
			wishlist_id,
			user_client_id,
			asset_name,
			serial_number,
			barcode,
			category_id,
			status_id,
			priority_level,
			price_estimate,
			notes
		FROM "asset_wishlist"
		WHERE user_client_id = ? AND wishlist_id = ? AND deleted_at IS NULL;
	`

	var asset assets.AssetWishlist

	err := r.db.Raw(selectQuery, clientID, assetWishlistID).Scan(&asset).Error
	if err != nil {
		log.Error().Uint("assetWishlistID", assetWishlistID).Str("clientID", clientID).Err(err).Msg("❌ Failed to retrieve asset wishlist")
		return nil, err
	}

	log.Info().Uint("assetWishlistID", assetWishlistID).Str("clientID", clientID).Msg("✅ Successfully retrieved asset wishlist")
	return &asset, nil
}

func (r assetWishlistRepository) GetAssetWishlistResponseByID(clientID string, assetWishlistID uint) (*response.AssetWishlistResponse, error) {
	selectQuery := `
        SELECT 
            aw.wishlist_id,
            aw.user_client_id,
            aw.asset_name,
            aw.serial_number,
            aw.barcode,
            ac.asset_category_id,
            ac.category_name,
            ac.description AS category_description,
            ass.asset_status_id,
            ass.status_name,
            ass.description AS status_description,
            aw.priority_level,
            aw.price_estimate,
            aw.notes
        FROM "asset_wishlist" aw
        INNER JOIN "asset_category" ac ON aw.category_id = ac.asset_category_id
        INNER JOIN "asset_status" ass ON aw.status_id = ass.asset_status_id
        WHERE aw.user_client_id = ? AND aw.wishlist_id = ? AND aw.deleted_at IS NULL
        ORDER BY aw.wishlist_id ASC;
    `

	row := r.db.Raw(selectQuery, clientID, assetWishlistID).Row()

	var asset response.AssetWishlistResponse
	var status response.AssetStatusResponse
	var category response.AssetCategoryResponse

	err := row.Scan(
		&asset.WishlistID,
		&asset.UserClientID,
		&asset.AssetName,
		&asset.SerialNumber,
		&asset.Barcode,
		&category.AssetCategoryID,
		&category.CategoryName,
		&category.Description,
		&status.AssetStatusID,
		&status.StatusName,
		&status.Description,
		&asset.PriorityLevel,
		&asset.PriceEstimate,
		&asset.Notes,
	)

	if err != nil {
		log.Error().Uint("assetWishlistID", assetWishlistID).Str("clientID", clientID).Err(err).Msg("❌ Failed to scan asset wishlist")
		return nil, err
	}

	asset.Status = status
	asset.Category = category

	log.Info().Uint("assetWishlistID", assetWishlistID).Str("clientID", clientID).Msg("✅ Successfully retrieved asset wishlist")
	return &asset, nil
}

func (r assetWishlistRepository) GetListAssetWishlist(clientID string, size, index int) ([]response.AssetWishlistResponse, error) {
	selectQuery := `
		SELECT
			aw.wishlist_id,
			aw.user_client_id,
			aw.asset_name,
			aw.serial_number,
			aw.barcode,
			ac.asset_category_id,
			ac.category_name,
			ac.description AS category_description,
			ass.asset_status_id,
			ass.status_name,
			ass.description AS status_description,
			aw.priority_level,
			aw.price_estimate,
			aw.notes
		FROM "asset_wishlist" aw
		INNER JOIN "asset_category" ac ON aw.category_id = ac.asset_category_id
		INNER JOIN "asset_status" ass ON aw.status_id = ass.asset_status_id
		WHERE aw.user_client_id = ? AND aw.deleted_at IS NULL
		ORDER BY aw.wishlist_id ASC
		LIMIT ? OFFSET ?;
	`

	var assetWishlists []response.AssetWishlistResponse

	rows, err := r.db.Raw(selectQuery, clientID, size, (index-1)*size).Rows()
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("❌ Failed to retrieve asset wishlist")
		return nil, err
	}

	for rows.Next() {
		var asset response.AssetWishlistResponse
		var status response.AssetStatusResponse
		var category response.AssetCategoryResponse

		err := rows.Scan(
			&asset.WishlistID,
			&asset.UserClientID,
			&asset.AssetName,
			&asset.SerialNumber,
			&asset.Barcode,
			&category.AssetCategoryID,
			&category.CategoryName,
			&category.Description,
			&status.AssetStatusID,
			&status.StatusName,
			&status.Description,
			&asset.PriorityLevel,
			&asset.PriceEstimate,
			&asset.Notes,
		)
		if err != nil {
			log.Error().Str("clientID", clientID).Err(err).Msg("❌ Failed to scan asset wishlist row")
			return nil, err
		}

		asset.Status = status
		asset.Category = category
		assetWishlists = append(assetWishlists, asset)
	}

	log.Info().Str("clientID", clientID).Msg("✅ Successfully retrieved asset wishlist")
	return assetWishlists, nil
}

func (r assetWishlistRepository) UpdateAssetWishlist(assetWishlist *assets.AssetWishlist) error {
	tx := r.db.Begin()
	defer tx.Rollback()

	if err := tx.Table(utils.TableAssetWishlistName).
		Where("wishlist_id = ? AND user_client_id = ? AND deleted_at IS NULL", assetWishlist.WishlistID, assetWishlist.UserClientID).
		Updates(map[string]interface{}{
			"asset_name":     assetWishlist.AssetName,
			"serial_number":  assetWishlist.SerialNumber,
			"barcode":        assetWishlist.Barcode,
			"category_id":    assetWishlist.CategoryID,
			"status_id":      assetWishlist.StatusID,
			"priority_level": assetWishlist.PriorityLevel,
			"price_estimate": assetWishlist.PriceEstimate,
			"notes":          assetWishlist.Notes,
			"updated_by":     assetWishlist.UpdatedBy,
			"updated_at":     gorm.Expr("NOW()"),
		}).Error; err != nil {
		return fmt.Errorf("failed to update asset wishlist: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r assetWishlistRepository) DeleteAssetWishlist(id uint, clientID string) error {
	tx := r.db.Begin()
	defer tx.Rollback()

	if err := tx.Table(utils.TableAssetWishlistName).
		Where("wishlist_id = ? AND user_client_id = ? AND deleted_at IS NULL", id, clientID).
		Updates(map[string]interface{}{
			"deleted_by": clientID,
			"deleted_at": gorm.Expr("NOW()"),
		}).Error; err != nil {
		return fmt.Errorf("failed to delete asset wishlist: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r assetWishlistRepository) GetListAssetWishlistCount(clientID string) (int64, error) {
	var count int64
	err := r.db.Table(utils.TableAssetWishlistName).
		Where("user_client_id = ? AND deleted_at IS NULL", clientID).
		Count(&count).Error

	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("❌ Failed to count asset wishlists")
		return 0, err
	}

	log.Info().Str("clientID", clientID).Msg("✅ Successfully counted asset wishlists")
	return count, nil
}
