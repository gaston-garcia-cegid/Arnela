package main

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/config"
	_ "github.com/gaston-garcia-cegid/arnela/backend/docs"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/handler"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/middleware"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository/postgres"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/service"
	"github.com/gaston-garcia-cegid/arnela/backend/pkg/cache"
	"github.com/gaston-garcia-cegid/arnela/backend/pkg/database"
	"github.com/gaston-garcia-cegid/arnela/backend/pkg/jwt"
	"github.com/gaston-garcia-cegid/arnela/backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Arnela CRM/CMS API
// @version         1.0
// @description     API for Arnela professional office management system
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@arnela.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	log := logger.NewLogger(cfg.App.Env)
	log.Info("Starting Arnela API", map[string]interface{}{
		"env":     cfg.App.Env,
		"port":    cfg.App.Port,
		"version": "1.0.0",
	})

	// Set Gin mode
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize PostgreSQL connection
	db, err := database.NewPostgresDB(cfg.Database.GetDSN())
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}
	defer db.Close()

	// Run migrations
	migrationsPath, _ := filepath.Abs("./migrations")
	if err := database.RunMigrations(db.DB.DB, migrationsPath); err != nil {
		log.Fatal("Failed to run migrations", err)
	}

	// Initialize Redis connection (optional for now)
	redisClient, err := cache.NewRedisClient(cfg.Redis.GetRedisAddress(), cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Warn("Redis connection failed, continuing without cache", map[string]interface{}{
			"error": err.Error(),
		})
		redisClient = nil
	} else {
		defer redisClient.Close()
		log.Info("Redis connected successfully", nil)
	}

	// Initialize JWT token manager
	tokenManager := jwt.NewTokenManager(cfg.JWT.Secret, "arnela-api")

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db.DB)
	clientRepo := postgres.NewClientRepository(db.DB)

	// Initialize services
	authService := service.NewAuthService(userRepo, tokenManager, 24*time.Hour)
	clientService := service.NewClientService(clientRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	clientHandler := handler.NewClientHandler(clientService)

	// Initialize Gin router
	router := gin.New()

	// Add middlewares
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware(log))

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "arnela-api",
			"version": "1.0.0",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(tokenManager))
		{
			protected.GET("/auth/me", authHandler.Me)

			// Client routes
			clients := protected.Group("/clients")
			{
				// Admin and Employee can manage all clients
				clients.POST("", middleware.RequireRole("admin", "employee"), clientHandler.CreateClient)
				clients.GET("", middleware.RequireRole("admin", "employee"), clientHandler.ListClients)
				clients.GET("/:id", middleware.RequireRole("admin", "employee"), clientHandler.GetClient)
				clients.PUT("/:id", middleware.RequireRole("admin", "employee"), clientHandler.UpdateClient)
				clients.DELETE("/:id", middleware.RequireRole("admin"), clientHandler.DeleteClient)

				// Client can only access their own data
				clients.GET("/me", middleware.RequireRole("client"), clientHandler.GetMyClient)
			}
		}
	}

	// Start server
	address := fmt.Sprintf("0.0.0.0:%s", cfg.App.Port)
	log.Info("🚀 Server starting", map[string]interface{}{
		"address": address,
		"env":     cfg.App.Env,
	})

	if err := router.Run(address); err != nil {
		log.Fatal("Failed to start server", err)
	}
}
