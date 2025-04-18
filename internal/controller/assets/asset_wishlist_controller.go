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

type AssetWishlistController interface {
	AddWishlistAsset(context *gin.Context)
	GetListWishlistAsset(context *gin.Context)
	GetWishlistAssetByID(context *gin.Context)
	UpdateWishlistAsset(context *gin.Context)
	DeleteWishlistAsset(context *gin.Context)
}

type assetWishlistController struct {
	AssetWishlistService assets.AssetWishlistService
	JWTService           utils.JWTService
}

func NewAssetWishlistController(AssetWishlistService assets.AssetWishlistService, JWTService utils.JWTService) AssetWishlistController {
	return assetWishlistController{AssetWishlistService: AssetWishlistService, JWTService: JWTService}
}

func (h assetWishlistController) AddWishlistAsset(c *gin.Context) {
	var req *request.AssetWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendResponse(c, 400, "Error", nil, err.Error())
		return
	}

	token, exist := utils.ExtractTokenClaims(c)
	if !exist {
		response.SendResponse(c, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	asset, err := h.AssetWishlistService.AddAssetWishlist(req, token.ClientID)
	if err != nil {
		response.SendResponse(c, 500, "Failed to add assets", nil, err.Error())
		return
	}
	response.SendResponse(c, 201, "Asset added successfully", asset, nil)
}

func (h assetWishlistController) GetListWishlistAsset(c *gin.Context) {

	token, exist := utils.ExtractTokenClaims(c)
	if !exist {
		response.SendResponse(c, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	wishlist, err := h.AssetWishlistService.GetAssetWishlist(token.ClientID)
	if err != nil {
		response.SendResponse(c, 500, "Failed to get assets wishlist", nil, err.Error())
		return
	}
	response.SendResponse(c, 200, "Success", wishlist, nil)
}

func (h assetWishlistController) GetWishlistAssetByID(c *gin.Context) {

	token, exist := utils.ExtractTokenClaims(c)
	if !exist {
		response.SendResponse(c, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	assetIDStr := c.Param("id")
	assetID, err := strconv.ParseUint(assetIDStr, 10, 32)
	if err != nil {
		response.SendResponse(c, 400, "Invalid asset ID", nil, err.Error())
		return
	}

	asset, err := h.AssetWishlistService.GetAssetWishlistByID(uint(assetID), token.ClientID)
	if err != nil {
		response.SendResponse(c, 500, "Failed to get assets wishlist", nil, err.Error())
		return
	}
	response.SendResponse(c, 200, "Success", asset, nil)
}

func (h assetWishlistController) UpdateWishlistAsset(c *gin.Context) {
	var req request.UpdateAssetWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendResponse(c, 400, "Error", nil, err.Error())
		return
	}

	token, exist := utils.ExtractTokenClaims(c)
	if !exist {
		response.SendResponse(c, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	assetIDStr := c.Param("id")
	assetID, err := strconv.ParseUint(assetIDStr, 10, 32)
	if err != nil {
		response.SendResponse(c, 400, "Invalid asset ID", nil, err.Error())
		return
	}

	result, err := h.AssetWishlistService.UpdateAssetWishlist(uint(assetID), req, token.ClientID)
	if err != nil {
		response.SendResponse(c, 500, "Failed to update assets", nil, err.Error())
		return
	}
	response.SendResponse(c, 201, "Asset update successfully", result, nil)
}

func (h assetWishlistController) DeleteWishlistAsset(c *gin.Context) {
	token, exist := utils.ExtractTokenClaims(c)
	if !exist {
		response.SendResponse(c, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	assetIDStr := c.Param("id")
	assetID, err := strconv.ParseUint(assetIDStr, 10, 32)
	if err != nil {
		response.SendResponse(c, 400, "Invalid asset ID", nil, err.Error())
		return
	}

	result, err := h.AssetWishlistService.DeleteAssetWishlist(uint(assetID), token.ClientID)
	if err != nil {
		response.SendResponse(c, 500, "Failed to delete assets", nil, err.Error())
		return
	}
	response.SendResponse(c, 200, "Asset deleted successfully", result, nil)
}
