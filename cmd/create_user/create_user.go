package main

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/onion-studio/onion-weekly/config"
	"github.com/onion-studio/onion-weekly/db"
	"github.com/onion-studio/onion-weekly/domain"
	"go.uber.org/fx"
)

func main() {
	var email, fullName string
	fmt.Println("Enter the email of new user:")
	if _, err := fmt.Scanln(&email); err != nil {
		panic(err)
	}
	fmt.Println("Enter the password of new user:")
	password, err := terminal.ReadPassword(0)
	if err != nil {
		panic(err)
	}
	fmt.Println("Enter the full name of new user:")
	if _, err := fmt.Scanln(&fullName); err != nil {
		panic(err)
	}
	app := fx.New(
		fx.Provide(
			config.NewAppConf,
			db.NewPgxPool,
			domain.NewUserService,
		),
		fx.Invoke(func(pgxPool *pgxpool.Pool, userService *domain.UserService) {
			ctx := context.Background()
			tx, err := pgxPool.Begin(ctx)
			if err != nil {
				panic(err)
			}
			_, _, _, err = userService.CreateUserWithEmailCredential(tx, domain.InputCreateUser{
				Email:    email,
				Password: string(password),
				FullName: fullName,
			})
			if err != nil {
				panic(err)
			}
			if err := tx.Commit(ctx); err != nil {
				panic(err)
			}
			fmt.Println("Successfully created.")
			os.Exit(0)
		}),
		fx.NopLogger,
	)

	app.Run()
}
