package repository

import (
	"asset-service/internal/dto/out"
	"asset-service/internal/models/asset"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type AssetRepository struct {
	DB *gorm.DB
}

const tableAssetName = "my-home.asset"

func NewAssetRepository(db *gorm.DB) *AssetRepository {
	return &AssetRepository{DB: db}
}

func (r AssetRepository) RegisterAsset(asset **asset.Asset) error {
	err := r.DB.Table(tableAssetName).Create(&asset).Error
	if err != nil {
		return err
	}
	return nil
}

func (r AssetRepository) AddAsset(asset *asset.Asset, maintenance *asset.AssetMaintenance) (*out.AssetResponse, error) {
	if asset == nil {
		return nil, errors.New("asset cannot be nil")
	}

	// Validate CategoryID and StatusID
	if asset.CategoryID == 0 || asset.StatusID == 0 {
		return nil, errors.New("category_id and status_id cannot be null or zero")
	}

	// Verify the existence of CategoryID and StatusID in a single query
	var exists int
	err := r.DB.Raw(
		`SELECT 
			(SELECT COUNT(*) FROM "my-home"."asset_category" WHERE asset_category_id = ?) + 
			(SELECT COUNT(*) FROM "my-home"."asset_status" WHERE asset_status_id = ?) AS exists`,
		asset.CategoryID, asset.StatusID).Scan(&exists).Error
	if err != nil {
		return nil, fmt.Errorf("failed to verify category_id or status_id: %w", err)
	}
	if exists != 2 {
		return nil, errors.New("invalid category_id or status_id")
	}

	// Start a transaction for atomic operations
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Insert the asset
	if err := tx.Table("my-home.asset").Create(asset).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create asset: %w", err)
	}

	// Set asset_id for maintenance and insert the maintenance record
	maintenance.AssetID = int(asset.AssetID)
	if err := tx.Table("my-home.asset_maintenance").Create(maintenance).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create maintenance record: %w", err)
	}

	// Retrieve the asset with maintenance information
	var result out.AssetResponse
	selectAssetQuery := `
		SELECT 
			a.asset_id,
			a.user_client_id,
			a.name,
			a.description,
			c.category_name,
			s.status_name,
			a.purchase_date,
			a.value,
			m.maintenance_date,
			m.maintenance_cost
		FROM "my-home"."asset" a
		INNER JOIN "my-home"."asset_category" c ON a.category_id = c.asset_category_id
		INNER JOIN "my-home"."asset_status" s ON a.status_id = s.asset_status_id
		LEFT JOIN "my-home"."asset_maintenance" m ON a.asset_id = m.asset_id
		WHERE a.asset_id = ?
	`
	if err := tx.Raw(selectAssetQuery, asset.AssetID).Scan(&result).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to retrieve asset after creation: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &result, nil
}

func (r AssetRepository) GetAssetByName(name string) (*asset.Asset, error) {
	var asset asset.Asset
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
			asset.value,
			maintenance.maintenance_date,
			maintenance.maintenance_cost
		FROM "my-home"."asset" asset
		INNER JOIN "my-home"."asset_category" category ON asset.category_id = category.asset_category_id
		INNER JOIN "my-home"."asset_status" status ON asset.status_id = status.asset_status_id
		LEFT JOIN "my-home"."asset_maintenance" maintenance ON asset.asset_id = maintenance.asset_id
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
