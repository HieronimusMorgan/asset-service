package routes

import (
	"asset-service/internal/controller"
	"asset-service/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AssetCategoryRoutes(r *gin.Engine, db *gorm.DB) {
	assetHandler := controller.NewAssetCategoryController(db)

	protected := r.Group("/asset-service/v1/asset/category")
	protected.Use(middleware.Middleware())
	{
		protected.POST("/add", assetHandler.AddAssetCategory)
		protected.POST("/update/:id", assetHandler.UpdateAssetCategory)
		protected.GET("", assetHandler.GetListAssetCategory)
		protected.GET("/:id", assetHandler.GetAssetCategoryById)
		protected.DELETE("/delete/:id", assetHandler.DeleteAssetCategory)
	}
}
