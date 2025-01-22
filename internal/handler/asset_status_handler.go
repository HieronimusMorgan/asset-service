package handler

import (
	"asset-service/internal/dto/in"
	"asset-service/internal/services"
	"asset-service/internal/utils"
	"asset-service/package/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetStatusHandler struct {
	AssetStatusService *services.AssetStatusService
}

func NewAssetStatusHandler(db *gorm.DB) *AssetStatusHandler {
	s := services.AddAssetStatus(db)
	return &AssetStatusHandler{AssetStatusService: s}
}

func (h AssetStatusHandler) AddAssetStatus(context *gin.Context) {
	var req in.AssetStatusRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, 400, "Error", nil, err.Error())
		return
	}

	token, err := utils.ExtractClaimsResponse(context)
	if err != nil {
		return
	}

	if err != nil {
		response.SendResponse(context, 401, "Error", nil, err.Error())
		return
	}
	assetStatus, err := h.AssetStatusService.AddAssetStatus(&req, token.ClientID)
	if err != nil {
		response.SendResponse(context, 500, "Failed to add asset status", nil, err)
		return
	}
	response.SendResponse(context, 201, "Asset status added successfully", assetStatus, nil)
}

func (h AssetStatusHandler) GetListAssetStatus(context *gin.Context) {
	assetStatus, err := h.AssetStatusService.GetAssetStatus()
	if err != nil {
		response.SendResponse(context, 500, "Failed to get asset status", nil, err)
		return
	}
	response.SendResponse(context, 200, "Success", assetStatus, nil)
}

func (h AssetStatusHandler) GetAssetStatusByID(context *gin.Context) {
	assetStatusID, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, 400, "Resource ID must be a number", nil, err)
		return
	}

	assetStatus, err := h.AssetStatusService.GetAssetStatusByID(assetStatusID)
	if err != nil {
		response.SendResponse(context, 500, "Failed to get asset status", nil, err)
		return
	}
	response.SendResponse(context, 200, "Success", assetStatus, nil)
}

func (h AssetStatusHandler) UpdateAssetStatus(context *gin.Context) {
	var req in.AssetStatusRequest
	assetStatusID, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, 400, "Resource ID must be a number", nil, err)
		return
	}

	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, 400, "Invalid request", nil, err)
		return
	}

	token, err := utils.ExtractClaimsResponse(context)
	if err != nil {
		return
	}

	assetStatus, err := h.AssetStatusService.UpdateAssetStatus(assetStatusID, &req, token.ClientID)
	if err != nil {
		response.SendResponse(context, 500, "Failed to update asset status", nil, err)
		return
	}
	response.SendResponse(context, 200, "Asset status updated successfully", assetStatus, nil)
}

func (h AssetStatusHandler) DeleteAssetStatus(context *gin.Context) {
	assetStatusID, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, 400, "Resource ID must be a number", nil, err)
		return
	}

	token, err := utils.ExtractClaimsResponse(context)
	if err != nil {
		return
	}

	err = h.AssetStatusService.DeleteAssetStatus(assetStatusID, token.ClientID)
	if err != nil {
		response.SendResponse(context, 500, "Failed to delete asset status", nil, err)
		return
	}
	response.SendResponse(context, 200, "Asset status deleted successfully", nil, nil)
}
