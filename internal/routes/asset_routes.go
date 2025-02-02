package routes

import (
	"asset-service/internal/handler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AssetRoutes(r *gin.Engine, db *gorm.DB) {
	assetHandler := handler.NewAssetHandler(db)

	public := r.Group("/asset-service/v1/asset")
	{
		public.POST("/add", assetHandler.AddAsset)
		public.POST("/update/:id", assetHandler.UpdateAsset)
		public.POST("/update-status/:id", assetHandler.UpdateAssetStatus)
		public.POST("/update-category/:id", assetHandler.UpdateAssetCategory)
		public.GET("", assetHandler.GetListAsset)
		public.GET("/:id", assetHandler.GetAssetById)
		public.DELETE("/delete/:id", assetHandler.DeleteAsset)

		//public.POST("/category/add", assetHandler.AddAssetCategory)
		//public.POST("maintenance-record/add", assetHandler.RegisterMaintenanceRecord)
		//public.GET("", assetHandler.GetAssets)
		//public.GET("/category", assetHandler.GetAssetCategories)
		//public.GET("/maintenance-record", assetHandler.GetMaintenanceRecords)
	}
}
