package api

import (
	"github.com/labstack/echo/v4"
)

type Job struct {
	Name string `json:"name"`
}

// JobGetHandler godoc
//
//	@Summary		Get all jobs
//	@Description	Get all jobs
//	@Tags			jobs
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]Job
//	@Router			/jobs [get]
func JobGetHandler(c echo.Context) error {
	return c.String(200, "You got a job!")
}
