package assets

import (
	"asset-service/internal/utils/redis"
	"asset-service/internal/utils/text"
	"errors"

	request "asset-service/internal/dto/in/assets"
	response "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	repository "asset-service/internal/repository/assets"
	"asset-service/internal/utils"

	"github.com/rs/zerolog/log"
)

type AssetCategoryService interface {
	AddAssetCategory(assetRequest *request.AssetCategoryRequest, credentialKey string, clientID string) (interface{}, error)
	UpdateAssetCategory(assetCategoryID uint, assetCategoryRequest *request.AssetCategoryRequest, clientID string, credentialKey string) (interface{}, error)
	GetListAssetCategory(clientID string, size int, index int) (interface{}, int64, error)
	GetAssetCategoryById(categoryID uint, clientID string) (interface{}, error)
	DeleteAssetCategory(categoryID uint, clientID string) error
}

type assetCategoryService struct {
	AssetCategoryRepository repository.AssetCategoryRepository
	AssetRepository         repository.AssetRepository
	AssetAuditLogRepository repository.AssetAuditLogRepository
	Redis                   redis.RedisService
}

func NewAssetCategoryService(
	assetCategoryRepository repository.AssetCategoryRepository,
	AssetRepository repository.AssetRepository,
	AssetAuditLogRepository repository.AssetAuditLogRepository,
	redis redis.RedisService) AssetCategoryService {
	return &assetCategoryService{
		AssetCategoryRepository: assetCategoryRepository,
		AssetRepository:         AssetRepository,
		AssetAuditLogRepository: AssetAuditLogRepository,
		Redis:                   redis,
	}
}

func (s *assetCategoryService) AddAssetCategory(assetRequest *request.AssetCategoryRequest, credentialKey string, clientID string) (interface{}, error) {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve data from Redis")
		return nil, err
	}

	err = text.CheckCredentialKey(s.Redis, credentialKey, data.ClientID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("credential key check failed")
		return nil, err
	}

	existingCategory, err := s.AssetCategoryRepository.GetAssetCategoryByNameAndClientID(assetRequest.CategoryName, data.ClientID)
	if existingCategory != nil {
		log.Warn().Str("category_name", assetRequest.CategoryName).Msg("Asset category already exists")
		return nil, errors.New("asset category already exists")
	}

	assetCategory := &assets.AssetCategory{
		UserClientID: clientID,
		CategoryName: assetRequest.CategoryName,
		Description:  assetRequest.Description,
		CreatedBy:    &data.ClientID,
		UpdatedBy:    &data.ClientID,
	}

	err = s.AssetCategoryRepository.AddAssetCategory(assetCategory)
	if err != nil {
		log.Error().Err(err).Str("category_name", assetRequest.CategoryName).Msg("Failed to add asset category")
		return nil, err
	}

	if auditErr := s.AssetAuditLogRepository.AfterCreateAssetCategory(assetCategory); auditErr != nil {
		log.Warn().Err(auditErr).
			Str("category_name", assetCategory.CategoryName).
			Msg("⚠ Audit log failed after adding asset category")
	}

	log.Info().Str("category_name", assetRequest.CategoryName).Msg("✅ Asset category added successfully")

	return response.AssetCategoryResponse{
		AssetCategoryID: assetCategory.AssetCategoryID,
		CategoryName:    assetCategory.CategoryName,
		Description:     assetCategory.Description,
	}, nil
}

func (s *assetCategoryService) UpdateAssetCategory(assetCategoryID uint, assetCategoryRequest *request.AssetCategoryRequest, clientID string, credentialKey string) (interface{}, error) {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve user from Redis")
		return nil, err
	}

	err = text.CheckCredentialKey(s.Redis, credentialKey, data.ClientID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("credential key check failed")
		return nil, err
	}

	oldAssetCategory, err := s.AssetCategoryRepository.GetAssetCategoryById(assetCategoryID, clientID)
	if err != nil {
		return nil, err
	}

	assetCategory, err := s.AssetCategoryRepository.GetAssetCategoryByIdAndNameNotExist(assetCategoryID, assetCategoryRequest.CategoryName)
	if err != nil {
		log.Warn().Uint("asset_category_id", assetCategoryID).Msg("Asset category not found or already exists")
		return nil, errors.New("asset category not found or already exists")
	}

	assetCategory.CategoryName = assetCategoryRequest.CategoryName
	assetCategory.Description = assetCategoryRequest.Description
	assetCategory.UpdatedBy = &data.ClientID

	err = s.AssetCategoryRepository.UpdateAssetCategory(assetCategory, clientID)
	if err != nil {
		log.Error().Err(err).Uint("asset_category_id", assetCategoryID).Msg("Failed to update asset category")
		return nil, err
	}

	if auditErr := s.AssetAuditLogRepository.AfterUpdateAssetCategory(oldAssetCategory, assetCategory); auditErr != nil {
		log.Warn().Err(auditErr).
			Uint("asset_category_id", assetCategory.AssetCategoryID).
			Msg("⚠ Audit log failed after updating asset category")
	}

	log.Info().Uint("asset_category_id", assetCategoryID).Msg("✅ Asset category updated successfully")

	return response.AssetCategoryResponse{
		AssetCategoryID: assetCategory.AssetCategoryID,
		CategoryName:    assetCategory.CategoryName,
		Description:     assetCategory.Description,
	}, nil
}

func (s *assetCategoryService) GetListAssetCategory(clientID string, size int, index int) (interface{}, int64, error) {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve user from Redis")
		return nil, 0, err
	}

	total, err := s.AssetCategoryRepository.GetCountAssetCategory(data.ClientID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve asset category count")
		return nil, 0, err
	}

	assetCategories, err := s.AssetCategoryRepository.GetListAssetCategory(data.ClientID, size, index)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve asset categories")
		return nil, total, err
	}

	var assetCategoriesResponse []response.AssetCategoryResponse
	for _, assetCategory := range assetCategories {
		assetCategoriesResponse = append(assetCategoriesResponse, response.AssetCategoryResponse{
			AssetCategoryID: assetCategory.AssetCategoryID,
			CategoryName:    assetCategory.CategoryName,
			Description:     assetCategory.Description,
		})
	}

	return assetCategoriesResponse, total, nil
}

func (s *assetCategoryService) GetAssetCategoryById(categoryID uint, clientID string) (interface{}, error) {
	assetCategory, err := s.AssetCategoryRepository.GetAssetCategoryById(categoryID, clientID)
	if err != nil {
		log.Warn().Uint("asset_category_id", categoryID).Msg("Asset category not found")
		return nil, errors.New("asset category not found")
	}

	return response.AssetCategoryResponse{
		AssetCategoryID: assetCategory.AssetCategoryID,
		CategoryName:    assetCategory.CategoryName,
		Description:     assetCategory.Description,
	}, nil
}

func (s *assetCategoryService) DeleteAssetCategory(categoryID uint, clientID string) error {
	data, err := redis.GetUserRedis(s.Redis, utils.User, clientID)
	if err != nil {
		log.Error().Str("clientID", clientID).Err(err).Msg("Failed to retrieve user from Redis")
		return err
	}

	assetCategory, err := s.AssetCategoryRepository.GetAssetCategoryById(categoryID, clientID)
	if err != nil {
		log.Warn().Uint("asset_category_id", categoryID).Msg("Asset category not found")
		return errors.New("asset category not found")
	}

	asset, err := s.AssetRepository.GetAssetByCategoryID(assetCategory.AssetCategoryID, clientID)
	if err != nil {
		log.Error().Err(err).Uint("asset_category_id", categoryID).Msg("Failed to get asset by category MaintenanceTypeID")
		return err
	}

	if len(asset) > 0 {
		log.Warn().Uint("asset_category_id", categoryID).Msg("Asset category is still in use by assets")
		return errors.New("asset category is still in use")
	}

	assetCategory.DeletedBy = &data.ClientID
	if err = s.AssetCategoryRepository.DeleteAssetCategory(assetCategory); err != nil {
		log.Error().Err(err).Uint("asset_category_id", categoryID).Msg("Failed to delete asset category")
		return err
	}

	if err = s.AssetAuditLogRepository.AfterDeleteAssetCategory(assetCategory); err != nil {
		log.Error().Err(err).Uint("asset_category_id", categoryID).Msg("Failed to audit log after deleting asset category")
		return err
	}

	log.Info().Uint("asset_category_id", categoryID).Msg("✅ Asset category deleted successfully")
	return nil
}
