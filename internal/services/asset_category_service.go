package services

import (
	"asset-service/internal/dto/in"
	"asset-service/internal/dto/out"
	"asset-service/internal/models/asset"
	"asset-service/internal/models/user"
	"asset-service/internal/repository"
	"asset-service/internal/utils"
	"errors"
	"gorm.io/gorm"
)

type AssetCategoryService struct {
	AssetCategoryRepository *repository.AssetCategoryRepository
}

func NewAssetCategoryService(db *gorm.DB) *AssetCategoryService {
	r := repository.NewAssetCategoryRepository(db)
	return &AssetCategoryService{AssetCategoryRepository: r}
}

func (s AssetCategoryService) AddAssetCategory(assetRequest *in.AssetCategoryRequest, clientID string) (interface{}, error) {
	var user = &user.User{}

	err := utils.GetDataFromRedis(utils.User, clientID, user)
	if err != nil {
		return nil, err
	}
	var assetCategory = &asset.AssetCategory{
		CategoryName: assetRequest.CategoryName,
		Description:  assetRequest.Description,
		CreatedBy:    user.FullName,
		UpdatedBy:    user.FullName,
	}

	err = s.AssetCategoryRepository.GetAssetCategoryByName(assetRequest.CategoryName)
	if err == nil {
		return nil, errors.New("asset category already exists")
	}

	err = s.AssetCategoryRepository.AddAssetCategory(&assetCategory)
	if err != nil {
		return nil, err
	}
	var assetCategoryResponse = out.AssetCategoryResponse{
		AssetCategoryID: assetCategory.AssetCategoryID,
		CategoryName:    assetCategory.CategoryName,
		Description:     assetCategory.Description,
	}
	return assetCategoryResponse, nil
}

func (s AssetCategoryService) UpdateAssetCategory(assetCategoryID uint, assetCategoryRequest *in.AssetCategoryRequest, clientID string) (interface{}, error) {
	var user = &user.User{}
	var assetCategory = &asset.AssetCategory{}

	err := utils.GetDataFromRedis(utils.User, clientID, user)
	if err != nil {
		return nil, err
	}

	assetCategory, err = s.AssetCategoryRepository.GetAssetCategoryByIdAndNameNotExist(assetCategoryID, assetCategoryRequest.CategoryName)
	if err != nil {
		return nil, errors.New("asset category not found or already exists")
	}

	assetCategory.CategoryName = assetCategoryRequest.CategoryName
	assetCategory.Description = assetCategoryRequest.Description
	assetCategory.UpdatedBy = user.FullName
	err = s.AssetCategoryRepository.UpdateAssetCategory(&assetCategory)
	if err != nil {
		return nil, err
	}
	return assetCategory, nil
}

func (s AssetCategoryService) GetListAssetCategory(clientID string) (interface{}, error) {
	var user = &user.User{}

	err := utils.GetDataFromRedis(utils.User, clientID, user)
	if err != nil {
		return nil, err
	}

	assetCategories, err := s.AssetCategoryRepository.GetListAssetCategory()
	if err != nil {
		return nil, err
	}
	var assetCategoriesResponse []out.AssetCategoryResponse

	for _, assetCategory := range assetCategories {
		assetCategoriesResponse = append(assetCategoriesResponse, out.AssetCategoryResponse{
			AssetCategoryID: assetCategory.AssetCategoryID,
			CategoryName:    assetCategory.CategoryName,
			Description:     assetCategory.Description,
		})
	}

	return assetCategoriesResponse, nil
}

func (s AssetCategoryService) GetAssetCategoryById(categoryID uint) (interface{}, error) {
	var assetCategory *asset.AssetCategory
	assetCategory, err := s.AssetCategoryRepository.GetAssetCategoryById(categoryID)
	if err != nil {
		return nil, errors.New("asset category not found")
	}

	var assetCategoryResponse = out.AssetCategoryResponse{
		AssetCategoryID: assetCategory.AssetCategoryID,
		CategoryName:    assetCategory.CategoryName,
		Description:     assetCategory.Description,
	}

	return assetCategoryResponse, nil
}

func (s AssetCategoryService) DeleteAssetCategory(categoryID uint, clientID string) error {
	var user = &user.User{}
	var assetCategory = &asset.AssetCategory{}

	err := utils.GetDataFromRedis(utils.User, clientID, user)
	if err != nil {
		return err
	}

	assetCategory, err = s.AssetCategoryRepository.GetAssetCategoryById(categoryID)
	if err != nil {
		return errors.New("asset category not found")
	}

	assetCategory.DeletedBy = &user.FullName
	err = s.AssetCategoryRepository.DeleteAssetCategory(assetCategory)
	if err != nil {
		return err
	}
	return nil
}
