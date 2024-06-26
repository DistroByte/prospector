package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health endpoint
//
//	@Summary		Health check
//	@Description	Check if the API is up and running
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	string	"OK"
//	@Router			/health [get]
func (c *Controller) Health(ctx *gin.Context) {
	ctx.String(http.StatusOK, "OK")
}
