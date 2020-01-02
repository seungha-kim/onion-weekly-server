package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"github.com/onion-studio/onion-weekly/config"
	"github.com/onion-studio/onion-weekly/db"
	"github.com/onion-studio/onion-weekly/domain"
	"go.uber.org/fx"
)

var pgxPool *pgxpool.Pool
var h *authHandler

func TestMain(m *testing.M) {
	fx.New(
		fx.Provide(
			db.NewPgxPool,
			config.NewTestAppConf,
			domain.NewUserService,
			NewAuthHandler,
		),
		fx.Populate(&h, &pgxPool),
		fx.NopLogger,
	)
	m.Run()
}

func TestAuthHandler_handlePostTestUser(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{}"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	db.RollbackForTest(pgxPool, func(tx pgx.Tx) {
		c.Set("tx", tx)
		if assert.NoError(t, h.handlePostUser(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			//assert.Equal(t, "", rec.Body.String())
		}
	})
}

func TestAuthHandler_handleAuth(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{}"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	db.RollbackForTest(pgxPool, func(tx pgx.Tx) {
		c.Set("tx", tx)

		if assert.NoError(t, h.handlePostUser(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			//assert.Equal(t, "", rec.Body.String())
		}
	})
}
