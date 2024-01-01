package routes

import (
	"github.com/gin-gonic/gin"

	controller "prospector/controllers"
	middleware "prospector/middleware"

	_ "prospector/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func CreateRoutes(r *gin.Engine, c *controller.Controller) {
	api := r.Group("/api")
	{
		api.GET("/health", c.Health)
		api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	middleware.BasicAuthMiddleware(r)

	authenticated := r.Group("/api/v1")
	{
		authenticated.GET("/auth", c.AuthHealth)
	}
}
