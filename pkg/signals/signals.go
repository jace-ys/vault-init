package signals

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func SetupSignalContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-sigc:
			cancel()
		case <-ctx.Done():
		}

		<-sigc
		os.Exit(1)
	}()

	return ctx, cancel
}
