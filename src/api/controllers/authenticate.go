package controllers

import (
	"prospector/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type GetUserNameResponse struct {
	UserID   string `json:"userID"`
	UserName string `json:"userName"`
}

// GetUserName endpoint
//
//	@Summary		Get user name
//	@Description	Get user name
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	GetUserNameResponse
//	@Router			/v1/user [get]
func (c *Controller) GetUserName(ctx *gin.Context) {
	claims := jwt.ExtractClaims(ctx)
	user, _ := ctx.Get(c.IdentityKey)
	ctx.JSON(200, gin.H{
		"userID":   claims[c.IdentityKey],
		"userName": user.(*middleware.User).Username,
	})
}

// Login endpoint
//
//	@Summary		Login
//	@Description	Login
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Body			{username: string, password: string}
//	@Success		200	{string}	string	middleware.AuthSucess
//	@Router			/login [post]
func (c *Controller) Login(ctx *gin.Context) {
	// pass to login middleware
	c.JWTMiddleware.LoginHandler(ctx)
}

// RefreshToken endpoint
//
//	@Summary		Refresh token
//	@Description	Refresh token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	string	middleware.AuthSucess
//	@Router			/v1/refresh [get]
func (c *Controller) RefreshToken(ctx *gin.Context) {
	c.JWTMiddleware.RefreshHandler(ctx)
}
