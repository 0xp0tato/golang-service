package util

import (
	"os"
	"os/signal"
	"syscall"
)

func GetSignalChannel() chan os.Signal {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	return signals

}
