package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
)

func (s *Server) Run() {
	httpServerCtx, httpServerStopCtx := context.WithCancel(context.Background())

	go func() {
		if err := s.app.Listen(s.serverAddr); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}()

	go func() {
		<-s.quit

		slog.Info(fmt.Sprintf("Shutting down %s gracefully...", s.serverAddr))

		shutdownCtx, cancel := context.WithTimeout(httpServerCtx, s.graceTimout)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()

			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				slog.Error(fmt.Sprintf("%s graceful shutdown timed out.. forcing exit.", s.serverAddr))
				os.Exit(1)
			}
		}()

		if errShutdown := s.app.Shutdown(); errShutdown != nil {
			slog.Error(errShutdown.Error())
			os.Exit(1)
		}

		slog.Info(fmt.Sprintf("%s gracefully stopped.", s.serverAddr))
		httpServerStopCtx()
	}()

	<-httpServerCtx.Done()
	s.done <- struct{}{}
}
