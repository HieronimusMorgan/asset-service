package assets

import (
	response "asset-service/internal/dto/out/assets"
	model "asset-service/internal/models/assets"
	"asset-service/internal/utils"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type AssetMaintenanceRecordRepository interface {
	AddAssetMaintenanceRecord(maintenance *model.AssetMaintenanceRecord) error
	GetMaintenanceByAssetID(assetID uint, clientID string) (*model.AssetMaintenanceRecord, error)
	GetMaintenanceRecordByMaintenanceID(maintenanceID uint, clientID string) (*response.AssetMaintenancesResponse, error)
	GetListMaintenanceByAssetID(assetID uint, clientID string) ([]response.AssetMaintenancesResponse, error)
	GetListMaintenance() ([]response.AssetMaintenancesResponse, error)
	GetListMaintenanceByClientID(clientID string) ([]response.AssetMaintenancesResponse, error)
	Update(maintenance *model.AssetMaintenanceRecord) error
	Delete(assetID uint, fullName string) error
	GetMaintenanceByTypeExist(clientID string, assetID int, typeID int) (model.AssetMaintenanceRecord, error)
}

type assetMaintenanceRecordRepository struct {
	db gorm.DB
}

func NewAssetMaintenanceRecordRepository(db gorm.DB) AssetMaintenanceRecordRepository {
	return assetMaintenanceRecordRepository{db: db}
}

func (r assetMaintenanceRecordRepository) AddAssetMaintenanceRecord(maintenance *model.AssetMaintenanceRecord) error {
	return r.db.Table(utils.TableAssetMaintenanceRecordName).Create(maintenance).Error
}

func (r assetMaintenanceRecordRepository) GetMaintenanceByAssetID(assetID uint, clientID string) (*model.AssetMaintenanceRecord, error) {

	var maintenance model.AssetMaintenanceRecord
	err := r.db.Table(utils.TableAssetMaintenanceRecordName).Where("asset_id = ? AND user_client_id = ? ", assetID, clientID).Order("maintenance_record_id ASC").Scan(&maintenance).Error
	return &maintenance, err
}

func (r assetMaintenanceRecordRepository) GetMaintenanceRecordByMaintenanceID(maintenanceID uint, clientID string) (*response.AssetMaintenancesResponse, error) {
	assetMaintenance := `
		SELECT am.id, am.user_client_id, am.asset_id, amt.type_id, amt.type_name, am.maintenance_date, am.maintenance_details, am.maintenance_cost, am.performed_by, am.interval_days, am.next_due_date
		FROM "asset-service"."asset_maintenance_record" am 
		LEFT JOIN "asset-service"."asset_maintenance_type" amt ON am.type_id = amt.type_id
		WHERE am.asset_id = ? AND am.user_client_id = ?
`

	rows := r.db.Raw(assetMaintenance, maintenanceID, clientID).Row()

	var maintenance response.AssetMaintenancesResponse
	var typeMaintenance response.MaintenanceTypeResponse

	err := rows.Scan(
		&maintenance.ID,
		&maintenance.UserClientID,
		&maintenance.AssetID,
		&typeMaintenance.TypeID,
		&typeMaintenance.TypeName,
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

func (r assetMaintenanceRecordRepository) GetListMaintenanceByAssetID(assetID uint, clientID string) ([]response.AssetMaintenancesResponse, error) {
	assetMaintenance := `
		SELECT am.id, am.user_client_id, am.asset_id, amt.type_id, amt.type_name, am.maintenance_date, am.maintenance_details, am.maintenance_cost, am.performed_by, am.interval_days, am.next_due_date
		FROM "asset-service"."asset_maintenance_record" am 
		LEFT JOIN "asset-service"."asset_maintenance_type" amt ON am.type_id = amt.type_id
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
			&typeMaintenance.TypeID,
			&typeMaintenance.TypeName,
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

func (r assetMaintenanceRecordRepository) GetListMaintenance() ([]response.AssetMaintenancesResponse, error) {
	assetMaintenance := `
		SELECT am.id, am.user_client_id, am.asset_id, amt.type_id, amt.type_name, am.maintenance_date, am.maintenance_details, am.maintenance_cost, am.performed_by, am.interval_days, am.next_due_date
		FROM "asset-service"."asset_maintenance_record" am 
		LEFT JOIN "asset-service"."asset_maintenance_type" amt ON am.type_id = amt.type_id
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
			&typeMaintenance.TypeID,
			&typeMaintenance.TypeName,
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

func (r assetMaintenanceRecordRepository) GetListMaintenanceByClientID(clientID string) ([]response.AssetMaintenancesResponse, error) {
	assetMaintenance := `
		SELECT am.id, am.user_client_id, am.asset_id, amt.type_id, amt.type_name, am.maintenance_date, am.maintenance_details, am.maintenance_cost, am.performed_by, am.interval_days, am.next_due_date
		FROM "asset-service"."asset_maintenance_record" am 
		LEFT JOIN "asset-service"."asset_maintenance_type" amt ON am.type_id = amt.type_id
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
			&typeMaintenance.TypeID,
			&typeMaintenance.TypeName,
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

func (r assetMaintenanceRecordRepository) Update(maintenance *model.AssetMaintenanceRecord) error {
	return r.db.Table(utils.TableAssetMaintenanceRecordName).Save(maintenance).Error
}

func (r assetMaintenanceRecordRepository) Delete(assetID uint, fullName string) error {
	if assetID != 0 { // Ensure it exists before deleting
		if err := r.db.Table("asset-service.asset_maintenance_record").Model(&model.AssetMaintenanceRecord{}).
			Where("asset_id = ?", assetID).
			Updates(map[string]interface{}{"deleted_by": fullName, "deleted_at": time.Now()}).
			Delete(&model.AssetMaintenanceRecord{}).Error; err != nil {
			return fmt.Errorf("failed to delete asset maintenance record: %w", err)
		}
	}
	return nil
}

func (r assetMaintenanceRecordRepository) GetMaintenanceByTypeExist(clientID string, assetID int, typeID int) (model.AssetMaintenanceRecord, error) {
	var maintenance model.AssetMaintenanceRecord
	err := r.db.Table(utils.TableAssetMaintenanceRecordName).Where("user_client_id = ? AND asset_id = ? AND type_id = ?", clientID, assetID, typeID).First(&maintenance).Error
	return maintenance, err
}
