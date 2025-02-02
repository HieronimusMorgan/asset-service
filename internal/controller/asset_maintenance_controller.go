package controller

import (
	"asset-service/internal/dto/in"
	"asset-service/internal/services"
	"asset-service/internal/utils"
	"asset-service/package/response"
	"gorm.io/gorm"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AssetMaintenanceController struct {
	Service *services.AssetMaintenanceService
}

func NewAssetMaintenanceController(db *gorm.DB) *AssetMaintenanceController {
	return &AssetMaintenanceController{Service: services.NewAssetMaintenanceService(db)}
}

func (c *AssetMaintenanceController) CreateMaintenance(ctx *gin.Context) {
	var maintenance in.AssetMaintenanceRequest
	if err := ctx.ShouldBindJSON(&maintenance); err != nil {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}

	token, err := utils.ExtractClaimsResponse(ctx)
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

func (c *AssetMaintenanceController) GetMaintenanceByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.SendResponse(ctx, http.StatusBadRequest, "Invalid maintenance ID", nil, err.Error())
		return
	}

	token, err := utils.ExtractClaimsResponse(ctx)
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

func (c *AssetMaintenanceController) GetMaintenancesByAssetID(ctx *gin.Context) {
	assetID, err := strconv.ParseUint(ctx.Param("asset_id"), 10, 32)
	if err != nil {
		response.SendResponse(ctx, http.StatusBadRequest, "Invalid asset ID", nil, err.Error())
		return
	}

	token, err := utils.ExtractClaimsResponse(ctx)
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

func (c *AssetMaintenanceController) UpdateMaintenance(ctx *gin.Context) {
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

func (c *AssetMaintenanceController) DeleteMaintenance(ctx *gin.Context) {
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
