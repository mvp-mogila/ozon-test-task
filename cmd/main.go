package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mvp-mogila/ozon-test-task/internal/app"
)

func main() {
	log.Println("Application init")
	a := app.NewApp()

	log.Println("Start application")
	go func() {
		err := a.Run()
		log.Fatal(err)
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGTERM, syscall.SIGINT)

	<-stopCh

	a.Stop()
	log.Println("Application stopped")
}
