package assets

import (
	"asset-service/internal/models/assets"
	"gorm.io/gorm"
)

type AssetStatusRepository interface {
	GetAssetStatusByName(name string) error
	AddAssetStatus(assetStatus **assets.AssetStatus) error
	GetAssetStatus() ([]assets.AssetStatus, error)
	GetAssetStatusByID(assetStatusID uint) (*assets.AssetStatus, error)
	UpdateAssetStatus(status *assets.AssetStatus) error
	DeleteAssetStatus(status *assets.AssetStatus) error
}

type assetStatusRepository struct {
	db gorm.DB
}

const tableAssetStatus = "my-home.asset_status"

func NewAssetStatusRepository(db gorm.DB) AssetStatusRepository {
	return assetStatusRepository{db: db}
}

func (r assetStatusRepository) GetAssetStatusByName(name string) error {
	var assetStatus assets.AssetStatus
	err := r.db.Table(tableAssetStatus).Where("status_name LIKE ?", name).First(&assetStatus).Error
	if err != nil {
		return err
	}
	return nil
}

func (r assetStatusRepository) AddAssetStatus(assetStatus **assets.AssetStatus) error {
	err := r.db.Table(tableAssetStatus).Create(assetStatus).Error
	if err != nil {
		return err
	}
	return nil
}

func (r assetStatusRepository) GetAssetStatus() ([]assets.AssetStatus, error) {
	var assetStatus []assets.AssetStatus
	err := r.db.Table(tableAssetStatus).Find(&assetStatus).Where("deleted_at IS NULL").Error
	if err != nil {
		return nil, err
	}
	return assetStatus, nil
}

func (r assetStatusRepository) GetAssetStatusByID(assetStatusID uint) (*assets.AssetStatus, error) {
	var assetStatus assets.AssetStatus
	err := r.db.Table(tableAssetStatus).Where("asset_status_id = ?", assetStatusID).First(&assetStatus).Error
	if err != nil {
		return nil, err
	}
	return &assetStatus, nil
}

func (r assetStatusRepository) UpdateAssetStatus(status *assets.AssetStatus) error {
	err := r.db.Table(tableAssetStatus).Save(status).Error
	if err != nil {
		return err
	}
	return nil
}

func (r assetStatusRepository) DeleteAssetStatus(status *assets.AssetStatus) error {
	err := r.db.Table(tableAssetStatus).Model(status).
		Update("deleted_by", status.DeletedBy).
		Delete(status).Error
	if err != nil {
		return err
	}
	return nil
}
