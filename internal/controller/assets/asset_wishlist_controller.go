package assets

import (
	request "asset-service/internal/dto/in/assets"
	responses "asset-service/internal/dto/out/assets"
	"asset-service/internal/services/assets"
	"asset-service/internal/utils"
	"asset-service/internal/utils/jwt"
	"asset-service/internal/utils/text"
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
	JWTService           jwt.Service
	IpCDN                string
}

func NewAssetWishlistController(AssetWishlistService assets.AssetWishlistService, JWTService jwt.Service, IpCDN string) AssetWishlistController {
	return assetWishlistController{AssetWishlistService: AssetWishlistService, JWTService: JWTService, IpCDN: IpCDN}
}

func (h assetWishlistController) AddWishlistAsset(c *gin.Context) {
	var req *request.AssetWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendResponse(c, 400, "Error", nil, err.Error())
		return
	}

	token, exist := jwt.ExtractTokenClaims(c)
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
	token, exist := jwt.ExtractTokenClaims(c)
	if !exist {
		response.SendResponseList(c, 400, "Error", response.PagedData{
			Total:     0,
			PageSize:  0,
			PageIndex: 0,
			Items:     nil,
		}, "Token not found")
		return
	}

	pageIndex, pageSize, err := utils.GetPageIndexPageSize(c)
	if err != nil {
		response.SendResponse(c, 400, "Invalid page index or page size", nil, err.Error())
		return
	}

	wishlist, total, err := h.AssetWishlistService.GetListAssetWishlist(token.ClientID, pageSize, pageIndex)
	if err != nil {
		response.SendResponseList(c, 500, "Invalid page_index", response.PagedData{
			Total:     total,
			PageSize:  pageSize,
			PageIndex: pageIndex,
			Items:     nil,
		}, "Failed to get assets wishlist")
		return
	}

	response.SendResponseList(c, 200, "Success", response.PagedData{
		Total:     total,
		PageSize:  pageSize,
		PageIndex: pageIndex,
		Items:     wishlist,
	}, nil)
}

func (h assetWishlistController) GetAssetWishlistByID(c *gin.Context) {
	id, err := utils.ConvertToUint(c.Param("id"))
	if err != nil {
		response.SendResponse(c, 400, "Error", nil, "Invalid asset ID")
		return
	}

	token, exist := jwt.ExtractTokenClaims(c)
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

	token, exist := jwt.ExtractTokenClaims(c)
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

	token, exist := jwt.ExtractTokenClaims(c)
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

func (h assetWishlistController) AddAssetWishlistToAsset(context *gin.Context) {
	id, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, 400, "Error", nil, "Invalid asset ID")
		return
	}

	err = context.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		response.SendResponse(context, 400, "Error parsing form", nil, err.Error())
		return
	}

	req := request.AssetRequest{
		SerialNumber:   text.GetOptionalString(context, "serial_number"),
		Name:           context.PostForm("name"),
		Description:    text.GetOptionalString(context, "description"),
		Barcode:        text.GetOptionalString(context, "barcode"),
		CategoryID:     utils.ParseFormUint(context, "category_id"),
		StatusID:       utils.ParseFormUint(context, "status_id"),
		PurchaseDate:   text.GetOptionalString(context, "purchase_date"),
		ExpiryDate:     text.GetOptionalString(context, "expiry_date"),
		WarrantyExpiry: text.GetOptionalString(context, "warranty_expiry_date"),
		Price:          utils.ParseFormFloat(context, "price"),
		Stock:          utils.ParseFormInt(context, "stock"),
		Notes:          text.GetOptionalString(context, "notes"),
	}

	// Extract token
	token, err := h.JWTService.ExtractClaims(context.GetHeader(utils.Authorization))
	if err != nil {
		response.SendResponse(context, 401, "Unauthorized", nil, err.Error())
		return
	}

	requestHeaderID := context.GetHeader("X-REQUEST-ID")

	// Extract files
	files := context.Request.MultipartForm.File["images"]
	if len(files) == 0 {
		response.SendResponse(context, 400, "No images uploaded", nil, "At least one image is required")
		return
	}

	var imageMetadata []responses.AssetImageResponse
	if len(files) != 0 {
		imageMetadata, err = uploadImagesToCDN(h.IpCDN, files, token.ClientID, context.GetHeader(utils.Authorization))
	}

	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	result, err := h.AssetWishlistService.AddAssetWishlistToAsset(id, &req, imageMetadata, token.ClientID, requestHeaderID)
	if err != nil {
		response.SendResponse(context, 500, "Failed to add asset wishlist to asset", nil, err.Error())
		return
	}
	response.SendResponse(context, 200, "Asset wishlist added to asset successfully", result, nil)
}
