package assets

import (
	response "asset-service/internal/dto/out/assets"
	model "asset-service/internal/models/assets"
	"asset-service/internal/utils"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type AssetMaintenanceRepository interface {
	AddAssetMaintenance(maintenance *model.AssetMaintenance) error
	GetMaintenanceByAssetID(assetID uint, clientID string) (*model.AssetMaintenance, error)
	GetMaintenanceByMaintenanceIDAndAssetID(maintenanceID uint, assetID uint, clientID string) (*model.AssetMaintenance, error)
	GetMaintenanceResponseByID(maintenanceID uint, clientID string) (*response.AssetMaintenancesResponse, error)
	GetListMaintenanceByAssetID(assetID uint, clientID string) ([]response.AssetMaintenancesResponse, error)
	GetListMaintenance() ([]response.AssetMaintenancesResponse, error)
	GetListMaintenanceByClientID(clientID string) ([]response.AssetMaintenancesResponse, error)
	Update(maintenance *model.AssetMaintenance) error
	Delete(assetID uint, fullName string) error
	GetMaintenanceByTypeExist(clientID string, assetID int, typeID int) (model.AssetMaintenance, error)
}

type assetMaintenanceRepository struct {
	db              gorm.DB
	assetRepository AssetRepository
}

func NewAssetMaintenanceRepository(db gorm.DB) AssetMaintenanceRepository {
	return assetMaintenanceRepository{db: db}
}

func (r assetMaintenanceRepository) AddAssetMaintenance(maintenance *model.AssetMaintenance) error {
	return r.db.Table(utils.TableAssetMaintenanceName).Create(maintenance).Error
}

func (r assetMaintenanceRepository) GetMaintenanceByAssetID(assetID uint, clientID string) (*model.AssetMaintenance, error) {
	var maintenance model.AssetMaintenance
	err := r.db.Table(utils.TableAssetMaintenanceName).Where("asset_id = ? AND user_client_id = ?", assetID, clientID).First(&maintenance).Error
	return &maintenance, err
}

func (r assetMaintenanceRepository) GetMaintenanceByMaintenanceIDAndAssetID(maintenanceID uint, assetID uint, clientID string) (*model.AssetMaintenance, error) {
	var maintenance model.AssetMaintenance
	err := r.db.Table(utils.TableAssetMaintenanceName).Where("maintenance_type_id = ? AND asset_id = ? AND user_client_id = ?", maintenanceID, assetID, clientID).First(&maintenance).Error
	return &maintenance, err
}

func (r assetMaintenanceRepository) GetMaintenanceResponseByID(maintenanceID uint, clientID string) (*response.AssetMaintenancesResponse, error) {

	assetMaintenance := `
		SELECT am.id, am.user_client_id, am.asset_id, amt.maintenance_type_id, amt.maintenance_type_name, am.maintenance_date, am.maintenance_details, am.maintenance_cost, am.performed_by, am.interval_days, am.next_due_date
		FROM "asset_maintenance" am 
		LEFT JOIN "asset_maintenance_type" amt ON am.maintenance_type_id = amt.maintenance_type_id
		WHERE am.id = ? AND am.user_client_id = ?
`

	rows := r.db.Raw(assetMaintenance, maintenanceID, clientID).Row()

	var maintenance response.AssetMaintenancesResponse
	var typeMaintenance response.MaintenanceTypeResponse

	err := rows.Scan(
		&maintenance.ID,
		&maintenance.UserClientID,
		&maintenance.AssetID,
		&typeMaintenance.MaintenanceTypeID,
		&typeMaintenance.MaintenanceTypeName,
		&maintenance.MaintenanceDate,
		&maintenance.MaintenanceDetails,
		&maintenance.MaintenanceCost,
		&maintenance.PerformedBy,
		&maintenance.IntervalDays,
		&maintenance.NextDueDate,
	)

	if err != nil {
		return nil, err
	}

	maintenance.Type = typeMaintenance

	return &maintenance, nil
}

func (r assetMaintenanceRepository) GetListMaintenanceByAssetID(assetID uint, clientID string) ([]response.AssetMaintenancesResponse, error) {
	assetMaintenance := `
		SELECT am.id, am.user_client_id, am.asset_id, amt.maintenance_type_id, amt.maintenance_type_name, am.maintenance_date, am.maintenance_details, am.maintenance_cost, am.performed_by, am.interval_days, am.next_due_date
		FROM "asset_maintenance" am 
		LEFT JOIN "asset_maintenance_type" amt ON am.maintenance_type_id = amt.maintenance_type_id
		WHERE am.asset_id = ? AND am.user_client_id = ?
`

	rows, err := r.db.Raw(assetMaintenance, assetID, clientID).Rows()
	if err != nil {
		return nil, err
	}

	var result []response.AssetMaintenancesResponse
	for rows.Next() {
		var maintenance response.AssetMaintenancesResponse
		var typeMaintenance response.MaintenanceTypeResponse

		err := rows.Scan(
			&maintenance.ID,
			&maintenance.UserClientID,
			&maintenance.AssetID,
			&typeMaintenance.MaintenanceTypeID,
			&typeMaintenance.MaintenanceTypeName,
			&maintenance.MaintenanceDate,
			&maintenance.MaintenanceDetails,
			&maintenance.MaintenanceCost,
			&maintenance.PerformedBy,
			&maintenance.IntervalDays,
			&maintenance.NextDueDate,
		)

		if err != nil {
			return nil, err
		}
		maintenance.Type = typeMaintenance
		result = append(result, maintenance)
	}

	return result, nil
}

func (r assetMaintenanceRepository) GetListMaintenance() ([]response.AssetMaintenancesResponse, error) {
	assetMaintenance := `
		SELECT am.id, am.user_client_id, am.asset_id, amt.maintenance_type_id, amt.maintenance_type_name, am.maintenance_date, am.maintenance_details, am.maintenance_cost, am.performed_by, am.interval_days, am.next_due_date
		FROM "asset_maintenance" am 
		LEFT JOIN "asset_maintenance_type" amt ON am.maintenance_type_id = amt.maintenance_type_id
`

	rows, err := r.db.Raw(assetMaintenance).Rows()
	if err != nil {
		return nil, err
	}

	var result []response.AssetMaintenancesResponse
	for rows.Next() {
		var maintenance response.AssetMaintenancesResponse
		var typeMaintenance response.MaintenanceTypeResponse

		err := rows.Scan(
			&maintenance.ID,
			&maintenance.UserClientID,
			&maintenance.AssetID,
			&typeMaintenance.MaintenanceTypeID,
			&typeMaintenance.MaintenanceTypeName,
			&maintenance.MaintenanceDate,
			&maintenance.MaintenanceDetails,
			&maintenance.MaintenanceCost,
			&maintenance.PerformedBy,
			&maintenance.IntervalDays,
			&maintenance.NextDueDate,
		)

		if err != nil {
			return nil, err
		}
		maintenance.Type = typeMaintenance
		result = append(result, maintenance)
	}

	return result, nil
}

func (r assetMaintenanceRepository) GetListMaintenanceByClientID(clientID string) ([]response.AssetMaintenancesResponse, error) {
	assetMaintenance := `
		SELECT am.id, am.user_client_id, am.asset_id, amt.maintenance_type_id, amt.maintenance_type_name, am.maintenance_date, am.maintenance_details, am.maintenance_cost, am.performed_by, am.interval_days, am.next_due_date
		FROM "asset_maintenance" am 
		LEFT JOIN "asset_maintenance_type" amt ON am.maintenance_type_id = amt.maintenance_type_id
		WHERE am.user_client_id = ?
`

	rows, err := r.db.Raw(assetMaintenance, clientID).Rows()
	if err != nil {
		return nil, err
	}

	var result []response.AssetMaintenancesResponse
	for rows.Next() {
		var maintenance response.AssetMaintenancesResponse
		var typeMaintenance response.MaintenanceTypeResponse

		err := rows.Scan(
			&maintenance.ID,
			&maintenance.UserClientID,
			&maintenance.AssetID,
			&typeMaintenance.MaintenanceTypeID,
			&typeMaintenance.MaintenanceTypeName,
			&maintenance.MaintenanceDate,
			&maintenance.MaintenanceDetails,
			&maintenance.MaintenanceCost,
			&maintenance.PerformedBy,
			&maintenance.IntervalDays,
			&maintenance.NextDueDate,
		)

		if err != nil {
			return nil, err
		}
		maintenance.Type = typeMaintenance
		result = append(result, maintenance)
	}

	return result, nil
}

func (r assetMaintenanceRepository) Update(maintenance *model.AssetMaintenance) error {
	return r.db.Table(utils.TableAssetMaintenanceName).Save(maintenance).Error
}

func (r assetMaintenanceRepository) Delete(assetID uint, fullName string) error {
	if assetID != 0 { // Ensure it exists before deleting
		if err := r.db.Table("asset_maintenance").Model(model.AssetMaintenance{}).
			Where("asset_id = ?", assetID).
			Updates(map[string]interface{}{"deleted_by": fullName, "deleted_at": time.Now()}).
			Delete(&model.AssetMaintenance{}).Error; err != nil {
			return fmt.Errorf("failed to delete asset maintenance: %w", err)
		}
	}
	return nil
}

func (r assetMaintenanceRepository) GetMaintenanceByTypeExist(clientID string, assetID int, typeID int) (model.AssetMaintenance, error) {
	var maintenance model.AssetMaintenance
	err := r.db.Table(utils.TableAssetMaintenanceName).Where("user_client_id = ? AND asset_id = ? AND maintenance_type_id = ?", clientID, assetID, typeID).First(&maintenance).Error
	return maintenance, err
}
