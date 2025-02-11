package assets

import (
	"asset-service/internal/models/assets"
	service "asset-service/internal/services/assets"
	"asset-service/internal/utils"
	"asset-service/package/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AssetMaintenanceTypeController interface {
	CreateMaintenanceType(ctx *gin.Context)
}

type assetMaintenanceTypeController struct {
	Service    service.AssetMaintenanceTypeService
	JWTService utils.JWTService
}

func NewAssetMaintenanceTypeController(Service service.AssetMaintenanceTypeService, JWTService utils.JWTService) AssetMaintenanceTypeController {
	return assetMaintenanceTypeController{Service: Service, JWTService: JWTService}
}

func (c assetMaintenanceTypeController) CreateMaintenanceType(ctx *gin.Context) {
	var maintenance assets.AssetMaintenanceType
	if err := ctx.ShouldBindJSON(&maintenance); err != nil {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}

	//token, err := utils.ExtractClaimsResponse(ctx)
	//if err != nil {
	//	return
	//}
	//
	//maintenances, err := c.Service
	//if err != nil {
	//	response.SendResponse(ctx, http.StatusInternalServerError, "Failed to create maintenance record", nil, err.Error())
	//	return
	//}
	//
	//response.SendResponse(ctx, http.StatusCreated, "Maintenance record created successfully", maintenances, nil)
}
