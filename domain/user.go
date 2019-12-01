package domain

import (
	"context"
	"fmt"
	"os"

	"github.com/onion-studio/onion-weekly/db"
	"golang.org/x/crypto/bcrypt"
)

// User represents a record from users table
type User struct {
	Id        UUID        `json:"id"`
	CreatedAt Timestamptz `json:"createdAt"`
}

type UserProfile struct {
	UserId   UUID   `json:"userId"`
	FullName string `json:"fullName"`
}

type EmailCredential struct {
	UserId         UUID        `json:"userId"`
	Email          string      `json:"email"`
	HashedPassword string      `json:"-"`
	CreatedAt      Timestamptz `json:"createdAt"`
}

type CreateUserWithEmailCredentialInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateUserWithEmailCredential creates user with email credential, and returns them.
func CreateUserWithEmailCredential(input CreateUserWithEmailCredentialInput) (user User, credential EmailCredential, err error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
	if err != nil {
		return
	}

	tx, err := db.Pool.Begin(context.Background())
	if err != nil {
		return
	}

	err = tx.
		QueryRow(
			context.Background(),
			`insert into users (id) 
			values (default) 
			returning id, created_at;`,
		).
		Scan(&user.Id, &user.CreatedAt)

	if err != nil {
		_ = tx.Rollback(context.Background())
		return
	}

	err = tx.
		QueryRow(
			context.Background(),
			`insert into email_credentials (user_id, email, hashed_password) 
			values ($1, $2, $3) 
			returning user_id, email, hashed_password, created_at;`,
			user.Id, input.Email, hashed).
		Scan(&credential.UserId, &credential.Email, &credential.HashedPassword, &credential.CreatedAt)

	if err != nil {
		_ = tx.Rollback(context.Background())
		return
	}

	err = tx.Commit(context.Background())

	return
}

// LoadFirstUser load the first user (for testing)
func LoadFirstUser() (user User, err error) {
	err = db.Pool.QueryRow(context.Background(), "select id from users").Scan(&user.Id)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	return
}
