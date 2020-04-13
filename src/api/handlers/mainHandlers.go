package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
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

// Root root
func Root(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// GetCat get
func GetCat(c echo.Context) error {
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

// AddCat add
func AddCat(c echo.Context) error { //fastest
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

// AddCat2 add2
func AddCat2(c echo.Context) error { //almost as fast,preferrable
	cat := Cat{}
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&cat)
	if err != nil {
		log.Printf("error is :%s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, cat)
}

// AddCat3 add3
func AddCat3(c echo.Context) error { //slower than the other two
	cat := Cat{}
	err := c.Bind(&cat)
	if err != nil {
		log.Printf("error is :%s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, cat)
}

//JwtPage jwt page
func JwtPage(c echo.Context) error {
	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	log.Println("username:", claims["name"], ",user id:", claims["jti"], ",raw token:", token.Raw)
	return c.String(http.StatusOK, "hello,this is jwt page.")
}

// Login login
func Login(c echo.Context) error {
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

// LoginPage login page
func LoginPage(c echo.Context) error {
	redirect := c.QueryParam("redirect")
	return c.HTML(http.StatusOK, fmt.Sprintf(`<form target="/login" method="POST"><input type="text" name="username"/><br><input type="password" name="password"/><input type="hidden" value="%s" name="redirect"/><input type="submit" /></form>`, redirect))
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
