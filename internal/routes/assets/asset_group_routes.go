package assets

import (
	"asset-service/config"
	"github.com/gin-gonic/gin"
)

func AssetGroupRoutes(r *gin.Engine, middleware config.Middleware, controller config.Controller) {
	assetGroup := r.Group("/v1/asset-group")
	assetGroup.Use(middleware.AssetMiddleware.HandlerAssetGroup())
	{
		assetGroup.POST("", controller.AssetGroupController.AddAssetGroup)
		assetGroup.GET("/add-invitation-token/:id", controller.AssetGroupController.AddInvitationTokenAssetGroup)
		assetGroup.GET("/remove-invitation-token/:id", controller.AssetGroupController.RemoveInvitationTokenAssetGroup)
		assetGroup.PUT("/:id", controller.AssetGroupController.UpdateAssetGroup)
		assetGroup.DELETE("/:id", controller.AssetGroupController.DeleteAssetGroup)
	}

	assetGroup.Use(middleware.AssetMiddleware.HandlerAsset())
	{
		assetGroup.GET("", controller.AssetGroupController.GetAssetGroup)
	}

	assetGroupAsset := r.Group("/v1/asset-group/asset")
	assetGroupAsset.Use(middleware.AssetMiddleware.HandlerAsset())
	{
		assetGroupAsset.GET("/:id", controller.AssetGroupController.GetListAssetGroupAsset)
		assetGroupAsset.POST("/add-stock", controller.AssetGroupController.AddStockAssetGroupAsset)
		assetGroupAsset.POST("/reduce-stock", controller.AssetGroupController.ReduceStockAssetGroupAsset)
	}

	assetPermission := r.Group("/v1/asset-group/permission")
	assetPermission.Use(middleware.AssetMiddleware.HandlerAsset())
	{
		assetPermission.POST("/add", controller.AssetGroupController.AddPermissionMemberAssetGroup)
		assetPermission.POST("/remove", controller.AssetGroupController.RemovePermissionMemberAssetGroup)
	}

	assetGroupMember := r.Group("/v1/asset-group/member")
	assetGroupMember.Use(middleware.AssetMiddleware.HandlerAsset())
	{
		assetGroupMember.POST("/add", controller.AssetGroupMemberController.InviteMemberAssetGroup)
		assetGroupMember.POST("/remove", controller.AssetGroupMemberController.RemoveMemberAssetGroup)
		assetGroupMember.GET("/:id", controller.AssetGroupMemberController.GetListMemberAssetGroup)
		assetGroupMember.DELETE("/:id", controller.AssetGroupMemberController.LeaveMemberAssetGroup)
	}
	assetGroupInvitation := r.Group("/v1/asset-group/invitation")
	assetGroupInvitation.Use(middleware.AssetMiddleware.HandlerAsset())
	{
		assetGroupInvitation.GET("/add-invitation-token/:id", controller.AssetGroupController.AddInvitationTokenAssetGroup)
		assetGroupInvitation.GET("/remove-invitation-token/:id", controller.AssetGroupController.RemoveInvitationTokenAssetGroup)
		assetGroupInvitation.PUT("/:id", controller.AssetGroupController.UpdateAssetGroup)
		assetGroupInvitation.GET("/:id", controller.AssetGroupController.GetAssetGroup)
		assetGroupInvitation.DELETE("/:id", controller.AssetGroupController.DeleteAssetGroup)
	}

	adminGroup := r.Group("/v1/admin/asset-group")
	adminGroup.Use(middleware.AdminMiddleware.HandlerAsset())
	{
		adminGroup.GET("/permission", controller.AssetGroupPermissionController.GetListAssetGroupPermission)
		adminGroup.GET("/permission/:id", controller.AssetGroupPermissionController.GetAssetGroupPermissionByID)
		adminGroup.POST("/permission", controller.AssetGroupPermissionController.AddAssetGroupPermission)
		adminGroup.PUT("/permission/:id", controller.AssetGroupPermissionController.UpdateAssetGroupPermission)
		adminGroup.DELETE("/permission/:id", controller.AssetGroupPermissionController.DeleteAssetGroupPermission)
	}
}
