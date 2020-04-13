package api

import (
	"api/handlers"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// MainGroup main endpoints
func MainGroup(e *echo.Echo) {
	e.GET("/root", handlers.Root)
	e.GET("/login", handlers.LoginPage)
	e.POST("/login", handlers.Login)
	e.GET("/cats/:datatype", handlers.GetCat)
	e.POST("/cat", handlers.AddCat)
	e.POST("/cat2", handlers.AddCat2)
	e.POST("/cat3", handlers.AddCat3)
	e.GET("/jwt", handlers.JwtPage, middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte("secret"),
		TokenLookup:   "cookie:JWTCookie", //to read from cookies, to use this, we should write cookie when user signed in.
	}))
}
