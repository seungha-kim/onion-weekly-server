package domain

import (
	"fmt"
	"testing"

	"github.com/onion-studio/onion-weekly/dto"

	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/onion-studio/onion-weekly/config"

	"github.com/onion-studio/onion-weekly/db"
)

var pgxPool *pgxpool.Pool
var appConf config.AppConf

func TestMain(m *testing.M) {
	appConf, _ = config.NewTestAppConf()
	pgxPool = db.NewPgxPool(appConf)
	m.Run()
	pgxPool.Close()
}

func TestCreateUserWithEmailCredential(t *testing.T) {
	srv := UserService{appConf: appConf}

	tests := []struct {
		name           string
		input          dto.CreateUserInput
		wantUser       dto.User
		wantCredential dto.EmailCredential
		wantProfile    dto.UserProfile
		wantErr        bool
	}{
		{
			name: "First Test",
			input: dto.CreateUserInput{
				Email:    "test@test.com",
				Password: "test",
				FullName: "Test Test",
			},
			wantCredential: dto.EmailCredential{
				Email:          "test@test.com",
				HashedPassword: "anything",
			},
			wantProfile: dto.UserProfile{
				FullName: "Test Test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		fmt.Println("what the")
		t.Run(tt.name, func(t *testing.T) {
			db.RollbackForTest(pgxPool, func(tx pgx.Tx) {
				_, gotCredential, gotProfile, err := srv.CreateUserWithEmailCredential(tx, tt.input)
				if (err != nil) != tt.wantErr {
					t.Errorf("CreateUserWithEmailCredential() error = %v, wantErr %v", err, tt.wantErr)
				}
				if err != nil {
					return
				}
				if gotCredential.Email != tt.wantCredential.Email {
					t.Errorf("CreateUserWithEmailCredential() gotCredential = %v, want %v", gotCredential, tt.wantCredential)
				}
				if gotProfile.FullName != tt.wantProfile.FullName {
					t.Errorf("CreateUserWithEmailCredential() gotProfile = %v, want %v", gotProfile, tt.wantProfile)
				}
			})
		})
	}
}

func TestCreateTokenByEmailCredential(t *testing.T) {
	srv := UserService{appConf: appConf}

	tests := []struct {
		name       string
		input      dto.CreatTokenByEmailCredentialInput
		wantOutput dto.Token
		wantErr    bool
	}{
		{
			name: "Happy",
			input: dto.CreatTokenByEmailCredentialInput{
				Email:    "test@test.com",
				Password: "test1234",
			},
			wantOutput: dto.Token{},
			wantErr:    false,
		},
		{
			name: "Wrong Password",
			input: dto.CreatTokenByEmailCredentialInput{
				Email:    "test@test.com",
				Password: "wrong_password",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db.RollbackForTest(pgxPool, func(tx pgx.Tx) {
				_, _, _, err := srv.CreateUserWithEmailCredential(
					tx,
					dto.CreateUserInput{
						Email:    "test@test.com",
						Password: "test1234",
						FullName: "Test Test",
					})
				_, err = srv.CreateTokenByEmailCredential(
					tx,
					tt.input)
				if (err != nil) != tt.wantErr {
					t.Errorf("CreateTokenByEmailCredential() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			})
		})
	}
}
