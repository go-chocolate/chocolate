package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chocolate/chocolate/example/internal/app"
)

func main() {
	defer app.Shutdown()

	go app.Run()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGHUP, syscall.SIGQUIT)
	<-ch
}
