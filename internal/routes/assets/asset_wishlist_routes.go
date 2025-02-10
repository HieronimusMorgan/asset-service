package assets

import (
	"asset-service/internal/controller/assets"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AssetWishlistRoutes(r *gin.Engine, db *gorm.DB) {
	assetHandler := assets.NewAssetWishlistController(db)

	public := r.Group("/asset-service/v1/asset-wishlist")
	{
		public.POST("/add", assetHandler.AddWishlistAsset)
		public.POST("/update/:id", assetHandler.UpdateWishlistAsset)
		public.GET("", assetHandler.GetListWishlistAsset)
		public.GET("/:id", assetHandler.GetWishlistAssetByID)
		public.DELETE("/delete/:id", assetHandler.DeleteWishlistAsset)
	}
}
