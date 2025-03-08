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
	"sync"
)

type AssetService interface {
	AddAsset(assetRequest *request.AssetRequest, images []response.AssetImageResponse, clientID string) (interface{}, error)
	UpdateAsset(assetID uint, assetRequest request.UpdateAssetRequest, clientID string) (interface{}, error)
	UpdateStockAsset(isAdded bool, assetID uint, stock struct {
		Stock  int     `json:"stock" binding:"required"`
		Reason *string `json:"reason"`
	}, clientID string) (interface{}, error)
	GetListAsset(clientID string) (interface{}, error)
	GetAssetByID(clientID string, assetID uint) (interface{}, error)
	UpdateAssetStatus(assetID uint, statusID uint, clientID string) error
	UpdateAssetCategory(assetID uint, categoryID uint, clientID string) error
	DeleteAsset(assetID uint, clientID string) error
}

type assetService struct {
	AssetRepository         repo.AssetRepository
	AssetCategoryRepository repo.AssetCategoryRepository
	AssetStatusRepository   repo.AssetStatusRepository
	AssetImageRepository    repo.AssetImageRepository
	Redis                   utils.RedisService
	AuditLogRepository      repo.AssetAuditLogRepository
	AssetTransaction        transaction.AssetTransactionRepository
	AssetStockRepository    repo.AssetStockRepository
}

func NewAssetService(assetRepository repo.AssetRepository,
	assetCategoryRepository repo.AssetCategoryRepository,
	assetStatusRepository repo.AssetStatusRepository,
	assetImageRepository repo.AssetImageRepository,
	log repo.AssetAuditLogRepository,
	redis utils.RedisService,
	assetTransaction transaction.AssetTransactionRepository,
	assetStockRepository repo.AssetStockRepository) AssetService {
	return assetService{
		AssetRepository:         assetRepository,
		AssetCategoryRepository: assetCategoryRepository,
		AssetStatusRepository:   assetStatusRepository,
		AssetImageRepository:    assetImageRepository,
		AuditLogRepository:      log,
		Redis:                   redis,
		AssetTransaction:        assetTransaction,
		AssetStockRepository:    assetStockRepository}
}

func (s assetService) AddAsset(assetRequest *request.AssetRequest, images []response.AssetImageResponse, clientID string) (interface{}, error) {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logError("GetRedisData", clientID, err, "Failed to get data from redis")
	}

	if exists, err := s.AssetRepository.AssetNameExists(assetRequest.Name, clientID); err != nil {
		return logError("AssetNameExists", clientID, err, "Failed to check asset name exists")
	} else if exists {
		return logError("AssetNameExists", clientID, errors.New("assets already exists"), "Asset name already exists")
	}

	if _, err := s.AssetCategoryRepository.GetAssetCategoryById(uint(assetRequest.CategoryID), clientID); err != nil {
		return logError("GetAssetCategoryById", clientID, errors.New("category not found"), "Failed to get asset category by ID")
	}

	if _, err := s.AssetStatusRepository.GetAssetStatusByID(uint(assetRequest.StatusID)); err != nil {
		return logError("GetAssetStatusByID", clientID, errors.New("status not found"), "Failed to get asset status by ID")
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
		CreatedBy:          data.ClientID,
		UpdatedBy:          data.ClientID,
	}

	if err := s.AssetRepository.AddAsset(asset); err != nil {
		return logError("AddAsset", clientID, err, "Failed to add asset")
	}

	if len(images) > 0 {
		var assetImages []assets.AssetImage
		for _, image := range images {
			assetImages := append(assetImages, assets.AssetImage{
				UserClientID: clientID,
				AssetID:      asset.AssetID,
				ImageURL:     image.ImageURL,
				CreatedBy:    data.ClientID,
				UpdatedBy:    data.ClientID,
			})
			if err := s.AssetImageRepository.AddAssetImage(assetImages); err != nil {
				return logError("AddAssetImage", clientID, err, "Failed to add asset image")
			}
		}
	}

	assetStock := &assets.AssetStock{
		AssetID:         asset.AssetID,
		UserClientID:    clientID,
		InitialQuantity: asset.Stock,
		LatestQuantity:  asset.Stock,
		ChangeType:      "INCREASE",
		Quantity:        asset.Stock,
		Reason:          nil,
		CreatedBy:       clientID,
	}

	if err := s.AssetStockRepository.AddAssetStock(assetStock); err != nil {
		return logError("AddAssetStock", clientID, err, "Failed to add asset stock")
	}
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := s.AuditLogRepository.AfterCreateAsset(asset); err != nil {
			log.Error().Str("clientID", clientID).Err(err).Msg("Failed to create asset log")
		}
	}()

	go func() {
		defer wg.Done()
		if err := s.AuditLogRepository.AfterCreateAssetStock(assetStock); err != nil {
			log.Error().Str("clientID", clientID).Err(err).Msg("Failed to create asset stock log")
		}
	}()

	result, err := s.AssetRepository.GetAssetResponseByID(clientID, asset.AssetID)
	if err != nil {
		return logError("GetAssetResponseByID", clientID, err, "Failed to get asset by ID")
	}

	assetImage, err := s.AssetImageRepository.GetAssetImageResponseByAssetID(asset.AssetID)
	if err != nil {
		return logError("GetAssetImageByAssetID", clientID, err, "Failed to get asset image by asset ID")
	}

	result.Images = *assetImage
	result.Stock = response.AssetStockResponse{
		StockID:         assetStock.StockID,
		AssetID:         assetStock.AssetID,
		InitialQuantity: assetStock.InitialQuantity,
		LatestQuantity:  assetStock.LatestQuantity,
		ChangeType:      assetStock.ChangeType,
		Quantity:        assetStock.Quantity,
		Reason:          nil,
	}

	log.Info().Str("key", "GetAssetResponseByID").Str("clientID", clientID).Fields(asset).Msg("Success to get asset by ID")
	return result, nil
}

func (s assetService) UpdateAsset(assetID uint, assetRequest request.UpdateAssetRequest, clientID string) (interface{}, error) {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logError("GetUserRedis", clientID, err, "Failed to get user redis")
	}

	oldAsset, err := s.AssetRepository.GetAsset(assetID, clientID)
	if err != nil {
		return logError("GetAsset", clientID, err, "Failed to get asset by ID")
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
		UpdatedBy:          data.ClientID,
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
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logError("GetUserRedis", clientID, err, "Failed to get user from Redis")
	}

	// Step 2: Retrieve asset and stock data
	asset, err := s.AssetRepository.GetAsset(assetID, clientID)
	if err != nil {
		return logError("GetAsset", clientID, err, "Failed to get asset by ID")
	}

	oldAssetStock, err := s.AssetStockRepository.GetAssetStockByAssetID(assetID, clientID)
	if err != nil {
		return logError("GetAssetStockByAssetID", clientID, err, "Failed to get asset stock by asset ID")
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
		UpdatedBy:       data.ClientID,
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

func (s assetService) GetListAsset(clientID string) (interface{}, error) {
	result, err := s.AssetRepository.GetListAssets(clientID)
	if err != nil {
		return logError("GetListAssets", clientID, err, "Failed to get list assets")
	}

	return result, nil
}

func (s assetService) GetAssetByID(clientID string, assetID uint) (interface{}, error) {
	_, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logError("GetUserRedis", clientID, err, "Failed to get user redis")
	}

	asset, err := s.AssetRepository.GetAssetResponseByID(clientID, assetID)
	if err != nil {
		return logError("GetAssetResponseByID", clientID, err, "Failed to get asset by ID")
	}

	assetImage, _ := s.AssetImageRepository.GetAssetImageResponseByAssetID(assetID)

	asset.Images = *assetImage

	return *asset, nil
}

func (s assetService) UpdateAssetStatus(assetID uint, statusID uint, clientID string) error {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logErrorWithNoReturn("GetUserRedis", clientID, err, "Failed to get user redis")
	}

	oldAsset, err := s.AssetRepository.GetAsset(assetID, clientID)
	if err != nil {
		return logErrorWithNoReturn("GetAsset", clientID, err, "Failed to get asset by ID")
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
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		return logErrorWithNoReturn("GetUserRedis", clientID, err, "Failed to get user redis")
	}

	oldAsset, err := s.AssetRepository.GetAsset(assetID, clientID)
	if err != nil {
		return logErrorWithNoReturn("GetAsset", clientID, err, "Failed to get asset by ID")
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
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
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
	return nil, err
}

func logErrorWithNoReturn(key, clientID string, err error, msg string) error {
	log.Error().Str("key", key).Str("clientID", clientID).Err(err).Msg(msg)
	return err
}
