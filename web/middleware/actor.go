package middleware

import (
	"context"

	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
	"github.com/onion-studio/onion-weekly/domain"
	"github.com/onion-studio/onion-weekly/dto"
)

var _ error = &ActorError{}

type ActorError struct{}

func (e *ActorError) Error() string {
	return "Cannot get uid"
}

func Actor(pgxPool *pgxpool.Pool, userService *domain.UserService) echo.MiddlewareFunc {
	// requires: em.JWT(h.appConf.Secret)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			uid, ok := claims["uid"]
			if !ok {
				return &ActorError{}
			}
			ctx := context.Background()
			tx, err := pgxPool.Begin(ctx)
			if err != nil {
				return err
			}
			defer tx.Rollback(ctx)
			input := dto.FindUserByIdInput{}
			err = input.Id.Set(uid)
			if err != nil {
				return err
			}
			actor, err := userService.FindUserById(tx, input)
			if err != nil {
				return err
			}
			c.Set("actor", actor)
			return next(c)
		}
	}
}
