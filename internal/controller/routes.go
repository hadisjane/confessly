package controller

import (
	"github.com/hadisjane/confessly/internal/configs"
	"github.com/hadisjane/confessly/internal/middleware"
	"github.com/hadisjane/confessly/logger"
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RunServer() error {
	// Set Gin mode based on configuration
	if configs.AppSettings.AppParams.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()

	// Add logging middleware
	r.Use(gin.LoggerWithWriter(logger.Info.Writer()))
	r.Use(gin.Recovery())

	// Health check endpoint
	r.GET("/", Ping)

	// Auth routes
	authG := r.Group("/auth")
	{
		authG.POST("/register", Register)
		authG.POST("/login", Login)
	}

	// Public routes (no auth required)
	public := r.Group("/public")
	public.Use(middleware.TryParseUserContext)
	public.Use(middleware.GuestUUIDMiddleware())
	{
		public.GET("/confessions", GetAllConfessions)
		public.GET("/confessions/:id", GetConfession)
		public.GET("/confessions/search", SearchConfessions)
		public.POST("/confessions", CreateConfession)
	}

	// API routes with authentication middleware
	apiG := r.Group("/api", middleware.CheckUserAuthentication)

	// Confession routes
	confessionsG := apiG.Group("/confessions")
	{
		confessionsG.PUT("/:id", UpdateConfession)
		confessionsG.DELETE("/:id", DeleteConfession)
		confessionsG.GET("/search", SearchConfessions)
	}

	// Report routes
	reportsG := apiG.Group("/reports")
	{
		reportsG.POST("", CreateReport)
	}

	// Admin routes
	adminG := apiG.Group("/admin", middleware.CheckAdminAuthentication)
	{
		adminG.GET("/reports", GetReports)
		adminG.PUT("/reports/:id", UpdateReport)
		adminG.GET("/reports/:id", GetReport)
		adminG.GET("/users", GetUsers)
		adminG.GET("/users/:id", GetUserByID)
		adminG.DELETE("/confessions/:id", DeleteConfessionByAdmin)
		adminG.POST("/users/:id/ban", BanUser)
		adminG.POST("/users/:id/unban", UnbanUser)
		adminG.GET("/guests", GetGuestUsers)
		adminG.GET("/guests/:uuid", GetGuestUser)
		adminG.POST("/guests/:uuid/ban", BanGuestUser)
		adminG.POST("/guests/:uuid/unban", UnbanGuestUser)
	}

	// Get server address from config
	serverAddr := ":" + configs.AppSettings.AppParams.PortRun
	if configs.AppSettings.AppParams.PortRun[0] == ':' {
		serverAddr = configs.AppSettings.AppParams.PortRun
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	logger.Info.Printf("Starting server on %s", serverAddr)

	// Start the server
	if err := r.Run(serverAddr); err != nil {
		logger.Error.Printf("Error running server: %v", err)
		return fmt.Errorf("failed to start server: %v", err)
	}

	return nil
}
