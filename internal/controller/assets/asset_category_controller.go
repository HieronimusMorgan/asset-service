package assets

import (
	request "asset-service/internal/dto/in/assets"
	"asset-service/internal/services/assets"
	"asset-service/internal/utils"
	"asset-service/package/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AssetCategoryController interface {
	AddAssetCategory(context *gin.Context)
	UpdateAssetCategory(context *gin.Context)
	GetAssetCategories(context *gin.Context)
	GetListAssetCategory(context *gin.Context)
	GetAssetCategoryById(context *gin.Context)
	DeleteAssetCategory(context *gin.Context)
}

type assetCategoryController struct {
	AssetCategoryService assets.AssetCategoryService
	JWTService           utils.JWTService
}

func NewAssetCategoryController(assetCategoryService assets.AssetCategoryService, jwtService utils.JWTService) AssetCategoryController {
	return assetCategoryController{AssetCategoryService: assetCategoryService, JWTService: jwtService}
}

func (h assetCategoryController) AddAssetCategory(context *gin.Context) {
	var req request.AssetCategoryRequest
	token, exist := utils.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
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

func (h assetCategoryController) UpdateAssetCategory(context *gin.Context) {
	var req request.AssetCategoryRequest
	assetCategoryID, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Resource MaintenanceTypeID must be a number", nil, err.Error())
		return
	}

	token, exist := utils.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
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

func (h assetCategoryController) GetAssetCategories(*gin.Context) {

}

func (h assetCategoryController) GetListAssetCategory(context *gin.Context) {
	token, exist := utils.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	assetCategories, err := h.AssetCategoryService.GetListAssetCategory(token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}
	response.SendResponse(context, http.StatusOK, "Asset categories retrieved successfully", assetCategories, nil)
}

func (h assetCategoryController) GetAssetCategoryById(context *gin.Context) {
	assetCategoryID, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Resource MaintenanceTypeID must be a number", nil, err.Error())
		return
	}

	token, exist := utils.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	assetCategory, err := h.AssetCategoryService.GetAssetCategoryById(assetCategoryID, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}
	response.SendResponse(context, http.StatusOK, "Asset category retrieved successfully", assetCategory, nil)
}

func (h assetCategoryController) DeleteAssetCategory(context *gin.Context) {
	assetCategoryID, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Resource MaintenanceTypeID must be a number", nil, err.Error())
		return
	}

	token, err := h.JWTService.ExtractClaims(context.GetHeader(utils.Authorization))
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
