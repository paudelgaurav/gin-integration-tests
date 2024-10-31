package main

import (
	"github.com/joho/godotenv"
	"github.com/paudelgaurav/gin-integration-tests/cmd"
)

func main() {

	_ = godotenv.Load()
	cmd.Execute()

}
