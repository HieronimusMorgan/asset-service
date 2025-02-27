package assets

import (
	request "asset-service/internal/dto/in/assets"
	response "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	repo "asset-service/internal/repository/assets"
	"asset-service/internal/repository/transaction"
	"asset-service/internal/utils"
	"errors"
	"github.com/rs/zerolog/log"
)

type AssetService interface {
	AddAsset(assetRequest *request.AssetRequest, images []response.AssetImageResponse, clientID string) (interface{}, error)
	UpdateAsset(assetID uint, assetRequest request.UpdateAssetRequest, clientID string) (interface{}, error)
	GetListAsset(clientID string) (interface{}, error)
	GetAssetByID(clientID string, assetID uint) (interface{}, error)
	UpdateAssetStatus(assetID uint, statusID uint, clientID string) error
	UpdateAssetCategory(assetID uint, categoryID uint, clientID string) error
	DeleteAsset(assetID uint, clientID string) error
}

type assetService struct {
	AssetRepository            repo.AssetRepository
	AssetCategoryRepository    repo.AssetCategoryRepository
	AssetStatusRepository      repo.AssetStatusRepository
	AssetMaintenanceRepository repo.AssetMaintenanceRepository
	AssetImageRepository       repo.AssetImageRepository
	Redis                      utils.RedisService
	AuditLogRepository         repo.AssetAuditLogRepository
	AssetTransaction           transaction.AssetTransactionRepository
}

func NewAssetService(assetRepository repo.AssetRepository,
	AssetCategoryRepository repo.AssetCategoryRepository,
	AssetStatusRepository repo.AssetStatusRepository,
	assetMaintenanceRepository repo.AssetMaintenanceRepository,
	AssetImageRepository repo.AssetImageRepository,
	log repo.AssetAuditLogRepository,
	redis utils.RedisService,
	AssetTransaction transaction.AssetTransactionRepository) AssetService {
	return assetService{
		AssetRepository:            assetRepository,
		AssetCategoryRepository:    AssetCategoryRepository,
		AssetStatusRepository:      AssetStatusRepository,
		AssetMaintenanceRepository: assetMaintenanceRepository,
		AssetImageRepository:       AssetImageRepository,
		AuditLogRepository:         log,
		Redis:                      redis,
		AssetTransaction:           AssetTransaction}
}

func (s assetService) AddAsset(assetRequest *request.AssetRequest, images []response.AssetImageResponse, clientID string) (interface{}, error) {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)

	if err != nil {
		log.Error().
			Str("key", "GetRedisData").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get data from redis")
		return nil, err
	}

	if check, err := s.AssetRepository.AssetNameExists(assetRequest.Name, clientID); err != nil {
		log.Error().
			Str("key", "AssetNameExists").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to check asset name exists")
		return nil, err
	} else if check {
		log.Error().
			Str("key", "AssetNameExists").
			Str("clientID", clientID).
			Err(err).
			Msg("Asset name already exists")
		return nil, errors.New("assets already exists")
	}

	if _, err := s.AssetCategoryRepository.GetAssetCategoryById(uint(assetRequest.CategoryID), clientID); err != nil {
		log.Error().
			Str("key", "GetAssetCategoryById").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset category by ID")
		return nil, errors.New("category not found")
	}

	if _, err := s.AssetStatusRepository.GetAssetStatusByID(uint(assetRequest.StatusID)); err != nil {
		log.Error().
			Str("key", "GetAssetStatusByID").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset status by ID")
		return nil, errors.New("status not found")
	}

	purchaseDate, err := utils.ParseOptionalDate(assetRequest.PurchaseDate)
	if err != nil {
		log.Error().
			Str("key", "ParseOptionalDate").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset status by ID")
	}

	expiryDate, err := utils.ParseOptionalDate(assetRequest.ExpiryDate)
	if err != nil {
		log.Error().
			Str("key", "ParseOptionalDate").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset status by ID")
	}

	warrantyExpiry, err := utils.ParseOptionalDate(assetRequest.WarrantyExpiry)
	if err != nil {
		log.Error().
			Str("key", "ParseOptionalDate").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset status by ID")
	}

	assetRequest.ConvertAssetRequestEmptyToNil()

	var asset = &assets.Asset{
		UserClientID:       clientID,
		SerialNumber:       assetRequest.SerialNumber,
		Name:               assetRequest.Name,
		Description:        assetRequest.Description,
		Barcode:            assetRequest.Barcode,
		CategoryID:         assetRequest.CategoryID,
		StatusID:           assetRequest.StatusID,
		PurchaseDate:       purchaseDate,
		ExpiryDate:         expiryDate,
		WarrantyExpiryDate: warrantyExpiry,
		Price:              assetRequest.Price,
		Stock:              assetRequest.Stock,
		Notes:              assetRequest.Notes,
		CreatedBy:          data.ClientID,
		UpdatedBy:          data.ClientID,
	}

	err = s.AssetRepository.AddAsset(asset)
	if err != nil {
		log.Error().
			Str("key", "AddAsset").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to add asset")
		return nil, err
	}

	if len(images) != 0 {
		for _, image := range images {
			var assetImage = &assets.AssetImage{
				UserClientID: clientID,
				AssetID:      asset.AssetID,
				ImageURL:     image.ImageURL,
				FileSize:     image.FileSize,
				FileType:     image.FileType,
				CreatedBy:    data.ClientID,
				UpdatedBy:    data.ClientID,
			}
			err = s.AssetImageRepository.AddAssetImage(assetImage)
		}

		if err != nil {
			log.Error().
				Str("key", "AddAssetImage").
				Str("clientID", clientID).
				Err(err).
				Msg("Failed to add asset image")
			return nil, err
		}
	}

	err = s.AuditLogRepository.AfterCreateAsset(asset)
	if err != nil {
		log.Error().
			Str("key", "AfterCreateAsset").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to create asset")
		return nil, err
	}

	result, err := s.AssetRepository.GetAssetResponseByID(clientID, asset.AssetID)
	if err != nil {
		log.Error().
			Str("key", "GetAssetResponseByID").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset by ID")
		return nil, err
	}

	log.Info().
		Str("key", "GetAssetResponseByID").
		Str("clientID", clientID).
		Fields(asset).
		Msg("Success to get asset by ID")
	return result, nil
}

func (s assetService) UpdateAsset(assetID uint, assetRequest request.UpdateAssetRequest, clientID string) (interface{}, error) {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetUserRedis").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get user redis")
		return nil, err
	}

	oldAsset, err := s.AssetRepository.GetAsset(assetID, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetAsset").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset by ID")
		return nil, err
	}

	purchaseDate, err := utils.ParseOptionalDate(assetRequest.PurchaseDate)
	if err != nil {
		log.Error().
			Str("key", "ParseOptionalDate").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset status by ID")
	}

	expiryDate, err := utils.ParseOptionalDate(assetRequest.ExpiryDate)
	if err != nil {
		log.Error().
			Str("key", "ParseOptionalDate ExpiryDate").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset status by ID")
	}

	warrantyExpiry, err := utils.ParseOptionalDate(assetRequest.WarrantyExpiryDate)
	if err != nil {
		log.Error().
			Str("key", "ParseOptionalDate WarrantyExpiryDate").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset status by ID")
	}

	var asset = &assets.Asset{
		AssetID:            assetID,
		UserClientID:       clientID,
		SerialNumber:       assetRequest.SerialNumber,
		Description:        assetRequest.Description,
		Barcode:            assetRequest.Barcode,
		CategoryID:         assetRequest.CategoryID,
		StatusID:           assetRequest.StatusID,
		Stock:              assetRequest.Stock,
		PurchaseDate:       purchaseDate,
		ExpiryDate:         expiryDate,
		WarrantyExpiryDate: warrantyExpiry,
		Price:              assetRequest.Price,
		Notes:              assetRequest.Notes,
		UpdatedBy:          data.ClientID,
	}

	err = s.AssetRepository.UpdateAsset(asset, clientID)
	if err != nil {
		log.Error().
			Str("key", "UpdateAsset").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to update asset")
		return nil, err
	}

	err = s.AuditLogRepository.AfterUpdateAsset(*oldAsset, asset)
	if err != nil {
		log.Error().
			Str("key", "AfterUpdateAsset").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to update asset")
		return nil, err
	}

	return asset, nil
}

func (s assetService) GetListAsset(clientID string) (interface{}, error) {
	result, err := s.AssetRepository.GetListAssets(clientID)
	if err != nil {
		log.Error().
			Str("key", "GetListAssets").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get list assets")
		return nil, err
	}

	return result, nil
}

func (s assetService) GetAssetByID(clientID string, assetID uint) (interface{}, error) {
	_, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetUserRedis").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get user redis")
		return nil, err
	}

	asset, err := s.AssetRepository.GetAssetResponseByID(clientID, assetID)
	if err != nil {
		log.Error().
			Str("key", "GetAssetResponseByID").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset by ID")
		return nil, err
	}

	return *asset, nil
}

func (s assetService) UpdateAssetStatus(assetID uint, statusID uint, clientID string) error {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetUserRedis").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get user redis")
		return err
	}

	oldAsset, err := s.AssetRepository.GetAsset(assetID, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetAsset").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset by ID")
		return err
	}

	asset, err := s.AssetRepository.UpdateAssetStatus(assetID, statusID, data.ClientID)
	if err != nil {
		log.Error().
			Str("key", "UpdateAssetStatus").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to update asset status")
		return err
	}

	err = s.AuditLogRepository.AfterUpdateAsset(*oldAsset, asset)
	if err != nil {
		log.Error().
			Str("key", "AfterUpdateAsset").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to update asset")
		return err
	}

	return nil
}

func (s assetService) UpdateAssetCategory(assetID uint, categoryID uint, clientID string) error {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetUserRedis").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get user redis")
		return err
	}

	oldAsset, err := s.AssetRepository.GetAsset(assetID, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetAsset").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset by ID")
		return err
	}

	asset, err := s.AssetRepository.UpdateAssetCategory(assetID, categoryID, data.ClientID)
	if err != nil {
		log.Error().
			Str("key", "UpdateAssetCategory").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to update asset category")
		return err
	}

	err = s.AuditLogRepository.AfterUpdateAsset(*oldAsset, asset)
	if err != nil {
		log.Error().
			Str("key", "AfterUpdateAsset").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to update asset")
		return err
	}
	return nil
}

func (s assetService) DeleteAsset(assetID uint, clientID string) error {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetUserRedis").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get user redis")
		return err
	}

	err = s.AssetTransaction.DeleteAsset(assetID, clientID, data.ClientID)
	if err != nil {
		log.Error().
			Str("key", "DeleteAsset").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to delete asset")
		return err
	}

	return nil
}
