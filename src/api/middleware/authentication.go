package middleware

import (
	"github.com/gin-gonic/gin"
)

// basic authenication middleware
func BasicAuthMiddleware(r *gin.Engine) {
	r.Use(gin.BasicAuth(gin.Accounts{
		"admin": "admin",
	}))
}
