package main

import (
	"log"
	"point-prevalence-survey/config"
	"point-prevalence-survey/database"
	"point-prevalence-survey/routes"

	"github.com/gin-gonic/gin"
)

// @title Point Prevalence Survey API
// @version 1.0
// @description API for managing Point Prevalence Survey data with PostgreSQL integration
// @host localhost:8080
// @BasePath /
func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	database.InitDB()
	database.Migrate()

	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Setup routes
	routes.SetupRoutes(r)

	// Start server
	port := cfg.GetPort()
	log.Printf("Server starting on port %s", port)
	log.Printf("Swagger documentation available at http://localhost%s/swagger/index.html", port)
	
	if err := r.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
