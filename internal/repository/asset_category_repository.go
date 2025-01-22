package repository

import (
	"asset-service/internal/models/assets"
	"errors"
	"gorm.io/gorm"
)

type AssetCategoryRepository struct {
	DB *gorm.DB
}

func NewAssetCategoryRepository(db *gorm.DB) *AssetCategoryRepository {
	return &AssetCategoryRepository{DB: db}
}

func (r AssetCategoryRepository) AddAssetCategory(assetCategory **assets.AssetCategory) error {
	err := r.DB.Table("my-home.asset_category").Create(assetCategory).Error
	if err != nil {
		return err
	}
	return nil
}

func (r AssetCategoryRepository) UpdateAssetCategory(assetCategory **assets.AssetCategory) error {
	err := r.DB.Table("my-home.asset_category").Save(assetCategory).Error
	if err != nil {
		return err
	}
	return nil
}

func (r AssetCategoryRepository) GetAssetCategoryByName(name string) error {
	var assetCategory assets.AssetCategory
	r.DB.Table("my-home.asset_category").Where("category_name = ?", name).First(&assetCategory)
	if assetCategory.AssetCategoryID != 0 {
		return nil
	}
	return errors.New("assets category not found")
}

func (r AssetCategoryRepository) GetAssetCategoryById(assetCategoryID uint) (*assets.AssetCategory, error) {
	var assetCategory assets.AssetCategory
	r.DB.Table("my-home.asset_category").Where("asset_category_id = ?", assetCategoryID).First(&assetCategory)
	if assetCategory.AssetCategoryID != 0 {
		return &assetCategory, nil
	}
	return nil, errors.New("assets category not found")
}

func (r AssetCategoryRepository) GetAssetCategoryByIdAndNameNotExist(assetCategoryID uint, categoryName string) (*assets.AssetCategory, error) {
	var assetCategory assets.AssetCategory
	r.DB.Table("my-home.asset_category").Where("asset_category_id = ? AND category_name NOT LIKE ?", assetCategoryID, categoryName).First(&assetCategory)
	if assetCategory.AssetCategoryID != 0 {
		return &assetCategory, nil
	}
	return nil, errors.New("assets category not found")
}

func (r AssetCategoryRepository) GetListAssetCategory() ([]assets.AssetCategory, error) {
	var assetCategories []assets.AssetCategory
	err := r.DB.Table("my-home.asset_category").Find(&assetCategories).Error
	if err != nil {
		return nil, err
	}
	return assetCategories, nil
}

func (r AssetCategoryRepository) DeleteAssetCategory(assetCategory *assets.AssetCategory) error {
	err := r.DB.Table("my-home.asset_category").Model(&assetCategory).
		Where("asset_category_id = ?", assetCategory.AssetCategoryID).
		Update("deleted_by", assetCategory.DeletedBy).
		Delete(&assetCategory).Error
	if err != nil {
		return err
	}
	return nil
}
