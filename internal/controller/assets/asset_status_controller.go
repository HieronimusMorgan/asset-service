package assets

import (
	assets2 "asset-service/internal/dto/in/assets"
	"asset-service/internal/services/assets"
	"asset-service/internal/utils"
	"asset-service/package/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AssetStatusController interface {
	AddAssetStatus(context *gin.Context)
	GetListAssetStatus(context *gin.Context)
	GetAssetStatusByID(context *gin.Context)
	UpdateAssetStatus(context *gin.Context)
	DeleteAssetStatus(context *gin.Context)
}

type assetStatusController struct {
	AssetStatusService assets.AssetStatusService
	JWTService         utils.JWTService
}

func NewAssetStatusController(AssetStatusService assets.AssetStatusService, JWTService utils.JWTService) AssetStatusController {
	return assetStatusController{AssetStatusService: AssetStatusService, JWTService: JWTService}
}

func (h assetStatusController) AddAssetStatus(context *gin.Context) {
	var req assets2.AssetStatusRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, 400, "Error", nil, err.Error())
		return
	}

	token, exist := utils.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	assetStatus, err := h.AssetStatusService.AddAssetStatus(&req, token.ClientID)
	if err != nil {
		response.SendResponse(context, 500, "Failed to add assets status", nil, err)
		return
	}
	response.SendResponse(context, 201, "Asset status added successfully", assetStatus, nil)
}

func (h assetStatusController) GetListAssetStatus(context *gin.Context) {
	assetStatus, err := h.AssetStatusService.GetAssetStatus()
	if err != nil {
		response.SendResponse(context, 500, "Failed to get assets status", nil, err)
		return
	}
	response.SendResponse(context, 200, "Success", assetStatus, nil)
}

func (h assetStatusController) GetAssetStatusByID(context *gin.Context) {
	assetStatusID, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, 400, "Resource ID must be a number", nil, err)
		return
	}

	assetStatus, err := h.AssetStatusService.GetAssetStatusByID(assetStatusID)
	if err != nil {
		response.SendResponse(context, 500, "Failed to get assets status", nil, err)
		return
	}
	response.SendResponse(context, 200, "Success", assetStatus, nil)
}

func (h assetStatusController) UpdateAssetStatus(context *gin.Context) {
	var req assets2.AssetStatusRequest
	assetStatusID, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, 400, "Resource ID must be a number", nil, err)
		return
	}

	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, 400, "Invalid request", nil, err)
		return
	}

	token, exist := utils.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	assetStatus, err := h.AssetStatusService.UpdateAssetStatus(assetStatusID, &req, token.ClientID)
	if err != nil {
		response.SendResponse(context, 500, "Failed to update assets status", nil, err)
		return
	}
	response.SendResponse(context, 200, "Asset status updated successfully", assetStatus, nil)
}

func (h assetStatusController) DeleteAssetStatus(context *gin.Context) {
	assetStatusID, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, 400, "Resource ID must be a number", nil, err)
		return
	}

	token, exist := utils.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	err = h.AssetStatusService.DeleteAssetStatus(assetStatusID, token.ClientID)
	if err != nil {
		response.SendResponse(context, 500, "Failed to delete assets status", nil, err)
		return
	}
	response.SendResponse(context, 200, "Asset status deleted successfully", nil, nil)
}
