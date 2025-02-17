package assets

import (
	response "asset-service/internal/dto/out/assets"
	model "asset-service/internal/models/assets"
	"asset-service/internal/utils"
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type AssetMaintenanceRecordRepository interface {
	AddAssetMaintenanceRecord(maintenance *model.AssetMaintenance) error
	GetMaintenanceByAssetID(assetID uint, clientID string) (*model.AssetMaintenance, error)
	GetMaintenanceByID(maintenanceID uint, clientID string) (*response.AssetMaintenanceResponse, error)
	GetListMaintenanceByAssetID(assetID uint, clientID string) ([]response.AssetMaintenancesResponse, error)
	GetListMaintenance() ([]response.AssetMaintenancesResponse, error)
	GetListMaintenanceByClientID(clientID string) ([]response.AssetMaintenancesResponse, error)
	Update(maintenance *model.AssetMaintenance) error
	Delete(assetID uint, fullName string) error
	GetMaintenanceByTypeExist(clientID string, assetID int, typeID int) (model.AssetMaintenance, error)
}

type assetMaintenanceRecordRepository struct {
	db gorm.DB
}

func NewAssetMaintenanceRecordRepository(db gorm.DB) AssetMaintenanceRecordRepository {
	return assetMaintenanceRecordRepository{db: db}
}

func (r assetMaintenanceRecordRepository) AddAssetMaintenanceRecord(maintenance *model.AssetMaintenance) error {
	return r.db.Table(utils.TableAssetMaintenanceRecordName).Create(maintenance).Error
}

func (r assetMaintenanceRecordRepository) GetMaintenanceByAssetID(assetID uint, clientID string) (*model.AssetMaintenance, error) {
	var maintenance model.AssetMaintenance
	err := r.db.Table(utils.TableAssetMaintenanceRecordName).Where("asset_id = ?", assetID).Order("maintenance_record_id ASC").Scan(&maintenance).Error
	return &maintenance, err
}

func (r assetMaintenanceRecordRepository) GetMaintenanceByID(maintenanceID uint, clientID string) (*response.AssetMaintenanceResponse, error) {
	var maintenance response.AssetMaintenanceResponse

	assetMaintenance := `
		SELECT am.id, am.asset_id, am.maintenance_details, am.maintenance_date, am.maintenance_cost 
		FROM "my-home"."asset_maintenance" am 
		WHERE am.id = ?
`

	err := r.db.Raw(assetMaintenance, maintenanceID).Scan(&maintenance).Error
	if err != nil {
		return nil, err
	}

	return &maintenance, nil
}

func (r assetMaintenanceRecordRepository) GetListMaintenanceByAssetID(assetID uint, clientID string) ([]response.AssetMaintenancesResponse, error) {
	assetMaintenance := `
		SELECT am.id, amt.type_name, am.maintenance_date, am.maintenance_details, am. maintenance_cost, am. performed_by, am.next_due_date 
		FROM "my-home"."asset_maintenance" am
		LEFT JOIN "my-home"."asset_maintenance_type" amt ON am.type_id = amt.type_id 
		WHERE am.asset_id = ? AND am.user_client_id = ?
`
	rows, err := r.db.Raw(assetMaintenance, assetID, clientID).Rows()
	if err != nil {
		return nil, err
	}

	var result []response.AssetMaintenancesResponse
	for rows.Next() {
		var asset response.AssetMaintenancesResponse

		var maintenanceDate sql.NullTime
		var nextDueDate sql.NullTime

		err := rows.Scan(
			&asset.ID,
			&asset.TypeName,
			&maintenanceDate,
			&asset.MaintenanceDetails,
			&asset.MaintenanceCost,
			&asset.PerformedBy,
			&nextDueDate,
		)

		if err != nil {
			return nil, err
		}
		if maintenanceDate.Valid {
			asset.MaintenanceDate = maintenanceDate.Time
		}
		if nextDueDate.Valid {
			asset.NextDueDate = &nextDueDate.Time
		}
		result = append(result, asset)
	}

	return result, nil
}

func (r assetMaintenanceRecordRepository) GetListMaintenance() ([]response.AssetMaintenancesResponse, error) {
	assetMaintenance := `
		SELECT am.id, amt.type_name, am.maintenance_date, am.maintenance_details, am. maintenance_cost, am. performed_by, am.next_due_date 
		FROM "my-home"."asset_maintenance" am
		LEFT JOIN "my-home"."asset_maintenance_type" amt ON am.type_id = amt.type_id
`
	rows, err := r.db.Raw(assetMaintenance).Rows()
	if err != nil {
		return nil, err
	}

	var result []response.AssetMaintenancesResponse
	for rows.Next() {
		var asset response.AssetMaintenancesResponse

		var maintenanceDate sql.NullTime
		var nextDueDate sql.NullTime

		err := rows.Scan(
			&asset.ID,
			&asset.TypeName,
			&maintenanceDate,
			&asset.MaintenanceDetails,
			&asset.MaintenanceCost,
			&asset.PerformedBy,
			&nextDueDate,
		)

		if err != nil {
			return nil, err
		}
		if maintenanceDate.Valid {
			asset.MaintenanceDate = maintenanceDate.Time
		}
		if nextDueDate.Valid {
			asset.NextDueDate = &nextDueDate.Time
		}
		result = append(result, asset)
	}

	return result, nil
}

func (r assetMaintenanceRecordRepository) GetListMaintenanceByClientID(clientID string) ([]response.AssetMaintenancesResponse, error) {
	assetMaintenance := `
		SELECT am.id, amt.type_name, am.maintenance_date, am.maintenance_details, am. maintenance_cost, am. performed_by, am.next_due_date 
		FROM "my-home"."asset_maintenance" am
		LEFT JOIN "my-home"."asset_maintenance_type" amt ON am.type_id = amt.type_id 
		WHERE am.user_client_id = ?
`
	rows, err := r.db.Raw(assetMaintenance, clientID).Rows()
	if err != nil {
		return nil, err
	}

	var result []response.AssetMaintenancesResponse
	for rows.Next() {
		var asset response.AssetMaintenancesResponse

		var maintenanceDate sql.NullTime
		var nextDueDate sql.NullTime

		err := rows.Scan(
			&asset.ID,
			&asset.TypeName,
			&maintenanceDate,
			&asset.MaintenanceDetails,
			&asset.MaintenanceCost,
			&asset.PerformedBy,
			&nextDueDate,
		)

		if err != nil {
			return nil, err
		}

		if maintenanceDate.Valid {
			asset.MaintenanceDate = maintenanceDate.Time
		}
		if nextDueDate.Valid {
			asset.NextDueDate = &nextDueDate.Time
		}
		result = append(result, asset)
	}

	return result, nil
}

func (r assetMaintenanceRecordRepository) Update(maintenance *model.AssetMaintenance) error {
	return r.db.Table(utils.TableAssetMaintenanceRecordName).Save(maintenance).Error
}

func (r assetMaintenanceRecordRepository) Delete(assetID uint, fullName string) error {
	if assetID != 0 { // Ensure it exists before deleting
		if err := r.db.Table("my-home.asset_maintenance_record").Model(&model.AssetMaintenance{}).
			Where("asset_id = ?", assetID).
			Updates(map[string]interface{}{"deleted_by": fullName, "deleted_at": time.Now()}).
			Delete(&model.AssetMaintenanceRecord{}).Error; err != nil {
			return fmt.Errorf("failed to delete asset maintenance record: %w", err)
		}
	}
	return nil
}

func (r assetMaintenanceRecordRepository) GetMaintenanceByTypeExist(clientID string, assetID int, typeID int) (model.AssetMaintenance, error) {
	var maintenance model.AssetMaintenance
	err := r.db.Table(utils.TableAssetMaintenanceRecordName).Where("user_client_id = ? AND asset_id = ? AND type_id = ?", clientID, assetID, typeID).First(&maintenance).Error
	return maintenance, err
}
