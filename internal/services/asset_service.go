package services

import (
	"asset-service/internal/dto/in"
	"asset-service/internal/dto/out"
	"asset-service/internal/models"
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
	var asset = &models.Asset{
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
	var user = &models.User{}
	err := utils.GetDataFromRedis(utils.User, clientID, user)

	if err != nil {
		return nil, err
	}

	_, err = s.AssetRepository.GetAssetByName(assetRequest.Name)
	if err == nil {
		return nil, errors.New("asset already exists")
	}
	layout := "2006-01-02" // Date-only format
	parsedDate, err := time.Parse(layout, assetRequest.PurchaseDate)
	if err != nil {
		return nil, fmt.Errorf("invalid purchase date format: %v", err)
	}

	var asset = &models.Asset{
		Name:         assetRequest.Name,
		UserClientID: clientID,
		Description:  assetRequest.Description,
		CategoryID:   assetRequest.CategoryID,
		StatusID:     assetRequest.StatusID,
		PurchaseDate: &parsedDate,
		Value:        assetRequest.Value,
		CreatedBy:    user.FullName,
		UpdatedBy:    user.FullName,
	}

	var result *out.AssetResponse
	result, err = s.AssetRepository.AddAsset(asset)
	if err != nil && result == nil {
		return nil, err
	}

	return result, nil
}

func (s AssetService) GetListAsset(clientID string) (interface{}, error) {
	var user = &models.User{}
	err := utils.GetDataFromRedis(utils.User, clientID, user)
	if err != nil {
		return nil, err
	}

	assets, err := s.AssetRepository.GetListAsset(clientID)
	if err != nil {
		return nil, err
	}

	var assetResponse []struct {
		ID           uint
		Name         string
		Description  string
		CategoryName string
		StatusName   string
		PurchaseDate string
		Value        float64
	}
	for _, asset := range assets {

		assetResponse = append(assetResponse, struct {
			ID           uint
			Name         string
			Description  string
			CategoryName string
			StatusName   string
			PurchaseDate string
			Value        float64
		}{
			ID:           asset.ID,
			Name:         asset.Name,
			Description:  asset.Description,
			CategoryName: asset.CategoryName,
			StatusName:   asset.StatusName,
			PurchaseDate: asset.PurchaseDate.Format("2006-01-02"),
			Value:        asset.Value,
		})
	}

	return assetResponse, nil
}
