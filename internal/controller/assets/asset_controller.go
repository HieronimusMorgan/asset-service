package assets

import (
	request "asset-service/internal/dto/in/assets"
	responses "asset-service/internal/dto/out/assets"
	"asset-service/internal/services/assets"
	"asset-service/internal/utils"
	"asset-service/package/response"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

type AssetController interface {
	AddAsset(context *gin.Context)
	UpdateAsset(context *gin.Context)
	UpdateAssetStatus(context *gin.Context)
	UpdateAssetCategory(context *gin.Context)
	AddStockAsset(context *gin.Context)
	ReduceStockAsset(context *gin.Context)
	GetListAsset(context *gin.Context)
	GetAssetById(context *gin.Context)
	DeleteAsset(context *gin.Context)
}

type assetController struct {
	AssetService assets.AssetService
	JWTService   utils.JWTService
	IpCDN        string
}

func NewAssetController(assetService assets.AssetService, jwtService utils.JWTService, IpCDN string) AssetController {
	return assetController{AssetService: assetService, JWTService: jwtService, IpCDN: IpCDN}
}

func (h assetController) AddAsset(context *gin.Context) {
	err := context.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		response.SendResponse(context, 400, "Error parsing form", nil, err.Error())
		return
	}

	req := request.AssetRequest{
		SerialNumber:   utils.GetOptionalString(context, "serial_number"),
		Name:           context.PostForm("name"),
		Description:    utils.GetOptionalString(context, "description"),
		Barcode:        utils.GetOptionalString(context, "barcode"),
		CategoryID:     utils.ParseFormUint(context, "category_id"),
		StatusID:       utils.ParseFormUint(context, "status_id"),
		PurchaseDate:   utils.GetOptionalString(context, "purchase_date"),
		ExpiryDate:     utils.GetOptionalString(context, "expiry_date"),
		WarrantyExpiry: utils.GetOptionalString(context, "warranty_expiry_date"),
		Price:          utils.ParseFormFloat(context, "price"),
		Stock:          utils.ParseFormInt(context, "stock"),
		Notes:          utils.GetOptionalString(context, "notes"),
	}

	// Extract token
	token, err := h.JWTService.ExtractClaims(context.GetHeader(utils.Authorization))
	if err != nil {
		response.SendResponse(context, 401, "Unauthorized", nil, err.Error())
		return
	}

	credentialKey := context.GetHeader("X-CREDENTIAL-KEY")
	if credentialKey == "" {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "CredentialKey not found")
		return
	}
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

	asset, err := h.AssetService.AddAsset(&req, imageMetadata, token.ClientID, credentialKey)
	if err != nil {
		response.SendResponse(context, 500, "Failed to add asset", nil, err.Error())
		return
	}

	// Return response with saved asset
	response.SendResponse(context, 201, "Asset added successfully", asset, nil)
}

func (h assetController) UpdateAsset(context *gin.Context) {
	var req request.UpdateAssetRequest
	assetIDStr := context.Param("id")
	assetID, err := strconv.ParseUint(assetIDStr, 10, 32)
	if err != nil {
		response.SendResponse(context, 400, "Invalid asset MaintenanceTypeID", nil, err.Error())
		return
	}
	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, 400, "Error", nil, err.Error())
		return
	}
	token, err := h.JWTService.ExtractClaims(context.GetHeader(utils.Authorization))
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
		response.SendResponse(context, 400, "Invalid asset MaintenanceTypeID", nil, err.Error())
		return
	}
	var req struct {
		StatusID uint `json:"status_id" binding:"required"`
	}
	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, 400, "Invalid request", nil, err.Error())
		return
	}
	token, err := h.JWTService.ExtractClaims(context.GetHeader(utils.Authorization))
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
		response.SendResponse(context, 400, "Invalid asset category MaintenanceTypeID", nil, err.Error())
		return
	}
	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, 400, "Invalid request", nil, err.Error())
		return
	}
	token, err := h.JWTService.ExtractClaims(context.GetHeader(utils.Authorization))
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

func (h assetController) AddStockAsset(context *gin.Context) {
	var req struct {
		Stock  int     `json:"stock" binding:"required"`
		Reason *string `json:"reason"`
	}
	assetIDStr := context.Param("id")
	assetID, err := strconv.ParseUint(assetIDStr, 10, 32)
	if err != nil {
		response.SendResponse(context, 400, "Invalid asset MaintenanceTypeID", nil, err.Error())
		return
	}
	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, 400, "Invalid request", nil, err.Error())
		return
	}
	token, err := h.JWTService.ExtractClaims(context.GetHeader(utils.Authorization))
	if err != nil {
		return
	}

	data, err := h.AssetService.UpdateStockAsset(true, uint(assetID), req, token.ClientID)
	if err != nil {
		response.SendResponse(context, 500, "Failed to update stock asset", nil, err.Error())
		return
	}

	response.SendResponse(context, 200, "Stock asset updated successfully", data, nil)
}

func (h assetController) ReduceStockAsset(context *gin.Context) {
	var req struct {
		Stock  int     `json:"stock" binding:"required"`
		Reason *string `json:"reason"`
	}
	assetIDStr := context.Param("id")
	assetID, err := strconv.ParseUint(assetIDStr, 10, 32)
	if err != nil {
		response.SendResponse(context, 400, "Invalid asset MaintenanceTypeID", nil, err.Error())
		return
	}
	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, 400, "Invalid request", nil, err.Error())
		return
	}
	token, err := h.JWTService.ExtractClaims(context.GetHeader(utils.Authorization))
	if err != nil {
		return
	}

	data, err := h.AssetService.UpdateStockAsset(false, uint(assetID), req, token.ClientID)
	if err != nil {
		response.SendResponse(context, 500, "Failed to update stock asset", nil, err.Error())
		return
	}

	response.SendResponse(context, 200, "Stock asset updated successfully", data, nil)
}

func (h assetController) GetListAsset(context *gin.Context) {

	token, err := h.JWTService.ExtractClaims(context.GetHeader(utils.Authorization))
	if err != nil {
		return
	}

	// Get pagination parameters
	pageSize, err := strconv.Atoi(context.DefaultQuery("page_size", "10"))
	if err != nil || pageSize <= 0 {
		response.SendResponse(context, 400, "Invalid page_size", nil, "page_size must be a positive integer")
		return
	}

	pageIndex, err := strconv.Atoi(context.DefaultQuery("page_index", "1"))
	if err != nil || pageIndex <= 0 {
		response.SendResponse(context, 400, "Invalid page_index", nil, "page_index must be a positive integer")
		return
	}

	asset, total, err := h.AssetService.GetListAsset(token.ClientID, pageIndex, pageSize)
	if err != nil {
		response.SendResponseList(context, 500, "Failed to get list assets", response.PagedData{
			Total:     total,
			PageIndex: pageIndex,
			PageSize:  pageSize,
			Items:     nil,
		}, err.Error())
		return
	}
	response.SendResponseList(context, 200, "Get list assets successfully", response.PagedData{
		Total:     total,
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Items:     asset,
	}, nil)
}

func (h assetController) GetAssetById(context *gin.Context) {

	assetID, err := utils.ConvertToUint(context.Param("id"))
	token, err := h.JWTService.ExtractClaims(context.GetHeader(utils.Authorization))
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
	token, err := h.JWTService.ExtractClaims(context.GetHeader(utils.Authorization))
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

func uploadImagesToCDN(ipCdn string, files []*multipart.FileHeader, clientID, authToken string) ([]responses.AssetImageResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add clientID field
	if err := writer.WriteField("client_id", clientID); err != nil {
		return nil, err
	}

	// Add files to form data
	for _, file := range files {
		part, err := writer.CreateFormFile("images", file.Filename)
		if err != nil {
			return nil, err
		}

		src, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()

		if _, err = io.Copy(part, src); err != nil {
			return nil, err
		}
	}

	// Close writer
	if err := writer.Close(); err != nil {
		return nil, err
	}

	// Create and send request
	req, err := http.NewRequest("POST", ipCdn+"/v1/upload", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", authToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to upload images: %s", resp.Status)
	}

	var res struct {
		Data []responses.AssetImageResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	for i := range res.Data {
		res.Data[i].ImageURL = ipCdn + "/v1" + res.Data[i].ImageURL
	}

	return res.Data, nil
}
