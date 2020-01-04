package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/onion-studio/onion-weekly/dto"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/onion-studio/onion-weekly/config"
	"github.com/onion-studio/onion-weekly/db"
	"github.com/onion-studio/onion-weekly/domain"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			config.NewAppConf,
			db.NewPgxPool,
			domain.NewUserService,
		),
		fx.Invoke(func(pgxPool *pgxpool.Pool, userService *domain.UserService) {
			var email, fullName string
			fmt.Println("Enter the email of new user:")
			if _, err := fmt.Scanln(&email); err != nil {
				log.Fatal(err)
			}
			fmt.Println("Enter the password of new user:")
			password, err := terminal.ReadPassword(0)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Enter the full name of new user:")
			if _, err := fmt.Scanln(&fullName); err != nil {
				log.Fatal(err)
			}

			ctx := context.Background()
			tx, err := pgxPool.Begin(ctx)
			if err != nil {
				log.Fatal(err)
			}
			_, _, _, err = userService.CreateUserWithEmailCredential(tx, dto.CreateUserInput{
				Email:    email,
				Password: string(password),
				FullName: fullName,
			})
			if err != nil {
				log.Fatal(err)
			}
			if err := tx.Commit(ctx); err != nil {
				log.Fatal(err)
			}
			fmt.Println("Successfully created.")
			os.Exit(0)
		}),
	)

	app.Run()
}
