package domain

import (
	"testing"

	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/onion-studio/onion-weekly/config"

	"github.com/onion-studio/onion-weekly/db"
)

var pgxPool *pgxpool.Pool
var appConf config.AppConf

func TestMain(m *testing.M) {
	appConf = config.LoadTestAppConf()
	pgxPool = db.CreatePool(appConf.PgURL)
	m.Run()
	pgxPool.Close()
}

func TestCreateUserWithEmailCredential(t *testing.T) {
	tests := []struct {
		name           string
		input          InputCreateUser
		wantUser       User
		wantCredential EmailCredential
		wantProfile    UserProfile
		wantErr        bool
	}{
		{
			name: "First Test",
			input: InputCreateUser{
				Email:    "test@test.com",
				Password: "test",
				FullName: "Test Test",
			},
			wantCredential: EmailCredential{
				Email:          "test@test.com",
				HashedPassword: "anything",
			},
			wantProfile: UserProfile{
				FullName: "Test Test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RunInTransaction(pgxPool, func(tx pgx.Tx) {
				_, gotCredential, gotProfile, err := CreateUserWithEmailCredential(appConf, tx, tt.input)
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
	tests := []struct {
		name       string
		input      InputCreatTokenByEmailCredential
		wantOutput OutputToken
		wantErr    bool
	}{
		{
			name: "Happy",
			input: InputCreatTokenByEmailCredential{
				Email:    "test@test.com",
				Password: "test1234",
			},
			wantOutput: OutputToken{},
			wantErr:    false,
		},
		{
			name: "Wrong Password",
			input: InputCreatTokenByEmailCredential{
				Email:    "test@test.com",
				Password: "wrong_password",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RunInTransaction(pgxPool, func(tx pgx.Tx) {
				_, _, _, err := CreateUserWithEmailCredential(
					appConf,
					tx,
					InputCreateUser{
						Email:    "test@test.com",
						Password: "test1234",
						FullName: "Test Test",
					})
				_, err = CreateTokenByEmailCredential(
					appConf,
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
