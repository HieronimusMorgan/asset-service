package assets

import (
	service "asset-service/internal/services/assets"
	"asset-service/internal/utils"
	"asset-service/internal/utils/jwt"
	"asset-service/package/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AssetMaintenanceRecordController interface {
	GetMaintenanceRecord(ctx *gin.Context)
	GetListMaintenancesRecordByAssetID(ctx *gin.Context)
	GetMaintenancesRecordByAssetIDAndMaintenanceID(ctx *gin.Context)
	GetMaintenanceRecordByRecordIDAndAssetIDAndMaintenanceID(ctx *gin.Context)
	GetMaintenanceRecordByID(ctx *gin.Context)
}

type assetMaintenanceRecordController struct {
	Service    service.AssetMaintenanceRecordService
	JWTService jwt.Service
}

func NewAssetMaintenanceRecordController(Service service.AssetMaintenanceRecordService, JWTService jwt.Service) AssetMaintenanceRecordController {
	return assetMaintenanceRecordController{Service: Service, JWTService: JWTService}
}

func (h assetMaintenanceRecordController) GetMaintenanceRecord(ctx *gin.Context) {
	assetIDStr := ctx.Query("asset_id")
	maintenanceIDStr := ctx.Query("maintenance_id")
	recordIDStr := ctx.Query("maintenance_record_id")

	if assetIDStr != "" && maintenanceIDStr != "" && recordIDStr != "" {
		h.GetMaintenanceRecordByRecordIDAndAssetIDAndMaintenanceID(ctx)
		return
	}

	if recordIDStr != "" {
		h.GetMaintenanceRecordByID(ctx)
		return
	}

	if assetIDStr != "" && maintenanceIDStr != "" {
		h.GetMaintenancesRecordByAssetIDAndMaintenanceID(ctx)
		return
	}

	if assetIDStr != "" {
		h.GetListMaintenancesRecordByAssetID(ctx)
		return
	}

	ctx.JSON(http.StatusBadRequest, gin.H{
		"status":  http.StatusBadRequest,
		"message": "❌ Invalid query parameters. Please provide valid 'asset_id', 'maintenance_id', or 'maintenance_record_id'.",
	})
}

func (h assetMaintenanceRecordController) GetMaintenanceRecordByRecordIDAndAssetIDAndMaintenanceID(ctx *gin.Context) {
	maintenanceRecordID, err := utils.ConvertToUint(ctx.Query("maintenance_record_id"))
	if err != nil {
		response.SendResponse(ctx, 400, "Error", nil, "Invalid maintenance record ID")
		return
	}

	assetID, err := utils.ConvertToUint(ctx.Query("asset_id"))
	if err != nil {
		response.SendResponse(ctx, 400, "Error", nil, "Invalid asset ID")
		return
	}

	maintenanceID, err := utils.ConvertToUint(ctx.Query("maintenance_id"))
	if err != nil {
		response.SendResponse(ctx, 400, "Error", nil, "Invalid maintenance ID")
		return
	}

	token, exist := jwt.ExtractTokenClaims(ctx)
	if !exist {
		response.SendResponse(ctx, 400, "Error", nil, "Token not found")
		return
	}

	result, err := h.Service.GetMaintenanceRecordByRecordIDAndAssetIDAndMaintenanceID(maintenanceRecordID, assetID, maintenanceID, token.ClientID)
	if err != nil {
		response.SendResponse(ctx, 500, "Failed to get maintenance record", nil, err.Error())
		return
	}
	response.SendResponse(ctx, 200, "Success", result, nil)
}

func (h assetMaintenanceRecordController) GetListMaintenancesRecordByAssetID(ctx *gin.Context) {
	assetID, err := utils.ConvertToUint(ctx.Query("asset_id"))
	if err != nil {
		response.SendResponse(ctx, 400, "Error", nil, "Invalid asset ID")
		return
	}

	token, exist := jwt.ExtractTokenClaims(ctx)
	if !exist {
		response.SendResponse(ctx, 400, "Error", nil, "Token not found")
		return
	}

	result, err := h.Service.GetListMaintenancesRecordByAssetID(assetID, token.ClientID)
	if err != nil {
		response.SendResponse(ctx, 500, "Failed to get maintenance record", nil, err.Error())
		return
	}
	response.SendResponse(ctx, 200, "Success", result, nil)
}

func (h assetMaintenanceRecordController) GetMaintenancesRecordByAssetIDAndMaintenanceID(ctx *gin.Context) {
	assetID, err := utils.ConvertToUint(ctx.Query("asset_id"))
	if err != nil {
		response.SendResponse(ctx, 400, "Error", nil, "Invalid asset ID")
		return
	}

	maintenanceID, err := utils.ConvertToUint(ctx.Query("maintenance_id"))
	if err != nil {
		response.SendResponse(ctx, 400, "Error", nil, "Invalid maintenance ID")
		return
	}

	token, exist := jwt.ExtractTokenClaims(ctx)
	if !exist {
		response.SendResponse(ctx, 400, "Error", nil, "Token not found")
		return
	}

	result, err := h.Service.GetMaintenancesRecordByAssetIDAndMaintenanceID(assetID, maintenanceID, token.ClientID)
	if err != nil {
		response.SendResponse(ctx, 500, "Failed to get maintenance record", nil, err.Error())
		return
	}
	response.SendResponse(ctx, 200, "Success", result, nil)
}

func (h assetMaintenanceRecordController) GetMaintenanceRecordByID(ctx *gin.Context) {
	maintenanceRecordID, err := utils.ConvertToUint(ctx.Query("maintenance_record_id"))
	if err != nil {
		response.SendResponse(ctx, 400, "Error", nil, "Invalid maintenance record ID")
		return
	}

	token, exist := jwt.ExtractTokenClaims(ctx)
	if !exist {
		response.SendResponse(ctx, 400, "Error", nil, "Token not found")
		return
	}

	result, err := h.Service.GetMaintenanceRecordByID(maintenanceRecordID, token.ClientID)
	if err != nil {
		response.SendResponse(ctx, 500, "Failed to get maintenance record", nil, err.Error())
		return
	}
	response.SendResponse(ctx, 200, "Success", result, nil)
}
