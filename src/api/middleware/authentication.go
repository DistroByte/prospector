package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shaj13/go-guardian/v2/auth"
	"github.com/shaj13/go-guardian/v2/auth/strategies/basic"
	"github.com/shaj13/go-guardian/v2/auth/strategies/ldap"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
	"github.com/shaj13/libcache"

	_ "github.com/shaj13/libcache/fifo"
)

var basicStrategy auth.Strategy
var strategy union.Union
var cache libcache.Cache

func SetupGoGuardian() {
	cfg := &ldap.Config{
		BaseDN:       "dc=example,dc=com",
		BindDN:       "cn=read-only-admin,dc=example,dc=com",
		Port:         "389",
		Host:         "ldap.forumsys.com",
		BindPassword: "password",
		Filter:       "(uid=%s)",
	}

	cache = libcache.FIFO.New(0)
	cache.SetTTL(time.Minute * 10)

	ldapStrategy := ldap.NewCached(cfg, cache)
	basicStrategy = basic.NewCached(validateBasicUser, cache)
	strategy = union.New(basicStrategy, ldapStrategy)
}

func validateBasicUser(ctx context.Context, r *http.Request, username, password string) (auth.Info, error) {
	fmt.Println("validateBasicUser", username, password)

	if username == "admin" && password == "admin" {
		return auth.NewDefaultUser("medium", "1", nil, nil), nil
	}

	return nil, fmt.Errorf("invalid basic auth credentials")
}

func AuthenticationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, user, err := strategy.AuthenticateRequest(ctx.Request)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid credentials",
			})
			return
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}
