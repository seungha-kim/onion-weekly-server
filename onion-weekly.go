package main

import (
	"github.com/onion-studio/onion-weekly/config"
	"github.com/onion-studio/onion-weekly/db"
	"github.com/onion-studio/onion-weekly/domain"
	"github.com/onion-studio/onion-weekly/web"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			config.NewAppConf,
			db.NewPgxPool,
			web.NewServer,
			web.NewAuthHandler,
			web.NewHelloHandler,
			domain.NewUserService,
		),
		fx.Invoke(func(_ *web.Server) {}),
	)
	app.Run()
}
