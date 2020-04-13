package middlewares

import (
	"errors"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// SetAdminMiddlewares set admin middlewares
func SetAdminMiddlewares(a *echo.Group) {
	// basic auth
	a.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "kin" && password == "jinyongnan" {
			return true, nil
		}
		return false, errors.New("failed")
	}))
}
