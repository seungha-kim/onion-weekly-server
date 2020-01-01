package web

import (
	"net/http"

	"github.com/onion-studio/onion-weekly/config"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo"
	em "github.com/labstack/echo/middleware"
	"github.com/onion-studio/onion-weekly/domain"
	"github.com/onion-studio/onion-weekly/web/middleware"
)

type authHandler struct {
	appConf     config.AppConf
	pgxPool     *pgxpool.Pool
	userService *domain.UserService
}

func NewAuthHandler(
	appConf config.AppConf,
	pgxPool *pgxpool.Pool,
	userService *domain.UserService,
) *authHandler {
	return &authHandler{
		appConf:     appConf,
		pgxPool:     pgxPool,
		userService: userService,
	}
}

func (h *authHandler) Register(g *echo.Group) {
	g.Use(middleware.Transaction(h.appConf, h.pgxPool))
	g.POST("", h.handleAuth)
	g.POST("/register", h.handlePostUser)
	// TODO: JWTWithConfig
	g.GET("/token", h.handleGetTokenPayload, em.JWT([]byte("mysecret")))
}

func (h *authHandler) handlePostUser(c echo.Context) (err error) {
	tx := c.Get("tx").(pgx.Tx)

	// TODO: recaptcha
	input := domain.InputCreateUser{}
	if err = c.Bind(&input); err != nil {
		return // TODO
	}

	user, _, _, err := h.userService.CreateUserWithEmailCredential(tx, input)
	if err != nil {
		switch err.(type) {
		case domain.DuplicateError:
			err = echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return
	}
	return c.JSON(200, user)
}

func (h *authHandler) handleAuth(c echo.Context) (err error) {
	tx := c.Get("tx").(pgx.Tx)

	input := domain.InputCreatTokenByEmailCredential{}
	if err = c.Bind(&input); err != nil {
		return // TODO
	}

	output, err := h.userService.CreateTokenByEmailCredential(tx, input)
	if err != nil {
		return
	}
	return c.JSON(200, output)
}

func (h *authHandler) handleGetTokenPayload(c echo.Context) (err error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return c.JSON(200, claims)
}
