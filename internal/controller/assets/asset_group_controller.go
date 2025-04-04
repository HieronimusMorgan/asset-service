package assets

import (
	request "asset-service/internal/dto/in/assets"
	"asset-service/internal/services/assets"
	"asset-service/internal/utils"
	"asset-service/package/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AssetGroupController interface {
	AddAssetGroup(context *gin.Context)
	UpdateAssetGroup(context *gin.Context)
	GetAssetGroupByID(context *gin.Context)
	AddMemberAssetGroup(context *gin.Context)
}

type assetGroupController struct {
	AssetGroupService assets.AssetGroupService
	JWTService        utils.JWTService
}

func NewAssetGroupController(AssetGroupService assets.AssetGroupService, JWTService utils.JWTService) AssetGroupController {
	return assetGroupController{AssetGroupService: AssetGroupService, JWTService: JWTService}
}

func (a assetGroupController) AddAssetGroup(context *gin.Context) {
	var req request.AssetGroupRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, exist := utils.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	err := a.AssetGroupService.AddAssetGroup(&req, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", err.Error(), err)
		return
	}

	response.SendResponse(context, http.StatusOK, "Success", nil, "Asset group created successfully")
}

func (a assetGroupController) UpdateAssetGroup(context *gin.Context) {
	var req request.AssetGroupRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	assetGroupID, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, 400, "Resource ID must be a number", nil, err)
		return
	}

	token, exist := utils.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	data, err := a.AssetGroupService.UpdateAssetGroup(assetGroupID, &req, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", err.Error(), err)
		return
	}

	response.SendResponse(context, http.StatusOK, "Asset group name updated successfully", data, nil)
}
func (a assetGroupController) GetAssetGroupByID(context *gin.Context) {
	assetGroupID, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, 400, "Resource ID must be a number", nil, err)
		return
	}

	token, exist := utils.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	data, err := a.AssetGroupService.GetAssetGroupDetailByID(assetGroupID, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", err.Error(), err)
		return
	}

	response.SendResponse(context, http.StatusOK, "Success", data, nil)
}

func (a assetGroupController) AddMemberAssetGroup(context *gin.Context) {
	var req request.AssetGroupMemberRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, exist := utils.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	err := a.AssetGroupService.AddMemberAssetGroup(&req, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", err.Error(), err)
		return
	}

	response.SendResponse(context, http.StatusOK, "Success", nil, "Member added to asset group successfully")
}
