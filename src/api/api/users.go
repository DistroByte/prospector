package api

import (
	"github.com/labstack/echo/v4"
)

type User struct {
	Name string `json:"name"`
}

// UserGetHandler godoc
//
//	@Summary		Get all users
//	@Description	Get all users
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]User
//	@Router			/users [get]
func UserGetHandler(c echo.Context) error {
	return c.String(200, "You got a user!")
}
