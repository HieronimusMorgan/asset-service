package main

import (
	"asset-service/config"
	"asset-service/internal/routes/assets"
	"log"
)

func main() {
	serverConfig, err := config.NewServerConfig()
	if err != nil {
		log.Fatalf("❌ Failed to initialize server: %v", err)
	}

	defer func() {
		sqlDB, _ := serverConfig.DB.DB()
		err := sqlDB.Close()
		if err != nil {
			return
		}
		log.Println("✅ Database connection closed")
	}()

	if err := serverConfig.Start(); err != nil {
		log.Fatalf("❌ Error starting server: %v", err)
	}

	engine := serverConfig.Gin

	assets.AssetCategoryRoutes(engine, serverConfig.Middleware, serverConfig.Controller.AssetCategory)
	assets.AssetStatusRoutes(engine, serverConfig.Middleware, serverConfig.Controller.AssetStatus)
	assets.AssetRoutes(engine, serverConfig.Middleware, serverConfig.Controller.Asset)
	assets.AssetWishlistRoutes(engine, serverConfig.Middleware, serverConfig.Controller.AssetWishlist)
	assets.AssetMaintenanceRoutes(engine, serverConfig.Middleware, serverConfig.Controller.AssetMaintenance)
	assets.AssetMaintenanceTypeRoutes(engine, serverConfig.Middleware, serverConfig.Controller.AssetMaintenanceType)
	assets.AssetGroupRoutes(engine, serverConfig.Middleware, serverConfig.Controller)
	// Run server
	log.Println("Starting server on :8081")
	err = engine.Run(":8081")
	if err != nil {
		return
	}
}
