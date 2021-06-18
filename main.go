package main

import (
	"github.com/joho/godotenv"
	"github.com/menarayanzshrestha/trello/bootstrap"
	"go.uber.org/fx"
)

func main() {
	godotenv.Load()

	fx.New(bootstrap.Module).Run()

}
