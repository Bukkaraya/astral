package main

import (
	"log"

	"clientgen-v2/internal/cli"
)

func main() {
	app := cli.NewApp()
	
	if err := app.Execute(); err != nil {
		log.Fatal(err)
	}
}