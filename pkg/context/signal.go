package context

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// WithSignalHandler wraps the context so a system Interrupt or SIGTERM signal cancels it
func WithSignalHandler(pCtx context.Context) context.Context {
	ctx, cancel := context.WithCancel(pCtx)
	go func() {
		sCh := make(chan os.Signal, 1)
		signal.Notify(sCh, os.Interrupt, syscall.SIGTERM)
		<-sCh
		cancel()
	}()
	return ctx
}

// Background wraps `context.Background()` with a signal handler
func Background() context.Context {
	return WithSignalHandler(context.Background())
}
