package services

import (
	"asset-service/internal/dto/in"
	"asset-service/internal/dto/out"
	"asset-service/internal/models/assets"
	"asset-service/internal/models/user"
	"asset-service/internal/repository"
	"asset-service/internal/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type AssetService struct {
	AssetRepository *repository.AssetRepository
}

func NewAssetService(db *gorm.DB) *AssetService {
	r := repository.NewAssetRepository(db)
	return &AssetService{AssetRepository: r}
}

func (s AssetService) RegisterAsset(s2 *struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"optional"`
	CategoryID  uint   `json:"category_id" binding:"required"`
}, userID uint) (interface{}, error) {
	var asset = &assets.Asset{
		Name:        s2.Name,
		Description: s2.Description,
	}
	err := s.AssetRepository.RegisterAsset(&asset)
	if err != nil {
		return nil, err
	}
	return asset, nil
}

func (s AssetService) AddAsset(assetRequest *in.AssetRequest, clientID string) (interface{}, error) {
	var user = &user.User{}
	err := utils.GetDataFromRedis(utils.User, clientID, user)

	if err != nil {
		return nil, err
	}

	_, err = s.AssetRepository.GetAssetByName(assetRequest.Name)
	if err == nil {
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

	var asset = &assets.Asset{
		Name:         assetRequest.Name,
		UserClientID: clientID,
		Description:  assetRequest.Description,
		CategoryID:   assetRequest.CategoryID,
		StatusID:     assetRequest.StatusID,
		PurchaseDate: &purchaseDate,
		ExpiryDate:   expiryDate,
		Value:        assetRequest.Value,
		CreatedBy:    user.FullName,
		UpdatedBy:    user.FullName,
	}

	var maintenanceDate *time.Time = nil
	if assetRequest.MaintenanceDate != "" {
		parsedMaintenanceDate, err := time.Parse(layout, assetRequest.MaintenanceDate)
		if err != nil {
			return nil, fmt.Errorf("invalid expiry date format: %v", err)
		}
		maintenanceDate = &parsedMaintenanceDate
	}

	var assetMaintenance = &assets.AssetMaintenance{
		MaintenanceDate:    *maintenanceDate,
		MaintenanceCost:    assetRequest.MaintenanceCost,
		MaintenanceDetails: nil,
		CreatedBy:          user.FullName,
		UpdatedBy:          user.FullName,
	}

	if assetRequest.MaintenanceDetail != "" {
		assetMaintenance.MaintenanceDetails = &assetRequest.MaintenanceDetail
	}

	var result *out.AssetResponse
	result, err = s.AssetRepository.AddAsset(asset, assetMaintenance)
	if err != nil && result == nil {
		return nil, err
	}

	return result, nil
}

func (s AssetService) GetListAsset(clientID string) (interface{}, error) {
	var user = &user.User{}
	err := utils.GetDataFromRedis(utils.User, clientID, user)
	if err != nil {
		return nil, err
	}

	assets, err := s.AssetRepository.GetListAsset(clientID)
	if err != nil {
		return nil, err
	}

	return assets, nil
}
