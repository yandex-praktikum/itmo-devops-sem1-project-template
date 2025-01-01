package main

import (
	"os"
	"os/signal"
	"syscall"

	"project_sem/internal/app"
)

func main() {
	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	app.Run(quitSignal)
}
