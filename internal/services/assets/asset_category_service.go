package assets

import (
	request "asset-service/internal/dto/in/assets"
	response "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	"asset-service/internal/models/user"
	repository "asset-service/internal/repository/assets"
	"asset-service/internal/utils"
	"errors"
)

type AssetCategoryService interface {
	AddAssetCategory(assetRequest *request.AssetCategoryRequest, clientID string) (interface{}, error)
	UpdateAssetCategory(assetCategoryID uint, assetCategoryRequest *request.AssetCategoryRequest, clientID string) (interface{}, error)
	GetListAssetCategory(clientID string) (interface{}, error)
	GetAssetCategoryById(categoryID uint) (interface{}, error)
	DeleteAssetCategory(categoryID uint, clientID string) error
}

type assetCategoryService struct {
	AssetCategoryRepository repository.AssetCategoryRepository
	Redis                   utils.RedisService
}

func NewAssetCategoryService(assetCategoryRepository repository.AssetCategoryRepository, redis utils.RedisService) AssetCategoryService {
	return assetCategoryService{AssetCategoryRepository: assetCategoryRepository, Redis: redis}
}

func (s assetCategoryService) AddAssetCategory(assetRequest *request.AssetCategoryRequest, clientID string) (interface{}, error) {
	var user = &user.User{}

	err := s.Redis.GetData(utils.User, clientID, user)
	if err != nil {
		return nil, err
	}
	var assetCategory = &assets.AssetCategory{
		CategoryName: assetRequest.CategoryName,
		Description:  assetRequest.Description,
		CreatedBy:    user.FullName,
		UpdatedBy:    user.FullName,
	}

	err = s.AssetCategoryRepository.GetAssetCategoryByName(assetRequest.CategoryName)
	if err == nil {
		return nil, errors.New("assets category already exists")
	}

	err = s.AssetCategoryRepository.AddAssetCategory(assetCategory)
	if err != nil {
		return nil, err
	}
	var assetCategoryResponse = response.AssetCategoryResponse{
		AssetCategoryID: assetCategory.AssetCategoryID,
		CategoryName:    assetCategory.CategoryName,
		Description:     assetCategory.Description,
	}
	return assetCategoryResponse, nil
}

func (s assetCategoryService) UpdateAssetCategory(assetCategoryID uint, assetCategoryRequest *request.AssetCategoryRequest, clientID string) (interface{}, error) {
	var user = &user.User{}
	var assetCategory = &assets.AssetCategory{}

	err := s.Redis.GetData(utils.User, clientID, user)
	if err != nil {
		return nil, err
	}

	assetCategory, err = s.AssetCategoryRepository.GetAssetCategoryByIdAndNameNotExist(assetCategoryID, assetCategoryRequest.CategoryName)
	if err != nil {
		return nil, errors.New("assets category not found or already exists")
	}

	assetCategory.CategoryName = assetCategoryRequest.CategoryName
	assetCategory.Description = assetCategoryRequest.Description
	assetCategory.UpdatedBy = user.FullName
	err = s.AssetCategoryRepository.UpdateAssetCategory(assetCategory)
	if err != nil {
		return nil, err
	}
	return assetCategory, nil
}

func (s assetCategoryService) GetListAssetCategory(clientID string) (interface{}, error) {
	var user = &user.User{}

	err := s.Redis.GetData(utils.User, clientID, user)
	if err != nil {
		return nil, err
	}

	assetCategories, err := s.AssetCategoryRepository.GetListAssetCategory()
	if err != nil {
		return nil, err
	}
	var assetCategoriesResponse []response.AssetCategoryResponse

	for _, assetCategory := range assetCategories {
		assetCategoriesResponse = append(assetCategoriesResponse, response.AssetCategoryResponse{
			AssetCategoryID: assetCategory.AssetCategoryID,
			CategoryName:    assetCategory.CategoryName,
			Description:     assetCategory.Description,
		})
	}

	return assetCategoriesResponse, nil
}

func (s assetCategoryService) GetAssetCategoryById(categoryID uint) (interface{}, error) {
	var assetCategory *assets.AssetCategory
	assetCategory, err := s.AssetCategoryRepository.GetAssetCategoryById(categoryID)
	if err != nil {
		return nil, errors.New("assets category not found")
	}

	var assetCategoryResponse = response.AssetCategoryResponse{
		AssetCategoryID: assetCategory.AssetCategoryID,
		CategoryName:    assetCategory.CategoryName,
		Description:     assetCategory.Description,
	}

	return assetCategoryResponse, nil
}

func (s assetCategoryService) DeleteAssetCategory(categoryID uint, clientID string) error {
	var user = &user.User{}
	var assetCategory = &assets.AssetCategory{}

	err := s.Redis.GetData(utils.User, clientID, user)
	if err != nil {
		return err
	}

	assetCategory, err = s.AssetCategoryRepository.GetAssetCategoryById(categoryID)
	if err != nil {
		return errors.New("assets category not found")
	}

	assetCategory.DeletedBy = &user.FullName
	err = s.AssetCategoryRepository.DeleteAssetCategory(assetCategory)
	if err != nil {
		return err
	}
	return nil
}
