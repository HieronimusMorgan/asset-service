package assets

import (
	request "asset-service/internal/dto/in/assets"
	response "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	"asset-service/internal/models/user"
	repo "asset-service/internal/repository/assets"
	"asset-service/internal/utils"
	"errors"
	"fmt"
	"time"
)

type AssetWishlistService interface {
	AddAssetWishlist(assetRequest *request.AssetWishlistRequest, clientID string) (interface{}, error)
	GetAssetWishlist(clientID string) (interface{}, error)
	UpdateAssetWishlist(assetID uint, assetRequest struct {
		Description  string  `json:"description"`
		PurchaseDate string  `json:"purchase_date" binding:"required"`
		CategoryID   int     `json:"category_id" binding:"required"`
		StatusID     int     `json:"status_id" binding:"required"`
		Price        float64 `json:"price" binding:"required"`
		IsWishlist   bool    `json:"is_wishlist" binding:"required"`
	}, clientID string) (interface{}, error)
	DeleteAssetWishlist(assetID uint, clientID string) (interface{}, error)
	GetAssetWishlistByID(assetID uint, clientID string) (interface{}, error)
}

type assetWishlistService struct {
	AssetWishlistRepository repo.AssetWishlistRepository
	AssetRepository         repo.AssetRepository
	Redis                   utils.RedisService
}

func NewAssetWishlistService(assetWishlistRepository repo.AssetWishlistRepository, assetRepository repo.AssetRepository, redis utils.RedisService) AssetWishlistService {
	return assetWishlistService{AssetWishlistRepository: assetWishlistRepository, AssetRepository: assetRepository, Redis: redis}
}

func (s assetWishlistService) AddAssetWishlist(assetRequest *request.AssetWishlistRequest, clientID string) (interface{}, error) {
	var user = &user.User{}
	err := s.Redis.GetData(utils.User, clientID, user)

	if err != nil {
		return nil, err
	}

	check, err := s.AssetRepository.AssetNameExists(assetRequest.Name, clientID)
	if err != nil {
		return nil, err
	}
	if check {
		return nil, errors.New("assets already exists")
	}

	layout := "2006-01-02" // Date-only format
	if assetRequest.PurchaseDate == "" {
		assetRequest.PurchaseDate = time.Now().Format("2006-01-02")
	}

	purchaseDate, err := time.Parse(layout, assetRequest.PurchaseDate)
	if err != nil {
		return nil, fmt.Errorf("invalid purchase date format: %v", err)
	}

	var asset = &assets.Asset{
		Name:         assetRequest.Name,
		UserClientID: clientID,
		Description:  assetRequest.Description,
		CategoryID:   assetRequest.CategoryID,
		StatusID:     assetRequest.StatusID,
		PurchaseDate: &purchaseDate,
		Price:        assetRequest.Price,
		IsWishlist:   assetRequest.IsWishlist,
		CreatedBy:    user.FullName,
		UpdatedBy:    user.FullName,
	}

	var result *response.AssetResponseList
	result, err = s.AssetWishlistRepository.AddAssetWishlist(asset)
	if err != nil && result == nil {
		return nil, err
	}

	return result, nil
}

func (s assetWishlistService) GetAssetWishlist(clientID string) (interface{}, error) {
	var user = &user.User{}
	err := s.Redis.GetData(utils.User, clientID, user)

	if err != nil {
		return nil, err
	}

	list, err := s.AssetWishlistRepository.GetAssetWishlistList(clientID)
	if err != nil {
		return nil, err
	}

	var assetResponse []response.AssetResponseList
	for _, asset := range list {
		assetResponse = append(assetResponse, response.AssetResponseList{
			ID:           asset.ID,
			ClientID:     asset.ClientID,
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

func (s assetWishlistService) UpdateAssetWishlist(assetID uint, assetRequest struct {
	Description  string  `json:"description"`
	PurchaseDate string  `json:"purchase_date" binding:"required"`
	CategoryID   int     `json:"category_id" binding:"required"`
	StatusID     int     `json:"status_id" binding:"required"`
	Price        float64 `json:"price" binding:"required"`
	IsWishlist   bool    `json:"is_wishlist" binding:"required"`
}, clientID string) (interface{}, error) {
	var user = &user.User{}
	err := s.Redis.GetData(utils.User, clientID, user)

	if err != nil {
		return nil, err
	}

	asset, err := s.AssetRepository.GetAsset(assetID, clientID)
	if err != nil {
		return nil, err
	}

	layout := "2006-01-02" // Date-only format
	purchaseDate, err := time.Parse(layout, assetRequest.PurchaseDate)
	if err != nil {
		return nil, fmt.Errorf("invalid purchase date format: %v", err)
	}

	asset.Description = assetRequest.Description
	asset.PurchaseDate = &purchaseDate
	asset.CategoryID = assetRequest.CategoryID
	asset.StatusID = assetRequest.StatusID
	asset.Price = assetRequest.Price
	asset.IsWishlist = assetRequest.IsWishlist
	asset.UpdatedBy = user.FullName

	var result *response.AssetResponseList
	result, err = s.AssetWishlistRepository.UpdateAssetWishlist(asset)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s assetWishlistService) DeleteAssetWishlist(assetID uint, clientID string) (interface{}, error) {
	var user = &user.User{}
	err := s.Redis.GetData(utils.User, clientID, user)

	if err != nil {
		return nil, err
	}

	asset, err := s.AssetRepository.GetAsset(assetID, clientID)
	if err != nil {
		return nil, err
	}

	err = s.AssetWishlistRepository.DeleteAssetWishlist(clientID, user.FullName, asset.AssetID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s assetWishlistService) GetAssetWishlistByID(assetID uint, clientID string) (interface{}, error) {
	var user = &user.User{}
	err := s.Redis.GetData(utils.User, clientID, user)

	if err != nil {
		return nil, err
	}

	asset, err := s.AssetWishlistRepository.GetAssetWishlistByID(clientID, assetID)
	if err != nil {
		return nil, err
	}

	return asset, nil
}
