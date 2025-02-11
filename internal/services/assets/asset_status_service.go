package assets

import (
	request "asset-service/internal/dto/in/assets"
	response "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	"asset-service/internal/models/user"
	repository "asset-service/internal/repository/assets"
	"asset-service/internal/utils"
	"errors"
)

type AssetStatusService interface {
	AddAssetStatus(assetStatusRequest *request.AssetStatusRequest, clientID string) (interface{}, error)
	GetAssetStatus() (interface{}, error)
	GetAssetStatusByID(assetStatusID uint) (interface{}, error)
	UpdateAssetStatus(assetStatusID uint, assetStatusRequest *request.AssetStatusRequest, clientID string) (interface{}, error)
	DeleteAssetStatus(assetStatusID uint, clientID string) error
}

type assetStatusService struct {
	AssetStatusRepository repository.AssetStatusRepository
	Redis                 utils.RedisService
}

func NewAssetStatusService(assetStatusRepository repository.AssetStatusRepository, redis utils.RedisService) AssetStatusService {
	return assetStatusService{AssetStatusRepository: assetStatusRepository, Redis: redis}
}

func (s assetStatusService) AddAssetStatus(assetStatusRequest *request.AssetStatusRequest, clientID string) (interface{}, error) {
	var user = &user.User{}

	err := s.Redis.GetData(utils.User, clientID, user)
	if err != nil {
		return nil, err
	}

	var assetStatus = &assets.AssetStatus{
		StatusName:  assetStatusRequest.StatusName,
		Description: assetStatusRequest.Description,
		CreatedBy:   user.FullName,
	}

	err = s.AssetStatusRepository.GetAssetStatusByName(assetStatusRequest.StatusName)
	if err == nil {
		return nil, errors.New("assets status already exists")
	}

	err = s.AssetStatusRepository.AddAssetStatus(&assetStatus)
	if err != nil {
		return nil, err
	}
	var assetStatusResponse = response.AssetStatusResponse{
		AssetStatusID: assetStatus.AssetStatusID,
		StatusName:    assetStatus.StatusName,
		Description:   assetStatus.Description,
	}

	return assetStatusResponse, nil
}

func (s assetStatusService) GetAssetStatus() (interface{}, error) {
	assetStatus, err := s.AssetStatusRepository.GetAssetStatus()
	if err != nil {
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
		return nil, err
	}
	var assetStatusResponse = response.AssetStatusResponse{
		AssetStatusID: assetStatus.AssetStatusID,
		StatusName:    assetStatus.StatusName,
		Description:   assetStatus.Description,
	}
	return assetStatusResponse, nil
}

func (s assetStatusService) UpdateAssetStatus(assetStatusID uint, assetStatusRequest *request.AssetStatusRequest, clientID string) (interface{}, error) {
	var user = &user.User{}

	err := s.Redis.GetData(utils.User, clientID, user)
	if err != nil {
		return nil, err
	}

	assetStatus, err := s.AssetStatusRepository.GetAssetStatusByID(assetStatusID)
	if err != nil {
		return nil, err
	}

	assetStatus.StatusName = assetStatusRequest.StatusName
	assetStatus.Description = assetStatusRequest.Description
	assetStatus.UpdatedBy = user.FullName

	err = s.AssetStatusRepository.UpdateAssetStatus(assetStatus)
	if err != nil {
		return nil, err
	}

	var assetStatusResponse = response.AssetStatusResponse{
		AssetStatusID: assetStatus.AssetStatusID,
		StatusName:    assetStatus.StatusName,
		Description:   assetStatus.Description,
	}

	return assetStatusResponse, nil
}

func (s assetStatusService) DeleteAssetStatus(assetStatusID uint, clientID string) error {
	var user = &user.User{}

	err := s.Redis.GetData(utils.User, clientID, user)
	if err != nil {
		return err
	}

	assetStatus, err := s.AssetStatusRepository.GetAssetStatusByID(assetStatusID)
	if err != nil {
		return err
	}

	assetStatus.DeletedBy = &user.FullName

	err = s.AssetStatusRepository.DeleteAssetStatus(assetStatus)
	if err != nil {
		return err
	}

	return nil
}
