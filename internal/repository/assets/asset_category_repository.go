package assets

import (
	"asset-service/internal/models/assets"
	"errors"
	"gorm.io/gorm"
)

type AssetCategoryRepository interface {
	AddAssetCategory(assetCategory *assets.AssetCategory) error
	UpdateAssetCategory(assetCategory *assets.AssetCategory) error
	GetAssetCategoryByName(name string) error
	GetAssetCategoryById(assetCategoryID uint) (*assets.AssetCategory, error)
	GetAssetCategoryByIdAndNameNotExist(assetCategoryID uint, categoryName string) (*assets.AssetCategory, error)
	GetListAssetCategory() ([]assets.AssetCategory, error)
	DeleteAssetCategory(assetCategory *assets.AssetCategory) error
}
type assetCategoryRepository struct {
	db       gorm.DB
	logAudit AssetAuditLogRepository
}

const tableAssetCategoryName = "my-home.asset_category"

func NewAssetCategoryRepository(db gorm.DB, logAudit AssetAuditLogRepository) AssetCategoryRepository {
	return &assetCategoryRepository{db: db, logAudit: logAudit}
}

func (r assetCategoryRepository) AddAssetCategory(assetCategory *assets.AssetCategory) error {
	err := r.db.Table(tableAssetCategoryName).Create(&assetCategory).Error
	if err != nil {
		return err
	}

	// Log audit
	err = r.logAudit.AfterCreateAssetCategory(assetCategory)
	if err != nil {
		return err
	}

	return nil
}

func (r assetCategoryRepository) UpdateAssetCategory(assetCategory *assets.AssetCategory) error {
	oldAssetCategory, err := r.GetAssetCategoryById(assetCategory.AssetCategoryID)
	if err != nil {
		return err
	}

	err = r.db.Table(tableAssetCategoryName).Save(assetCategory).Error
	if err != nil {
		return err
	}

	err = r.logAudit.AfterUpdateAssetCategory(oldAssetCategory, assetCategory)
	if err != nil {
		return err
	}

	return nil
}

func (r assetCategoryRepository) GetAssetCategoryByName(name string) error {
	var assetCategory assets.AssetCategory
	r.db.Table(tableAssetCategoryName).Where("category_name = ?", name).First(&assetCategory)
	if assetCategory.AssetCategoryID != 0 {
		return nil
	}
	return errors.New("assets category not found")
}

func (r assetCategoryRepository) GetAssetCategoryById(assetCategoryID uint) (*assets.AssetCategory, error) {
	var assetCategory assets.AssetCategory
	r.db.Table(tableAssetCategoryName).Where("asset_category_id = ?", assetCategoryID).First(&assetCategory)
	if assetCategory.AssetCategoryID != 0 {
		return &assetCategory, nil
	}
	return nil, errors.New("assets category not found")
}

func (r assetCategoryRepository) GetAssetCategoryByIdAndNameNotExist(assetCategoryID uint, categoryName string) (*assets.AssetCategory, error) {
	var assetCategory assets.AssetCategory
	r.db.Table(tableAssetCategoryName).Where("asset_category_id = ? AND category_name NOT LIKE ?", assetCategoryID, categoryName).First(&assetCategory)
	if assetCategory.AssetCategoryID != 0 {
		return &assetCategory, nil
	}
	return nil, errors.New("assets category not found")
}

func (r assetCategoryRepository) GetListAssetCategory() ([]assets.AssetCategory, error) {
	var assetCategories []assets.AssetCategory
	err := r.db.Table(tableAssetCategoryName).Find(&assetCategories).Error
	if err != nil {
		return nil, err
	}
	return assetCategories, nil
}

func (r assetCategoryRepository) DeleteAssetCategory(assetCategory *assets.AssetCategory) error {
	err := r.db.Table(tableAssetCategoryName).Model(&assetCategory).
		Where("asset_category_id = ?", assetCategory.AssetCategoryID).
		Update("deleted_by", assetCategory.DeletedBy).
		Delete(&assetCategory).Error
	if err != nil {
		return err
	}
	return nil
}
