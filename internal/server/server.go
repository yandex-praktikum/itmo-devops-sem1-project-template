package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app         *fiber.App
	done        chan struct{}
	serverAddr  string
	graceTimout time.Duration
	quit        <-chan struct{}
}

func New(app *fiber.App, done chan struct{}, serverAddr string, graceTimout time.Duration, quit <-chan struct{}) *Server {
	return &Server{
		app:         app,
		done:        done,
		serverAddr:  serverAddr,
		graceTimout: graceTimout,
		quit:        quit,
	}
}
