package assets

import (
	"asset-service/internal/models/assets"
	"asset-service/internal/utils"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// AssetCategoryRepository defines the interface
type AssetCategoryRepository interface {
	AddAssetCategory(assetCategory *assets.AssetCategory) error
	UpdateAssetCategory(assetCategory *assets.AssetCategory, clientID string) error
	GetAssetCategoryByNameAndClientID(name, clientID string) (*assets.AssetCategory, error)
	GetAssetCategoryById(assetCategoryID uint, clientID string) (*assets.AssetCategory, error)
	GetAssetCategoryByIdAndNameNotExist(assetCategoryID uint, categoryName string) (*assets.AssetCategory, error)
	GetListAssetCategory(clientID string) ([]assets.AssetCategory, error)
	DeleteAssetCategory(assetCategory *assets.AssetCategory) error
}

// assetCategoryRepository implementation
type assetCategoryRepository struct {
	db       gorm.DB
	logAudit AssetAuditLogRepository
}

// NewAssetCategoryRepository initializes the repository
func NewAssetCategoryRepository(db gorm.DB, logAudit AssetAuditLogRepository) AssetCategoryRepository {
	return &assetCategoryRepository{db: db, logAudit: logAudit}
}

// AddAssetCategory inserts a new asset category and logs audit
func (r *assetCategoryRepository) AddAssetCategory(assetCategory *assets.AssetCategory) error {
	err := r.db.Table(utils.TableAssetCategoryName).Create(&assetCategory).Error
	if err != nil {
		log.Error().Err(err).
			Str("category_name", assetCategory.CategoryName).
			Msg("❌ Failed to add asset category")
		return err
	}

	log.Info().
		Str("category_name", assetCategory.CategoryName).
		Msg("✅ Asset category added successfully")
	return nil
}

// UpdateAssetCategory modifies an existing asset category and logs changes
func (r *assetCategoryRepository) UpdateAssetCategory(assetCategory *assets.AssetCategory, clientID string) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error().Msg("🔥 Panic occurred, rolling back transaction!")
		}
	}()

	err := tx.Table(utils.TableAssetCategoryName).
		Where("asset_category_id = ? AND user_client_id = ?", assetCategory.AssetCategoryID, clientID).
		Updates(assetCategory).Error
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).
			Uint("asset_category_id", assetCategory.AssetCategoryID).
			Msg("❌ Failed to update asset category")
		return err
	}

	tx.Commit()
	log.Info().
		Uint("asset_category_id", assetCategory.AssetCategoryID).
		Msg("✅ Asset category updated successfully")
	return nil
}

// GetAssetCategoryByName fetches a category by name
func (r *assetCategoryRepository) GetAssetCategoryByNameAndClientID(name, clientID string) (*assets.AssetCategory, error) {
	var assetCategory assets.AssetCategory
	err := r.db.Table(utils.TableAssetCategoryName).
		Where("category_name = ? AND user_client_id = ?", name, clientID).
		First(&assetCategory).Error
	if err != nil {
		log.Warn().
			Str("category_name", name).
			Str("user_client_id", clientID).
			Msg("⚠ Asset category not found")
		return nil, err
	}
	return &assetCategory, nil
}

// GetAssetCategoryById retrieves a category by ID
func (r *assetCategoryRepository) GetAssetCategoryById(assetCategoryID uint, clientID string) (*assets.AssetCategory, error) {
	var assetCategory assets.AssetCategory
	err := r.db.Table(utils.TableAssetCategoryName).
		Where("asset_category_id = ? AND user_client_id = ?", assetCategoryID, clientID).
		First(&assetCategory).Error
	if err != nil {
		log.Warn().
			Uint("asset_category_id", assetCategoryID).
			Msg("⚠ Asset category not found")
		return nil, err
	}
	return &assetCategory, nil
}

// GetAssetCategoryByIdAndNameNotExist checks if category ID and name do not match
func (r *assetCategoryRepository) GetAssetCategoryByIdAndNameNotExist(assetCategoryID uint, categoryName string) (*assets.AssetCategory, error) {
	var assetCategory assets.AssetCategory
	err := r.db.Table(utils.TableAssetCategoryName).
		Where("asset_category_id = ? AND category_name NOT LIKE ?", assetCategoryID, categoryName).
		First(&assetCategory).Error
	if err != nil {
		return nil, err
	}
	return &assetCategory, nil
}

// GetListAssetCategory retrieves all categories
func (r *assetCategoryRepository) GetListAssetCategory(clientID string) ([]assets.AssetCategory, error) {
	var assetCategories []assets.AssetCategory
	err := r.db.Table(utils.TableAssetCategoryName).
		Where("user_client_id = ?", clientID).
		Find(&assetCategories).Error
	if err != nil {
		log.Error().Err(err).Msg("❌ Failed to retrieve asset categories")
		return nil, err
	}
	return assetCategories, nil
}

// DeleteAssetCategory marks a category as deleted
func (r *assetCategoryRepository) DeleteAssetCategory(assetCategory *assets.AssetCategory) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error().Msg("🔥 Panic occurred, rolling back transaction!")
		}
	}()

	err := tx.Table(utils.TableAssetCategoryName).
		Where("asset_category_id = ?", assetCategory.AssetCategoryID).
		Update("deleted_by", assetCategory.DeletedBy).
		Delete(&assetCategory).Error
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).
			Uint("asset_category_id", assetCategory.AssetCategoryID).
			Msg("❌ Failed to delete asset category")
		return err
	}

	tx.Commit()
	log.Info().
		Uint("asset_category_id", assetCategory.AssetCategoryID).
		Msg("✅ Asset category deleted successfully")
	return nil
}
