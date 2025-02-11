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

type AssetService interface {
	AddAsset(assetRequest *request.AssetRequest, clientID string) (interface{}, error)
	UpdateAsset(assetID uint, assetRequest struct {
		Description  string  `json:"description"`
		PurchaseDate string  `json:"purchase_date" binding:"required"`
		ExpiryDate   string  `json:"expiry_date"`
		Price        float64 `json:"price" binding:"required"`
	}, clientID string) (interface{}, error)
	GetListAsset(clientID string) (interface{}, error)
	GetAssetByID(clientID string, assetID uint) (response.AssetResponseList, error)
	UpdateAssetStatus(assetID uint, statusID uint, clientID string) error
	UpdateAssetCategory(assetID uint, categoryID uint, clientID string) error
	DeleteAsset(assetID uint, clientID string) error
}

type assetService struct {
	AssetRepository            repo.AssetRepository
	AssetMaintenanceRepository repo.AssetMaintenanceRepository
	Redis                      utils.RedisService
}

func NewAssetService(assetRepository repo.AssetRepository, assetMaintenanceRepository repo.AssetMaintenanceRepository, redis utils.RedisService) AssetService {
	return assetService{AssetRepository: assetRepository, AssetMaintenanceRepository: assetMaintenanceRepository, Redis: redis}
}

func (s assetService) AddAsset(assetRequest *request.AssetRequest, clientID string) (interface{}, error) {
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
	purchaseDate, err := time.Parse(layout, assetRequest.PurchaseDate)
	if err != nil {
		return nil, fmt.Errorf("invalid purchase date format: %v", err)
	}

	var expiryDate *time.Time = nil
	if assetRequest.ExpiryDate != "" {
		parsedExpiryDate, err := time.Parse(layout, assetRequest.ExpiryDate)
		if err != nil {
			return nil, fmt.Errorf("invalid expiry date format: %v", err)
		}
		expiryDate = &parsedExpiryDate
	}
	var warrantyExpiry *time.Time = nil
	if assetRequest.WarrantyExpiry != "" {
		parsedExpiryDate, err := time.Parse(layout, assetRequest.WarrantyExpiry)
		if err != nil {
			return nil, fmt.Errorf("invalid expiry date format: %v", err)
		}
		warrantyExpiry = &parsedExpiryDate
	}

	var asset = &assets.Asset{
		UserClientID:       clientID,
		AssetCode:          assetRequest.AssetCode,
		Name:               assetRequest.Name,
		Description:        assetRequest.Description,
		Barcode:            assetRequest.Barcode,
		CategoryID:         assetRequest.CategoryID,
		StatusID:           assetRequest.StatusID,
		PurchaseDate:       &purchaseDate,
		ExpiryDate:         expiryDate,
		WarrantyExpiryDate: warrantyExpiry,
		Price:              assetRequest.Price,
		Stock:              assetRequest.Stock,
		CreatedBy:          user.FullName,
		UpdatedBy:          user.FullName,
	}

	var result *response.AssetResponseList
	result, err = s.AssetRepository.AddAsset(asset)
	if err != nil && result == nil {
		return nil, err
	}

	return result, nil
}

func (s assetService) UpdateAsset(assetID uint, assetRequest struct {
	Description  string  `json:"description"`
	PurchaseDate string  `json:"purchase_date" binding:"required"`
	ExpiryDate   string  `json:"expiry_date"`
	Price        float64 `json:"price" binding:"required"`
}, clientID string) (interface{}, error) {
	var user = &user.User{}
	err := s.Redis.GetData(utils.User, clientID, user)

	if err != nil {
		return nil, err
	}

	layout := "2006-01-02" // Date-only format
	purchaseDate, err := time.Parse(layout, assetRequest.PurchaseDate)
	if err != nil {
		return nil, fmt.Errorf("invalid purchase date format: %v", err)
	}

	var expiryDate *time.Time = nil
	if assetRequest.ExpiryDate != "" {
		parsedExpiryDate, err := time.Parse(layout, assetRequest.ExpiryDate)
		if err != nil {
			return nil, fmt.Errorf("invalid expiry date format: %v", err)
		}
		expiryDate = &parsedExpiryDate
	}

	var asset = &assets.Asset{
		AssetID:      assetID,
		Description:  assetRequest.Description,
		PurchaseDate: &purchaseDate,
		ExpiryDate:   expiryDate,
		Price:        assetRequest.Price,
		UpdatedBy:    user.FullName,
	}

	var result *response.AssetResponse
	result, err = s.AssetRepository.UpdateAsset(asset, clientID)
	if err != nil && result == nil {
		return nil, err
	}

	return result, nil
}

func (s assetService) GetListAsset(clientID string) (interface{}, error) {
	var user = &user.User{}
	err := s.Redis.GetData(utils.User, clientID, user)
	if err != nil {
		return nil, err
	}

	result, err := s.AssetRepository.GetListAsset(clientID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s assetService) GetAssetByID(clientID string, assetID uint) (response.AssetResponseList, error) {
	var user = &user.User{}
	err := s.Redis.GetData(utils.User, clientID, user)
	if err != nil {
		return response.AssetResponseList{}, err
	}

	asset, err := s.AssetRepository.GetAssetByID(clientID, assetID)
	if err != nil {
		return response.AssetResponseList{}, err
	}

	return *asset, nil
}

func (s assetService) UpdateAssetStatus(assetID uint, statusID uint, clientID string) error {
	var user = &user.User{}
	err := s.Redis.GetData(utils.User, clientID, user)
	if err != nil {
		return err
	}

	err = s.AssetRepository.UpdateAssetStatus(assetID, statusID, user.ClientID, user.FullName)
	if err != nil {
		return err
	}

	return nil
}

func (s assetService) UpdateAssetCategory(assetID uint, categoryID uint, clientID string) error {
	var user = &user.User{}
	err := s.Redis.GetData(utils.User, clientID, user)
	if err != nil {
		return err
	}

	err = s.AssetRepository.UpdateAssetCategory(assetID, categoryID, user.ClientID, user.FullName)
	if err != nil {
		return err
	}

	return nil

}

func (s assetService) DeleteAsset(assetID uint, clientID string) error {
	var user = &user.User{}
	err := s.Redis.GetData(utils.User, clientID, user)
	if err != nil {
		return err
	}

	err = s.AssetRepository.DeleteAsset(assetID, user.ClientID, user.FullName)
	if err != nil {
		return err
	}

	return nil
}
