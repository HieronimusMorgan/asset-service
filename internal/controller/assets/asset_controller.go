package assets

import (
	assets2 "asset-service/internal/dto/in/assets"
	"asset-service/internal/services/assets"
	"asset-service/internal/utils"
	"asset-service/package/response"
	"github.com/gin-gonic/gin"
	"strconv"
)

type AssetController interface {
	AddAsset(context *gin.Context)
	UpdateAsset(context *gin.Context)
	UpdateAssetStatus(context *gin.Context)
	UpdateAssetCategory(context *gin.Context)
	GetListAsset(context *gin.Context)
	GetAssetById(context *gin.Context)
	DeleteAsset(context *gin.Context)
}

type assetController struct {
	AssetService assets.AssetService
	JWTService   utils.JWTService
}

func NewAssetController(assetService assets.AssetService, jwtService utils.JWTService) AssetController {
	return assetController{AssetService: assetService, JWTService: jwtService}
}

func (h assetController) AddAsset(context *gin.Context) {
	var req *assets2.AssetRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, 400, "Error", nil, err.Error())
		return
	}
	token, err := h.JWTService.ExtractClaims(context.GetHeader("Authorization"))
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

func (h assetController) UpdateAsset(context *gin.Context) {
	var req struct {
		Description  string  `json:"description"`
		PurchaseDate string  `json:"purchase_date" binding:"required"`
		ExpiryDate   string  `json:"expiry_date"`
		Price        float64 `json:"price" binding:"required"`
	}
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
	token, err := h.JWTService.ExtractClaims(context.GetHeader("Authorization"))
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

func (h assetController) UpdateAssetStatus(context *gin.Context) {
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
	token, err := h.JWTService.ExtractClaims(context.GetHeader("Authorization"))
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

func (h assetController) UpdateAssetCategory(context *gin.Context) {
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
	token, err := h.JWTService.ExtractClaims(context.GetHeader("Authorization"))
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

func (h assetController) GetListAsset(context *gin.Context) {

	token, err := h.JWTService.ExtractClaims(context.GetHeader("Authorization"))
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

func (h assetController) GetAssetById(context *gin.Context) {

	assetID, err := utils.ConvertToUint(context.Param("id"))
	token, err := h.JWTService.ExtractClaims(context.GetHeader("Authorization"))
	if err != nil {
		return
	}

	asset, err := h.AssetService.GetAssetByID(token.ClientID, assetID)
	if err != nil {
		response.SendResponse(context, 500, "Failed to get detail assets", nil, err.Error())
		return
	}
	response.SendResponse(context, 200, "Get detail assets successfully", asset, nil)

}

func (h assetController) DeleteAsset(context *gin.Context) {
	assetID, err := utils.ConvertToUint(context.Param("id"))
	token, err := h.JWTService.ExtractClaims(context.GetHeader("Authorization"))
	if err != nil {
		return
	}

	err = h.AssetService.DeleteAsset(assetID, token.ClientID)
	if err != nil {
		response.SendResponse(context, 500, "Failed to delete asset", nil, err.Error())
		return
	}

	response.SendResponse(context, 200, "Asset deleted successfully", nil, nil)
}
