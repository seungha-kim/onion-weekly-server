package config

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConf struct {
	Test       bool
	Secret     string
	Debug      bool
	BcryptCost int
	PgURL      string
	Port       int64
}

type Error struct {
	message string
}

func (c Error) Error() string {
	return c.message
}

func NewAppConf() (AppConf, error) {
	wd, _ := os.Getwd()
	_ = godotenv.Load(path.Join(wd, "main.env"))
	return loadAppConf(false)
}

func NewTestAppConf() (AppConf, error) {
	wd, _ := os.Getwd()
	_ = godotenv.Load(path.Join(wd, "..", "test.env"))
	return loadAppConf(true)
}

func loadEnvs() (map[string]string, error) {
	defaults := map[string]string{
		"PORT":        "1323",
		"DEBUG":       "",
		"BCRYPT_COST": "14",
	}

	envs := map[string]string{
		"PORT":        "",
		"DEBUG":       "",
		"PG_URL":      "",
		"BCRYPT_COST": "",
		"SECRET":      "",
	}

	// check existence
	for k := range envs {
		fromEnv := os.Getenv(k)
		defaultValue, ok := defaults[k]
		if !ok && fromEnv == "" {
			return nil, Error{message: fmt.Sprintf("Cannot get %s", k)}
		}
		if fromEnv != "" {
			envs[k] = fromEnv
		} else {
			envs[k] = defaultValue
		}
	}

	return envs, nil
}

func loadAppConf(test bool) (AppConf, error) {
	envs, err := loadEnvs()

	port, err := strconv.ParseInt(envs["PORT"], 10, 64)
	if err != nil {
		return AppConf{}, err
	}

	bcryptCost, err := strconv.ParseInt(envs["BCRYPT_COST"], 10, 64)

	return AppConf{
		Test:       test,
		BcryptCost: int(bcryptCost),
		PgURL:      envs["PG_URL"],
		Debug:      envs["DEBUG"] != "",
		Port:       port,
		Secret:     envs["SECRET"],
	}, nil
}
