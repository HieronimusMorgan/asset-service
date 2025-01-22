package handler

import (
	"asset-service/internal/dto/in"
	"asset-service/internal/services"
	"asset-service/internal/utils"
	"asset-service/package/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		response.SendResponse(context, 500, "Failed to add asset", nil, err.Error())
		return
	}
	response.SendResponse(context, 201, "Asset added successfully", asset, nil)
}

func (h AssetHandler) GetListAsset(context *gin.Context) {

	token, err := utils.ExtractClaimsResponse(context)
	if err != nil {
		return
	}

	asset, err := h.AssetService.GetListAsset(token.ClientID)
	if err != nil {
		response.SendResponse(context, 500, "Failed to get list asset", nil, err.Error())
		return
	}
	response.SendResponse(context, 200, "Get list asset successfully", asset, nil)
}

func (h AssetHandler) UpdateAsset(context *gin.Context) {

}

func (h AssetHandler) GetAssetById(context *gin.Context) {

}

func (h AssetHandler) DeleteAsset(context *gin.Context) {

}
