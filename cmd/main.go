package main

import (
	"log"

	"github.com/mvp-mogila/ozon-test-task/internal/app"
)

func main() {
	log.Println("Application init")
	app := app.NewApp()

	log.Println("Start application")
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
