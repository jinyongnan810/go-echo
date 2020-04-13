package api

import (
	"api/handlers"

	"github.com/labstack/echo"
)

// AdminGroup admin endpoints
func AdminGroup(a *echo.Group) {
	a.GET("/main", handlers.AdminMain)
}
