package assets

import (
	"asset-service/internal/controller/assets"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AssetRoutes(r *gin.Engine, db *gorm.DB) {
	assetHandler := assets.NewAssetController(db)

	public := r.Group("/asset-service/v1/asset")
	{
		public.POST("/add", assetHandler.AddAsset)
		//public.POST("/wishlist", assetHandler.AddWishlistAsset)
		public.POST("/update/:id", assetHandler.UpdateAsset)
		public.POST("/update-status/:id", assetHandler.UpdateAssetStatus)
		public.POST("/update-category/:id", assetHandler.UpdateAssetCategory)
		public.GET("", assetHandler.GetListAsset)
		public.GET("/:id", assetHandler.GetAssetById)
		public.DELETE("/delete/:id", assetHandler.DeleteAsset)
	}
}
