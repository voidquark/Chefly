package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"chefly/config"
	"chefly/database"
	"chefly/handlers"
	"chefly/middleware"
	"chefly/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	Version   = "production"
	BuildTime = "unknown"
)

//go:embed frontend/dist/*
var frontendFS embed.FS

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.InitDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize audit logger
	auditLogger := services.NewAuditLogger(
		cfg.AuditLogEnabled,
		cfg.AuditLogLevel,
		cfg.AuditLogFormat,
	)

	// Log system startup
	auditLogger.Info("system.startup", "Application started", nil)

	// Start automatic token cleanup goroutine (runs every hour)
	go cleanupExpiredTokens(db, auditLogger)

	// Setup Gin router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())

	// Security headers middleware (first)
	router.Use(middleware.SecurityHeaders())

	// Audit middleware
	router.Use(middleware.AuditMiddleware(auditLogger))

	// CORS middleware
	corsConfig := cors.DefaultConfig()
	if cfg.Environment == "production" {
		// In production, restrict to specific origins
		corsConfig.AllowOrigins = []string{"http://localhost:5173", "http://localhost:8080"}
	} else {
		// In development, allow all origins
		corsConfig.AllowAllOrigins = true
	}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = false // Must be false when AllowAllOrigins is true

	router.Use(cors.New(corsConfig))

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, cfg.JWTSecret, cfg.RegistrationEnabled)
	recipeHandler := handlers.NewRecipeHandler(db, cfg.ClaudeAPIKey, cfg.ClaudeModel, cfg.OpenAIAPIKey, cfg.OpenAIModel, cfg.RecipeGenerationLimit, auditLogger)
	shoppingListHandler := handlers.NewShoppingListHandler(db)
	adminHandler := handlers.NewAdminHandler(db, auditLogger)

	// Public routes
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken) // Public - no auth required
		}

		// Public recipe view (for sharing)
		api.GET("/recipes/shared/:id", recipeHandler.GetPublicRecipe)

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// Auth routes (requires authentication)
			protected.POST("/auth/logout", authHandler.Logout)


			// User routes
			protected.GET("/user/profile", authHandler.GetProfile)
			protected.PUT("/user/profile", authHandler.UpdateProfile)
			protected.GET("/user/stats", authHandler.GetStats)

			// Recipe routes
			recipes := protected.Group("/recipes")
			{
				recipes.POST("/generate", recipeHandler.GenerateRecipe)
				recipes.GET("", recipeHandler.GetRecipes)
				recipes.GET("/:id", recipeHandler.GetRecipe)
				recipes.DELETE("/:id", recipeHandler.DeleteRecipe)
				recipes.POST("/:id/favorite", recipeHandler.ToggleFavorite)
			}

			// Filter options routes
			filters := protected.Group("/filters")
			{
				filters.GET("/countries", recipeHandler.GetCountries)
				filters.GET("/meats", recipeHandler.GetMeatTypes)
				filters.GET("/ingredients", recipeHandler.GetIngredients)
			}

			// Shopping list routes
			shoppingList := protected.Group("/shopping-list")
			{
				shoppingList.GET("", shoppingListHandler.GetShoppingList)
				shoppingList.POST("/add-recipe", shoppingListHandler.AddRecipeToShoppingList)
				shoppingList.POST("/:id/toggle", shoppingListHandler.ToggleItemChecked)
				shoppingList.DELETE("/:id", shoppingListHandler.DeleteItem)
				shoppingList.DELETE("/clear/checked", shoppingListHandler.ClearCheckedItems)
				shoppingList.DELETE("/clear/all", shoppingListHandler.ClearAllItems)
			}

			// Admin routes (requires admin privileges)
			admin := protected.Group("/admin")
			admin.Use(middleware.AdminOnly())
			{
				admin.GET("/users", adminHandler.GetAllUsers)
				admin.DELETE("/users/:id", adminHandler.DeleteUser)
				admin.GET("/stats", adminHandler.GetAdminStats)
				admin.PUT("/users/:id/recipe-limit", adminHandler.UpdateUserRecipeLimit)
			}
		}
	}

	// Health check (must be before NoRoute)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Serve uploaded images (recipe images, thumbnails)
	router.Static("/uploads", "./uploads")

	// Serve embedded frontend for SPA routing
	// Note: Frontend is always embedded in the binary after build
	distFS, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		log.Fatalf("Failed to get frontend dist: %v", err)
	}

	// Serve static assets (JS, CSS, images, etc.) from the assets subdirectory
	assetsFS, err := fs.Sub(distFS, "assets")
	if err != nil {
		log.Fatalf("Failed to get assets directory: %v", err)
	}
	router.StaticFS("/assets", http.FS(assetsFS))

	// Serve other static files (like vite.svg)
	router.GET("/vite.svg", func(c *gin.Context) {
		c.FileFromFS("vite.svg", http.FS(distFS))
	})

	// Custom NoRoute handler for SPA routing
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// If path starts with /api, return JSON 404 (API route not found)
		if len(path) >= 4 && path[:4] == "/api" {
			c.JSON(404, gin.H{
				"error":   "Not Found",
				"message": "API endpoint not found",
				"path":    path,
			})
			return
		}

		// For all other routes, serve the frontend SPA (React Router handles it)
		// Serve index.html for all non-API routes (SPA will handle routing)
		indexHTML, err := distFS.Open("index.html")
		if err != nil {
			c.String(500, "Failed to load application")
			return
		}
		defer indexHTML.Close()

		stat, err := indexHTML.Stat()
		if err != nil {
			c.String(500, "Failed to load application")
			return
		}

		c.DataFromReader(200, stat.Size(), "text/html; charset=utf-8", indexHTML, nil)
	})

	// Start server
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Chefly server starting on port %s\n", port)
	fmt.Printf("📊 Database: %s\n", cfg.DBPath)
	fmt.Printf("🌍 Environment: %s\n", cfg.Environment)
	fmt.Printf("📝 Audit Logging: %v (level: %s, format: %s)\n", cfg.AuditLogEnabled, cfg.AuditLogLevel, cfg.AuditLogFormat)
	fmt.Printf("👥 Registration: %v\n", cfg.RegistrationEnabled)
	fmt.Printf("🍳 Recipe Generation Limit: %s\n", cfg.RecipeGenerationLimit)
	fmt.Printf("📌 Version: %s (built: %s)\n", Version, BuildTime)

	// Create HTTP server with timeouts
	// Note: WriteTimeout must be long enough for AI recipe generation (typically 30-60 seconds)
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 120 * time.Second, // Allow up to 2 minutes for AI recipe generation
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🌐 Server starting on port %s", port)
		log.Printf("📱 Application: http://localhost:%s", port)
		log.Printf("🏥 Health check: http://localhost:%s/health", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal (Ctrl+C or container stop)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Graceful shutdown with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("⚠️  Server forced to shutdown: %v", err)
	}

	// Close database connection
	if err := db.Close(); err != nil {
		log.Printf("⚠️  Error closing database: %v", err)
	}

	log.Println("✅ Server stopped cleanly")
}

// cleanupExpiredTokens runs in a goroutine and cleans up expired/revoked tokens every hour
func cleanupExpiredTokens(db *sql.DB, logger *services.AuditLogger) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	// Run cleanup immediately on startup
	performCleanup(db, logger)

	// Then run every hour
	for range ticker.C {
		performCleanup(db, logger)
	}
}

// performCleanup deletes expired and revoked refresh tokens
func performCleanup(db *sql.DB, logger *services.AuditLogger) {
	result, err := db.Exec(`
		DELETE FROM refresh_tokens
		WHERE revoked = 1 OR expires_at < datetime('now')
	`)

	if err != nil {
		if logger != nil {
			logger.Error("token.cleanup.failed", "Failed to cleanup expired tokens", err, nil)
		}
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 && logger != nil {
		logger.Info("token.cleanup.success", fmt.Sprintf("Cleaned up %d expired/revoked tokens", rowsAffected), nil)
	}
}
