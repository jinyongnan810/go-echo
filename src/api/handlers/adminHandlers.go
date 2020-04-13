package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

// AdminMain admin main page
func AdminMain(c echo.Context) error {
	return c.String(http.StatusOK, "hello,this is admin main page.")
}
