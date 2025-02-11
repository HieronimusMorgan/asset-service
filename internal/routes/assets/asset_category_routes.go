package assets

import (
	"asset-service/internal/controller/assets"
	"asset-service/internal/middleware"
	"github.com/gin-gonic/gin"
)

func AssetCategoryRoutes(r *gin.Engine, middleware middleware.AuthMiddleware, controller assets.AssetCategoryController) {

	protected := r.Group("/asset-service/v1/asset/category")
	protected.Use(middleware.Handler())
	{
		protected.POST("/add", controller.AddAssetCategory)
		protected.POST("/update/:id", controller.UpdateAssetCategory)
		protected.GET("", controller.GetListAssetCategory)
		protected.GET("/:id", controller.GetAssetCategoryById)
		protected.DELETE("/delete/:id", controller.DeleteAssetCategory)
	}
}
