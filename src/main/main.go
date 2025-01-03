package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"project_sem/src/archive"
	"project_sem/src/db"
	"strconv"
)

type PostResponse struct {
	TotalItems      int `json:"total_items"`
	TotalCategories int `json:"total_categories"`
	TotalPrice      int `json:"total_price"`
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		prices, err := db.GetPrices()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		z, err := archive.WritePricesIntoZip(prices)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		r.Header.Set("Content-Type", "application/zip")
		r.Header.Set("Content-Disposition", "attachment; filename=\"prices.zip\"")
		r.Header.Set("Content-Length", strconv.Itoa(len(z)))

		_, err = fmt.Fprint(w, string(z))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	if r.Method == "POST" {
		file, handler, err := r.FormFile("file")

		if err != nil {
			log.Println(err)
		}

		defer func() {
			if err := file.Close(); err != nil {
				log.Println(err)
			}
		}()

		t := r.URL.Query().Get("type")

		var prices []db.Price
		if t == "tar" {
			prices, err = archive.GetPricesFromTar(file)
		} else {
			prices, err = archive.GetPricesFromZip(file, handler.Size)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		err = db.InsertPrices(prices)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		sP, cC, err := db.GetSumPriceAndCountCategories()

		body, err := json.Marshal(PostResponse{
			TotalItems:      len(prices),
			TotalCategories: cC,
			TotalPrice:      sP,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = fmt.Fprint(w, string(body))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}
}

func main() {
	http.HandleFunc("/api/v0/prices", httpHandler)
	log.Println(http.ListenAndServe(":8080", nil))
}
