package main

import (
	"os"

	"github.com/onion-studio/onion-weekly/config"
	"github.com/onion-studio/onion-weekly/db"
	"github.com/onion-studio/onion-weekly/domain"
	"github.com/onion-studio/onion-weekly/web"
	"go.uber.org/fx"
)

func main() {
	opts := []fx.Option{
		fx.Provide(
			config.NewAppConf,
			db.NewPgxPool,
			web.NewServer,
			web.NewAuthHandler,
			web.NewHelloHandler,
			domain.NewUserService,
		),
		fx.Invoke(func(_ *web.Server) {}),
	}
	if os.Getenv("DEBUG") != "" {
		opts = append(opts, fx.NopLogger)
	}
	app := fx.New(opts...)
	app.Run()
}
