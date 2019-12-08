package web

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/onion-studio/onion-weekly/domain"
	"net/http"
)

func RegisterAuthHandlers(g *echo.Group) {
	g.POST("/register", handlePostTestUser)
	g.POST("/session", handlePostTestToken)
	// TODO: JWTWithConfig
	g.GET("/token", handleGetTokenPayload, middleware.JWT([]byte("mysecret")))
}

func handlePostTestUser(c echo.Context) (err error) {
	input := domain.InputCreateUser{}
	if err = c.Bind(&input); err != nil {
		return // TODO
	}
	user, _, _, err := domain.CreateUserWithEmailCredential(input)
	if err != nil {
		switch err.(type) {
		case domain.DuplicateError:
			err = echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return
	}
	return c.JSON(200, user)
}

func handlePostTestToken(c echo.Context) (err error) {
	input := domain.InputCreatTokenByEmailCredential{}
	if err = c.Bind(&input); err != nil {
		return // TODO
	}
	outputToken, err := domain.CreateTokenByEmailCredential(input)
	if err != nil {
		switch err.(type) {

		}
		return
	}
	switch err.(type) {
	case domain.DuplicateError:
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(200, outputToken)
}

func handleGetTokenPayload(c echo.Context) (err error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return c.JSON(200, claims)
}
