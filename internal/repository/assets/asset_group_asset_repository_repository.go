package assets

import (
	"asset-service/internal/models/assets"
	"asset-service/internal/utils"
	"gorm.io/gorm"
)

type AssetGroupAssetRepository interface {
	AddAssetGroupAsset(asset *assets.AssetGroupAsset) error
	UpdateAssetGroupAsset(asset *assets.AssetGroupAsset) error
	GetAssetGroupAssetByID(assetGroupID uint) (*assets.AssetGroupAsset, error)
	GetListAssetGroupAssetByID(assetGroupID uint) ([]assets.AssetGroupAsset, error)
	DeleteAssetGroupAsset(assetGroupID uint) error
}

type assetGroupAssetRepository struct {
	db    gorm.DB
	audit AssetAuditLogRepository
}

func NewAssetGroupAssetRepository(db gorm.DB, audit AssetAuditLogRepository) AssetGroupAssetRepository {
	return assetGroupAssetRepository{db: db, audit: audit}
}

func (r assetGroupAssetRepository) AddAssetGroupAsset(asset *assets.AssetGroupAsset) error {
	return r.db.Table(utils.TableAssetGroupAssetName).Create(asset).Error
}

func (r assetGroupAssetRepository) UpdateAssetGroupAsset(asset *assets.AssetGroupAsset) error {
	return r.db.Table(utils.TableAssetGroupAssetName).Save(asset).Error
}

func (r assetGroupAssetRepository) GetAssetGroupAssetByID(assetGroupID uint) (*assets.AssetGroupAsset, error) {
	var asset assets.AssetGroupAsset
	if err := r.db.Table(utils.TableAssetGroupAssetName).First(&asset, assetGroupID).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r assetGroupAssetRepository) GetListAssetGroupAssetByID(assetGroupID uint) ([]assets.AssetGroupAsset, error) {
	var groupAssets []assets.AssetGroupAsset
	if err := r.db.Table(utils.TableAssetGroupAssetName).
		Where("asset_group_id = ?", assetGroupID).
		Order("asset_group_id ASC").
		Find(&groupAssets).
		Error; err != nil {
		return nil, err
	}
	return groupAssets, nil
}

func (r assetGroupAssetRepository) DeleteAssetGroupAsset(assetGroupID uint) error {
	return r.db.Table(utils.TableAssetGroupAssetName).Delete(&assets.AssetGroupAsset{}, assetGroupID).Error
}
