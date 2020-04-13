package api

import (
	"api/handlers"

	"github.com/labstack/echo"
)

// UserGroup user endpoints
func UserGroup(u *echo.Group) {
	u.GET("/main", handlers.UserMain)
}
