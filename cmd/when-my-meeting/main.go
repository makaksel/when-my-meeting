package main

import (
	"context"
	"log"
	"makaksel/when-my-meeting/internal/app"
	"os"
	"os/signal"
	"syscall"

	_ "time/tzdata"
)

var version = "dev"

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	a, err := app.New(cancel)
	if err != nil {
		log.Fatalf("app init: %v", err)
	}

	if err := a.Run(ctx); err != nil {
		log.Fatalf("app run: %v", err)
	}
}
