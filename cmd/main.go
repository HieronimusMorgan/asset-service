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

	// Ensure database connection closes when the server shuts down
	defer func() {
		sqlDB, _ := serverConfig.DB.DB()
		sqlDB.Close()
		log.Println("✅ Database connection closed")
	}()

	// Start server config (Ensure everything is ready)
	if err := serverConfig.Start(); err != nil {
		log.Fatalf("❌ Error starting server: %v", err)
	}

	// Initialize Router
	engine := serverConfig.Gin

	// Register routes
	assets.AssetCategoryRoutes(engine, serverConfig.Middleware.AuthMiddleware, serverConfig.Controller.AssetCategory)
	assets.AssetStatusRoutes(engine, serverConfig.Middleware.AuthMiddleware, serverConfig.Controller.AssetStatus)
	assets.AssetRoutes(engine, serverConfig.Middleware.AuthMiddleware, serverConfig.Controller.Asset)
	assets.AssetWishlistRoutes(engine, serverConfig.Middleware.AuthMiddleware, serverConfig.Controller.AssetWishlist)
	assets.AssetMaintenanceRoutes(engine, serverConfig.Middleware.AuthMiddleware, serverConfig.Controller.AssetMaintenance)

	// Run server
	log.Println("Starting server on :8081")
	err = engine.Run(":8081")
	if err != nil {
		return
	}
}
