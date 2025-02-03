package repository

import (
	"asset-service/internal/dto/out"
	"asset-service/internal/models/assets"
	"gorm.io/gorm"
)

type AssetMaintenanceRepository struct {
	DB              *gorm.DB
	assetRepository *AssetRepository
}

const tableAssetMaintenanceName = "my-home.asset_maintenance"

func NewAssetMaintenanceRepository(db *gorm.DB) *AssetMaintenanceRepository {
	return &AssetMaintenanceRepository{DB: db, assetRepository: NewAssetRepository(db)}
}

func (r *AssetMaintenanceRepository) Create(clientID string, maintenance *assets.AssetMaintenance) error {
	asset, err := r.assetRepository.GetAssetByID(clientID, uint(maintenance.AssetID))
	if err != nil {
		return err
	}

	if asset == nil {
		return gorm.ErrRecordNotFound
	}

	return r.DB.Table(tableAssetMaintenanceName).Create(maintenance).Error
}

func (r *AssetMaintenanceRepository) GetByID(maintenanceID uint, clientID string) (*out.AssetMaintenanceResponse, error) {
	var maintenance out.AssetMaintenanceResponse

	assetMaintenance := `
		SELECT am.id, am.asset_id, am.maintenance_details, am.maintenance_date, am.maintenance_cost FROM "my-home"."asset_maintenance" am
		LEFT JOIN "my-home"."asset" a ON am.asset_id = a.asset_id
		         WHERE am.id = ? AND a.user_client_id = ?
`

	err := r.DB.Raw(assetMaintenance, maintenanceID, clientID).Scan(&maintenance).Error
	if err != nil {
		return nil, err
	}

	return &maintenance, nil
}

func (r *AssetMaintenanceRepository) GetByAssetID(assetID uint, clientID string) (*assets.AssetMaintenance, error) {
	var maintenance assets.AssetMaintenance

	assetMaintenance := `
		SELECT * FROM "my-home"."asset_maintenance" am
		LEFT JOIN "my-home"."asset" a ON am.asset_id = a.asset_id
		         WHERE a.asset_id = ? AND a.user_client_id = ?
`

	err := r.DB.Raw(assetMaintenance, assetID, clientID).Scan(&maintenance).Error
	if err != nil {
		return nil, err
	}

	return &maintenance, nil
}

func (r *AssetMaintenanceRepository) GetList() ([]out.AssetMaintenanceResponse, error) {
	var maintenances []out.AssetMaintenanceResponse

	assetMaintenance := `
		SELECT am.id, am.asset_id, am.maintenance_details, am.maintenance_date, am.maintenance_cost FROM "my-home"."asset_maintenance" am
`

	err := r.DB.Raw(assetMaintenance).Scan(&maintenances).Error
	if err != nil {
		return nil, err
	}

	return maintenances, nil
}

func (r *AssetMaintenanceRepository) Update(maintenance *assets.AssetMaintenance) error {
	return r.DB.Table(tableAssetMaintenanceName).Save(maintenance).Error
}

func (r *AssetMaintenanceRepository) Delete(maintenanceID uint) error {
	return r.DB.Table(tableAssetMaintenanceName).Delete(&assets.AssetMaintenance{}, maintenanceID).Error
}
