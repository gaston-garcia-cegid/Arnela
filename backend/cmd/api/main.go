package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/config"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/handler"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/middleware"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository/postgres"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/service"
	"github.com/gaston-garcia-cegid/arnela/backend/pkg/database"
	"github.com/gaston-garcia-cegid/arnela/backend/pkg/jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/gaston-garcia-cegid/arnela/backend/docs"
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
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// ✅ PASO 1: Crear conexión inicial SOLO para migraciones
	log.Println("[DEBUG] Creating temporary database connection for migrations...")
	dbForMigrations, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database for migrations: %v", err)
	}

	// ✅ PASO 2: Ejecutar migraciones con la conexión temporal
	migrationsPath, _ := filepath.Abs("./migrations")
	if err := database.RunMigrations(dbForMigrations.DB, migrationsPath); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// ✅ PASO 3: Cerrar la conexión temporal
	if err := dbForMigrations.Close(); err != nil {
		log.Printf("[WARN] Error closing migrations database connection: %v", err)
	}
	log.Println("[DEBUG] Migrations connection closed")

	// ✅ PASO 4: Crear NUEVA conexión para la aplicación
	log.Println("[DEBUG] Creating application database connection...")
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Printf("[DEBUG] Main - db pointer: %p", db)

	// ✅ PASO 5: Verificar que la nueva conexión está sana
	if err := database.HealthCheck(db); err != nil {
		log.Fatalf("Database health check failed: %v", err)
	}

	stats := db.Stats()
	log.Printf("Database connection healthy - Pool stats: Open=%d, Idle=%d, InUse=%d",
		stats.OpenConnections, stats.Idle, stats.InUse)

	// Goroutine para mantener pool activo
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if err := database.HealthCheck(db); err != nil {
				log.Printf("[ERROR] Periodic health check failed: %v", err)
			} else {
				stats := db.Stats()
				log.Printf("[DEBUG] Periodic health check - Pool stats: Open=%d, Idle=%d, InUse=%d",
					stats.OpenConnections, stats.Idle, stats.InUse)
			}
		}
	}()

	// Initialize JWT token manager
	tokenManager := jwt.NewTokenManager(cfg.JWT.Secret, "arnela-api")

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	clientRepo := postgres.NewClientRepository(db)
	appointmentRepo := postgres.NewAppointmentRepository(db)
	employeeRepo := postgres.NewEmployeeRepository(db)
	statsRepo := postgres.NewStatsRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, clientRepo, tokenManager, cfg.JWT.TokenExpiry)
	clientService := service.NewClientService(clientRepo, userRepo)
	appointmentService := service.NewAppointmentService(appointmentRepo, clientRepo, employeeRepo)
	employeeService := service.NewEmployeeService(employeeRepo, userRepo)
	statsService := service.NewStatsService(statsRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	clientHandler := handler.NewClientHandler(clientService)
	appointmentHandler := handler.NewAppointmentHandler(appointmentService)
	employeeHandler := handler.NewEmployeeHandler(employeeService)
	statsHandler := handler.NewStatsHandler(statsService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenManager)
	authMiddleware.SetClientRepo(clientRepo)

	// Setup Gin router
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// CORS Configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		if err := database.HealthCheck(db); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "unhealthy",
				"database": "disconnected",
				"error":    err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "healthy",
			"database": "connected",
		})
	})

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", authMiddleware.RequireAuth(), authHandler.Me)
		}

		// Client routes (authenticated)
		clients := v1.Group("/clients")
		clients.Use(authMiddleware.RequireAuth())
		{
			// Client self-service (any authenticated user can get their own client info)
			clients.GET("/me", clientHandler.GetMyClient)

			// Admin/Employee only routes
			clients.POST("", authMiddleware.RequireRole("admin", "employee"), clientHandler.CreateClient)
			clients.GET("", authMiddleware.RequireRole("admin", "employee"), clientHandler.ListClients)
			clients.GET("/:id", authMiddleware.RequireRole("admin", "employee"), clientHandler.GetClient)
			clients.PUT("/:id", authMiddleware.RequireRole("admin", "employee"), clientHandler.UpdateClient)
			clients.DELETE("/:id", authMiddleware.RequireRole("admin"), clientHandler.DeleteClient)
		}

		// Appointment routes (authenticated)
		appointments := v1.Group("/appointments")
		appointments.Use(authMiddleware.RequireAuth())
		{
			// Public endpoints (all authenticated users)
			appointments.GET("/therapists", appointmentHandler.GetTherapists)
			appointments.GET("/available-slots", appointmentHandler.GetAvailableSlots)
			appointments.POST("", appointmentHandler.CreateAppointment)
			appointments.GET("/:id", appointmentHandler.GetAppointment)
			appointments.PUT("/:id", appointmentHandler.UpdateAppointment)
			appointments.POST("/:id/cancel", appointmentHandler.CancelAppointment)

			// Client-specific endpoint
			appointments.GET("/me", appointmentHandler.GetMyAppointments)

			// Admin/Employee only routes
			appointments.GET("", authMiddleware.RequireRole("admin", "employee"), appointmentHandler.ListAppointments)
			appointments.POST("/:id/confirm", authMiddleware.RequireRole("admin", "employee"), appointmentHandler.ConfirmAppointment)
		}

		// Employee routes (authenticated)
		employees := v1.Group("/employees")
		employees.Use(authMiddleware.RequireAuth())
		{
			// Public endpoints (clients can see employees and specialties for booking)
			employees.GET("", employeeHandler.ListEmployees)
			employees.GET("/:id", employeeHandler.GetEmployee)
			employees.GET("/specialty/:specialty", employeeHandler.GetEmployeesBySpecialty)

			// Admin only routes
			employees.POST("", authMiddleware.RequireRole("admin"), employeeHandler.CreateEmployee)
			employees.PUT("/:id", authMiddleware.RequireRole("admin"), employeeHandler.UpdateEmployee)
			employees.DELETE("/:id", authMiddleware.RequireRole("admin"), employeeHandler.DeleteEmployee)
		}

		// Stats routes (authenticated)
		stats := v1.Group("/stats")
		stats.Use(authMiddleware.RequireAuth())
		{
			// Admin/Employee only routes
			stats.GET("/dashboard", authMiddleware.RequireRole("admin", "employee"), statsHandler.GetDashboardStats)
		}
	}

	// Setup graceful shutdown
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Gracefully shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	// Cerrar la DB al final
	if err := db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	} else {
		log.Println("Database connection closed successfully")
	}

	log.Println("Server exited gracefully")
}
