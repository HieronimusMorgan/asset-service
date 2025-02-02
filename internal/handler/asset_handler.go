package handler

import (
	"asset-service/internal/dto/in"
	"asset-service/internal/services"
	"asset-service/internal/utils"
	"asset-service/package/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

type AssetHandler struct {
	AssetService *services.AssetService
}

func NewAssetHandler(db *gorm.DB) *AssetHandler {
	s := services.NewAssetService(db)
	return &AssetHandler{AssetService: s}
}

func (h AssetHandler) AddAsset(context *gin.Context) {
	var req *in.AssetRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, 400, "Error", nil, err.Error())
		return
	}
	token, err := utils.ExtractClaimsResponse(context)
	if err != nil {
		return
	}

	asset, err := h.AssetService.AddAsset(req, token.ClientID)
	if err != nil {
		response.SendResponse(context, 500, "Failed to add assets", nil, err.Error())
		return
	}
	response.SendResponse(context, 201, "Asset added successfully", asset, nil)
}

func (h AssetHandler) UpdateAsset(context *gin.Context) {
	var req *in.AssetRequest
	assetIDStr := context.Param("id")
	assetID, err := strconv.ParseUint(assetIDStr, 10, 32)
	if err != nil {
		response.SendResponse(context, 400, "Invalid asset ID", nil, err.Error())
		return
	}
	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, 400, "Error", nil, err.Error())
		return
	}
	token, err := utils.ExtractClaimsResponse(context)
	if err != nil {
		return
	}

	asset, err := h.AssetService.UpdateAsset(uint(assetID), req, token.ClientID)
	if err != nil {
		response.SendResponse(context, 500, "Failed to update assets", nil, err.Error())
		return
	}
	response.SendResponse(context, 201, "Asset update successfully", asset, nil)

}

func (h AssetHandler) UpdateAssetStatus(context *gin.Context) {
	assetIDStr := context.Param("id")
	assetID, err := strconv.ParseUint(assetIDStr, 10, 32)
	if err != nil {
		response.SendResponse(context, 400, "Invalid asset ID", nil, err.Error())
		return
	}
	var req struct {
		StatusID uint `json:"status_id" binding:"required"`
	}
	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, 400, "Invalid request", nil, err.Error())
		return
	}
	token, err := utils.ExtractClaimsResponse(context)
	if err != nil {
		return
	}

	err = h.AssetService.UpdateAssetStatus(uint(assetID), req.StatusID, token.ClientID)
	if err != nil {
		response.SendResponse(context, 500, "Failed to update asset status", nil, err.Error())
		return
	}

	response.SendResponse(context, 200, "Asset status updated successfully", nil, nil)
}

func (h AssetHandler) UpdateAssetCategory(context *gin.Context) {
	var req struct {
		CategoryID uint `json:"category_id" binding:"required"`
	}
	assetIDStr := context.Param("id")
	assetID, err := strconv.ParseUint(assetIDStr, 10, 32)
	if err != nil {
		response.SendResponse(context, 400, "Invalid asset category ID", nil, err.Error())
		return
	}
	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, 400, "Invalid request", nil, err.Error())
		return
	}
	token, err := utils.ExtractClaimsResponse(context)
	if err != nil {
		return
	}

	err = h.AssetService.UpdateAssetCategory(uint(assetID), req.CategoryID, token.ClientID)
	if err != nil {
		response.SendResponse(context, 500, "Failed to update asset category", nil, err.Error())
		return
	}

	response.SendResponse(context, 200, "Asset category updated successfully", nil, nil)
}

func (h AssetHandler) GetListAsset(context *gin.Context) {

	token, err := utils.ExtractClaimsResponse(context)
	if err != nil {
		return
	}

	asset, err := h.AssetService.GetListAsset(token.ClientID)
	if err != nil {
		response.SendResponse(context, 500, "Failed to get list assets", nil, err.Error())
		return
	}
	response.SendResponse(context, 200, "Get list assets successfully", asset, nil)
}

func (h AssetHandler) GetAssetById(context *gin.Context) {

	assetID, err := utils.ConvertToUint(context.Param("id"))
	token, err := utils.ExtractClaimsResponse(context)
	if err != nil {
		return
	}

	asset, err := h.AssetService.GetAssetByID(token.ClientID, assetID)
	if err != nil {
		response.SendResponse(context, 500, "Failed to get list assets", nil, err.Error())
		return
	}
	response.SendResponse(context, 200, "Get list assets successfully", asset, nil)

}

func (h AssetHandler) DeleteAsset(context *gin.Context) {

}
