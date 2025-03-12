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
	"github.com/rs/zerolog/log"
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
	// Parse multipart form (for images + JSON fields)
	err := context.Request.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		response.SendResponse(context, 400, "Error parsing form", nil, err.Error())
		return
	}

	// Extract form data (JSON fields)
	req := request.AssetRequest{
		SerialNumber:   getOptionalString(context, "serial_number"),
		Name:           context.PostForm("name"),
		Description:    getOptionalString(context, "description"),
		Barcode:        getOptionalString(context, "barcode"),
		CategoryID:     parseFormInt(context, "category_id"),
		StatusID:       parseFormInt(context, "status_id"),
		PurchaseDate:   getOptionalString(context, "purchase_date"),
		ExpiryDate:     getOptionalString(context, "expiry_date"),
		WarrantyExpiry: getOptionalString(context, "warranty_expiry_date"),
		Price:          parseFormFloat(context, "price"),
		Stock:          parseFormInt(context, "stock"),
		Notes:          getOptionalString(context, "notes"),
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

	asset, err := h.AssetService.AddAsset(&req, imageMetadata, token.ClientID, requestHeaderID)
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
		response.SendResponse(context, 400, "Invalid asset ID", nil, err.Error())
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
		response.SendResponse(context, 400, "Invalid asset category ID", nil, err.Error())
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
		response.SendResponse(context, 400, "Invalid asset ID", nil, err.Error())
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
		response.SendResponse(context, 400, "Invalid asset ID", nil, err.Error())
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

	asset, err := h.AssetService.GetListAsset(token.ClientID)
	if err != nil {
		response.SendResponse(context, 500, "Failed to get list assets", nil, err.Error())
		return
	}
	response.SendResponse(context, 200, "Get list assets successfully", asset, nil)
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

func uploadImagesToCDN(ipCdn string, files []*multipart.FileHeader, clientID string, authToken string) ([]responses.AssetImageResponse, error) {
	var _ []responses.AssetImageResponse
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	log.Info().Msgf("Uploading %d images", len(files))
	// Add clientID field
	_ = writer.WriteField("client_id", clientID)

	// Add files to form data
	for _, file := range files {
		part, err := writer.CreateFormFile("images", file.Filename)
		if err != nil {
			return nil, err
		}

		// Open file
		src, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()

		// Copy file content to form
		_, err = io.Copy(part, src)
		if err != nil {
			return nil, err
		}
	}

	// Close writer
	writer.Close()

	// Send request to `cdn-service`
	log.Log().Msgf("ipCdn: %s", ipCdn)
	req, err := http.NewRequest("POST", ipCdn+"/v1/upload", body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create request")
		return nil, err
	}

	// Set headers
	req.Header.Set("Authorization", authToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upload images")
		return nil, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Failed to upload images: %s", resp.Status)
	}

	// Parse response JSON
	var res struct {
		Data []responses.AssetImageResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	var result struct {
		Data []responses.AssetImageResponse `json:"data"`
	}
	for _, img := range res.Data {
		result.Data = append(result.Data, responses.AssetImageResponse{
			ImageURL: ipCdn + "/v1" + img.ImageURL,
		})
	}

	return result.Data, nil
}

// Get an optional string field from the form
func getOptionalString(context *gin.Context, field string) *string {
	val := context.PostForm(field)
	if val == "" {
		return nil
	}
	return &val
}

// Parse an integer from the form data
func parseFormInt(context *gin.Context, field string) int {
	val := context.PostForm(field)
	if val == "" {
		return 0
	}
	intVal, _ := strconv.Atoi(val)
	return intVal
}

// Parse a float from the form data
func parseFormFloat(context *gin.Context, field string) float64 {
	val := context.PostForm(field)
	if val == "" {
		return 0.0
	}
	floatVal, _ := strconv.ParseFloat(val, 64)
	return floatVal
}
