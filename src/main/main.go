package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Cat : the cat
type Cat struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// JWTClaims create jwt token
type JWTClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

func main() {
	e := echo.New()
	e.Use(ServerHeader)
	g := e.Group("/admin")
	u := e.Group("/user")
	u.Use(CheckLogin)
	e.GET("/jwt", jwtPage, middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte("secret"),
		TokenLookup:   "cookie:JWTCookie", //to read from cookies, to use this, we should write cookie when user signed in.
	}))

	f, err := os.Create("test.log") // import "os"
	if err != nil {
		log.Printf("error:%s", err)
	}
	defer f.Close()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${status} ${method} ${host}${path} ${latency_human}` + "\n",
		Output: f,
	}))
	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "kin" && password == "jinyongnan" {
			return true, nil
		}
		return false, errors.New("failed")
	}))
	g.GET("/main", adminMain)
	u.GET("/main", userMain)

	e.GET("/", root)
	e.GET("/login", loginPage)
	e.POST("/login", login)
	e.GET("/cats/:datatype", getCat)
	e.POST("/cat", addCat)
	e.POST("/cat2", addCat2)
	e.POST("/cat3", addCat3)
	e.Logger.Fatal(e.Start(":1323"))
}
func root(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
func getCat(c echo.Context) error {
	catName := c.QueryParam("name")
	catType := c.QueryParam("type")
	dataType := c.Param("datatype")
	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("name:%s,type:%s", catName, catType))
	} else if dataType == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"name": catName,
			"type": catType,
		})
	}
	return c.JSON(http.StatusBadRequest, map[string]string{
		"error": "datatype not specified",
	})
}
func addCat(c echo.Context) error { //fastest
	cat := Cat{}
	defer c.Request().Body.Close()
	b, err := ioutil.ReadAll(c.Request().Body) //import "io/ioutil"
	if err != nil {
		log.Printf("error is :%s", err)
		return c.String(http.StatusInternalServerError, "")
	}
	err = json.Unmarshal(b, &cat)
	if err != nil {
		log.Printf("error is :%s", err)
		return c.String(http.StatusInternalServerError, "")
	}
	return c.JSON(http.StatusOK, cat)
}
func addCat2(c echo.Context) error { //almost as fast,preferrable
	cat := Cat{}
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&cat)
	if err != nil {
		log.Printf("error is :%s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, cat)
}
func addCat3(c echo.Context) error { //slower than the other two
	cat := Cat{}
	err := c.Bind(&cat)
	if err != nil {
		log.Printf("error is :%s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, cat)
}

func adminMain(c echo.Context) error {
	return c.String(http.StatusOK, "hello,this is admin main page.")
}
func userMain(c echo.Context) error {
	return c.String(http.StatusOK, "hello,this is user main page.")
}
func jwtPage(c echo.Context) error {
	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	log.Println("username:", claims["name"], ",user id:", claims["jti"], ",raw token:", token.Raw)
	return c.String(http.StatusOK, "hello,this is jwt page.")
}
func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	redirect := c.FormValue("redirect")
	if redirect == "" {
		redirect = "/"
	}
	if username == "kin" && password == "jinyongnan" {
		cookie := new(http.Cookie)
		cookie.Name = "SessionID"
		cookie.Value = "some hash"
		cookie.Expires = time.Now().Add(2 * time.Hour)
		c.SetCookie(cookie)
		token, err := createJWTToken()
		if err != nil {
			log.Println("error:", err)
			return c.String(http.StatusInternalServerError, "token create failed")
		}
		jwtCookie := new(http.Cookie)
		jwtCookie.Name = "JWTCookie"
		jwtCookie.Value = token
		jwtCookie.Expires = time.Now().Add(2 * time.Hour)
		c.SetCookie(jwtCookie)
		return c.JSON(http.StatusOK, map[string]string{
			"message": "login success",
			"token":   token,
		})
		// return c.Redirect(http.StatusMovedPermanently, redirect)
	}
	return c.JSON(http.StatusOK, map[string]string{
		"error": "failed",
	})
}
func loginPage(c echo.Context) error {
	redirect := c.QueryParam("redirect")
	return c.HTML(http.StatusOK, fmt.Sprintf(`<form target="/input" method="POST"><input type="text" name="username"/><br><input type="password" name="password"/><input type="hidden" value="%s" name="redirect"/><input type="submit" /></form>`, redirect))
}

func createJWTToken() (string, error) { //import jwt "github.com/dgrijalva/jwt-go"
	claims := JWTClaims{
		"test",
		jwt.StandardClaims{
			Id:        "user_id",
			ExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
		},
	}
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token, err := rawToken.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return token, err
}

////////////middlewares////////////

// ServerHeader add header to all responses
func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "server kin")
		return next(c)
	}
}

// CheckLogin check user login
func CheckLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("SessionID")
		if err == nil && cookie != nil && cookie.Value == "some hash" {
			return next(c)
		}
		return c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/login?redirect=%s", c.Path()))
	}
}
