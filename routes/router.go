package routes

import (
	"github.com/adarsh-jaiss/assingment/api/controllers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Routes(rateLimitCtrl *controllers.RateLimitController, productCtrl *controllers.ProductController) *gin.Engine {
	r := gin.Default()
	
	// Swagger Documentation
	r.GET("/api/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")


	api.POST("/request", rateLimitCtrl.HandleRequest)
	api.GET("/stats", rateLimitCtrl.GetStats)

	api.POST("/products", productCtrl.CreateProduct)
	api.GET("/products", productCtrl.GetProducts)
	api.GET("/products/:id", productCtrl.GetProductByID)
	api.POST("/products/:id/media", productCtrl.AppendMedia)

	return r
}
