package domain

import (
	"context"
	"strings"

	"github.com/onion-studio/onion-weekly/dto"

	"github.com/onion-studio/onion-weekly/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	appConf config.AppConf
}

func NewUserService(appConf config.AppConf) *UserService {
	return &UserService{appConf: appConf}
}

// language=PostgreSQL
const sqlInsertUsers = `
insert into users (id)
values (default) 
returning id, created_at;`

// language=PostgreSQL
const sqlInsertEmailCredentials = `
insert into email_credentials (user_id, email, hashed_password) 
values ($1, $2, $3) 
returning user_id, email, hashed_password, created_at;`

// language=PostgreSQL
const sqlInsertUserProfiles = `
insert into user_profiles (user_id, full_name)
values ($1, $2)
returning user_id, full_name;`

// CreateUserWithEmailCredential creates user with email credential, and returns them.
func (srv *UserService) CreateUserWithEmailCredential(
	tx pgx.Tx,
	input dto.CreateUserInput,
) (user dto.User, credential dto.EmailCredential, profile dto.UserProfile, err error) {
	ctx := context.Background()

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), srv.appConf.BcryptCost)
	if err != nil {
		return
	}

	err = tx.
		QueryRow(ctx, sqlInsertUsers).
		Scan(&user.Id, &user.CreatedAt)

	if err != nil {
		_ = tx.Rollback(context.Background())
		return
	}

	err = tx.
		QueryRow(ctx, sqlInsertEmailCredentials, user.Id, input.Email, hashed).
		Scan(&credential.UserId, &credential.Email, &credential.HashedPassword, &credential.CreatedAt)

	if err != nil {
		_ = tx.Rollback(context.Background())

		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"email_credentials_email_key\"") {
			err = DuplicateError{fieldName: "email"}
		}

		return
	}

	err = tx.
		QueryRow(ctx, sqlInsertUserProfiles, user.Id, input.FullName).
		Scan(&profile.UserId, &profile.FullName)

	return
}

// language=PostgreSQL
const sqlSelectEmailCredentialsByEmail = `
select e.user_id, e.email, e.hashed_password, e.created_at 
from email_credentials e
where e.email = $1;`

func (srv *UserService) CreateTokenByEmailCredential(
	tx pgx.Tx,
	input dto.CreatTokenByEmailCredentialInput,
) (output dto.Token, err error) {
	ctx := context.Background()
	ec := dto.EmailCredential{}

	err = tx.
		QueryRow(ctx, sqlSelectEmailCredentialsByEmail, input.Email).
		Scan(&ec.UserId, &ec.Email, &ec.HashedPassword, &ec.CreatedAt)

	if err = bcrypt.CompareHashAndPassword([]byte(ec.HashedPassword), []byte(input.Password)); err != nil {
		return
	}
	var buf []byte
	buf, err = ec.UserId.EncodeText(nil, buf)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": string(buf),
	})

	tokenString, err := token.SignedString(srv.appConf.Secret)

	if err != nil {
		return
	}
	output.Token = tokenString
	return
}
