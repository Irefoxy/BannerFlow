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
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		a.Run()
	}()

	<-stopChan
	a.Stop()
}
