package config

import (
	controller "asset-service/internal/controller/assets"
	"asset-service/internal/middleware"
	repository "asset-service/internal/repository/assets"
	"asset-service/internal/repository/transaction"
	users "asset-service/internal/repository/users"
	services "asset-service/internal/services/assets"
	"asset-service/internal/utils"
	controllercron "asset-service/internal/utils/cron/controller"
	repositorycron "asset-service/internal/utils/cron/repository"
	"asset-service/internal/utils/cron/service"
	nt "asset-service/internal/utils/nats"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func NewServerConfig() (*ServerConfig, error) {
	cfg := LoadConfig()
	redisClient := InitRedis(cfg)
	redisService := utils.NewRedisService(*redisClient)
	db := InitDatabase(cfg)
	engine := InitGin()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("🛑 Shutting down gracefully...")

		// Close database and Redis before exiting
		CloseDatabase(db)
		CloseRedis(redisClient)

		os.Exit(0)
	}()

	server := &ServerConfig{
		Gin:        engine,
		Config:     cfg,
		DB:         db,
		Redis:      redisService,
		JWTService: utils.NewJWTService(cfg.JWTSecret),
	}

	server.initNats()
	server.initRepository()
	server.initTransaction()
	server.initServices()
	server.initController()
	server.initMiddleware()
	server.initCron()
	return server, nil
}

// InitGin initializes the Gin engine with appropriate configurations
func InitGin() *gin.Engine {
	// Set Gin mode based on environment
	if ginMode := gin.Mode(); ginMode != gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
		logrus.Warn("⚠ Running in DEBUG mode. Use `GIN_MODE=release` in production.")
	} else {
		logrus.Info("✅ Running in RELEASE mode.")
	}

	// Create a new Gin router
	engine := gin.New()

	// Middleware
	engine.Use(gin.Recovery()) // Handles panics and prevents crashes
	engine.Use(gin.Logger())   // Logs HTTP requests

	// Security Headers (Prevents Clickjacking & XSS Attacks)
	engine.Use(func(c *gin.Context) {
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Next()
	})

	logrus.Info("🚀 Gin HTTP server initialized successfully")
	return engine
}

// initRepository initializes database access objects (Repository)
func (s *ServerConfig) initRepository() {
	s.Repository = Repository{
		UserRepository:                       users.NewUserRepository(*s.DB),
		UserSettingRepository:                users.NewUserSettingRepository(*s.DB),
		AssetAuditLog:                        repository.NewAssetAuditLogRepository(*s.DB),
		AssetCategory:                        repository.NewAssetCategoryRepository(*s.DB, s.Repository.AssetAuditLog),
		AssetMaintenanceType:                 repository.NewAssetMaintenanceTypeRepository(*s.DB),
		AssetRepository:                      repository.NewAssetRepository(*s.DB, s.Repository.AssetAuditLog),
		AssetStatusRepository:                repository.NewAssetStatusRepository(*s.DB),
		AssetWishlistRepository:              repository.NewAssetWishlistRepository(*s.DB, s.Repository.AssetAuditLog, s.Repository.AssetRepository),
		AssetMaintenance:                     repository.NewAssetMaintenanceRepository(*s.DB),
		AssetMaintenanceRecord:               repository.NewAssetMaintenanceRecordRepository(*s.DB),
		AssetImageRepository:                 repository.NewAssetImageRepository(*s.DB),
		AssetStockRepository:                 repository.NewAssetStockRepository(*s.DB),
		AssetGroupRepository:                 repository.NewAssetGroupRepository(*s.DB, s.Repository.AssetAuditLog),
		AssetGroupAssetRepository:            repository.NewAssetGroupAssetRepository(*s.DB, s.Repository.AssetAuditLog),
		AssetGroupMemberRepository:           repository.NewAssetGroupMemberRepository(*s.DB, s.Repository.AssetAuditLog),
		AssetGroupMemberPermissionRepository: repository.NewAssetGroupMemberPermissionRepository(*s.DB, repository.NewAssetAuditLogRepository(*s.DB)),
		AssetGroupPermissionRepository:       repository.NewAssetGroupPermissionRepository(*s.DB, repository.NewAssetAuditLogRepository(*s.DB)),
		AssetGroupInvitation:                 repository.NewAssetGroupInvitationRepository(*s.DB),
	}
}

// initServices initializes the application services
func (s *ServerConfig) initServices() {
	s.Services = Services{
		AssetCategory: services.NewAssetCategoryService(
			s.Repository.AssetCategory,
			s.Repository.AssetRepository,
			s.Repository.AssetAuditLog,
			s.Redis),
		AssetMaintenance: services.NewAssetMaintenanceService(
			s.Repository.AssetMaintenance,
			s.Repository.AssetRepository,
			s.Repository.AssetMaintenanceRecord,
			s.Repository.AssetAuditLog,
			s.Redis),
		AssetMaintenanceType: services.NewAssetMaintenanceTypeService(
			s.Repository.AssetMaintenanceType,
			s.Repository.AssetMaintenance,
			s.Redis),
		Asset: services.NewAssetService(
			s.Repository.UserRepository,
			s.Repository.AssetRepository,
			s.Repository.AssetCategory,
			s.Repository.AssetStatusRepository,
			s.Repository.AssetImageRepository,
			s.Repository.AssetAuditLog,
			s.Repository.AssetGroupMemberRepository,
			s.Repository.AssetGroupAssetRepository,
			s.Redis,
			s.Transaction.AssetTransactionRepository,
			s.Repository.AssetStockRepository),
		AssetStatus: services.NewAssetStatusService(
			s.Repository.AssetStatusRepository,
			s.Repository.AssetAuditLog,
			s.Redis),
		AssetWishlist: services.NewAssetWishlistService(
			s.Repository.AssetWishlistRepository,
			s.Repository.AssetCategory,
			s.Repository.AssetStatusRepository,
			s.Repository.AssetRepository,
			s.Repository.AssetAuditLog,
			s.Redis),
		AssetImage: services.NewAssetImageService(
			s.Repository.AssetImageRepository,
			s.Repository.AssetRepository,
			s.Redis,
			s.Nats.NatsService),
		AssetGroupAssetService: services.NewAssetGroupAssetService(
			s.Repository.AssetGroupAssetRepository,
			s.Repository.AssetRepository,
			s.Repository.AssetAuditLog,
			s.Redis,
		),
		AssetGroupPermissionService: services.NewAssetGroupPermissionService(
			s.Repository.AssetGroupPermissionRepository,
			s.Repository.AssetRepository,
			s.Repository.AssetAuditLog,
			s.Redis),
		AssetGroupMemberService: services.NewAssetGroupMemberService(
			s.Repository.UserRepository,
			s.Repository.UserSettingRepository,
			s.Repository.AssetGroupRepository,
			s.Repository.AssetGroupMemberPermissionRepository,
			s.Repository.AssetGroupMemberRepository,
			s.Repository.AssetRepository,
			s.Repository.AssetGroupInvitation,
			s.Repository.AssetAuditLog,
			s.Redis),
		AssetGroupService: services.NewAssetGroupService(
			s.Repository.UserRepository,
			s.Repository.AssetGroupRepository,
			s.Repository.AssetGroupPermissionRepository,
			s.Repository.AssetGroupMemberPermissionRepository,
			s.Repository.AssetGroupMemberRepository,
			s.Repository.AssetGroupAssetRepository,
			s.Repository.AssetRepository,
			s.Repository.AssetStockRepository,
			s.Repository.AssetAuditLog,
			s.Redis),
	}
}

// initNats initializes the application services
func (s *ServerConfig) initNats() {
	s.Nats = Nats{
		NatsService: nt.NewNatsService(s.Config.NatsUrl),
	}
}

// Start initializes everything and returns an error if something fails
func (s *ServerConfig) Start() error {
	log.Println("✅ Server configuration initialized successfully!")
	return nil
}

func (s *ServerConfig) initController() {
	s.Controller = Controller{
		AssetCategory:                  controller.NewAssetCategoryController(s.Services.AssetCategory, s.JWTService),
		AssetMaintenance:               controller.NewAssetMaintenanceController(s.Services.AssetMaintenance, s.JWTService),
		AssetMaintenanceType:           controller.NewAssetMaintenanceTypeController(s.Services.AssetMaintenanceType, s.JWTService),
		Asset:                          controller.NewAssetController(s.Services.Asset, s.JWTService, s.Config.CdnUrl),
		AssetStatus:                    controller.NewAssetStatusController(s.Services.AssetStatus, s.JWTService),
		AssetWishlist:                  controller.NewAssetWishlistController(s.Services.AssetWishlist, s.JWTService),
		AssetGroupController:           controller.NewAssetGroupController(s.Services.AssetGroupService, s.JWTService),
		AssetGroupMemberController:     controller.NewAssetGroupMemberController(s.Services.AssetGroupMemberService, s.JWTService),
		AssetGroupPermissionController: controller.NewAssetGroupPermissionController(s.Services.AssetGroupPermissionService, s.JWTService),
	}
}

func (s *ServerConfig) initMiddleware() {
	s.Middleware = Middleware{
		AuthMiddleware:  middleware.NewAuthMiddleware(s.JWTService),
		AdminMiddleware: middleware.NewAdminMiddleware(s.JWTService),
	}
}

func (s *ServerConfig) initTransaction() {
	s.Transaction = Transaction{
		AssetTransactionRepository: transaction.NewAssetTransactionRepository(*s.DB,
			s.Repository.AssetRepository,
			s.Repository.AssetCategory,
			s.Repository.AssetStatusRepository,
			s.Repository.AssetMaintenance,
			s.Repository.AssetMaintenanceRecord,
			s.Repository.AssetAuditLog),
	}
}

func (s *ServerConfig) initCron() {
	s.Cron = Cron{
		CronRepository: repositorycron.NewCronRepository(*s.DB),
		CronService:    service.NewCronService(*s.DB, repositorycron.NewCronRepository(*s.DB), s.Services.AssetMaintenance, s.Services.AssetImage),
		CronController: controllercron.NewCronJobController(service.NewCronService(*s.DB, repositorycron.NewCronRepository(*s.DB), s.Services.AssetMaintenance, nil)),
	}
	s.Cron.CronService.Start()
}
