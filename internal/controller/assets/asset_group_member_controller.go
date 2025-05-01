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

type AssetGroupMemberController interface {
	InviteMemberAssetGroup(context *gin.Context)
	RemoveMemberAssetGroup(context *gin.Context)
	GetListMemberAssetGroup(context *gin.Context)
	LeaveMemberAssetGroup(context *gin.Context)
}

type assetGroupMemberController struct {
	AssetGroupMemberService assets.AssetGroupMemberService
	JWTService              jwt.Service
}

func NewAssetGroupMemberController(AssetGroupMemberService assets.AssetGroupMemberService, JWTService jwt.Service) AssetGroupMemberController {
	return assetGroupMemberController{AssetGroupMemberService: AssetGroupMemberService, JWTService: JWTService}
}

func (a assetGroupMemberController) InviteMemberAssetGroup(context *gin.Context) {
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

	err := a.AssetGroupMemberService.AddAssetGroupMember(&req, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", err.Error(), err)
		return
	}

	response.SendResponse(context, http.StatusOK, "Member added to asset group successfully", nil, nil)
}

func (a assetGroupMemberController) RemoveMemberAssetGroup(context *gin.Context) {
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

	err := a.AssetGroupMemberService.RemoveMemberAssetGroup(req, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", err.Error(), err)
		return
	}

	response.SendResponse(context, http.StatusOK, "Member removed from asset group successfully", nil, nil)
}

func (a assetGroupMemberController) GetListMemberAssetGroup(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Asset group MaintenanceTypeID is required")
		return
	}

	assetGroupID, err := utils.ConvertToUint(id)
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Invalid asset group MaintenanceTypeID")
		return
	}
	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	data, err := a.AssetGroupMemberService.GetListAssetGroupMember(assetGroupID, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", err.Error(), err)
		return
	}

	response.SendResponse(context, http.StatusOK, "List of members in asset group", data, nil)
}

func (a assetGroupMemberController) LeaveMemberAssetGroup(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Asset group MaintenanceTypeID is required")
		return
	}

	assetGroupID, err := utils.ConvertToUint(id)
	if err != nil {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Invalid asset group MaintenanceTypeID")
		return
	}

	token, exist := jwt.ExtractTokenClaims(context)
	if !exist {
		response.SendResponse(context, http.StatusBadRequest, "Error", nil, "Token not found")
		return
	}

	err = a.AssetGroupMemberService.LeaveMemberAssetGroup(assetGroupID, token.ClientID)
	if err != nil {
		response.SendResponse(context, http.StatusInternalServerError, "Error", err.Error(), err)
		return
	}

	response.SendResponse(context, http.StatusOK, "Member left asset group successfully", nil, nil)
}
