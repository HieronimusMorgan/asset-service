package assets

import (
	request "asset-service/internal/dto/in/assets"
	service "asset-service/internal/services/assets"
	"asset-service/internal/utils"
	"asset-service/internal/utils/jwt"
	"asset-service/package/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AssetMaintenanceTypeController interface {
	CreateMaintenanceType(ctx *gin.Context)
	GetMaintenanceByID(context *gin.Context)
	GetListMaintenanceType(context *gin.Context)
}

type assetMaintenanceTypeController struct {
	Service    service.AssetMaintenanceTypeService
	JWTService jwt.Service
}

func NewAssetMaintenanceTypeController(Service service.AssetMaintenanceTypeService, JWTService jwt.Service) AssetMaintenanceTypeController {
	return assetMaintenanceTypeController{Service: Service, JWTService: JWTService}
}

func (c assetMaintenanceTypeController) CreateMaintenanceType(ctx *gin.Context) {
	var maintenance *request.AssetMaintenanceTypeRequest
	if err := ctx.ShouldBindJSON(&maintenance); err != nil {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}

	credentialKey := ctx.GetHeader(utils.XCredentialKey)
	if credentialKey == "" {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, "credential key not found")
		return
	}

	token, exist := jwt.ExtractTokenClaims(ctx)
	if !exist {
		response.SendResponse(ctx, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	maintenances, err := c.Service.AddMaintenanceType(maintenance, token.ClientID, credentialKey)
	if err != nil {
		response.SendResponse(ctx, http.StatusInternalServerError, "Failed to create maintenance record", nil, err.Error())
		return
	}

	response.SendResponse(ctx, http.StatusCreated, "Maintenance record created successfully", maintenances, nil)

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

func (c assetMaintenanceTypeController) GetMaintenanceByID(context *gin.Context) {
	assetID, err := utils.ConvertToUint(context.Param("id"))
	token, err := c.JWTService.ExtractClaims(context.GetHeader(utils.Authorization))
	if err != nil {
		return
	}

	maintenance, err := c.Service.GetMaintenanceTypeByID(assetID, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusNotFound, "Maintenance not found", nil, err.Error())
		return
	}

	if maintenance == nil {
		response.SendResponse(context, http.StatusNotFound, "Maintenance not found", nil, nil)
		return
	}

	response.SendResponse(context, http.StatusOK, "Success", maintenance, nil)

}

func (c assetMaintenanceTypeController) GetListMaintenanceType(context *gin.Context) {
	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	maintenanceTypes, err := c.Service.GetListMaintenanceType(token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}
	if maintenanceTypes == nil {
		response.SendResponse(context, http.StatusNotFound, "Maintenance types not found", nil, nil)
		return
	}
	response.SendResponse(context, http.StatusOK, "Success", maintenanceTypes, nil)
}
