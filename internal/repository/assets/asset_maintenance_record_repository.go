package assets

import (
	out "asset-service/internal/dto/out/assets"
	response "asset-service/internal/dto/out/assets"
	model "asset-service/internal/models/assets"
	"asset-service/internal/utils"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type AssetMaintenanceRecordRepository interface {
	AddAssetMaintenanceRecord(maintenance *model.AssetMaintenanceRecord) error
	GetCountTotalMaintenanceRecordByAssetID(assetID uint, clientID string) (int64, error)
	GetMaintenanceRecordByAssetID(assetID uint, clientID string) (*model.AssetMaintenanceRecord, error)
	GetMaintenanceRecordByMaintenanceID(maintenanceID uint, clientID string) (*response.AssetMaintenancesResponse, error)
	GetListMaintenanceRecordByAssetIDAndMaintenanceID(assetID, maintenanceID uint, clientID string) (*[]out.AssetMaintenanceRecordResponse, error)
	GetMaintenanceRecordByRecordIDAndAssetIDAndMaintenanceID(maintenanceRecordID, assetID, maintenanceID uint, clientID string) (interface{}, error)
	GetListMaintenance() ([]response.AssetMaintenancesResponse, error)
	GetListMaintenanceByClientID(clientID string) ([]response.AssetMaintenancesResponse, error)
	GetListMaintenanceRecordByAssetID(assetID uint, clientID string) (*[]out.AssetMaintenanceRecordResponse, error)
	Update(maintenance *model.AssetMaintenanceRecord) error
	Delete(assetID uint) error
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

func (r assetMaintenanceRecordRepository) GetCountTotalMaintenanceRecordByAssetID(assetID uint, clientID string) (int64, error) {
	var count int64
	err := r.db.Table(utils.TableAssetMaintenanceRecordName).Where("asset_id = ? AND user_client_id = ?", assetID, clientID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r assetMaintenanceRecordRepository) GetMaintenanceRecordByAssetID(assetID uint, clientID string) (*model.AssetMaintenanceRecord, error) {

	var maintenance model.AssetMaintenanceRecord
	err := r.db.Table(utils.TableAssetMaintenanceRecordName).Where("asset_id = ? AND user_client_id = ? ", assetID, clientID).Order("maintenance_record_id ASC").Scan(&maintenance).Error
	return &maintenance, err
}

func (r assetMaintenanceRecordRepository) GetMaintenanceRecordByMaintenanceID(maintenanceID uint, clientID string) (*response.AssetMaintenancesResponse, error) {
	assetMaintenance := `
		SELECT am.id, am.user_client_id, am.asset_id, amt.maintenance_type_id, amt.maintenance_type_name, am.maintenance_date, am.maintenance_details, am.maintenance_cost, am.performed_by, am.interval_days, am.next_due_date
		FROM "asset_maintenance_record" am 
		LEFT JOIN "asset_maintenance_type" amt ON am.maintenance_type_id = amt.maintenance_type_id
		WHERE am.asset_id = ? AND am.user_client_id = ?
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

func (r assetMaintenanceRecordRepository) GetListMaintenanceRecordByAssetIDAndMaintenanceID(assetID, maintenanceID uint, clientID string) (*[]out.AssetMaintenanceRecordResponse, error) {
	query := `
		SELECT 
			am.maintenance_record_id,
			amt.maintenance_type_name,
			am.maintenance_details,
			am.maintenance_date,
			am.maintenance_cost,
			am.performed_by,
			am.interval_days,
			am.next_due_date
		FROM asset_maintenance_record am
		LEFT JOIN asset_maintenance_type amt ON am.maintenance_type_id = amt.maintenance_type_id
		WHERE am.asset_id = ? AND am.user_client_id = ? AND am.maintenance_id  = ?
	`

	type maintenanceRow struct {
		MaintenanceRecordID uint
		MaintenanceTypeName string
		MaintenanceDetails  *string
		MaintenanceDate     *time.Time
		MaintenanceCost     float64
		PerformedBy         *string
		IntervalDays        *int
		NextDueDate         *time.Time
	}

	var rows []maintenanceRow
	if err := r.db.Raw(query, assetID, clientID, maintenanceID).Scan(&rows).Error; err != nil {
		return nil, err
	}

	var result []out.AssetMaintenanceRecordResponse
	for _, row := range rows {
		result = append(result, out.AssetMaintenanceRecordResponse{
			MaintenanceRecordID: row.MaintenanceRecordID,
			MaintenanceTypeName: row.MaintenanceTypeName,
			MaintenanceDetails:  row.MaintenanceDetails,
			MaintenanceDate:     row.MaintenanceDate,
			MaintenanceCost:     row.MaintenanceCost,
			PerformedBy:         row.PerformedBy,
			IntervalDays:        row.IntervalDays,
			NextDueDate:         row.NextDueDate,
		})
	}

	return &result, nil
}

func (r assetMaintenanceRecordRepository) GetListMaintenance() ([]response.AssetMaintenancesResponse, error) {
	assetMaintenance := `
		SELECT am.id, am.user_client_id, am.asset_id, amt.maintenance_type_id, amt.maintenance_type_name, am.maintenance_date, am.maintenance_details, am.maintenance_cost, am.performed_by, am.interval_days, am.next_due_date
		FROM "asset_maintenance_record" am 
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

func (r assetMaintenanceRecordRepository) GetListMaintenanceByClientID(clientID string) ([]response.AssetMaintenancesResponse, error) {
	assetMaintenance := `
		SELECT am.id, am.user_client_id, am.asset_id, amt.maintenance_type_id, amt.maintenance_type_name, am.maintenance_date, am.maintenance_details, am.maintenance_cost, am.performed_by, am.interval_days, am.next_due_date
		FROM "asset_maintenance_record" am 
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

func (r assetMaintenanceRecordRepository) GetListMaintenanceRecordByAssetID(assetID uint, clientID string) (*[]out.AssetMaintenanceRecordResponse, error) {
	query := `
		SELECT 
			am.maintenance_record_id,
			amt.maintenance_type_name,
			am.maintenance_details,
			am.maintenance_date,
			am.maintenance_cost,
			am.performed_by,
			am.interval_days,
			am.next_due_date
		FROM asset_maintenance_record am
		LEFT JOIN asset_maintenance_type amt ON am.maintenance_type_id = amt.maintenance_type_id
		WHERE am.asset_id = ? AND am.user_client_id = ?
	`

	type maintenanceRow struct {
		MaintenanceRecordID uint
		MaintenanceTypeName string
		MaintenanceDetails  *string
		MaintenanceDate     *time.Time
		MaintenanceCost     float64
		PerformedBy         *string
		IntervalDays        *int
		NextDueDate         *time.Time
	}

	var rows []maintenanceRow
	if err := r.db.Raw(query, assetID, clientID).Scan(&rows).Error; err != nil {
		return nil, err
	}

	var result []out.AssetMaintenanceRecordResponse
	for _, row := range rows {
		result = append(result, out.AssetMaintenanceRecordResponse{
			MaintenanceRecordID: row.MaintenanceRecordID,
			MaintenanceTypeName: row.MaintenanceTypeName,
			MaintenanceDetails:  row.MaintenanceDetails,
			MaintenanceDate:     row.MaintenanceDate,
			MaintenanceCost:     row.MaintenanceCost,
			PerformedBy:         row.PerformedBy,
			IntervalDays:        row.IntervalDays,
			NextDueDate:         row.NextDueDate,
		})
	}

	return &result, nil
}

func (r assetMaintenanceRecordRepository) Update(maintenance *model.AssetMaintenanceRecord) error {
	return r.db.Table(utils.TableAssetMaintenanceRecordName).Save(maintenance).Error
}

func (r assetMaintenanceRecordRepository) Delete(assetID uint) error {
	if assetID != 0 { // Ensure it exists before deleting
		if err := r.db.Unscoped().Table("asset_maintenance_record").Model(&model.AssetMaintenanceRecord{}).
			Where("asset_id = ?", assetID).
			Delete(&model.AssetMaintenanceRecord{}).Error; err != nil {
			return fmt.Errorf("failed to delete asset maintenance record: %w", err)
		}
	}
	return nil
}

func (r assetMaintenanceRecordRepository) GetMaintenanceByTypeExist(clientID string, assetID int, typeID int) (model.AssetMaintenanceRecord, error) {
	var maintenance model.AssetMaintenanceRecord
	err := r.db.Table(utils.TableAssetMaintenanceRecordName).Where("user_client_id = ? AND asset_id = ? AND maintenance_type_id = ?", clientID, assetID, typeID).First(&maintenance).Error
	return maintenance, err
}

func (r assetMaintenanceRecordRepository) GetMaintenanceRecordByRecordIDAndAssetIDAndMaintenanceID(maintenanceRecordID, assetID, maintenanceID uint, clientID string) (interface{}, error) {
	query := `
		SELECT 
			am.maintenance_record_id,
			amt.maintenance_type_name,
			am.maintenance_details,
			am.maintenance_date,
			am.maintenance_cost,
			am.performed_by,
			am.interval_days,
			am.next_due_date
		FROM asset_maintenance_record am
		LEFT JOIN asset_maintenance_type amt ON am.maintenance_type_id = amt.maintenance_type_id
		WHERE am.asset_id = ? AND am.user_client_id = ? AND am.maintenance_id  = ? AND am.maintenance_record_id = ?
	`

	type maintenanceRow struct {
		MaintenanceRecordID uint
		MaintenanceTypeName string
		MaintenanceDetails  *string
		MaintenanceDate     *time.Time
		MaintenanceCost     float64
		PerformedBy         *string
		IntervalDays        *int
		NextDueDate         *time.Time
	}

	var rows []maintenanceRow
	if err := r.db.Raw(query, assetID, clientID, maintenanceID, maintenanceRecordID).Scan(&rows).Error; err != nil {
		return nil, err
	}

	var result []out.AssetMaintenanceRecordResponse
	for _, row := range rows {
		result = append(result, out.AssetMaintenanceRecordResponse{
			MaintenanceRecordID: row.MaintenanceRecordID,
			MaintenanceTypeName: row.MaintenanceTypeName,
			MaintenanceDetails:  row.MaintenanceDetails,
			MaintenanceDate:     row.MaintenanceDate,
			MaintenanceCost:     row.MaintenanceCost,
			PerformedBy:         row.PerformedBy,
			IntervalDays:        row.IntervalDays,
			NextDueDate:         row.NextDueDate,
		})
	}

	return &result, nil
}
