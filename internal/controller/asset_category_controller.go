package controller

import (
	"asset-service/internal/dto/in"
	"asset-service/internal/services"
	"asset-service/internal/utils"
	"asset-service/package/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type AssetCategoryController struct {
	AssetCategoryService *services.AssetCategoryService
}

func NewAssetCategoryController(db *gorm.DB) *AssetCategoryController {
	s := services.NewAssetCategoryService(db)
	return &AssetCategoryController{AssetCategoryService: s}
}

func (h AssetCategoryController) AddAssetCategory(context *gin.Context) {
	var req in.AssetCategoryRequest
	token, err := utils.ExtractClaimsResponse(context)
	if err != nil {
		return
	}

	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}

	assetCategory, err := h.AssetCategoryService.AddAssetCategory(&req, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}
	response.SendResponse(context, http.StatusOK, "Asset category added successfully", assetCategory, nil)
}

func (h AssetCategoryController) UpdateAssetCategory(context *gin.Context) {
	var req in.AssetCategoryRequest
	assetCategoryID, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Resource ID must be a number", nil, err.Error())
		return
	}

	token, err := utils.ExtractClaimsResponse(context)
	if err != nil {
		return
	}

	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Invalid request", nil, err.Error())
		return
	}

	assetCategory, err := h.AssetCategoryService.UpdateAssetCategory(assetCategoryID, &req, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}
	response.SendResponse(context, http.StatusOK, "Asset category updated successfully", assetCategory, nil)

}

func (h AssetCategoryController) GetAssetCategories(context *gin.Context) {

}

func (h AssetCategoryController) GetListAssetCategory(context *gin.Context) {
	token, err := utils.ExtractClaimsResponse(context)
	if err != nil {
		return
	}

	assetCategories, err := h.AssetCategoryService.GetListAssetCategory(token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}
	response.SendResponse(context, http.StatusOK, "Asset categories retrieved successfully", assetCategories, nil)
}

func (h AssetCategoryController) GetAssetCategoryById(context *gin.Context) {
	assetCategoryID, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Resource ID must be a number", nil, err.Error())
		return
	}

	assetCategory, err := h.AssetCategoryService.GetAssetCategoryById(assetCategoryID)
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}
	response.SendResponse(context, http.StatusOK, "Asset category retrieved successfully", assetCategory, nil)
}

func (h AssetCategoryController) DeleteAssetCategory(context *gin.Context) {
	assetCategoryID, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Resource ID must be a number", nil, err.Error())
		return
	}

	token, err := utils.ExtractClaimsResponse(context)
	if err != nil {
		return
	}

	err = h.AssetCategoryService.DeleteAssetCategory(assetCategoryID, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}
	response.SendResponse(context, http.StatusOK, "Asset category deleted successfully", nil, nil)
}
