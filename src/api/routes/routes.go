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

	c := controller.Controller{Client: &controller.DefaultNomadClient{}, IdentityKey: identityKey, JWTMiddleware: middleware.AuthMiddleware(identityKey)}

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	docs.SwaggerInfo.Title = "Prospector API"
	docs.SwaggerInfo.Description = "API backend for Prospector"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "prospector.ie"
	docs.SwaggerInfo.Schemes = []string{"https"}
	docs.SwaggerInfo.BasePath = "/api"

	//	@securityDefinitions.apikey	BearerAuth
	//	@in							header
	//	@name						Authorization

	api := r.Group("/api")
	{
		api.GET("/health", c.Health)
		api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		api.POST("/login", c.Login)
		api.GET("/refresh", c.RefreshToken)
	}

	println("identityKey: %#v\n", identityKey)

	r.NoRoute(c.JWTMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		println("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	println("valid identityKey: %#v\n", identityKey)

	authenticated := r.Group("/api/v1")
	authenticated.Use(c.JWTMiddleware.MiddlewareFunc())
	{
		authenticated.GET("/user", c.GetUserName)
		authenticated.GET("/jobs", c.GetJobs)
		authenticated.GET("/jobs/:id", c.GetJob)
		authenticated.POST("/jobs", c.CreateJob)
		authenticated.DELETE("/jobs/:id", c.DeleteJob)
	}
}
