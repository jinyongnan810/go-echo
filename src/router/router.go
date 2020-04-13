package router

import (
	"api"
	"api/middlewares"

	"github.com/labstack/echo"
)

// New new echo
func New() *echo.Echo {
	e := echo.New()

	// create groups
	a := e.Group("/admin")
	u := e.Group("/user")

	// set middlewares
	middlewares.SetMainMiddlewares(e)
	middlewares.SetUserMiddlewares(u)
	middlewares.SetAdminMiddlewares(a)

	// set endpoints
	api.MainGroup(e)
	api.UserGroup(u)
	api.AdminGroup(a)

	return e
}
