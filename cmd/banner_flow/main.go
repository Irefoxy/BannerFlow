package main

import (
	"BannerFlow/internal/app"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	stopChan := make(chan os.Signal, 1)
	a := app.NewApp()

	go func() {
		a.Run()
	}()

	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan
	a.Stop()
}
