package main

import (
	"fmt"
	"os"
	"path"
)

func main() {
	wd, _ := os.Getwd()
	fmt.Println(path.Join(wd, "test.env"))
}
