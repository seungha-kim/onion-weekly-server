package web

import (
	"github.com/labstack/echo"
	"github.com/onion-studio/onion-weekly/domain"
	"net/http"
)

func PostTestUserHandler(c echo.Context) (err error) {
	input := domain.InputCreateUser{}
	if err = c.Bind(&input); err != nil {
		return // TODO
	}
	user, _, _, err := domain.CreateUserWithEmailCredential(input)
	if err != nil {
		switch err.(type) {
		case domain.DuplicateError:
			err = echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return
	}
	return c.JSON(200, user)
}

func PostTestTokenHandler(c echo.Context) (err error) {
	input := domain.InputCreatTokenByEmailCredential{}
	if err = c.Bind(&input); err != nil {
		return // TODO
	}
	outputToken, err := domain.CreateTokenByEmailCredential(input)
	if err != nil {
		switch err.(type) {

		}
		return
	}
	switch err.(type) {
	case domain.DuplicateError:
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(200, outputToken)
}
