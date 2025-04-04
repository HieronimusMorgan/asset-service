package assets

import (
	"asset-service/config"
	"asset-service/internal/controller/assets"
	"github.com/gin-gonic/gin"
)

func AssetGroupRoutes(r *gin.Engine, middleware config.Middleware, assetGroupController assets.AssetGroupController, assetGroupPermission assets.AssetGroupPermissionController) {

	routerGroup := r.Group("/v1/asset-group")
	routerGroup.Use(middleware.AuthMiddleware.Handler())
	{
		routerGroup.POST("", assetGroupController.AddAssetGroup)
		routerGroup.PUT("/:id", assetGroupController.UpdateAssetGroup)
		routerGroup.GET("/:id", assetGroupController.GetAssetGroupByID)
		routerGroup.POST("/member", assetGroupController.AddMemberAssetGroup)
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
