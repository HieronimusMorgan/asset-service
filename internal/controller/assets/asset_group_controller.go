package assets

import (
	request "asset-service/internal/dto/in/assets"
	"asset-service/internal/services/assets"
	"asset-service/internal/utils"
	"asset-service/internal/utils/jwt"
	"asset-service/package/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AssetGroupController interface {
	AddAssetGroup(context *gin.Context)
	UpdateAssetGroup(context *gin.Context)
	GetAssetGroup(context *gin.Context)
	DeleteAssetGroup(context *gin.Context)

	AddInvitationTokenAssetGroup(context *gin.Context)
	RemoveInvitationTokenAssetGroup(context *gin.Context)

	InviteMemberAssetGroup(context *gin.Context)
	RemoveMemberAssetGroup(context *gin.Context)

	AddPermissionMemberAssetGroup(context *gin.Context)
	RemovePermissionMemberAssetGroup(context *gin.Context)

	GetListAssetGroupAsset(context *gin.Context)
	AddStockAssetGroupAsset(context *gin.Context)
	ReduceStockAssetGroupAsset(context *gin.Context)
}

type assetGroupController struct {
	AssetGroupService assets.AssetGroupService
	JWTService        jwt.Service
}

func NewAssetGroupController(AssetGroupService assets.AssetGroupService, JWTService jwt.Service) AssetGroupController {
	return assetGroupController{AssetGroupService: AssetGroupService, JWTService: JWTService}
}

func (a assetGroupController) AddAssetGroup(context *gin.Context) {
	var req request.AssetGroupRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	credentialKey := context.GetHeader(utils.XCredentialKey)
	if credentialKey == "" {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "credential key not found")
		return
	}

	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	data, err := a.AssetGroupService.AddAssetGroup(&req, token.ClientID, credentialKey)
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}

	response.SendResponse(context, http.StatusOK, "Asset group created successfully", data, nil)
}

func (a assetGroupController) AddInvitationTokenAssetGroup(context *gin.Context) {
	assetGroupID, err := utils.ConvertToUint(context.Param("id"))
	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	data, err := a.AssetGroupService.AddInvitationAssetGroup(assetGroupID, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", err.Error(), err)
		return
	}

	response.SendResponse(context, http.StatusOK, "Invitation created successfully", data, nil)
}

func (a assetGroupController) RemoveInvitationTokenAssetGroup(context *gin.Context) {
	assetGroupID, err := utils.ConvertToUint(context.Param("id"))
	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	err = a.AssetGroupService.RemoveInvitationAssetGroup(assetGroupID, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", err.Error(), err)
		return
	}

	response.SendResponse(context, http.StatusOK, "Invitation created successfully", nil, nil)
}

func (a assetGroupController) UpdateAssetGroup(context *gin.Context) {
	var req request.AssetGroupRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assetGroupID, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, 400, "Resource MaintenanceTypeID must be a number", nil, err)
		return
	}

	credentialKey := context.GetHeader(utils.XCredentialKey)
	if credentialKey == "" {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "credential key not found")
		return
	}

	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	data, err := a.AssetGroupService.UpdateAssetGroup(assetGroupID, &req, token.ClientID, credentialKey)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", err.Error(), err)
		return
	}

	response.SendResponse(context, http.StatusOK, "Asset group name updated successfully", data, nil)
}

func (a assetGroupController) GetAssetGroup(context *gin.Context) {
	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	data, err := a.AssetGroupService.GetAssetGroupDetail(token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, err.Error())
		return
	}

	response.SendResponse(context, http.StatusOK, "Success", data, nil)
}

func (a assetGroupController) DeleteAssetGroup(context *gin.Context) {
	assetGroupID, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, 400, "Resource MaintenanceTypeID must be a number", nil, err)
		return
	}

	credentialKey := context.GetHeader(utils.XCredentialKey)
	if credentialKey == "" {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "credential key not found")
		return
	}

	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	err = a.AssetGroupService.DeleteAssetGroup(assetGroupID, token.ClientID, credentialKey)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", err.Error(), err)
		return
	}

	response.SendResponse(context, http.StatusOK, "Asset group deleted successfully", nil, nil)
}

func (a assetGroupController) InviteMemberAssetGroup(context *gin.Context) {
	var req request.AssetGroupMemberRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	err := a.AssetGroupService.InviteMemberAssetGroup(&req, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", err.Error(), err)
		return
	}

	response.SendResponse(context, http.StatusOK, "Member added to asset group successfully", nil, nil)
}

func (a assetGroupController) RemoveMemberAssetGroup(context *gin.Context) {
	var req request.AssetGroupMemberRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	err := a.AssetGroupService.RemoveMemberAssetGroup(req, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", err.Error(), err)
		return
	}

	response.SendResponse(context, http.StatusOK, "Member removed from asset group successfully", nil, nil)
}

func (a assetGroupController) AddPermissionMemberAssetGroup(context *gin.Context) {
	var req request.ChangeAssetGroupPermissionRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	err := a.AssetGroupService.AddPermissionMemberAssetGroup(&req, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", err.Error(), err)
		return
	}

	response.SendResponse(context, http.StatusOK, "Permission added to member successfully", nil, nil)
}

func (a assetGroupController) RemovePermissionMemberAssetGroup(context *gin.Context) {
	var req request.ChangeAssetGroupPermissionRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	err := a.AssetGroupService.RemovePermissionMemberAssetGroup(&req, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", err.Error(), err)
		return
	}
	response.SendResponse(context, http.StatusOK, "Permission removed from member successfully", nil, nil)
}

func (a assetGroupController) GetListAssetGroupAsset(context *gin.Context) {
	assetGroupID, err := utils.ConvertToUint(context.Param("id"))
	if err != nil {
		response.SendResponse(context, 400, "Resource MaintenanceTypeID must be a number", nil, err)
		return
	}

	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	pageIndex, pageSize, err := utils.GetPageIndexPageSize(context)
	if err != nil {
		response.SendResponse(context, 400, "Invalid page index or page size", nil, err.Error())
		return
	}

	data, total, err := a.AssetGroupService.GetListAssetGroupAsset(assetGroupID, pageIndex, pageSize, token.ClientID)
	if err != nil {
		response.SendResponseList(context, 500, "Failed to get list assets", response.PagedData{
			Total:     total,
			PageIndex: pageIndex,
			PageSize:  pageSize,
			Items:     nil,
		}, err.Error())
		return
	}
	response.SendResponseList(context, 200, "Get list assets group asset successfully", response.PagedData{
		Total:     total,
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Items:     data,
	}, nil)
}

func (a assetGroupController) AddStockAssetGroupAsset(context *gin.Context) {
	var req request.ChangeAssetStockRequest

	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, 400, "Invalid request", nil, err.Error())
		return
	}
	token, err := a.JWTService.ExtractClaims(context.GetHeader(utils.Authorization))
	if err != nil {
		return
	}

	data, err := a.AssetGroupService.UpdateStockAssetGroupAsset(true, req, token.ClientID)
	if err != nil {
		response.SendResponse(context, 500, "Failed to update stock asset", nil, err.Error())
		return
	}

	response.SendResponse(context, 200, "Stock asset updated successfully", data, nil)
}

func (a assetGroupController) ReduceStockAssetGroupAsset(context *gin.Context) {
	var req request.ChangeAssetStockRequest

	if err := context.ShouldBindJSON(&req); err != nil {
		response.SendResponse(context, 400, "Invalid request", nil, err.Error())
		return
	}

	token, err := a.JWTService.ExtractClaims(context.GetHeader(utils.Authorization))
	if err != nil {
		return
	}

	data, err := a.AssetGroupService.UpdateStockAssetGroupAsset(false, req, token.ClientID)
	if err != nil {
		response.SendResponse(context, 500, "Failed to update stock asset", nil, err.Error())
		return
	}

	response.SendResponse(context, 200, "Stock asset updated successfully", data, nil)
}
