package assets

import (
	request "asset-service/internal/dto/in/assets"
	response "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	repo "asset-service/internal/repository/assets"
	"asset-service/internal/utils"
	"errors"
	"github.com/rs/zerolog/log"
)

type AssetWishlistService interface {
	AddAssetWishlist(assetRequest *request.AssetWishlistRequest, clientID string) (interface{}, error)
	GetAssetWishlist(clientID string) (interface{}, error)
	UpdateAssetWishlist(assetID uint, assetRequest request.UpdateAssetWishlistRequest, clientID string) (interface{}, error)
	DeleteAssetWishlist(assetID uint, clientID string) (interface{}, error)
	GetAssetWishlistByID(assetID uint, clientID string) (interface{}, error)
}

type assetWishlistService struct {
	AssetWishlistRepository repo.AssetWishlistRepository
	AssetRepository         repo.AssetRepository
	AssetCategoryRepository repo.AssetCategoryRepository
	AssetStatusRepository   repo.AssetStatusRepository
	AuditLogRepository      repo.AssetAuditLogRepository
	Redis                   utils.RedisService
}

func NewAssetWishlistService(
	assetWishlistRepository repo.AssetWishlistRepository,
	AssetCategoryRepository repo.AssetCategoryRepository,
	AssetStatusRepository repo.AssetStatusRepository,
	assetRepository repo.AssetRepository,
	AuditLogRepository repo.AssetAuditLogRepository,
	redis utils.RedisService) AssetWishlistService {
	return assetWishlistService{
		AssetWishlistRepository: assetWishlistRepository,
		AssetRepository:         assetRepository,
		AssetCategoryRepository: AssetCategoryRepository,
		AssetStatusRepository:   AssetStatusRepository,
		AuditLogRepository:      AuditLogRepository,
		Redis:                   redis,
	}
}

func (s assetWishlistService) AddAssetWishlist(assetRequest *request.AssetWishlistRequest, clientID string) (interface{}, error) {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetUserRedis").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get user redis")
		return nil, err
	}

	check, err := s.AssetRepository.AssetNameExists(assetRequest.Name, clientID)
	if err != nil {
		log.Error().
			Str("key", "AssetNameExists").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to check asset name exists")
		return nil, err
	}
	if check {
		log.Error().
			Str("key", "AssetNameExists").
			Str("clientID", clientID).
			Err(err).
			Msg("Asset name already exists")
		return nil, errors.New("assets already exists")
	}

	purchaseDate, err := utils.ParseOptionalDate(assetRequest.PurchaseDate)
	if err != nil {
		log.Error().
			Str("key", "ParseOptionalDate").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset status by ID")
	}

	var asset = &assets.Asset{
		Name:         assetRequest.Name,
		UserClientID: clientID,
		Description:  assetRequest.Description,
		CategoryID:   assetRequest.CategoryID,
		StatusID:     assetRequest.StatusID,
		PurchaseDate: purchaseDate,
		Price:        assetRequest.Price,
		IsWishlist:   assetRequest.IsWishlist,
		CreatedBy:    data.ClientID,
		UpdatedBy:    data.ClientID,
	}

	err = s.AssetWishlistRepository.AddAssetWishlist(asset)
	if err != nil {
		log.Error().
			Str("key", "AddAssetWishlist").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to add asset wishlist")
		return nil, err
	}

	assetResult, err := s.AssetWishlistRepository.GetAssetWishlistByID(clientID, asset.AssetID)
	if err != nil {
		log.Error().
			Str("key", "GetAssetByID").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get asset by ID")
		return nil, err
	}

	return assetResult, nil
}

func (s assetWishlistService) GetAssetWishlist(clientID string) (interface{}, error) {
	list, err := s.AssetWishlistRepository.GetAssetWishlistList(clientID)
	if err != nil {
		return nil, err
	}

	var assetResponse []response.AssetResponse
	for _, asset := range list {
		assetResponse = append(assetResponse, response.AssetResponse{
			AssetID:      asset.AssetID,
			UserClientID: asset.UserClientID,
			Name:         asset.Name,
			Description:  asset.Description,
			Status:       response.AssetStatusResponse{AssetStatusID: asset.Status.AssetStatusID, StatusName: asset.Status.StatusName},
			Category:     response.AssetCategoryResponse{AssetCategoryID: asset.Category.AssetCategoryID, CategoryName: asset.Category.CategoryName},
			PurchaseDate: asset.PurchaseDate,
			Price:        asset.Price,
		})
	}

	return assetResponse, nil
}

func (s assetWishlistService) UpdateAssetWishlist(assetID uint, assetRequest request.UpdateAssetWishlistRequest, clientID string) (interface{}, error) {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetUserRedis").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get user redis")
		return nil, err
	}

	if _, err := s.AssetCategoryRepository.GetAssetCategoryById(uint(assetRequest.CategoryID)); err != nil {
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

	asset, err := s.AssetRepository.GetAsset(assetID, clientID)
	if err != nil {
		log.Error().Uint("assetID", assetID).Str("clientID", clientID).Err(err).Msg("❌ Asset not found")
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

	updateData := map[string]interface{}{
		"description":   assetRequest.Description,
		"image_url":     assetRequest.ImageUrl,
		"category_id":   assetRequest.CategoryID,
		"status_id":     assetRequest.StatusID,
		"price":         assetRequest.Price,
		"purchase_date": purchaseDate,
		"notes":         assetRequest.Notes,
		"is_wishlist":   assetRequest.IsWishlist,
		"updated_by":    data.ClientID,
	}

	err = s.AssetWishlistRepository.UpdateAssetWishlist(asset.AssetID, updateData)
	if err != nil {
		return nil, err
	}

	wishlist, err := s.AssetRepository.GetAsset(assetID, clientID)
	if err != nil {
		return nil, err
	}

	err = s.AuditLogRepository.AfterUpdateAsset(*asset, wishlist)
	if err != nil {
		return nil, err
	}

	return wishlist, nil
}

func (s assetWishlistService) DeleteAssetWishlist(assetID uint, clientID string) (interface{}, error) {
	data, err := utils.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().
			Str("key", "GetUserRedis").
			Str("clientID", clientID).
			Err(err).
			Msg("Failed to get user redis")
		return nil, err
	}

	asset, err := s.AssetRepository.GetAsset(assetID, clientID)
	if err != nil {
		return nil, err
	}

	err = s.AssetWishlistRepository.DeleteAssetWishlist(clientID, data.ClientID, asset.AssetID)
	if err != nil {
		return nil, err
	}

	err = s.AuditLogRepository.AfterDeleteAsset(asset)
	return nil, nil
}

func (s assetWishlistService) GetAssetWishlistByID(assetID uint, clientID string) (interface{}, error) {
	asset, err := s.AssetWishlistRepository.GetAssetWishlistByID(clientID, assetID)
	if err != nil {
		return nil, err
	}

	return asset, nil
}
