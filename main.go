package main

import (
	"github.com/joho/godotenv"
	"github.com/paudelgaurav/gin-boilerplate/cmd"
)

func main() {

	_ = godotenv.Load()
	cmd.Execute()

}
