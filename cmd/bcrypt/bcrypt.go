package main

import (
	"fmt"
	"os"

	"github.com/onion-studio/onion-weekly/config"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	appConf := config.NewAppConf()
	input := os.Args[1]
	hashed, err := bcrypt.GenerateFromPassword([]byte(input), appConf.BcryptCost)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(string(hashed))
	}
}
