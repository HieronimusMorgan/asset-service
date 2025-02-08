package assets

import (
	"asset-service/internal/models/assets"
	service "asset-service/internal/services/assets"
	"asset-service/package/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type AssetMaintenanceTypeController struct {
	Service *service.AssetMaintenanceTypeService
}

func NewAssetMaintenanceTypeController(db *gorm.DB) *AssetMaintenanceTypeController {
	return &AssetMaintenanceTypeController{Service: service.NewAssetMaintenanceTypeService(db)}
}

func (c *AssetMaintenanceTypeController) CreateMaintenanceType(ctx *gin.Context) {
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
