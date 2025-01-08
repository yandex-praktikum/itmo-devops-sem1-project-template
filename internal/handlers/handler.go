package handlers

import (
	"github.com/gorilla/mux"
	"project_sem/internal/db"
)

type Handler struct {
	pricesRepo *db.PricesRepo
}

func (h *Handler) InitRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/v0/prices", h.uploadItems).Methods("POST")
	router.HandleFunc("/api/v0/prices", h.getItems).Methods("GET")
	return router
}

func NewHandler(repo *db.PricesRepo) *Handler {
	return &Handler{pricesRepo: repo}
}
