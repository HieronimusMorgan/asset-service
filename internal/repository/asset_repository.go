package repository

import (
	"asset-service/internal/dto/out"
	"asset-service/internal/models"
	"gorm.io/gorm"
)

type AssetRepository struct {
	DB *gorm.DB
}

const tableAssetName = "asset-service.asset"

func NewAssetRepository(db *gorm.DB) *AssetRepository {
	return &AssetRepository{DB: db}
}

func (r AssetRepository) RegisterAsset(asset **models.Asset) error {
	err := r.DB.Table(tableAssetName).Create(&asset).Error
	if err != nil {
		return err
	}
	return nil
}

func (r AssetRepository) AddAsset(asset *models.Asset) (*out.AssetResponse, error) {
	err := r.DB.Table("asset-service.asset").Create(asset).Error
	if err != nil {
		return nil, err
	}

	selectQuery := `
        SELECT 
            a.asset_id,
            a.user_client_id,
            a.name,
            a.description,
            c.category_name,
            s.status_name,
            a.purchase_date,
            a.value
        
        FROM "asset-service"."asset" a
        INNER JOIN "asset-service"."asset_category" c ON a.category_id = c.asset_category_id
        INNER JOIN "asset-service"."asset_status" s ON a.status_id = s.asset_status_id
        WHERE a.asset_id = ?
    `
	var result out.AssetResponse
	err = r.DB.Raw(selectQuery, asset.AssetID).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r AssetRepository) GetAssetByName(name string) (*models.Asset, error) {
	var asset models.Asset
	err := r.DB.Table(tableAssetName).Where("name LIKE ?", name).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r AssetRepository) GetListAsset(clientID string) ([]out.AssetResponse, error) {
	selectQuery := `
		SELECT 
			asset.asset_id,
			asset.name,
			asset.description,
			category.category_name,
			status.status_name,
			asset.purchase_date,
			asset.value
		
		FROM "asset-service"."asset" asset
		INNER JOIN "asset-service"."asset_category" category ON asset.category_id = category.asset_category_id
		INNER JOIN "asset-service"."asset_status" status ON asset.status_id = status.asset_status_id
		WHERE asset.user_client_id = ? AND asset.deleted_at IS NULL
		ORDER BY asset.name DESC
	`
	var result []out.AssetResponse
	err := r.DB.Raw(selectQuery, clientID).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}
