package assets

import (
	request "asset-service/internal/dto/in/assets"
	model "asset-service/internal/models/assets"
	repo "asset-service/internal/repository/assets"
	"asset-service/internal/utils"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type AssetMaintenanceTypeService interface {
	GetMaintenanceTypeByID(maintenanceTypeID uint, clientID string) (interface{}, error)
	GetListMaintenanceType(clientID string) ([]model.AssetMaintenanceType, error)
	AddMaintenanceType(maintenanceType *request.AssetMaintenanceTypeRequest, clientID string) (interface{}, error)
	UpdateMaintenanceType(id uint, clientID string, maintenanceType *request.AssetMaintenanceTypeRequest) (interface{}, error)
	DeleteMaintenanceType(maintenanceTypeID uint, clientID string) error
}
type assetMaintenanceTypeService struct {
	AssetMaintenanceTypeRepository repo.AssetMaintenanceTypeRepository
	AssetMaintenanceRepository     repo.AssetMaintenanceRepository
	Redis                          utils.RedisService
}

func NewAssetMaintenanceTypeService(
	assetMaintenanceTypeRepository repo.AssetMaintenanceTypeRepository,
	assetMaintenanceRepository repo.AssetMaintenanceRepository,
	redis utils.RedisService) AssetMaintenanceTypeService {
	return assetMaintenanceTypeService{
		AssetMaintenanceTypeRepository: assetMaintenanceTypeRepository,
		AssetMaintenanceRepository:     assetMaintenanceRepository,
		Redis:                          redis}
}

func (s assetMaintenanceTypeService) AddMaintenanceType(maintenanceType *request.AssetMaintenanceTypeRequest, clientID string) (interface{}, error) {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetUserRedis").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get user redis")
		return nil, err
	}

	maintenanceTypeRecord := &model.AssetMaintenanceType{
		UserClientID:        clientID,
		MaintenanceTypeName: maintenanceType.MaintenanceTypeName,
		Description:         maintenanceType.Description,
		CreatedBy:           data.ClientID,
	}

	if err = s.AssetMaintenanceTypeRepository.AddAssetMaintenanceType(maintenanceTypeRecord, clientID); err != nil {
		log.Error().
			Str("key", "AddAssetMaintenanceType").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to add asset maintenance type")
		return nil, err
	}

	return maintenanceType, nil
}

func (s assetMaintenanceTypeService) GetMaintenanceTypeByID(maintenanceTypeID uint, clientID string) (interface{}, error) {
	maintenanceType, err := s.AssetMaintenanceTypeRepository.GetAssetMaintenanceTypeByID(maintenanceTypeID, clientID)
	if err != nil {
		return nil, err
	}

	if maintenanceType.MaintenanceTypeID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return maintenanceType, nil
}

func (s assetMaintenanceTypeService) GetListMaintenanceType(clientID string) ([]model.AssetMaintenanceType, error) {
	maintenanceTypes, err := s.AssetMaintenanceTypeRepository.GetAssetMaintenanceType(clientID)
	if err != nil {
		return nil, err
	}

	return maintenanceTypes, nil
}

func (s assetMaintenanceTypeService) UpdateMaintenanceType(id uint, clientID string, maintenanceType *request.AssetMaintenanceTypeRequest) (interface{}, error) {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetUserRedis").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get user redis")
		return nil, err
	}

	exist, err := s.AssetMaintenanceTypeRepository.GetAssetMaintenanceTypeByName(maintenanceType.MaintenanceTypeName, clientID)
	if err != nil {
		return nil, err
	}

	if exist.MaintenanceTypeID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	maintenanceTypeRecord := &model.AssetMaintenanceType{
		MaintenanceTypeID:   id,
		MaintenanceTypeName: maintenanceType.MaintenanceTypeName,
		Description:         maintenanceType.Description,
		UpdatedBy:           data.ClientID,
	}

	err = s.AssetMaintenanceTypeRepository.UpdateAssetMaintenanceType(maintenanceTypeRecord, clientID)
	if err != nil {
		return nil, err
	}

	return maintenanceTypeRecord, nil
}

func (s assetMaintenanceTypeService) DeleteMaintenanceType(maintenanceTypeID uint, clientID string) error {
	if err := s.AssetMaintenanceTypeRepository.DeleteAssetMaintenanceTypeByID(maintenanceTypeID, clientID); err != nil {
		return err
	}

	return nil
}
