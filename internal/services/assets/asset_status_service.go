package assets

import (
	request "asset-service/internal/dto/in/assets"
	response "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	repository "asset-service/internal/repository/assets"
	"asset-service/internal/utils"
	"errors"
	"github.com/rs/zerolog/log"
)

type AssetStatusService interface {
	AddAssetStatus(assetStatusRequest *request.AssetStatusRequest, clientID string) (interface{}, error)
	GetAssetStatus() (interface{}, error)
	GetAssetStatusByID(assetStatusID uint) (interface{}, error)
	UpdateAssetStatus(assetStatusID uint, assetStatusRequest *request.AssetStatusRequest, clientID string) (interface{}, error)
	DeleteAssetStatus(assetStatusID uint, clientID string) error
}

type assetStatusService struct {
	AssetStatusRepository   repository.AssetStatusRepository
	AssetAuditLogRepository repository.AssetAuditLogRepository
	Redis                   utils.RedisService
}

func NewAssetStatusService(
	assetStatusRepository repository.AssetStatusRepository,
	AssetAuditLogRepository repository.AssetAuditLogRepository,
	redis utils.RedisService) AssetStatusService {
	return assetStatusService{
		AssetStatusRepository:   assetStatusRepository,
		AssetAuditLogRepository: AssetAuditLogRepository,
		Redis:                   redis}
}

func (s assetStatusService) AddAssetStatus(assetStatusRequest *request.AssetStatusRequest, clientID string) (interface{}, error) {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)

	var assetStatus = &assets.AssetStatus{
		StatusName:  assetStatusRequest.StatusName,
		Description: assetStatusRequest.Description,
		CreatedBy:   data.ClientID,
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

func (s assetStatusService) GetAssetStatus() (interface{}, error) {
	assetStatus, err := s.AssetStatusRepository.GetAssetStatus()
	if err != nil {
		log.Error().
			Str("key", "GetAssetStatus").
			Err(err).
			Msg("Failed to get asset status")
		return nil, err
	}
	var assetStatusResponse []response.AssetStatusResponse
	for _, status := range assetStatus {
		assetStatusResponse = append(assetStatusResponse, response.AssetStatusResponse{
			AssetStatusID: status.AssetStatusID,
			StatusName:    status.StatusName,
			Description:   status.Description,
		})
	}
	return assetStatusResponse, nil
}

func (s assetStatusService) GetAssetStatusByID(assetStatusID uint) (interface{}, error) {
	assetStatus, err := s.AssetStatusRepository.GetAssetStatusByID(assetStatusID)
	if err != nil {
		log.Error().
			Str("key", "GetAssetStatusByID").
			Err(err).
			Msg("Failed to get asset status by ID")
		return nil, err
	}
	return response.AssetStatusResponse{
		AssetStatusID: assetStatus.AssetStatusID,
		StatusName:    assetStatus.StatusName,
		Description:   assetStatus.Description,
	}, nil
}

func (s assetStatusService) UpdateAssetStatus(assetStatusID uint, assetStatusRequest *request.AssetStatusRequest, clientID string) (interface{}, error) {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)

	assetStatus, err := s.AssetStatusRepository.GetAssetStatusByID(assetStatusID)
	if err != nil {
		log.Error().
			Str("key", "GetAssetStatusByID").
			Err(err).
			Msg("Failed to get asset status by ID")
		return nil, err
	}
	var oldAsset = assetStatus
	assetStatus.StatusName = assetStatusRequest.StatusName
	assetStatus.Description = assetStatusRequest.Description
	assetStatus.UpdatedBy = data.ClientID

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
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	assetStatus, err := s.AssetStatusRepository.GetAssetStatusByID(assetStatusID)
	if err != nil {
		log.Error().
			Str("key", "GetAssetStatusByID").
			Err(err).
			Msg("Failed to get asset status by ID")
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
