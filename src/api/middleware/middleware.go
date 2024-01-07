package middleware

import (
	"github.com/gin-gonic/gin"
)

// CreateMiddlewares creates the middlewares for the server
func CreateStandardMiddlewares(r *gin.Engine) {
	// Add the metrics middleware
	MetricsMiddleware(r)
}

func CreateAuthMiddlewares(r *gin.Engine) {
	// Setup the authentication middleware
	SetupGoGuardian()
}
