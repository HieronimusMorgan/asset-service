package assets

import (
	request "asset-service/internal/dto/in/assets"
	response "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	repository "asset-service/internal/repository/assets"
	"asset-service/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

type AssetMaintenanceService interface {
	AddAssetMaintenance(maintenance request.AssetMaintenanceRequest, clientID string) (*assets.AssetMaintenance, error)
	GetMaintenanceByID(maintenanceID uint, clientID string) (interface{}, error)
	UpdateMaintenance(clientID string, maintenance *request.AssetMaintenanceRequest) error
	DeleteMaintenance(maintenanceID uint, clientID string) error
	GetMaintenancesByAssetID(assetID uint, clientID string) (interface{}, error)
	PerformMaintenance(assetPerform request.AssetMaintenancePerformRequest, clientID string) (interface{}, error)
	PerformMaintenanceCheck() error
}

type assetMaintenanceService struct {
	AssetMaintenanceRepository repository.AssetMaintenanceRepository
	AssetRepository            repository.AssetRepository
	AssetMaintenanceRecord     repository.AssetMaintenanceRecordRepository
	AssetAuditLogRepository    repository.AssetAuditLogRepository
	Redis                      utils.RedisService
}

func NewAssetMaintenanceService(
	AssetMaintenance repository.AssetMaintenanceRepository,
	assetRepository repository.AssetRepository,
	AssetMaintenanceRecord repository.AssetMaintenanceRecordRepository,
	AssetAuditLogRepository repository.AssetAuditLogRepository,
	RedisService utils.RedisService) AssetMaintenanceService {
	return assetMaintenanceService{
		AssetMaintenanceRepository: AssetMaintenance,
		AssetRepository:            assetRepository,
		AssetMaintenanceRecord:     AssetMaintenanceRecord,
		AssetAuditLogRepository:    AssetAuditLogRepository,
		Redis:                      RedisService}
}

func (s assetMaintenanceService) AddAssetMaintenance(maintenance request.AssetMaintenanceRequest, clientID string) (*assets.AssetMaintenance, error) {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetUserRedis").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get user redis")
		return nil, err
	}

	if maintenance.TypeID == 0 {
		log.Error().
			Str("key", "maintenance type id").
			Str("clientID", clientID).
			Err(err).
			Msg("maintenance type ID is required")
		return nil, fmt.Errorf("maintenance type ID is required")
	}

	if asset, err := s.AssetRepository.GetAssetResponseByID(clientID, uint(maintenance.AssetID)); err != nil || asset.AssetID == 0 {
		log.Error().
			Str("key", "GetAssetResponseByID").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset by ID")
		return nil, err
	}

	maintenanceExist, _ := s.AssetMaintenanceRepository.GetMaintenanceByTypeExist(clientID, maintenance.AssetID, maintenance.TypeID)
	if maintenanceExist.ID != 0 {
		log.Error().
			Str("key", "GetMaintenanceByTypeExist").
			Str("clientID", clientID).
			Err(err).
			Msg("maintenance type already exist")
		return nil, fmt.Errorf("maintenance type already exist")
	}

	maintenanceDate, err := utils.ParseOptionalDate(&maintenance.MaintenanceDate)
	if err != nil {
		log.Error().
			Str("key", "ParseOptionalDate").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset status by ID")
	}

	var nextDueDate *time.Time
	if !maintenanceDate.IsZero() && maintenance.IntervalDays != nil && *maintenance.IntervalDays != 0 {
		nextDueDate, err = utils.CalculateNextDueDate(maintenanceDate, maintenance.IntervalDays)
	}

	maintenances := assets.AssetMaintenance{
		AssetID:            maintenance.AssetID,
		TypeID:             maintenance.TypeID,
		UserClientID:       clientID,
		MaintenanceDetails: maintenance.MaintenanceDetails,
		MaintenanceCost:    maintenance.MaintenanceCost,
		MaintenanceDate:    maintenanceDate,
		IntervalDays:       maintenance.IntervalDays,
		NextDueDate:        nextDueDate,
		PerformedBy:        &data.FullName,
		CreatedBy:          data.FullName,
		UpdatedBy:          data.FullName,
	}

	if maintenance.MaintenanceDetails == nil {
		maintenances.MaintenanceDetails = nil
	}

	if err = s.AssetMaintenanceRepository.AddAssetMaintenance(&maintenances); err != nil {
		log.Error().
			Str("key", "AddAssetMaintenance").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to add asset maintenance")
		return nil, err
	}

	if auditErr := s.AssetAuditLogRepository.AfterCreateAssetMaintenance(&maintenances); auditErr != nil {
		log.Error().
			Str("key", "AfterCreateAssetMaintenance").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to log audit after create asset maintenance")
		return nil, err
	}

	return &maintenances, nil
}

func (s assetMaintenanceService) GetMaintenanceByID(maintenanceID uint, clientID string) (interface{}, error) {
	maintenance, err := s.AssetMaintenanceRepository.GetMaintenanceResponseByID(maintenanceID, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetMaintenanceResponseByID").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get maintenance by ID")
		return nil, err
	}

	if maintenance.ID == 0 {
		log.Error().
			Str("key", "maintenance.ID").
			Str("clientID", clientID).
			Err(err).
			Msg("maintenance ID not found")
		return nil, gorm.ErrRecordNotFound
	}

	return maintenance, nil
}

func (s assetMaintenanceService) UpdateMaintenance(clientID string, maintenance *request.AssetMaintenanceRequest) error {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetUserRedis").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get user redis")
		return err
	}

	if maintenance.TypeID == 0 {
		log.Error().
			Str("key", "maintenance type id").
			Str("clientID", clientID).
			Err(err).
			Msg("maintenance type ID is required")
		return fmt.Errorf("maintenance type ID is required")
	}

	if asset, err := s.AssetRepository.GetAssetResponseByID(clientID, uint(maintenance.AssetID)); err != nil || asset.AssetID == 0 {
		log.Error().
			Str("key", "GetAssetResponseByID").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset by ID")
		return err
	}

	maintenanceExist, _ := s.AssetMaintenanceRepository.GetMaintenanceByTypeExist(clientID, maintenance.AssetID, maintenance.TypeID)
	if maintenanceExist.ID != 0 {
		log.Error().
			Str("key", "GetMaintenanceByTypeExist").
			Str("clientID", clientID).
			Err(err).
			Msg("maintenance type already exist")
		return fmt.Errorf("maintenance type already exist")
	}

	maintenanceDate, err := utils.ParseOptionalDate(&maintenance.MaintenanceDate)
	if err != nil {
		log.Error().
			Str("key", "ParseOptionalDate").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to parse maintenance date")
	}

	var nextDueDate *time.Time
	if !maintenanceDate.IsZero() && maintenance.IntervalDays != nil && *maintenance.IntervalDays != 0 {
		nextDueDate, err = utils.CalculateNextDueDate(maintenanceDate, maintenance.IntervalDays)
	}

	maintenances := assets.AssetMaintenance{
		AssetID:            maintenance.AssetID,
		TypeID:             maintenance.TypeID,
		UserClientID:       clientID,
		MaintenanceDetails: maintenance.MaintenanceDetails,
		MaintenanceCost:    maintenance.MaintenanceCost,
		MaintenanceDate:    maintenanceDate,
		IntervalDays:       maintenance.IntervalDays,
		NextDueDate:        nextDueDate,
		PerformedBy:        &data.FullName,
		UpdatedBy:          data.FullName,
	}

	if maintenance.MaintenanceDetails == nil {
		maintenances.MaintenanceDetails = nil
	}

	if err = s.AssetMaintenanceRepository.Update(&maintenances); err != nil {
		log.Error().
			Str("key", "Update").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to update maintenance")
	}

	if auditErr := s.AssetAuditLogRepository.AfterUpdateAssetMaintenance(maintenanceExist, &maintenances); auditErr != nil {
		log.Error().
			Str("key", "AfterUpdateAssetMaintenance").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to log audit after update asset maintenance")
		return auditErr
	}

	return nil
}

func (s assetMaintenanceService) PerformMaintenance(assetPerform request.AssetMaintenancePerformRequest, clientID string) (interface{}, error) {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetUserRedis").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get user redis")
		return nil, err
	}

	maintenance, err := s.AssetMaintenanceRepository.GetMaintenanceByID(uint(assetPerform.MaintenanceID), clientID)
	if err != nil {
		log.Error().
			Str("key", "GetMaintenanceResponseByID").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get maintenance by ID")
		return nil, err
	}

	if maintenance.ID == 0 {
		log.Error().
			Str("key", "maintenance.ID").
			Str("clientID", clientID).
			Err(err).
			Msg("maintenance ID not found")
		return nil, gorm.ErrRecordNotFound
	}

	if asset, err := s.AssetRepository.GetAssetByID(clientID, uint(assetPerform.AssetID)); err != nil || asset.AssetID == 0 {
		log.Error().
			Str("key", "GetAssetResponseByID").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset by ID")
		return nil, err
	}

	var maintenanceDate = time.Now()

	var nextDueDate *time.Time
	if maintenance.IntervalDays != nil && *maintenance.IntervalDays != 0 {
		nextDueDate, err = utils.CalculateNextDueDate(&maintenanceDate, maintenance.IntervalDays)
	} else {
		nextDueDate = nil
	}

	maintenances := assets.AssetMaintenance{
		ID:                 maintenance.ID,
		AssetID:            maintenance.AssetID,
		TypeID:             maintenance.TypeID,
		UserClientID:       clientID,
		MaintenanceDetails: maintenance.MaintenanceDetails,
		MaintenanceCost:    maintenance.MaintenanceCost,
		MaintenanceDate:    &maintenanceDate,
		IntervalDays:       maintenance.IntervalDays,
		NextDueDate:        nextDueDate,
		PerformedBy:        &data.FullName,
		UpdatedBy:          data.FullName,
	}

	if maintenance.MaintenanceDetails == nil {
		maintenances.MaintenanceDetails = nil
	}

	if err = s.AssetMaintenanceRepository.Update(&maintenances); err != nil {
		log.Error().
			Str("key", "Update").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to update maintenance")
	}

	var assetMaintenanceRecord = assets.AssetMaintenanceRecord{
		AssetID:            maintenance.AssetID,
		TypeID:             maintenance.TypeID,
		UserClientID:       clientID,
		MaintenanceDetails: maintenance.MaintenanceDetails,
		MaintenanceCost:    maintenance.MaintenanceCost,
		MaintenanceDate:    maintenance.MaintenanceDate,
		IntervalDays:       maintenance.IntervalDays,
		NextDueDate:        maintenance.NextDueDate,
		PerformedBy:        maintenance.PerformedBy,
		CreatedBy:          data.FullName,
		UpdatedBy:          data.FullName,
	}

	if err = s.AssetMaintenanceRecord.AddAssetMaintenanceRecord(&assetMaintenanceRecord); err != nil {
		log.Error().
			Str("key", "AddAssetMaintenanceRecord").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to add asset maintenance record")
		return nil, err
	}

	if auditErr := s.AssetAuditLogRepository.AfterUpdateAssetMaintenance(*maintenance, &maintenances); auditErr != nil {
		log.Error().
			Str("key", "AfterUpdateAssetMaintenance").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to log audit after update asset maintenance")
	}

	if auditErr := s.AssetAuditLogRepository.AfterCreateAssetMaintenanceRecord(&assetMaintenanceRecord); auditErr != nil {
		log.Error().
			Str("key", "AfterCreateAssetMaintenanceRecord").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to log audit after create asset maintenance record")
	}

	return response.AssetMaintenancesResponse{
		ID:           maintenances.ID,
		UserClientID: maintenances.UserClientID,
		AssetID:      maintenances.AssetID,
		Type: response.MaintenanceTypeResponse{
			TypeID: maintenances.TypeID,
		},
		MaintenanceDate:    (*response.DateOnly)(maintenances.MaintenanceDate),
		MaintenanceDetails: maintenances.MaintenanceDetails,
		MaintenanceCost:    maintenances.MaintenanceCost,
		PerformedBy:        maintenances.PerformedBy,
		IntervalDays:       maintenances.IntervalDays,
		NextDueDate:        (*response.DateOnly)(maintenances.NextDueDate),
	}, nil

}

func (s assetMaintenanceService) DeleteMaintenance(maintenanceID uint, clientID string) error {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetUserRedis").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get user redis")
		return err
	}

	return s.AssetMaintenanceRepository.Delete(maintenanceID, data.ClientID)
}

func (s assetMaintenanceService) GetMaintenancesByAssetID(assetID uint, clientID string) (interface{}, error) {
	result, err := s.AssetMaintenanceRepository.GetListMaintenanceByAssetID(assetID, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetListMaintenanceByAssetID").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get maintenance by asset ID")
		return nil, err
	}

	return result, nil
}

func (s assetMaintenanceService) PerformMaintenanceCheck() error {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	maintenance, err := s.AssetMaintenanceRepository.GetListMaintenance()
	if err != nil {
		log.Error().
			Str("key", "GetListMaintenance").
			Err(err).
			Msg("Failed to get maintenance list")
		return err
	}

	for _, m := range maintenance {
		if m.MaintenanceDate == nil || m.NextDueDate == nil {
			continue
		}
		jsonPretty, err := json.Marshal(m)
		if err != nil {
			log.Error().
				Str("key", "MarshalIndent").
				Err(err).
				Msg("Failed to marshal maintenance data")
			return err
		}
		log.Info().
			Str("method", "Get Data").RawJSON("data", jsonPretty).Send()
	}
	return nil
}
