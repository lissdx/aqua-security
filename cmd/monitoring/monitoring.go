package main

import (
	"context"
	"github.com/lissdx/aqua-security/internal/initializers/invokers"
	"github.com/lissdx/aqua-security/internal/initializers/providers"
	"go.uber.org/fx"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	kill := make(chan os.Signal, 1)
	signal.Notify(kill, os.Interrupt, syscall.SIGTERM)

	app := fx.New(
		fx.Provide(providers.Dependencies()...),
		fx.Invoke(invokers.Invokers()...),
	)

	err := app.Start(ctx)
	if err != nil {
		log.Fatalln("Error starting app:", err)
	}

	<-kill
	_ = app.Stop(ctx)
	cancel()
}
