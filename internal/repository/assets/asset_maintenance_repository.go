package assets

import (
	assets2 "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	"gorm.io/gorm"
)

type AssetMaintenanceRepository interface {
	Create(clientID string, maintenance *assets.AssetMaintenance) error
	GetByID(maintenanceID uint, clientID string) (*assets2.AssetMaintenanceResponse, error)
	GetByAssetID(assetID uint, clientID string) (*assets.AssetMaintenance, error)
	GetList() ([]assets2.AssetMaintenanceResponse, error)
	Update(maintenance *assets.AssetMaintenance) error
	Delete(maintenanceID uint) error
}
type assetMaintenanceRepository struct {
	db              gorm.DB
	assetRepository AssetRepository
}

const tableAssetMaintenanceName = "my-home.asset_maintenance"

func NewAssetMaintenanceRepository(db gorm.DB, assetRepository AssetRepository) AssetMaintenanceRepository {
	return assetMaintenanceRepository{db: db, assetRepository: assetRepository}
}

func (r assetMaintenanceRepository) Create(clientID string, maintenance *assets.AssetMaintenance) error {
	asset, err := r.assetRepository.GetAssetByID(clientID, uint(maintenance.AssetID))
	if err != nil {
		return err
	}

	if asset == nil {
		return gorm.ErrRecordNotFound
	}

	return r.db.Table(tableAssetMaintenanceName).Create(maintenance).Error
}

func (r assetMaintenanceRepository) GetByID(maintenanceID uint, clientID string) (*assets2.AssetMaintenanceResponse, error) {
	var maintenance assets2.AssetMaintenanceResponse

	assetMaintenance := `
		SELECT am.id, am.asset_id, am.maintenance_details, am.maintenance_date, am.maintenance_cost FROM "my-home"."asset_maintenance" am
		LEFT JOIN "my-home"."asset" a ON am.asset_id = a.asset_id
		         WHERE am.id = ? AND a.user_client_id = ?
`

	err := r.db.Raw(assetMaintenance, maintenanceID, clientID).Scan(&maintenance).Error
	if err != nil {
		return nil, err
	}

	return &maintenance, nil
}

func (r assetMaintenanceRepository) GetByAssetID(assetID uint, clientID string) (*assets.AssetMaintenance, error) {
	var maintenance assets.AssetMaintenance

	assetMaintenance := `
		SELECT * FROM "my-home"."asset_maintenance" am
		LEFT JOIN "my-home"."asset" a ON am.asset_id = a.asset_id
		         WHERE a.asset_id = ? AND a.user_client_id = ?
`

	err := r.db.Raw(assetMaintenance, assetID, clientID).Scan(&maintenance).Error
	if err != nil {
		return nil, err
	}

	return &maintenance, nil
}

func (r assetMaintenanceRepository) GetList() ([]assets2.AssetMaintenanceResponse, error) {
	var maintenances []assets2.AssetMaintenanceResponse

	assetMaintenance := `
		SELECT am.id, am.asset_id, am.maintenance_details, am.maintenance_date, am.maintenance_cost FROM "my-home"."asset_maintenance" am
`

	err := r.db.Raw(assetMaintenance).Scan(&maintenances).Error
	if err != nil {
		return nil, err
	}

	return maintenances, nil
}

func (r assetMaintenanceRepository) Update(maintenance *assets.AssetMaintenance) error {
	return r.db.Table(tableAssetMaintenanceName).Save(maintenance).Error
}

func (r assetMaintenanceRepository) Delete(maintenanceID uint) error {
	return r.db.Table(tableAssetMaintenanceName).Delete(&assets.AssetMaintenance{}, maintenanceID).Error
}
