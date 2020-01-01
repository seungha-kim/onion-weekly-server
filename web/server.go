package web

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.uber.org/fx"

	"github.com/labstack/echo"
	em "github.com/labstack/echo/middleware"
	"github.com/onion-studio/onion-weekly/config"
)

type Server struct {
	echo    *echo.Echo
	appConf config.AppConf
}

func NewServer(
	lc fx.Lifecycle,
	appConf config.AppConf,
	helloHandler *helloHandler,
	authHandler *authHandler,
) *Server {
	s := &Server{echo: echo.New(), appConf: appConf}
	s.echo.Pre(em.RemoveTrailingSlashWithConfig(em.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))
	s.echo.Use(em.Logger())
	s.echo.Use(em.Recover())
	authHandler.Register(s.echo.Group("/auth"))
	helloHandler.Register(s.echo.Group("/hello"))
	s.echo.Debug = appConf.Debug

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Fatal(s.echo.Start(fmt.Sprintf(":%d", appConf.Port)))
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return s.echo.Shutdown(ctx)
		},
	})

	return s
}
