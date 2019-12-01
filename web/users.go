package web

import (
	"github.com/labstack/echo"
	"github.com/onion-studio/onion-weekly/domain"
)

// GetFirstUserHandler haha
func GetFirstUserHandler(c echo.Context) error {
	user, err := domain.LoadFirstUser()
	if err != nil {
		return c.JSON(500, nil)
	}

	return c.JSON(200, user)
}

func PostTestUserHandler(c echo.Context) error {
	input := domain.CreateUserWithEmailCredentialInput{
		Email:    "hello@example.com",
		Password: "asdf11234",
	}
	_, cred, err := domain.CreateUserWithEmailCredential(input)
	if err != nil {
		return err
	}
	return c.JSON(200, cred)
}