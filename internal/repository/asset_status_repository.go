package repository

import (
	"asset-service/internal/models/asset"
	"gorm.io/gorm"
)

type AssetStatusRepository struct {
	DB *gorm.DB
}

func (r AssetStatusRepository) GetAssetStatusByName(name string) error {
	var assetStatus asset.AssetStatus
	err := r.DB.Table("my-home.asset_status").Where("status_name LIKE ?", name).First(&assetStatus).Error
	if err != nil {
		return err
	}
	return nil
}

func (r AssetStatusRepository) AddAssetStatus(assetStatus **asset.AssetStatus) error {
	err := r.DB.Table("my-home.asset_status").Create(assetStatus).Error
	if err != nil {
		return err
	}
	return nil
}

func (r AssetStatusRepository) GetAssetStatus() ([]asset.AssetStatus, error) {
	var assetStatus []asset.AssetStatus
	err := r.DB.Table("my-home.asset_status").Find(&assetStatus).Where("deleted_at IS NULL").Error
	if err != nil {
		return nil, err
	}
	return assetStatus, nil
}

func (r AssetStatusRepository) GetAssetStatusByID(assetStatusID uint) (*asset.AssetStatus, error) {
	var assetStatus asset.AssetStatus
	err := r.DB.Table("my-home.asset_status").Where("asset_status_id = ?", assetStatusID).First(&assetStatus).Error
	if err != nil {
		return nil, err
	}
	return &assetStatus, nil
}

func (r AssetStatusRepository) UpdateAssetStatus(status *asset.AssetStatus) error {
	err := r.DB.Table("my-home.asset_status").Save(status).Error
	if err != nil {
		return err
	}
	return nil
}

func (r AssetStatusRepository) DeleteAssetStatus(status *asset.AssetStatus) error {
	err := r.DB.Table("my-home.asset_status").Model(status).
		Update("deleted_by", status.DeletedBy).
		Delete(status).Error
	if err != nil {
		return err
	}
	return nil
}

func NewAssetStatusRepository(db *gorm.DB) *AssetStatusRepository {
	return &AssetStatusRepository{DB: db}
}
