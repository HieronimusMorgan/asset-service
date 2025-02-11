package assets

import (
	assets2 "asset-service/internal/dto/in/assets"
	"asset-service/internal/services/assets"
	"asset-service/internal/utils"
	"asset-service/package/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AssetMaintenanceController interface {
	CreateMaintenance(ctx *gin.Context)
	GetMaintenanceByID(ctx *gin.Context)
	GetMaintenancesByAssetID(ctx *gin.Context)
	UpdateMaintenance(ctx *gin.Context)
	DeleteMaintenance(ctx *gin.Context)
}

type assetMaintenanceController struct {
	Service    assets.AssetMaintenanceService
	JWTService utils.JWTService
}

func NewAssetMaintenanceController(Service assets.AssetMaintenanceService, JWTService utils.JWTService) AssetMaintenanceController {
	return assetMaintenanceController{Service: Service, JWTService: JWTService}
}

func (c assetMaintenanceController) CreateMaintenance(ctx *gin.Context) {
	var maintenance assets2.AssetMaintenanceRequest
	if err := ctx.ShouldBindJSON(&maintenance); err != nil {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}

	token, err := c.JWTService.ExtractClaims(ctx.GetHeader("Authorization"))
	if err != nil {
		return
	}

	maintenances, err := c.Service.CreateMaintenance(maintenance, token.ClientID)
	if err != nil {
		response.SendResponse(ctx, http.StatusInternalServerError, "Failed to create maintenance record", nil, err.Error())
		return
	}

	response.SendResponse(ctx, http.StatusCreated, "Maintenance record created successfully", maintenances, nil)
}

func (c assetMaintenanceController) GetMaintenanceByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.SendResponse(ctx, http.StatusBadRequest, "Invalid maintenance ID", nil, err.Error())
		return
	}

	token, err := c.JWTService.ExtractClaims(ctx.GetHeader("Authorization"))
	if err != nil {
		return
	}

	maintenance, err := c.Service.GetMaintenanceByID(uint(id), token.ClientID)
	if err != nil {
		response.SendResponse(ctx, http.StatusNotFound, "Maintenance not found", nil, err.Error())
		return
	}

	if maintenance == nil {
		response.SendResponse(ctx, http.StatusNotFound, "Maintenance not found", nil, nil)
		return
	}

	response.SendResponse(ctx, http.StatusOK, "Success", maintenance, nil)
}

func (c assetMaintenanceController) GetMaintenancesByAssetID(ctx *gin.Context) {
	assetID, err := strconv.ParseUint(ctx.Param("asset_id"), 10, 32)
	if err != nil {
		response.SendResponse(ctx, http.StatusBadRequest, "Invalid asset ID", nil, err.Error())
		return
	}

	token, err := c.JWTService.ExtractClaims(ctx.GetHeader("Authorization"))
	if err != nil {
		return
	}

	maintenances, err := c.Service.GetMaintenancesByAssetID(uint(assetID), token.ClientID)
	if err != nil {
		response.SendResponse(ctx, http.StatusInternalServerError, "Failed to get maintenance records", nil, err.Error())
		return
	}
	response.SendResponse(ctx, http.StatusOK, "Success", maintenances, nil)
}

func (c assetMaintenanceController) UpdateMaintenance(ctx *gin.Context) {
	//var maintenance in.AssetMaintenanceRequest
	//if err := ctx.ShouldBindJSON(&maintenance); err != nil {
	//	response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, err.Error())
	//	return
	//}
	//if err := c.Service.UpdateMaintenance(&maintenance); err != nil {
	//	response.SendResponse(ctx, http.StatusInternalServerError, "Failed to update maintenance record", nil, err.Error())
	//	return
	//}
	//response.SendResponse(ctx, http.StatusOK, "Maintenance record updated successfully", maintenance, nil)
}

func (c assetMaintenanceController) DeleteMaintenance(ctx *gin.Context) {
	//id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	//if err != nil {
	//	response.SendResponse(ctx, http.StatusBadRequest, "Invalid maintenance ID", nil, err.Error())
	//	return
	//}
	//if err := c.Service.DeleteMaintenance(uint(id)); err != nil {
	//	response.SendResponse(ctx, http.StatusInternalServerError, "Failed to delete maintenance record", nil, err.Error())
	//	return
	//}
	//response.SendResponse(ctx, http.StatusNoContent, "Maintenance record deleted successfully", nil, nil)
}
