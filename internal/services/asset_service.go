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

func (s AssetService) AddAsset(assetRequest *in.AssetRequest, clientID string) (interface{}, error) {
	var user = &user.User{}
	err := utils.GetDataFromRedis(utils.User, clientID, user)

	if err != nil {
		return nil, err
	}

	_, err = s.AssetRepository.GetAssetByNameAndClientID(assetRequest.Name, clientID)
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

func (s AssetService) UpdateAsset(assetID uint, assetRequest struct {
	Description  string  `json:"description"`
	PurchaseDate string  `json:"purchase_date" binding:"required"`
	ExpiryDate   string  `json:"expiry_date"`
	Value        float64 `json:"value" binding:"required"`
}, clientID string) (interface{}, error) {
	var user = &user.User{}
	err := utils.GetDataFromRedis(utils.User, clientID, user)

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
		Value:        assetRequest.Value,
		UpdatedBy:    user.FullName,
	}

	var result *out.AssetResponse
	result, err = s.AssetRepository.UpdateAsset(asset, clientID)
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

	result, err := s.AssetRepository.GetListAsset(clientID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s AssetService) GetAssetByID(clientID string, assetID uint) (out.AssetResponseList, error) {
	var user = &user.User{}
	err := utils.GetDataFromRedis(utils.User, clientID, user)
	if err != nil {
		return out.AssetResponseList{}, err
	}

	asset, err := s.AssetRepository.GetAssetByID(clientID, assetID)
	if err != nil {
		return out.AssetResponseList{}, err
	}

	return *asset, nil
}

func (s AssetService) UpdateAssetStatus(assetID uint, statusID uint, clientID string) error {
	var user = &user.User{}
	err := utils.GetDataFromRedis(utils.User, clientID, user)
	if err != nil {
		return err
	}

	err = s.AssetRepository.UpdateAssetStatus(assetID, statusID, user.ClientID, user.FullName)
	if err != nil {
		return err
	}

	return nil
}

func (s AssetService) UpdateAssetCategory(assetID uint, categoryID uint, clientID string) error {
	var user = &user.User{}
	err := utils.GetDataFromRedis(utils.User, clientID, user)
	if err != nil {
		return err
	}

	err = s.AssetRepository.UpdateAssetCategory(assetID, categoryID, user.ClientID, user.FullName)
	if err != nil {
		return err
	}

	return nil

}

func (s AssetService) DeleteAsset(assetID uint, clientID string) error {
	var user = &user.User{}
	err := utils.GetDataFromRedis(utils.User, clientID, user)
	if err != nil {
		return err
	}

	err = s.AssetRepository.DeleteAsset(assetID, user.ClientID, user.FullName)
	if err != nil {
		return err
	}

	return nil
}
