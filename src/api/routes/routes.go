package routes

import (
	"github.com/gin-gonic/gin"

	controller "prospector/controllers"
	"prospector/middleware"

	_ "prospector/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func CreateRoutes(r *gin.Engine) {

	c := controller.Controller{Client: &controller.DefaultNomadClient{}}

	api := r.Group("/api")
	{
		api.GET("/health", c.Health)
		api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		api.GET("/jobs", c.GetJobs)
		api.POST("/jobs", c.CreateJob)
	}

	authenticated := r.Group("/api/v1")
	authenticated.Use(middleware.AuthenticationMiddleware())
	{
		authenticated.GET("/auth", c.AuthHealth)
	}
}
