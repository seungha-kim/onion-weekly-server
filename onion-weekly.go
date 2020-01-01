package main

import (
	"github.com/onion-studio/onion-weekly/config"
	"github.com/onion-studio/onion-weekly/db"
	"github.com/onion-studio/onion-weekly/domain"
	"github.com/onion-studio/onion-weekly/web"
	"go.uber.org/fx"
)

func main() {
	appConf := config.NewAppConf()
	opts := []fx.Option{
		fx.Provide(
			func() config.AppConf { return appConf },
			db.NewPgxPool,
			web.NewServer,
			web.NewAuthHandler,
			web.NewHelloHandler,
			domain.NewUserService,
		),
		fx.Invoke(func(_ *web.Server) {}),
	}
	if !appConf.Debug {
		opts = append(opts, fx.NopLogger)
	}
	app := fx.New(opts...)
	app.Run()
}
