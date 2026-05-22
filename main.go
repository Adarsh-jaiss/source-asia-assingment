package main

import (
	"log"

	_ "github.com/adarsh-jaiss/assingment/docs" // Import generated swagger docs
	"github.com/adarsh-jaiss/assingment/api/controllers"

	"github.com/adarsh-jaiss/assingment/api/repository"
	"github.com/adarsh-jaiss/assingment/routes"
)

// @title Source Asia Backend API
// @version 1.0
// @description Rate-limited API and Product Catalog for the Source Asia assignment.
// @host localhost:8080
// @BasePath /
func main() {
	// Initialize repositories
	rateLimitRepo := repository.NewRateLimitRepo()
	productRepo := repository.NewProductRepo()

	// Initialize controllers
	rateLimitCtrl := controllers.NewRateLimitController(rateLimitRepo)
	productCtrl := controllers.NewProductController(productRepo)

	// Setup routes
	ginEngine := routes.Routes(rateLimitCtrl, productCtrl)

	log.Println("Starting server on :8080")
	if err := ginEngine.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
