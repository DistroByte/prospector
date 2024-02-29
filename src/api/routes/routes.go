package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	controller "prospector/controllers"
	"prospector/middleware"

	docs "prospector/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func CreateRoutes(r *gin.Engine) {

	c := controller.Controller{Client: &controller.DefaultNomadClient{}}

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	docs.SwaggerInfo.Title = "Prospector API"
	docs.SwaggerInfo.Description = "API backend for Prospector"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "prospector.ie"
	docs.SwaggerInfo.Schemes = []string{"https"}
	docs.SwaggerInfo.BasePath = "/api"

	api := r.Group("/api")
	{
		api.GET("/health", c.Health)
		api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		api.GET("/jobs", c.GetJobs)
		api.GET("/jobs/:id", c.GetJob)
		api.POST("/jobs", c.CreateJob)
		api.DELETE("/jobs/:id", c.DeleteJob)
	}

	// @securityDefinitions.basic	BasicAuth
	authenticated := r.Group("/api/v1")
	authenticated.Use(middleware.AuthenticationMiddleware())
	{
		authenticated.GET("/auth", c.AuthHealth)
	}
}
