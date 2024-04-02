package middleware

import (
	"fmt"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type User struct {
	Username string
	FistName string
	LastName string
	Email    string
}

type AuthSucess struct {
	Token  string `json:"token"`
	Expire string `json:"expire"`
	Code   int    `json:"code"`
}

func AuthMiddleware(identityKey string) *jwt.GinJWTMiddleware {
	jwtMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "prospector",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour * 8,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				Username: claims[identityKey].(string),
			}
		},
		Authenticator: Authenticate,
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*User); ok && (v.Username != "test") {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",

		TimeFunc: time.Now,
	})

	if err != nil {
		fmt.Println("error creating jwt middleware")
	}

	return jwtMiddleware
}

func Authenticate(c *gin.Context) (interface{}, error) {
	var loginVals login
	if err := c.ShouldBind(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	username := loginVals.Username
	password := loginVals.Password

	println("Authenticating with local")

	if (username == "admin" && password == "admin") || (username == "test" && password == "test") {
		return &User{
			Username: username,
			FistName: "Test",
			LastName: "User",
		}, nil
	}

	println("Authenticating with LDAP")

	user, err := AuthenticateLdap(username, password)
	if err != nil {
		return nil, err
	}

	return user, nil
}
