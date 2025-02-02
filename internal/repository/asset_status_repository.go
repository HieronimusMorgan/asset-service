package repository

import (
	"asset-service/internal/models/assets"
	"gorm.io/gorm"
)

type AssetStatusRepository struct {
	DB *gorm.DB
}

func NewAssetStatusRepository(db *gorm.DB) *AssetStatusRepository {
	return &AssetStatusRepository{DB: db.Table("my-home.asset_status")}
}

func (r AssetStatusRepository) GetAssetStatusByName(name string) error {
	var assetStatus assets.AssetStatus
	err := r.DB.Where("status_name LIKE ?", name).First(&assetStatus).Error
	if err != nil {
		return err
	}
	return nil
}

func (r AssetStatusRepository) AddAssetStatus(assetStatus **assets.AssetStatus) error {
	err := r.DB.Create(assetStatus).Error
	if err != nil {
		return err
	}
	return nil
}

func (r AssetStatusRepository) GetAssetStatus() ([]assets.AssetStatus, error) {
	var assetStatus []assets.AssetStatus
	err := r.DB.Find(&assetStatus).Where("deleted_at IS NULL").Error
	if err != nil {
		return nil, err
	}
	return assetStatus, nil
}

func (r AssetStatusRepository) GetAssetStatusByID(assetStatusID uint) (*assets.AssetStatus, error) {
	var assetStatus assets.AssetStatus
	err := r.DB.Where("asset_status_id = ?", assetStatusID).First(&assetStatus).Error
	if err != nil {
		return nil, err
	}
	return &assetStatus, nil
}

func (r AssetStatusRepository) UpdateAssetStatus(status *assets.AssetStatus) error {
	err := r.DB.Save(status).Error
	if err != nil {
		return err
	}
	return nil
}

func (r AssetStatusRepository) DeleteAssetStatus(status *assets.AssetStatus) error {
	err := r.DB.Model(status).
		Update("deleted_by", status.DeletedBy).
		Delete(status).Error
	if err != nil {
		return err
	}
	return nil
}
