package assets

import (
	request "asset-service/internal/dto/in/assets"
	"asset-service/internal/services/assets"
	"asset-service/internal/utils"
	"asset-service/package/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AssetWishlistController interface {
	AddWishlistAsset(context *gin.Context)
	GetListAssetWishlist(context *gin.Context)
	GetAssetWishlistByID(context *gin.Context)
	UpdateAssetWishlist(context *gin.Context)
	DeleteAssetWishlist(context *gin.Context)
	AddAssetWishlistToAsset(context *gin.Context)
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

func (h assetWishlistController) GetListAssetWishlist(c *gin.Context) {
	token, exist := utils.ExtractTokenClaims(c)
	if !exist {
		response.SendResponse(c, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	asset, err := h.AssetWishlistService.GetListAssetWishlist(token.ClientID)
	if err != nil {
		response.SendResponse(c, 500, "Failed to get assets", nil, err.Error())
		return
	}
	response.SendResponse(c, 200, "Success", asset, nil)
}

func (h assetWishlistController) GetAssetWishlistByID(c *gin.Context) {
	id, err := utils.ConvertToUint(c.Param("id"))
	if err != nil {
		response.SendResponse(c, 400, "Error", nil, "Invalid asset ID")
		return
	}

	token, exist := utils.ExtractTokenClaims(c)
	if !exist {
		response.SendResponse(c, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	asset, err := h.AssetWishlistService.GetAssetWishlistByID(id, token.ClientID)
	if err != nil {
		response.SendResponse(c, 500, "Failed to get assets", nil, err.Error())
		return
	}
	response.SendResponse(c, 200, "Success", asset, nil)
}

func (h assetWishlistController) UpdateAssetWishlist(c *gin.Context) {
	id, err := utils.ConvertToUint(c.Param("id"))
	if err != nil {
		response.SendResponse(c, 400, "Error", nil, "Invalid asset ID")
		return
	}

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

	asset, err := h.AssetWishlistService.UpdateAssetWishlist(id, req, token.ClientID)
	if err != nil {
		response.SendResponse(c, 500, "Failed to update assets", nil, err.Error())
		return
	}
	response.SendResponse(c, 200, "Asset updated successfully", asset, nil)
}

func (h assetWishlistController) DeleteAssetWishlist(c *gin.Context) {
	id, err := utils.ConvertToUint(c.Param("id"))
	if err != nil {
		response.SendResponse(c, 400, "Error", nil, "Invalid asset ID")
		return
	}

	token, exist := utils.ExtractTokenClaims(c)
	if !exist {
		response.SendResponse(c, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	err = h.AssetWishlistService.DeleteAssetWishlist(id, token.ClientID)
	if err != nil {
		response.SendResponse(c, 500, "Failed to delete assets", nil, err.Error())
		return
	}
	response.SendResponse(c, 200, "Asset wishlist deleted successfully", nil, nil)
}

func (h assetWishlistController) AddAssetWishlistToAsset(c *gin.Context) {
	id, err := utils.ConvertToUint(c.Param("id"))
	if err != nil {
		response.SendResponse(c, 400, "Error", nil, "Invalid asset ID")
		return
	}

	token, exist := utils.ExtractTokenClaims(c)
	if !exist {
		response.SendResponse(c, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	result, err := h.AssetWishlistService.AddAssetWishlistToAsset(id, token.ClientID)
	if err != nil {
		response.SendResponse(c, 500, "Failed to add asset wishlist to asset", nil, err.Error())
		return
	}
	response.SendResponse(c, 200, "Asset wishlist added to asset successfully", result, nil)
}
