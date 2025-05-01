package assets

import (
	request "asset-service/internal/dto/in/assets"
	response "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	repo "asset-service/internal/repository/assets"
	"asset-service/internal/repository/transaction"
	repouser "asset-service/internal/repository/users"
	"asset-service/internal/utils"
	"asset-service/internal/utils/redis"
	"asset-service/internal/utils/text"
	"errors"
	"github.com/rs/zerolog/log"
)

type AssetService interface {
	AddAsset(assetRequest *request.AssetRequest, images []response.AssetImageResponse, clientID, requestHeaderID string) (interface{}, error)
	UpdateAsset(assetID uint, assetRequest request.UpdateAssetRequest, clientID string, credentialKey string) (interface{}, error)
	UpdateStockAsset(isAdded bool, assetID uint, stock struct {
		Stock  int     `json:"stock" binding:"required"`
		Reason *string `json:"reason"`
	}, clientID string) (interface{}, error)
	UpdateImageAsset(assetID uint, clientID string, metadata []response.AssetImageResponse) error
	GetListAsset(clientID string, index int, size int) (interface{}, int64, error)
	GetAssetByID(clientID string, assetID uint) (interface{}, error)
	UpdateAssetStatus(assetID uint, statusID uint, clientID string) error
	UpdateAssetCategory(assetID uint, categoryID uint, clientID string) error
	DeleteAsset(assetID uint, clientID string) error
}

type assetService struct {
	UserRepository             repouser.UserRepository
	AssetRepository            repo.AssetRepository
	AssetCategoryRepository    repo.AssetCategoryRepository
	AssetStatusRepository      repo.AssetStatusRepository
	AssetImageRepository       repo.AssetImageRepository
	Redis                      redis.RedisService
	AuditLogRepository         repo.AssetAuditLogRepository
	AssetGroupMemberRepository repo.AssetGroupMemberRepository
	AssetGroupAssetRepository  repo.AssetGroupAssetRepository
	AssetTransaction           transaction.AssetTransactionRepository
	AssetStockRepository       repo.AssetStockRepository
}

func NewAssetService(userRepository repouser.UserRepository,
	assetRepository repo.AssetRepository,
	assetCategoryRepository repo.AssetCategoryRepository,
	assetStatusRepository repo.AssetStatusRepository,
	assetImageRepository repo.AssetImageRepository,
	log repo.AssetAuditLogRepository,
	assetGroupMemberRepository repo.AssetGroupMemberRepository,
	assetGroupAssetRepository repo.AssetGroupAssetRepository,
	redis redis.RedisService,
	assetTransaction transaction.AssetTransactionRepository,
	assetStockRepository repo.AssetStockRepository) AssetService {
	return assetService{
		UserRepository:             userRepository,
		AssetRepository:            assetRepository,
		AssetCategoryRepository:    assetCategoryRepository,
		AssetStatusRepository:      assetStatusRepository,
		AssetImageRepository:       assetImageRepository,
		AuditLogRepository:         log,
		AssetGroupMemberRepository: assetGroupMemberRepository,
		AssetGroupAssetRepository:  assetGroupAssetRepository,
		Redis:                      redis,
		AssetTransaction:           assetTransaction,
		AssetStockRepository:       assetStockRepository}
}

func (s assetService) AddAsset(assetRequest *request.AssetRequest, images []response.AssetImageResponse, clientID, credentialKey string) (interface{}, error) {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logError("GetRedisData", clientID, err, "Failed to get data from redis")
	}

	user, err := s.UserRepository.GetUserByClientID(data.ClientID)
	if err != nil {
		return logError("GetUserByClientID", clientID, err, "Failed to get user by client MaintenanceTypeID")
	}

	err = text.CheckCredentialKey(s.Redis, credentialKey, data.ClientID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("credential key check failed")
		return nil, err
	}

	if exists, err := s.AssetRepository.AssetNameExists(assetRequest.Name, clientID); err != nil {
		return logError("AssetNameExists", clientID, err, "Failed to check asset name exists")
	} else if exists {
		return logError("AssetNameExists", clientID, errors.New("assets already exists"), "Asset name already exists")
	}

	if _, err := s.AssetCategoryRepository.GetAssetCategoryById(uint(assetRequest.CategoryID), clientID); err != nil {
		return logError("GetAssetCategoryById", clientID, errors.New("category not found"), "Failed to get asset category by MaintenanceTypeID")
	}

	if _, err := s.AssetStatusRepository.GetAssetStatusByID(uint(assetRequest.StatusID)); err != nil {
		return logError("GetAssetStatusByID", clientID, errors.New("status not found"), "Failed to get asset status by MaintenanceTypeID")
	}

	purchaseDate, _ := utils.ParseOptionalDate(assetRequest.PurchaseDate)
	expiryDate, _ := utils.ParseOptionalDate(assetRequest.ExpiryDate)
	warrantyExpiry, _ := utils.ParseOptionalDate(assetRequest.WarrantyExpiry)

	assetRequest.ConvertAssetRequestEmptyToNil()

	asset := &assets.Asset{
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
		CreatedBy:          &data.ClientID,
		UpdatedBy:          &data.ClientID,
	}

	if err := s.AssetRepository.AddAsset(asset, images); err != nil {
		return logError("AddAsset", clientID, err, "Failed to add asset")
	}

	// check user have asset group
	assetGroupMember, _ := s.AssetGroupMemberRepository.GetAssetGroupMemberByUserID(user.UserID)
	if assetGroupMember != nil && assetGroupMember.UserID != 0 {
		assetGroupAsset := &assets.AssetGroupAsset{
			AssetID:      asset.AssetID,
			AssetGroupID: assetGroupMember.AssetGroupID,
			UserID:       user.UserID,
			CreatedBy:    &data.ClientID,
		}

		if err := s.AssetGroupAssetRepository.AddAssetGroupAsset(assetGroupAsset); err != nil {
			return logError("AddAssetGroupAsset", clientID, err, "Failed to add asset group asset")
		}
	}

	assetImage, err := s.AssetImageRepository.GetAssetImageResponseByAssetID(asset.AssetID)
	if err != nil {
		return logError("GetAssetImageResponseByAssetID", clientID, err, "Failed to get asset image by MaintenanceTypeID")
	}

	result, err := s.AssetRepository.GetAssetResponseByID(clientID, asset.AssetID)
	if err != nil {
		return logError("GetAssetResponseByID", clientID, err, "Failed to get asset by MaintenanceTypeID")
	}

	result.Images = *assetImage

	log.Info().Str("key", "GetAssetResponseByID").Str("clientID", clientID).Fields(asset).Msg("Success to get asset by MaintenanceTypeID")
	return result, nil
}

func (s assetService) UpdateAsset(assetID uint, assetRequest request.UpdateAssetRequest, clientID string, credentialKey string) (interface{}, error) {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logError("GetUserRedis", clientID, err, "Failed to get user redis")
	}

	err = text.CheckCredentialKey(s.Redis, credentialKey, data.ClientID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("credential key check failed")
		return nil, err
	}

	oldAsset, err := s.AssetRepository.GetAsset(assetID, clientID)
	if err != nil {
		return logError("GetAsset", clientID, err, "Failed to get asset by MaintenanceTypeID")
	}

	purchaseDate, _ := utils.ParseOptionalDate(assetRequest.PurchaseDate)
	expiryDate, _ := utils.ParseOptionalDate(assetRequest.ExpiryDate)
	warrantyExpiry, _ := utils.ParseOptionalDate(assetRequest.WarrantyExpiryDate)

	asset := &assets.Asset{
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
		UpdatedBy:          &data.ClientID,
	}

	if err := s.AssetRepository.UpdateAsset(asset, clientID); err != nil {
		return logError("UpdateAsset", clientID, err, "Failed to update asset")
	}

	if err := s.AuditLogRepository.AfterUpdateAsset(*oldAsset, asset); err != nil {
		return logError("AfterUpdateAsset", clientID, err, "Failed to update asset log")
	}

	return asset, nil
}

func (s assetService) UpdateStockAsset(isAdded bool, assetID uint, stock struct {
	Stock  int     `json:"stock" binding:"required"`
	Reason *string `json:"reason"`
}, clientID string) (interface{}, error) {
	// Step 1: Fetch user data from Redis
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logError("GetUserRedis", clientID, err, "Failed to get user from Redis")
	}

	// Step 2: Retrieve asset and stock data
	asset, err := s.AssetRepository.GetAsset(assetID, clientID)
	if err != nil {
		return logError("GetAsset", clientID, err, "Failed to get asset by MaintenanceTypeID")
	}

	oldAssetStock, err := s.AssetStockRepository.GetAssetStockByAssetID(assetID, clientID)
	if err != nil {
		return logError("GetAssetStockByAssetID", clientID, err, "Failed to get asset stock by asset MaintenanceTypeID")
	}

	log.Info().
		Uint("assetID", assetID).
		Int("Previous Stock", oldAssetStock.LatestQuantity).
		Msg("Retrieved current stock")

	// Step 3: Determine new stock quantity
	var stockType string
	var latestQuantity int

	if isAdded {
		stockType = "INCREASE"
		latestQuantity = oldAssetStock.LatestQuantity + stock.Stock
	} else {
		stockType = "DECREASE"
		if oldAssetStock.LatestQuantity < stock.Stock {
			return logError("UpdateStockAsset", clientID, errors.New("insufficient stock"), "Stock cannot be negative")
		}
		latestQuantity = oldAssetStock.LatestQuantity - stock.Stock
	}

	// Step 4: Create stock update struct
	newAssetStock := &assets.AssetStock{
		AssetID:         asset.AssetID,
		UserClientID:    data.ClientID,
		InitialQuantity: oldAssetStock.InitialQuantity,
		LatestQuantity:  latestQuantity,
		Quantity:        stock.Stock,
		ChangeType:      stockType,
		Reason:          stock.Reason,
		UpdatedBy:       &data.ClientID,
	}

	// Step 5: Update stock in a transaction
	err = s.AssetStockRepository.UpdateAssetStock(newAssetStock, clientID)
	if err != nil {
		return logError("UpdateAssetStock", clientID, err, "Failed to update asset stock")
	}

	log.Info().
		Uint("assetID", assetID).
		Int("Updated Stock", latestQuantity).
		Str("Change Type", stockType).
		Msg("Stock updated successfully")

	// Step 6: Log stock change in audit log
	err = s.AuditLogRepository.AfterUpdateAssetStock(*oldAssetStock, newAssetStock)
	if err != nil {
		return logError("AfterUpdateAssetStock", clientID, err, "Failed to update asset stock log")
	}

	// Step 7: Return updated stock response
	return response.AssetStockResponse{
		StockID:         newAssetStock.StockID,
		AssetID:         newAssetStock.AssetID,
		InitialQuantity: newAssetStock.InitialQuantity,
		LatestQuantity:  newAssetStock.LatestQuantity,
		ChangeType:      newAssetStock.ChangeType,
		Quantity:        newAssetStock.Quantity,
	}, nil
}

func (s assetService) UpdateImageAsset(assetID uint, clientID string, metadata []response.AssetImageResponse) error {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logErrorWithNoReturn("GetUserRedis", clientID, err, "Failed to get user redis")
	}

	err = s.AssetImageRepository.UpdateAssetImage(assetID, metadata, data.ClientID)
	if err != nil {
		return logErrorWithNoReturn("UpdateAssetImage", clientID, err, "Failed to update asset image")
	}

	return nil
}

func (s assetService) GetListAsset(clientID string, index int, size int) (interface{}, int64, error) {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logListError("GetUserRedis", clientID, err, "Failed to get user from Redis")
	}

	result, err := s.AssetRepository.GetListAssets(data.ClientID, index, size)
	if err != nil {
		return logListError("GetListAssets", clientID, err, "Failed to get list assets")
	}

	assetCount, err := s.AssetRepository.GetCountAsset(data.ClientID)
	if err != nil {
		return logListError("GetCountAssets", clientID, err, "Failed to get count assets")
	}

	return result, assetCount, nil
}

func (s assetService) GetAssetByID(clientID string, assetID uint) (interface{}, error) {
	_, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logError("GetUserRedis", clientID, err, "Failed to get user redis")
	}

	asset, err := s.AssetRepository.GetAssetResponseByID(clientID, assetID)
	if err != nil {
		return logError("GetAssetResponseByID", clientID, err, "Failed to get asset by MaintenanceTypeID")
	}

	assetImage, _ := s.AssetImageRepository.GetAssetImageResponseByAssetID(assetID)

	asset.Images = *assetImage

	return *asset, nil
}

func (s assetService) UpdateAssetStatus(assetID uint, statusID uint, clientID string) error {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logErrorWithNoReturn("GetUserRedis", clientID, err, "Failed to get user redis")
	}

	oldAsset, err := s.AssetRepository.GetAsset(assetID, clientID)
	if err != nil {
		return logErrorWithNoReturn("GetAsset", clientID, err, "Failed to get asset by MaintenanceTypeID")
	}

	asset, err := s.AssetRepository.UpdateAssetStatus(assetID, statusID, data.ClientID)
	if err != nil {
		return logErrorWithNoReturn("UpdateAssetStatus", clientID, err, "Failed to update asset status")
	}

	err = s.AuditLogRepository.AfterUpdateAsset(*oldAsset, asset)
	if err != nil {
		return logErrorWithNoReturn("AfterUpdateAsset", clientID, err, "Failed to update asset")
	}

	return nil
}

func (s assetService) UpdateAssetCategory(assetID uint, categoryID uint, clientID string) error {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logErrorWithNoReturn("GetUserRedis", clientID, err, "Failed to get user redis")
	}

	oldAsset, err := s.AssetRepository.GetAsset(assetID, clientID)
	if err != nil {
		return logErrorWithNoReturn("GetAsset", clientID, err, "Failed to get asset by MaintenanceTypeID")
	}

	asset, err := s.AssetRepository.UpdateAssetCategory(assetID, categoryID, data.ClientID)
	if err != nil {
		return logErrorWithNoReturn("UpdateAssetCategory", clientID, err, "Failed to update asset category")
	}

	err = s.AuditLogRepository.AfterUpdateAsset(*oldAsset, asset)
	if err != nil {
		return logErrorWithNoReturn("AfterUpdateAsset", clientID, err, "Failed to update asset")
	}
	return nil
}

func (s assetService) DeleteAsset(assetID uint, clientID string) error {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logErrorWithNoReturn("GetUserRedis", clientID, err, "Failed to get user redis")
	}

	err = s.AssetTransaction.DeleteAsset(assetID, clientID, data.ClientID)
	if err != nil {
		return logErrorWithNoReturn("DeleteAsset", clientID, err, "Failed to delete asset")
	}

	return nil
}

func logError(key, clientID string, err error, msg string) (interface{}, error) {
	log.Error().Str("key", key).Str("clientID", clientID).Err(err).Msg(msg)
	if err == nil {
		err = errors.New(msg)
	}
	return nil, err
}

func logListError(key, clientID string, err error, msg string) (interface{}, int64, error) {
	log.Error().Str("key", key).Str("clientID", clientID).Err(err).Msg(msg)
	if err == nil {
		err = errors.New(msg)
	}
	return nil, 0, err
}

func logErrorWithNoReturn(key, clientID string, err error, msg string) error {
	log.Error().Str("key", key).Str("clientID", clientID).Err(err).Msg(msg)
	if err == nil {
		err = errors.New(msg)
	}
	return err
}
