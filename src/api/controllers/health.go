package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary		Health check
// @Description	Check if the API is up and running
// @Tags			health
// @Accept			json
// @Produce		json
// @Security		None
// @Success		200	{object}	Message
// @Router			/health [get]
func (c *Controller) Health(ctx *gin.Context) {
	ctx.String(http.StatusOK, "OK\n")
}

// @Summary		Authenticated Health check
// @Description	Check if the API is up and running
// @Tags			health
// @Accept			json
// @Produce		json
// @Security		BasicAuth
// @Success		200	{object}	Message
// @Router			/v1/auth [get]
func (c *Controller) AuthHealth(ctx *gin.Context) {
	ctx.String(http.StatusOK, "OK\n")
}
