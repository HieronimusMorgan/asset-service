package assets

import (
	"asset-service/internal/services/assets"
	"asset-service/internal/utils"
	"asset-service/internal/utils/jwt"
	"asset-service/package/response"
	"github.com/gin-gonic/gin"

	request "asset-service/internal/dto/in/assets"
	"net/http"
)

type AssetGroupPermissionController interface {
	AddAssetGroupPermission(context *gin.Context)
	UpdateAssetGroupPermission(context *gin.Context)
	GetAssetGroupPermissionByID(context *gin.Context)
	GetListAssetGroupPermission(context *gin.Context)
	DeleteAssetGroupPermission(context *gin.Context)
}

type assetGroupPermissionController struct {
	AssetGroupPermissionService assets.AssetGroupPermissionService
	JWTService                  jwt.Service
}

func NewAssetGroupPermissionController(AssetGroupPermissionService assets.AssetGroupPermissionService, JWTService jwt.Service) AssetGroupPermissionController {
	return assetGroupPermissionController{AssetGroupPermissionService: AssetGroupPermissionService, JWTService: JWTService}
}

func (a assetGroupPermissionController) AddAssetGroupPermission(context *gin.Context) {
	var req request.AssetGroupPermissionRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	err := a.AssetGroupPermissionService.AddAssetGroupPermission(&req, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", err.Error(), err)
		return
	}

	response.SendResponse(context, http.StatusOK, "Success", nil, "Asset group permission created successfully")
}

func (a assetGroupPermissionController) UpdateAssetGroupPermission(context *gin.Context) {
	var req request.AssetGroupPermissionRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Resource MaintenanceTypeID must be a number", nil, err.Error())
		return
	}

	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	err = a.AssetGroupPermissionService.UpdateAssetGroupPermission(id, &req, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", err.Error(), err)
		return
	}

	response.SendResponse(context, http.StatusOK, "Success", nil, "Asset group updated successfully")
}

func (a assetGroupPermissionController) GetAssetGroupPermissionByID(context *gin.Context) {
	id, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Resource MaintenanceTypeID must be a number", nil, err.Error())
		return
	}

	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	data, err := a.AssetGroupPermissionService.GetAssetGroupPermissionById(id, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", err.Error(), err)
		return
	}

	response.SendResponse(context, http.StatusOK, "Success", data, "Asset group permission found")
}

func (a assetGroupPermissionController) GetListAssetGroupPermission(context *gin.Context) {
	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	data, err := a.AssetGroupPermissionService.GetListAssetGroupPermission(token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", err.Error(), err)
		return
	}

	response.SendResponse(context, http.StatusOK, "Success", data, nil)
}

func (a assetGroupPermissionController) DeleteAssetGroupPermission(context *gin.Context) {
	id, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Resource MaintenanceTypeID must be a number", nil, err.Error())
		return
	}

	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	err = a.AssetGroupPermissionService.DeleteAssetGroupPermission(id, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", err.Error(), err)
		return
	}

	response.SendResponse(context, http.StatusOK, "Success", nil, "Asset group permission deleted successfully")
}
