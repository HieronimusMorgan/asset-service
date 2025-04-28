package assets

import (
	request "asset-service/internal/dto/in/assets"
	"asset-service/internal/services/assets"
	"asset-service/internal/utils"
	"asset-service/package/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

	credentialKey := context.GetHeader("X-CREDENTIAL-KEY")
	if credentialKey == "" {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "CredentialKey not found")
		return
	}

	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}

	assetCategory, err := h.AssetCategoryService.AddAssetCategory(&req, credentialKey, token.ClientID)
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
	pageSize, err := strconv.Atoi(context.DefaultQuery("page_size", "10"))
	if err != nil || pageSize <= 0 {
		response.SendResponse(context, 400, "Invalid page_size", nil, "page_size must be a positive integer")
		return
	}

	pageIndex, err := strconv.Atoi(context.DefaultQuery("page_index", "1"))
	if err != nil || pageIndex <= 0 {
		response.SendResponse(context, 400, "Invalid page_index", nil, "page_index must be a positive integer")
		return
	}

	token, exist := utils.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	assetCategories, total, err := h.AssetCategoryService.GetListAssetCategory(token.ClientID, pageSize, pageIndex)
	if err != nil {
		response.SendResponseList(context, 500, "Failed to get list assets", response.PagedData{
			Total:     total,
			PageIndex: pageIndex,
			PageSize:  pageSize,
			Items:     nil,
		}, err.Error())
		return
	}

	response.SendResponseList(context, 200, "List of asset categories", response.PagedData{
		Total:     total,
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Items:     assetCategories,
	}, nil)
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
