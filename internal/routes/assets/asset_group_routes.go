package assets

import (
	"asset-service/config"
	"asset-service/internal/controller/assets"
	"github.com/gin-gonic/gin"
)

func AssetGroupRoutes(r *gin.Engine, middleware config.Middleware, assetGroupController assets.AssetGroupController, assetGroupPermission assets.AssetGroupPermissionController) {
	assetGroup := r.Group("/v1/asset-group")
	assetGroup.Use(middleware.AuthMiddleware.Handler())
	{
		assetGroup.POST("", assetGroupController.AddAssetGroup)
		assetGroup.PUT("/:id", assetGroupController.UpdateAssetGroup)
		assetGroup.GET("/:id", assetGroupController.GetAssetGroupByID)
		assetGroup.DELETE("/:id", assetGroupController.DeleteAssetGroup)
	}

	assetGroupAsset := r.Group("/v1/asset-group/asset")
	assetGroupAsset.Use(middleware.AuthMiddleware.Handler())
	{
		assetGroupAsset.GET("/:id", assetGroupController.GetListAssetGroupAsset)
		assetGroupAsset.POST("/add-stock", assetGroupController.AddStockAssetGroupAsset)
		assetGroupAsset.POST("/reduce-stock", assetGroupController.ReduceStockAssetGroupAsset)
	}

	assetPermission := r.Group("/v1/asset-group/permission")
	assetPermission.Use(middleware.AuthMiddleware.Handler())
	{
		assetPermission.POST("/add", assetGroupController.AddPermissionMemberAssetGroup)
		assetPermission.POST("/remove", assetGroupController.RemovePermissionMemberAssetGroup)
	}

	assetGroupMember := r.Group("/v1/asset-group/member")
	assetGroupMember.Use(middleware.AuthMiddleware.Handler())
	{
		assetGroupMember.POST("/add", assetGroupController.AddMemberAssetGroup)
		assetGroupMember.POST("/remove", assetGroupController.RemoveMemberAssetGroup)
	}

	adminGroup := r.Group("/v1/admin/asset-group")
	adminGroup.Use(middleware.AdminMiddleware.Handler())
	{
		adminGroup.GET("/permission", assetGroupPermission.GetListAssetGroupPermission)
		adminGroup.GET("/permission/:id", assetGroupPermission.GetAssetGroupPermissionByID)
		adminGroup.POST("/permission", assetGroupPermission.AddAssetGroupPermission)
		adminGroup.PUT("/permission/:id", assetGroupPermission.UpdateAssetGroupPermission)
		adminGroup.DELETE("/permission/:id", assetGroupPermission.DeleteAssetGroupPermission)
	}
}
