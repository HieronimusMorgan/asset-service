package assets

import (
	request "asset-service/internal/dto/in/assets"
	"asset-service/internal/models/assets"
	repository "asset-service/internal/repository/assets"
	"asset-service/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type AssetMaintenanceService interface {
	AddAssetMaintenance(maintenance request.AssetMaintenanceRequest, clientID string) (*assets.AssetMaintenance, error)
	GetMaintenanceByID(maintenanceID uint, clientID string) (interface{}, error)
	UpdateMaintenance(maintenance *request.AssetMaintenanceRequest) error
	DeleteMaintenance(maintenanceID uint, clientID string) error
	GetMaintenancesByAssetID(assetID uint, clientID string) (interface{}, error)
	PerformMaintenanceCheck() error
}

type assetMaintenanceService struct {
	AssetMaintenanceRepository repository.AssetMaintenanceRepository
	AssetRepository            repository.AssetRepository
	Redis                      utils.RedisService
}

func NewAssetMaintenanceService(AssetMaintenance repository.AssetMaintenanceRepository, assetRepository repository.AssetRepository, RedisService utils.RedisService) AssetMaintenanceService {
	return assetMaintenanceService{AssetMaintenanceRepository: AssetMaintenance, AssetRepository: assetRepository, Redis: RedisService}
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
		return nil, fmt.Errorf("maintenance type ID is required")
	}

	if asset, err := s.AssetRepository.GetAssetByID(clientID, uint(maintenance.AssetID)); err != nil || asset.AssetID == 0 {
		return nil, err
	}

	maintenanceExist, _ := s.AssetMaintenanceRepository.GetMaintenanceByTypeExist(clientID, maintenance.AssetID, maintenance.TypeID)
	if maintenanceExist.ID != 0 {
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

	nextDueDate, err := utils.ParseOptionalDate(&maintenance.NextDueDate)
	if err != nil {
		log.Error().
			Str("key", "ParseOptionalDate").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset status by ID")
	}

	maintenanceRecord := assets.AssetMaintenance{
		AssetID:            maintenance.AssetID,
		TypeID:             maintenance.TypeID,
		UserClientID:       clientID,
		MaintenanceDetails: maintenance.MaintenanceDetails,
		MaintenanceCost:    maintenance.MaintenanceCost,
		MaintenanceDate:    maintenanceDate,
		NextDueDate:        nextDueDate,
		PerformedBy:        &data.ClientID,
		CreatedBy:          data.ClientID,
		UpdatedBy:          data.ClientID,
	}

	if maintenance.MaintenanceDetails == nil {
		maintenanceRecord.MaintenanceDetails = nil
	}

	err = s.AssetMaintenanceRepository.AddAssetMaintenance(&maintenanceRecord)
	if err != nil {
		return nil, err
	}

	return &maintenanceRecord, nil
}

func (s assetMaintenanceService) GetMaintenanceByID(maintenanceID uint, clientID string) (interface{}, error) {
	maintenance, err := s.AssetMaintenanceRepository.GetMaintenanceByID(maintenanceID, clientID)
	if err != nil {
		return nil, err
	}

	if maintenance.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return maintenance, nil
}

func (s assetMaintenanceService) UpdateMaintenance(maintenance *request.AssetMaintenanceRequest) error {
	return nil
}

func (s assetMaintenanceService) DeleteMaintenance(maintenanceID uint, clientID string) error {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return err
	}

	return s.AssetMaintenanceRepository.Delete(maintenanceID, data.ClientID)
}

func (s assetMaintenanceService) GetMaintenancesByAssetID(assetID uint, clientID string) (interface{}, error) {
	result, err := s.AssetMaintenanceRepository.GetListMaintenanceByAssetID(assetID, clientID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s assetMaintenanceService) PerformMaintenanceCheck() error {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	maintenance, err := s.AssetMaintenanceRepository.GetListMaintenance()
	if err != nil {
		return err
	}

	for _, m := range maintenance {
		if m.MaintenanceDate.IsZero() {
			continue
		}
		jsonPretty, err := json.MarshalIndent(m, "", "  ")
		if err != nil {
			return err
		}
		log.Info().
			Str("method", "Get Data").RawJSON("data", jsonPretty).Send()
	}
	return nil
}
