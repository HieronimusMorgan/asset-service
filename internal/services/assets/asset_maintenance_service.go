package assets

import (
	request "asset-service/internal/dto/in/assets"
	response "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	repository "asset-service/internal/repository/assets"
	"asset-service/internal/utils"
	"asset-service/internal/utils/redis"
	"asset-service/internal/utils/text"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

type AssetMaintenanceService interface {
	AddAssetMaintenance(maintenance request.AssetMaintenanceRequest, clientID string, credentialKey string) (*assets.AssetMaintenance, error)
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
	Redis                      redis.RedisService
}

func NewAssetMaintenanceService(
	AssetMaintenance repository.AssetMaintenanceRepository,
	assetRepository repository.AssetRepository,
	AssetMaintenanceRecord repository.AssetMaintenanceRecordRepository,
	AssetAuditLogRepository repository.AssetAuditLogRepository,
	RedisService redis.RedisService) AssetMaintenanceService {
	return assetMaintenanceService{
		AssetMaintenanceRepository: AssetMaintenance,
		AssetRepository:            assetRepository,
		AssetMaintenanceRecord:     AssetMaintenanceRecord,
		AssetAuditLogRepository:    AssetAuditLogRepository,
		Redis:                      RedisService}
}

func (s assetMaintenanceService) AddAssetMaintenance(maintenance request.AssetMaintenanceRequest, clientID string, credentialKey string) (*assets.AssetMaintenance, error) {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetUserRedis").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get user redis")
		return nil, err
	}

	err = text.CheckCredentialKey(s.Redis, credentialKey, data.ClientID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("credential key check failed")
		return nil, err
	}

	if maintenance.MaintenanceTypeID == 0 {
		log.Error().
			Str("key", "maintenance type id").
			Str("clientID", clientID).
			Err(err).
			Msg("maintenance type MaintenanceTypeID is required")
		return nil, fmt.Errorf("maintenance type MaintenanceTypeID is required")
	}

	if asset, err := s.AssetRepository.GetAssetResponseByID(clientID, uint(maintenance.AssetID)); err != nil || asset.AssetID == 0 {
		log.Error().
			Str("key", "GetAssetResponseByID").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset by MaintenanceTypeID")
		return nil, err
	}

	maintenanceExist, _ := s.AssetMaintenanceRepository.GetMaintenanceByTypeExist(clientID, maintenance.AssetID, maintenance.MaintenanceTypeID)
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
			Msg("Failed to get asset status by MaintenanceTypeID")
	}

	var nextDueDate *time.Time
	if !maintenanceDate.IsZero() && maintenance.IntervalDays != nil && *maintenance.IntervalDays != 0 {
		nextDueDate, err = utils.CalculateNextDueDate(maintenanceDate, maintenance.IntervalDays)
	}

	maintenances := assets.AssetMaintenance{
		AssetID:            maintenance.AssetID,
		MaintenanceTypeID:  maintenance.MaintenanceTypeID,
		UserClientID:       clientID,
		MaintenanceDetails: maintenance.MaintenanceDetails,
		MaintenanceCost:    maintenance.MaintenanceCost,
		MaintenanceDate:    maintenanceDate,
		IntervalDays:       maintenance.IntervalDays,
		NextDueDate:        nextDueDate,
		PerformedBy:        &data.ClientID,
		CreatedBy:          &data.ClientID,
		UpdatedBy:          &data.ClientID,
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
			Msg("Failed to get maintenance by MaintenanceTypeID")
		return nil, err
	}

	if maintenance.ID == 0 {
		log.Error().
			Str("key", "maintenance.MaintenanceTypeID").
			Str("clientID", clientID).
			Err(err).
			Msg("maintenance MaintenanceTypeID not found")
		return nil, gorm.ErrRecordNotFound
	}

	return maintenance, nil
}

func (s assetMaintenanceService) UpdateMaintenance(clientID string, maintenance *request.AssetMaintenanceRequest) error {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetUserRedis").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get user redis")
		return err
	}

	if maintenance.MaintenanceTypeID == 0 {
		log.Error().
			Str("key", "maintenance type id").
			Str("clientID", clientID).
			Err(err).
			Msg("maintenance type MaintenanceTypeID is required")
		return fmt.Errorf("maintenance type MaintenanceTypeID is required")
	}

	if asset, err := s.AssetRepository.GetAssetResponseByID(clientID, uint(maintenance.AssetID)); err != nil || asset.AssetID == 0 {
		log.Error().
			Str("key", "GetAssetResponseByID").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset by MaintenanceTypeID")
		return err
	}

	maintenanceExist, _ := s.AssetMaintenanceRepository.GetMaintenanceByTypeExist(clientID, maintenance.AssetID, maintenance.MaintenanceTypeID)
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
		MaintenanceTypeID:  maintenance.MaintenanceTypeID,
		UserClientID:       clientID,
		MaintenanceDetails: maintenance.MaintenanceDetails,
		MaintenanceCost:    maintenance.MaintenanceCost,
		MaintenanceDate:    maintenanceDate,
		IntervalDays:       maintenance.IntervalDays,
		NextDueDate:        nextDueDate,
		PerformedBy:        &data.ClientID,
		UpdatedBy:          &data.ClientID,
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
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetUserRedis").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get user redis")
		return nil, err
	}

	maintenance, err := s.AssetMaintenanceRepository.GetMaintenanceByMaintenanceIDAndAssetID(uint(assetPerform.MaintenanceID), uint(assetPerform.AssetID), clientID)
	if err != nil {
		log.Error().
			Str("key", "GetMaintenanceResponseByID").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get maintenance by MaintenanceTypeID")
		return nil, errors.New("maintenance asset not found")
	}

	if maintenance.ID == 0 {
		log.Error().
			Str("key", "maintenance.MaintenanceTypeID").
			Str("clientID", clientID).
			Err(err).
			Msg("maintenance MaintenanceTypeID not found")
		return nil, errors.New("maintenance asset not found")
	}

	if asset, err := s.AssetRepository.GetAssetByID(clientID, uint(assetPerform.AssetID)); err != nil || asset.AssetID == 0 {
		log.Error().
			Str("key", "GetAssetResponseByID").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset by MaintenanceTypeID")
		return nil, errors.New("asset not found")
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
		MaintenanceTypeID:  maintenance.MaintenanceTypeID,
		UserClientID:       clientID,
		MaintenanceDetails: maintenance.MaintenanceDetails,
		MaintenanceCost:    maintenance.MaintenanceCost,
		MaintenanceDate:    &maintenanceDate,
		IntervalDays:       maintenance.IntervalDays,
		NextDueDate:        nextDueDate,
		PerformedBy:        &data.ClientID,
		UpdatedBy:          &data.ClientID,
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
		MaintenanceID:      int(maintenance.ID),
		MaintenanceTypeID:  maintenance.MaintenanceTypeID,
		UserClientID:       clientID,
		MaintenanceDetails: maintenance.MaintenanceDetails,
		MaintenanceCost:    maintenance.MaintenanceCost,
		MaintenanceDate:    maintenance.MaintenanceDate,
		IntervalDays:       maintenance.IntervalDays,
		NextDueDate:        maintenance.NextDueDate,
		PerformedBy:        maintenance.PerformedBy,
		CreatedBy:          &data.ClientID,
		UpdatedBy:          &data.ClientID,
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
			MaintenanceTypeID: maintenances.MaintenanceTypeID,
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
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
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
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetUserRedis").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get user redis")
		return nil, err
	}

	asset, err := s.AssetRepository.GetAssetResponseByID(data.ClientID, assetID)
	if err != nil || asset.AssetID == 0 {
		log.Error().
			Str("key", "GetAssetByID").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset by MaintenanceTypeID")
		return nil, errors.New("asset not found")
	}

	// Check if the asset belongs to the client
	if asset.UserClientID != data.ClientID {
		log.Error().
			Str("key", "GetAssetByID").
			Str("clientID", data.ClientID).
			Err(err).
			Msg("Asset does not belong to the client")
		return nil, errors.New("asset does not belong to the client")
	}
	// Check if the asset belongs to the client

	result, err := s.AssetMaintenanceRepository.GetListMaintenanceByAssetID(assetID, data.ClientID)
	if err != nil {
		log.Error().
			Str("key", "GetListMaintenanceByAssetID").
			Str("clientID", data.ClientID).
			Err(err).
			Msg("Failed to get maintenance by asset MaintenanceTypeID")
		return nil, err
	}

	return struct {
		Asset       response.AssetResponse               `json:"asset"`
		Maintenance []response.AssetMaintenancesResponse `json:"maintenances"`
	}{
		Asset:       *asset,
		Maintenance: result,
	}, nil
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
