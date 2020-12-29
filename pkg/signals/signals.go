package signals

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func SetupSignalHandler() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}

		<-c
		os.Exit(1)
	}()

	return ctx, cancel
}
