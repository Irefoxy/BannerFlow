package main

import (
	"BannerFlow/internal/app"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	stopChan := make(chan os.Signal, 1)
	a := app.NewApp()
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		a.Run()
	}()

	go func() {
		time.Sleep(5 * time.Second)
		stopChan <- syscall.SIGTERM
	}()
	<-stopChan
	a.Stop()
}
