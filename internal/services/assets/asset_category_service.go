package assets

import (
	assets3 "asset-service/internal/dto/in/assets"
	assets4 "asset-service/internal/dto/out/assets"
	"asset-service/internal/models/assets"
	"asset-service/internal/models/user"
	assets2 "asset-service/internal/repository/assets"
	"asset-service/internal/utils"
	"errors"
	"gorm.io/gorm"
)

type AssetCategoryService struct {
	AssetCategoryRepository *assets2.AssetCategoryRepository
}

func NewAssetCategoryService(db *gorm.DB) *AssetCategoryService {
	r := assets2.NewAssetCategoryRepository(db)
	return &AssetCategoryService{AssetCategoryRepository: r}
}

func (s AssetCategoryService) AddAssetCategory(assetRequest *assets3.AssetCategoryRequest, clientID string) (interface{}, error) {
	var user = &user.User{}

	err := utils.GetDataFromRedis(utils.User, clientID, user)
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
	var assetCategoryResponse = assets4.AssetCategoryResponse{
		AssetCategoryID: assetCategory.AssetCategoryID,
		CategoryName:    assetCategory.CategoryName,
		Description:     assetCategory.Description,
	}
	return assetCategoryResponse, nil
}

func (s AssetCategoryService) UpdateAssetCategory(assetCategoryID uint, assetCategoryRequest *assets3.AssetCategoryRequest, clientID string) (interface{}, error) {
	var user = &user.User{}
	var assetCategory = &assets.AssetCategory{}

	err := utils.GetDataFromRedis(utils.User, clientID, user)
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
	var assetCategoriesResponse []assets4.AssetCategoryResponse

	for _, assetCategory := range assetCategories {
		assetCategoriesResponse = append(assetCategoriesResponse, assets4.AssetCategoryResponse{
			AssetCategoryID: assetCategory.AssetCategoryID,
			CategoryName:    assetCategory.CategoryName,
			Description:     assetCategory.Description,
		})
	}

	return assetCategoriesResponse, nil
}

func (s AssetCategoryService) GetAssetCategoryById(categoryID uint) (interface{}, error) {
	var assetCategory *assets.AssetCategory
	assetCategory, err := s.AssetCategoryRepository.GetAssetCategoryById(categoryID)
	if err != nil {
		return nil, errors.New("assets category not found")
	}

	var assetCategoryResponse = assets4.AssetCategoryResponse{
		AssetCategoryID: assetCategory.AssetCategoryID,
		CategoryName:    assetCategory.CategoryName,
		Description:     assetCategory.Description,
	}

	return assetCategoryResponse, nil
}

func (s AssetCategoryService) DeleteAssetCategory(categoryID uint, clientID string) error {
	var user = &user.User{}
	var assetCategory = &assets.AssetCategory{}

	err := utils.GetDataFromRedis(utils.User, clientID, user)
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
