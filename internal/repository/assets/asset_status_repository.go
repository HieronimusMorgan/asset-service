package assets

import (
	response "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	"asset-service/internal/utils"
	"gorm.io/gorm"
)

type AssetStatusRepository interface {
	GetAssetStatusByName(name string) error
	GetCountAssetStatus() (int64, error)
	AddAssetStatus(assetStatus **assets.AssetStatus) error
	GetAssetStatus(size, index int) ([]response.AssetStatusResponse, error)
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
	err := r.db.Table(utils.TableAssetStatusName).Where("status_name LIKE ?", name).First(&assetStatus).Error
	if err != nil {
		return err
	}
	return nil
}

func (r assetStatusRepository) GetCountAssetStatus() (int64, error) {
	var count int64
	err := r.db.Table(utils.TableAssetStatusName).Where("deleted_at IS NULL").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r assetStatusRepository) AddAssetStatus(assetStatus **assets.AssetStatus) error {
	err := r.db.Table(utils.TableAssetStatusName).Create(assetStatus).Error
	if err != nil {
		return err
	}
	return nil
}

func (r assetStatusRepository) GetAssetStatus(size, index int) ([]response.AssetStatusResponse, error) {
	var assetStatus []response.AssetStatusResponse
	offset := (index - 1) * size
	err := r.db.Table(utils.TableAssetStatusName).
		Where("deleted_at IS NULL").
		Limit(size).
		Offset(offset).
		Find(&assetStatus).Error
	if err != nil {
		return nil, err
	}
	return assetStatus, nil
}

func (r assetStatusRepository) GetAssetStatusByID(assetStatusID uint) (*assets.AssetStatus, error) {
	var assetStatus assets.AssetStatus
	err := r.db.Table(utils.TableAssetStatusName).Where("asset_status_id = ?", assetStatusID).First(&assetStatus).Error
	if err != nil {
		return nil, err
	}
	return &assetStatus, nil
}

func (r assetStatusRepository) UpdateAssetStatus(status *assets.AssetStatus) error {
	err := r.db.Table(utils.TableAssetStatusName).Save(status).Error
	if err != nil {
		return err
	}
	return nil
}

func (r assetStatusRepository) DeleteAssetStatus(status *assets.AssetStatus) error {
	err := r.db.Table(utils.TableAssetStatusName).Model(status).
		Update("deleted_by", status.DeletedBy).
		Delete(status).Error
	if err != nil {
		return err
	}
	return nil
}
