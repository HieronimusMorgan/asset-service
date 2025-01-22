package repository

import (
	"asset-service/internal/models/asset"
	"errors"
	"gorm.io/gorm"
)

type AssetCategoryRepository struct {
	DB *gorm.DB
}

func NewAssetCategoryRepository(db *gorm.DB) *AssetCategoryRepository {
	return &AssetCategoryRepository{DB: db}
}

func (r AssetCategoryRepository) AddAssetCategory(assetCategory **asset.AssetCategory) error {
	err := r.DB.Table("my-home.asset_category").Create(assetCategory).Error
	if err != nil {
		return err
	}
	return nil
}

func (r AssetCategoryRepository) UpdateAssetCategory(assetCategory **asset.AssetCategory) error {
	err := r.DB.Table("my-home.asset_category").Save(assetCategory).Error
	if err != nil {
		return err
	}
	return nil
}

func (r AssetCategoryRepository) GetAssetCategoryByName(name string) error {
	var assetCategory asset.AssetCategory
	r.DB.Table("my-home.asset_category").Where("category_name = ?", name).First(&assetCategory)
	if assetCategory.AssetCategoryID != 0 {
		return nil
	}
	return errors.New("asset category not found")
}

func (r AssetCategoryRepository) GetAssetCategoryById(assetCategoryID uint) (*asset.AssetCategory, error) {
	var assetCategory asset.AssetCategory
	r.DB.Table("my-home.asset_category").Where("asset_category_id = ?", assetCategoryID).First(&assetCategory)
	if assetCategory.AssetCategoryID != 0 {
		return &assetCategory, nil
	}
	return nil, errors.New("asset category not found")
}

func (r AssetCategoryRepository) GetAssetCategoryByIdAndNameNotExist(assetCategoryID uint, categoryName string) (*asset.AssetCategory, error) {
	var assetCategory asset.AssetCategory
	r.DB.Table("my-home.asset_category").Where("asset_category_id = ? AND category_name NOT LIKE ?", assetCategoryID, categoryName).First(&assetCategory)
	if assetCategory.AssetCategoryID != 0 {
		return &assetCategory, nil
	}
	return nil, errors.New("asset category not found")
}

func (r AssetCategoryRepository) GetListAssetCategory() ([]asset.AssetCategory, error) {
	var assetCategories []asset.AssetCategory
	err := r.DB.Table("my-home.asset_category").Find(&assetCategories).Error
	if err != nil {
		return nil, err
	}
	return assetCategories, nil
}

func (r AssetCategoryRepository) DeleteAssetCategory(assetCategory *asset.AssetCategory) error {
	err := r.DB.Table("my-home.asset_category").Model(&assetCategory).
		Where("asset_category_id = ?", assetCategory.AssetCategoryID).
		Update("deleted_by", assetCategory.DeletedBy).
		Delete(&assetCategory).Error
	if err != nil {
		return err
	}
	return nil
}
