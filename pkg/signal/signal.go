package signal

import (
	"os"
	"os/signal"
	"syscall"
)

func SetupSignalHandler() (stopCh <-chan struct{}) {
	var (
		onlyOneSignalHandler = make(chan struct{})
		shutdownSignals      = []os.Signal{os.Interrupt, syscall.SIGTERM}
	)

	close(onlyOneSignalHandler) // panics when called twice

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}
