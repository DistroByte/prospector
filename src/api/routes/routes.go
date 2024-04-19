package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	controller "prospector/controllers"
	"prospector/middleware"

	docs "prospector/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Route(r *gin.Engine, identityKey string) {

	c := controller.Controller{
		Client: &controller.DefaultNomadClient{
			URL: "http://zeus.internal:4646/v1",
		},
		IdentityKey:   identityKey,
		JWTMiddleware: middleware.AuthMiddleware(identityKey),
	}

	cProxy := controller.NomadProxyController{
		Client: &controller.DefaultNomadClient{
			URL: "http://zeus.internal:4646/v1",
		},
	}

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	r.Use(cors.New(config))

	docs.SwaggerInfo.Title = "Prospector API"
	docs.SwaggerInfo.Description = "API backend for Prospector"
	docs.SwaggerInfo.Version = "1.0"
	if gin.Mode() == gin.ReleaseMode {
		docs.SwaggerInfo.Host = "prospector.ie"
		docs.SwaggerInfo.Schemes = []string{"https"}
	} else {
		docs.SwaggerInfo.Host = "localhost:3434"
		docs.SwaggerInfo.Schemes = []string{"http"}
	}
	docs.SwaggerInfo.BasePath = "/api"

	//	@securityDefinitions.apikey	BearerAuth
	//	@in							header
	//	@name						Authorization

	api := r.Group("/api")
	{
		api.GET("/health", c.Health)
		api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		api.POST("/login", c.Login)

		// serve static files from ./vm-config
		api.GET("/vm-config/*filepath", func(c *gin.Context) {
			filepath := c.Param("filepath")
			c.File("vm-config/" + filepath)
		})
	}

	r.NoRoute(c.JWTMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found", "data": claims})
	})

	authenticated := r.Group("/api/v1")
	authenticated.GET("/refresh", c.RefreshToken)
	authenticated.Use(c.JWTMiddleware.MiddlewareFunc())
	{
		authenticated.GET("/user", c.GetUserName)

		jobs := authenticated.Group("/jobs")
		jobs.GET("", c.GetJobs)
		jobs.POST("", c.CreateJob)
		jobs.GET("/:id", c.GetJob)
		jobs.DELETE("/:id", c.DeleteJob)
		jobs.GET("/:id/components", c.GetComponents)
		jobs.GET("/:id/logs", cProxy.StreamLogs)
		jobs.PUT("/:id/restart", c.RestartJob)
		jobs.POST("/:id/start", c.StartJob)
		jobs.PUT("/:id/component/:component/restart", c.RestartAlloc)

		resources := authenticated.Group("/resources")
		resources.GET("", c.GetUserUsedResources)
		resources.GET("/allocated", c.GetUserAllocatedResources)
		resources.GET("/:id", c.GetJobUsedResources)
		resources.GET("/:id/allocated", c.GetJobAllocatedResources)
		resources.GET("/:id/:component", c.GetComponentUsedResources)
		resources.GET("/:id/:component/allocated", c.GetComponentAllocatedResources)
	}
}
