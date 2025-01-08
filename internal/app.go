package internal

import (
	"log"
	"net/http"
	"project_sem/internal/db"
	"project_sem/internal/handlers"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run() {
	conn := db.GetDB()
	repo := db.NewPricesRepo(conn)
	handler := handlers.NewHandler(repo)
	s.httpServer = &http.Server{
		Addr:    ":8080",
		Handler: handler.InitRoutes(),
	}
	log.Fatal(s.httpServer.ListenAndServe())
}
