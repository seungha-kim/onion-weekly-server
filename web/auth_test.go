package web

import (
	"testing"

	"github.com/onion-studio/onion-weekly/config"
	"github.com/onion-studio/onion-weekly/db"
	"github.com/onion-studio/onion-weekly/domain"
	"go.uber.org/fx"
)

var h *authHandler

func TestMain(m *testing.M) {
	fx.New(
		fx.Provide(
			db.NewPgxPool,
			config.NewTestAppConf,
			domain.NewUserService,
			NewAuthHandler,
		),
		fx.Populate(&h),
	)
	m.Run()
}

func TestAuthHandler_handlePostTestUser(t *testing.T) {
	//fmt.Printf("authHandler %v\n", h)
	//e := echo.New()
	//req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{}"))
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	//rec := httptest.NewRecorder()
	//c := e.NewContext(req, rec)
	//group := e.Group("/auth")
	//h.Register(group)
	//
	//if assert.NoError(t, h.handlePostUser(c)) {
	//	assert.Equal(t, http.StatusCreated, rec.Code)
	//	assert.Equal(t, "", rec.Body.String())
	//}
	// 포기! handler를 직접 호출하면 미들웨어 체인을 통과 못하니까,
	// group에다가 register하고 group가지고 테스트를 해보고싶은데 어떻게 하는 건지 모르겠다.
}
