package web

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/onion-studio/onion-weekly/config"
	"github.com/onion-studio/onion-weekly/domain"
	m "github.com/onion-studio/onion-weekly/web/middleware"
)

func RegisterAuthHandlers(g *echo.Group) {
	g.Use(m.Transaction)
	g.POST("/register", handlePostTestUser)
	g.POST("/session", handlePostTestToken)
	// TODO: JWTWithConfig
	g.GET("/token", handleGetTokenPayload, middleware.JWT([]byte("mysecret")))
}

func handlePostTestUser(c echo.Context) (err error) {
	appConf := c.Get("appConf").(config.AppConf)
	tx := c.Get("tx").(pgx.Tx)

	input := domain.InputCreateUser{}
	if err = c.Bind(&input); err != nil {
		return // TODO
	}

	user, _, _, err := domain.CreateUserWithEmailCredential(appConf, tx, input)
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
	appConf := c.Get("appConf").(config.AppConf)
	tx := c.Get("tx").(pgx.Tx)

	input := domain.InputCreatTokenByEmailCredential{}
	if err = c.Bind(&input); err != nil {
		return // TODO
	}

	outputToken, err := domain.CreateTokenByEmailCredential(appConf, tx, input)
	if err != nil {
		switch err.(type) {

		}
		return
	}
	return c.JSON(200, outputToken)
}

func handleGetTokenPayload(c echo.Context) (err error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return c.JSON(200, claims)
}
