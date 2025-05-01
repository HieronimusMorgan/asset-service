package assets

import (
	request "asset-service/internal/dto/in/assets"
	response "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	repository "asset-service/internal/repository/assets"
	"asset-service/internal/utils"
	"asset-service/internal/utils/redis"
	"asset-service/internal/utils/text"
	"errors"
	"github.com/rs/zerolog/log"
)

type AssetStatusService interface {
	AddAssetStatus(assetStatusRequest *request.AssetStatusRequest, clientID string, credentialKey string) (interface{}, error)
	GetAssetStatus(clientID string, size, index int) (interface{}, int64, error)
	GetAssetStatusByID(assetStatusID uint) (interface{}, error)
	UpdateAssetStatus(assetStatusID uint, assetStatusRequest *request.AssetStatusRequest, clientID string) (interface{}, error)
	DeleteAssetStatus(assetStatusID uint, clientID string) error
}

type assetStatusService struct {
	AssetStatusRepository   repository.AssetStatusRepository
	AssetAuditLogRepository repository.AssetAuditLogRepository
	Redis                   redis.RedisService
}

func NewAssetStatusService(
	assetStatusRepository repository.AssetStatusRepository,
	AssetAuditLogRepository repository.AssetAuditLogRepository,
	redis redis.RedisService) AssetStatusService {
	return assetStatusService{
		AssetStatusRepository:   assetStatusRepository,
		AssetAuditLogRepository: AssetAuditLogRepository,
		Redis:                   redis}
}

func (s assetStatusService) AddAssetStatus(assetStatusRequest *request.AssetStatusRequest, clientID string, credentialKey string) (interface{}, error) {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)

	err = text.CheckCredentialKey(s.Redis, credentialKey, data.ClientID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("credential key check failed")
		return nil, err
	}

	var assetStatus = &assets.AssetStatus{
		StatusName:  assetStatusRequest.StatusName,
		Description: assetStatusRequest.Description,
		CreatedBy:   &data.ClientID,
	}

	if err = s.AssetStatusRepository.GetAssetStatusByName(assetStatusRequest.StatusName); err == nil {
		log.Error().
			Str("key", "GetAssetStatusByName").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset status by name")
		return nil, errors.New("assets status already exists")
	}

	if err = s.AssetStatusRepository.AddAssetStatus(&assetStatus); err != nil {
		log.Error().
			Str("key", "AddAssetStatus").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to add asset status")
		return nil, err
	}

	if err = s.AssetAuditLogRepository.AfterCreateAssetStatus(assetStatus); err != nil {
		log.Error().
			Str("key", "AfterCreateAssetStatus").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to log audit after create asset status")
		return nil, err
	}

	return response.AssetStatusResponse{
		AssetStatusID: assetStatus.AssetStatusID,
		StatusName:    assetStatus.StatusName,
		Description:   assetStatus.Description,
	}, nil
}

func (s assetStatusService) GetAssetStatus(clientID string, size, index int) (interface{}, int64, error) {
	_, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetRedisData").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to retrieve data from Redis")
		return nil, 0, err
	}
	total, err := s.AssetStatusRepository.GetCountAssetStatus()
	if err != nil {
		log.Error().
			Str("key", "GetCountAssetStatus").
			Err(err).
			Msg("Failed to get count of asset status")
		return nil, 0, err
	}

	assetStatus, err := s.AssetStatusRepository.GetAssetStatus(size, index)
	if err != nil {
		log.Error().
			Str("key", "GetAssetStatus").
			Err(err).
			Msg("Failed to get asset status")
		return nil, total, err
	}

	return assetStatus, total, nil
}

func (s assetStatusService) GetAssetStatusByID(assetStatusID uint) (interface{}, error) {
	assetStatus, err := s.AssetStatusRepository.GetAssetStatusByID(assetStatusID)
	if err != nil {
		log.Error().
			Str("key", "GetAssetStatusByID").
			Err(err).
			Msg("Failed to get asset status by MaintenanceTypeID")
		return nil, err
	}
	return response.AssetStatusResponse{
		AssetStatusID: assetStatus.AssetStatusID,
		StatusName:    assetStatus.StatusName,
		Description:   assetStatus.Description,
	}, nil
}

func (s assetStatusService) UpdateAssetStatus(assetStatusID uint, assetStatusRequest *request.AssetStatusRequest, clientID string) (interface{}, error) {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)

	assetStatus, err := s.AssetStatusRepository.GetAssetStatusByID(assetStatusID)
	if err != nil {
		log.Error().
			Str("key", "GetAssetStatusByID").
			Err(err).
			Msg("Failed to get asset status by MaintenanceTypeID")
		return nil, err
	}
	var oldAsset = assetStatus
	assetStatus.StatusName = assetStatusRequest.StatusName
	assetStatus.Description = assetStatusRequest.Description
	assetStatus.UpdatedBy = &data.ClientID

	if err = s.AssetStatusRepository.UpdateAssetStatus(assetStatus); err != nil {
		log.Error().
			Str("key", "UpdateAssetStatus").
			Err(err).
			Msg("Failed to update asset status")
		return nil, err
	}

	err = s.AssetAuditLogRepository.AfterUpdateAssetStatus(*oldAsset, assetStatus)
	if err != nil {
		log.Error().
			Str("key", "AfterUpdateAssetStatus").
			Err(err).
			Msg("Failed to log audit after update asset status")
		return nil, err
	}

	return response.AssetStatusResponse{
		AssetStatusID: assetStatus.AssetStatusID,
		StatusName:    assetStatus.StatusName,
		Description:   assetStatus.Description,
	}, nil
}

func (s assetStatusService) DeleteAssetStatus(assetStatusID uint, clientID string) error {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	assetStatus, err := s.AssetStatusRepository.GetAssetStatusByID(assetStatusID)
	if err != nil {
		log.Error().
			Str("key", "GetAssetStatusByID").
			Err(err).
			Msg("Failed to get asset status by MaintenanceTypeID")
		return err
	}

	assetStatus.DeletedBy = &data.ClientID

	if err = s.AssetStatusRepository.DeleteAssetStatus(assetStatus); err != nil {
		log.Error().
			Str("key", "DeleteAssetStatus").
			Err(err).
			Msg("Failed to delete asset status")
		return err
	}

	if err = s.AssetAuditLogRepository.AfterDeleteAssetStatus(assetStatus); err != nil {
		log.Error().
			Str("key", "AfterDeleteAssetStatus").
			Err(err).
			Msg("Failed to log audit after delete asset status")
		return err
	}
	return nil
}
