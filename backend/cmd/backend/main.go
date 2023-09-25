package main

import (
	"context"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/app"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	s := app.New(
		app.WithPort(config.Port),
	)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		<-sigChan
		if err := s.Stop(ctx); err != nil {
			log.Fatalln("Error while shutting down server", err)
		}
		cancel()
	}()

	err = s.Start()
	if err != nil {
		log.Fatalln("[Error] failed to start Gin server due to:", err)
	}
}
