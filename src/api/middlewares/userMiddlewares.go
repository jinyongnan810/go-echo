package middlewares

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

// SetUserMiddlewares set user middlewares
func SetUserMiddlewares(u *echo.Group) {
	// check login
	u.Use(checkLogin)
}

// checkLogin check user login
func checkLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("SessionID")
		if err == nil && cookie != nil && cookie.Value == "some hash" {
			return next(c)
		}
		return c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/login?redirect=%s", c.Path()))
	}
}
