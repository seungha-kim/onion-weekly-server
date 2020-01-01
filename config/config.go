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

func NewAppConf() AppConf {
	wd, _ := os.Getwd()
	_ = godotenv.Load(path.Join(wd, "main.env"))
	return loadAppConf(false)
}

func NewTestAppConf() AppConf {
	wd, _ := os.Getwd()
	_ = godotenv.Load(path.Join(wd, "..", "test.env"))
	return loadAppConf(true)
}

func loadEnvs() map[string]string {
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
			panic(fmt.Sprintf("Cannot get %s", k))
		}
		if fromEnv != "" {
			envs[k] = fromEnv
		} else {
			envs[k] = defaultValue
		}
	}

	return envs
}

func loadAppConf(test bool) AppConf {
	envs := loadEnvs()

	port, err := strconv.ParseInt(envs["PORT"], 10, 64)
	if err != nil {
		panic(err)
	}

	bcryptCost, err := strconv.ParseInt(envs["BCRYPT_COST"], 10, 64)

	return AppConf{
		Test:       test,
		BcryptCost: int(bcryptCost),
		PgURL:      envs["PG_URL"],
		Debug:      envs["DEBUG"] != "",
		Port:       port,
		Secret:     envs["SECRET"],
	}
}
