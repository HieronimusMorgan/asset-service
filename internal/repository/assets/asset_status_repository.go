package assets

import (
	"asset-service/internal/models/assets"
	"asset-service/internal/utils"
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

func NewAssetStatusRepository(db gorm.DB) AssetStatusRepository {
	return assetStatusRepository{db: db}
}

func (r assetStatusRepository) GetAssetStatusByName(name string) error {
	var assetStatus assets.AssetStatus
	err := r.db.Table(utils.TableAssetStatus).Where("status_name LIKE ?", name).First(&assetStatus).Error
	if err != nil {
		return err
	}
	return nil
}

func (r assetStatusRepository) AddAssetStatus(assetStatus **assets.AssetStatus) error {
	err := r.db.Table(utils.TableAssetStatus).Create(assetStatus).Error
	if err != nil {
		return err
	}
	return nil
}

func (r assetStatusRepository) GetAssetStatus() ([]assets.AssetStatus, error) {
	var assetStatus []assets.AssetStatus
	err := r.db.Table(utils.TableAssetStatus).Find(&assetStatus).Where("deleted_at IS NULL").Error
	if err != nil {
		return nil, err
	}
	return assetStatus, nil
}

func (r assetStatusRepository) GetAssetStatusByID(assetStatusID uint) (*assets.AssetStatus, error) {
	var assetStatus assets.AssetStatus
	err := r.db.Table(utils.TableAssetStatus).Where("asset_status_id = ?", assetStatusID).First(&assetStatus).Error
	if err != nil {
		return nil, err
	}
	return &assetStatus, nil
}

func (r assetStatusRepository) UpdateAssetStatus(status *assets.AssetStatus) error {
	err := r.db.Table(utils.TableAssetStatus).Save(status).Error
	if err != nil {
		return err
	}
	return nil
}

func (r assetStatusRepository) DeleteAssetStatus(status *assets.AssetStatus) error {
	err := r.db.Table(utils.TableAssetStatus).Model(status).
		Update("deleted_by", status.DeletedBy).
		Delete(status).Error
	if err != nil {
		return err
	}
	return nil
}
