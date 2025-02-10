package assets

import (
	request "asset-service/internal/dto/in/assets"
	"asset-service/internal/services/assets"
	"asset-service/internal/utils"
	"asset-service/package/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

type AssetWishlistController struct {
	AssetWishlistService *assets.AssetWishlistService
}

func NewAssetWishlistController(db *gorm.DB) *AssetWishlistController {
	s := assets.NewAssetWishlistService(db)
	return &AssetWishlistController{AssetWishlistService: s}
}

func (a AssetWishlistController) AddWishlistAsset(c *gin.Context) {
	var req *request.AssetWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendResponse(c, 400, "Error", nil, err.Error())
		return
	}
	token, err := utils.ExtractClaimsResponse(c)
	if err != nil {
		return
	}

	asset, err := a.AssetWishlistService.AddAssetWishlist(req, token.ClientID)
	if err != nil {
		response.SendResponse(c, 500, "Failed to add assets", nil, err.Error())
		return
	}
	response.SendResponse(c, 201, "Asset added successfully", asset, nil)
}

func (a AssetWishlistController) GetListWishlistAsset(c *gin.Context) {
	token, err := utils.ExtractClaimsResponse(c)
	if err != nil {
		return
	}

	wishlist, err := a.AssetWishlistService.GetAssetWishlist(token.ClientID)
	if err != nil {
		response.SendResponse(c, 500, "Failed to get assets wishlist", nil, err.Error())
		return
	}
	response.SendResponse(c, 200, "Success", wishlist, nil)
}

func (a AssetWishlistController) GetWishlistAssetByID(c *gin.Context) {
	token, err := utils.ExtractClaimsResponse(c)
	if err != nil {
		return
	}

	assetIDStr := c.Param("id")
	assetID, err := strconv.ParseUint(assetIDStr, 10, 32)
	if err != nil {
		response.SendResponse(c, 400, "Invalid asset ID", nil, err.Error())
		return
	}

	asset, err := a.AssetWishlistService.GetAssetWishlistByID(uint(assetID), token.ClientID)
	if err != nil {
		response.SendResponse(c, 500, "Failed to get assets wishlist", nil, err.Error())
		return
	}
	response.SendResponse(c, 200, "Success", asset, nil)
}

func (a AssetWishlistController) UpdateWishlistAsset(c *gin.Context) {
	var req struct {
		Description  string  `json:"description"`
		PurchaseDate string  `json:"purchase_date" binding:"required"`
		CategoryID   int     `json:"category_id" binding:"required"`
		StatusID     int     `json:"status_id" binding:"required"`
		Price        float64 `json:"price" binding:"required"`
		IsWishlist   bool    `json:"is_wishlist" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendResponse(c, 400, "Error", nil, err.Error())
		return
	}
	token, err := utils.ExtractClaimsResponse(c)
	if err != nil {
		return
	}

	assetIDStr := c.Param("id")
	assetID, err := strconv.ParseUint(assetIDStr, 10, 32)
	if err != nil {
		response.SendResponse(c, 400, "Invalid asset ID", nil, err.Error())
		return
	}

	result, err := a.AssetWishlistService.UpdateAssetWishlist(uint(assetID), req, token.ClientID)
	if err != nil {
		response.SendResponse(c, 500, "Failed to update assets", nil, err.Error())
		return
	}
	response.SendResponse(c, 201, "Asset update successfully", result, nil)
}

func (a AssetWishlistController) DeleteWishlistAsset(c *gin.Context) {
	token, err := utils.ExtractClaimsResponse(c)
	if err != nil {
		return
	}

	assetIDStr := c.Param("id")
	assetID, err := strconv.ParseUint(assetIDStr, 10, 32)
	if err != nil {
		response.SendResponse(c, 400, "Invalid asset ID", nil, err.Error())
		return
	}

	result, err := a.AssetWishlistService.DeleteAssetWishlist(uint(assetID), token.ClientID)
	if err != nil {
		response.SendResponse(c, 500, "Failed to delete assets", nil, err.Error())
		return
	}
	response.SendResponse(c, 200, "Asset deleted successfully", result, nil)
}
