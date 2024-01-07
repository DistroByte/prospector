package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *Controller) Authenticate(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "authenticate",
	})
}
