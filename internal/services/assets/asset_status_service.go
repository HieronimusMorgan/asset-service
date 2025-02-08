package assets

import (
	assets3 "asset-service/internal/dto/in/assets"
	assets4 "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	"asset-service/internal/models/user"
	assets2 "asset-service/internal/repository/assets"
	"asset-service/internal/utils"
	"errors"
	"gorm.io/gorm"
)

type AssetStatusService struct {
	AssetStatusRepository *assets2.AssetStatusRepository
}

func (s AssetStatusService) AddAssetStatus(assetStatusRequest *assets3.AssetStatusRequest, clientID string) (interface{}, error) {
	var user = &user.User{}

	err := utils.GetDataFromRedis(utils.User, clientID, user)
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
	var assetStatusResponse = assets4.AssetStatusResponse{
		AssetStatusID: assetStatus.AssetStatusID,
		StatusName:    assetStatus.StatusName,
		Description:   assetStatus.Description,
	}

	return assetStatusResponse, nil
}

func (s AssetStatusService) GetAssetStatus() (interface{}, error) {
	assetStatus, err := s.AssetStatusRepository.GetAssetStatus()
	if err != nil {
		return nil, err
	}
	var assetStatusResponse []assets4.AssetStatusResponse
	for _, status := range assetStatus {
		assetStatusResponse = append(assetStatusResponse, assets4.AssetStatusResponse{
			AssetStatusID: status.AssetStatusID,
			StatusName:    status.StatusName,
			Description:   status.Description,
		})
	}
	return assetStatusResponse, nil
}

func (s AssetStatusService) GetAssetStatusByID(assetStatusID uint) (interface{}, error) {
	assetStatus, err := s.AssetStatusRepository.GetAssetStatusByID(assetStatusID)
	if err != nil {
		return nil, err
	}
	var assetStatusResponse = assets4.AssetStatusResponse{
		AssetStatusID: assetStatus.AssetStatusID,
		StatusName:    assetStatus.StatusName,
		Description:   assetStatus.Description,
	}
	return assetStatusResponse, nil
}

func (s AssetStatusService) UpdateAssetStatus(assetStatusID uint, assetStatusRequest *assets3.AssetStatusRequest, clientID string) (interface{}, error) {
	var user = &user.User{}

	err := utils.GetDataFromRedis(utils.User, clientID, user)
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

	var assetStatusResponse = assets4.AssetStatusResponse{
		AssetStatusID: assetStatus.AssetStatusID,
		StatusName:    assetStatus.StatusName,
		Description:   assetStatus.Description,
	}

	return assetStatusResponse, nil
}

func (s AssetStatusService) DeleteAssetStatus(assetStatusID uint, clientID string) error {
	var user = &user.User{}

	err := utils.GetDataFromRedis(utils.User, clientID, user)
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

func AddAssetStatus(db *gorm.DB) *AssetStatusService {
	assetStatusRepo := assets2.NewAssetStatusRepository(db)
	return &AssetStatusService{AssetStatusRepository: assetStatusRepo}
}
