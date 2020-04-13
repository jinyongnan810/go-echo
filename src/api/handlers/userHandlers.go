package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

// UserMain user main page
func UserMain(c echo.Context) error {
	return c.String(http.StatusOK, "hello,this is user main page.")
}
