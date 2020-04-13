package middlewares

import (
	"log"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// SetMainMiddlewares set main middlewares
func SetMainMiddlewares(e *echo.Echo) {
	// static middleware
	e.Static("/", "static")
	// logger middleware
	f, err := os.Create("test.log") // import "os"
	if err != nil {
		log.Printf("error:%s", err)
	}
	defer f.Close()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${status} ${method} ${host}${path} ${latency_human}` + "\n",
		Output: f,
	}))

	///// custom middlewares //////
	e.Use(serverHeader)
}

// serverHeader add header to all responses
func serverHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "server kin")
		return next(c)
	}
}
