package handlers

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	fp "path/filepath"
	"project_sem/internal/models"
	"strconv"
	"time"
)

func generateFileName() string {
	return fmt.Sprintf("data-%d.csv", time.Now().UnixNano())
}

func (h *Handler) uploadItems(response http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Printf("An error occured on reading request body - %s", err)
		http.Error(response, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		log.Printf("An error occured on reading zip file - %s", err)
		http.Error(response, "Invalid request body", http.StatusBadRequest)
		return
	}

	zipFile := zipReader.File[0]
	f, err := zipFile.Open()
	if err != nil {
		log.Printf("An error occured on open zip file - %s", err)
		http.Error(response, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer f.Close()

	var items []models.Item
	reader := csv.NewReader(f)
	firstLine := true
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if firstLine {
			firstLine = false
			continue
		}
		var item models.Item
		itemId, err := strconv.Atoi(record[0])
		if err != nil {
			log.Printf("Invalid ID format for row - %s", record)
			http.Error(response, "Invalid request body", http.StatusBadRequest)
			return
		}
		item.Id = itemId
		item.Name = record[1]
		item.Category = record[2]
		itemPrice, err := strconv.ParseFloat(record[3], 32)
		if err != nil {
			log.Printf("Invalid price for row - %s", record)
			http.Error(response, "Invalid request body", http.StatusBadRequest)
			return
		}
		item.Price = itemPrice
		item.CreateDate = record[4]
		items = append(items, item)
	}

	result, err := h.pricesRepo.AddItems(items)
	if err != nil {
		log.Printf("An error occurd on adding items - %s", err)
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(result)
	if err != nil {
		log.Printf("An error occurd on adding items - %s", err)
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Write(responseBody)
}

func (h *Handler) getItems(response http.ResponseWriter, request *http.Request) {
	items, err := h.pricesRepo.GetItems()
	if err != nil {
		log.Printf("An error occurd on adding items - %s", err)
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Printf("Got %d items from db", len(items))

	dir := os.TempDir()
	csvFileName := generateFileName()
	csvFilePath := fp.Join(dir, csvFileName)
	csvFile, err := os.Create(csvFilePath)
	if err != nil {
		log.Printf("An error occurd on creating file %s - %s", csvFilePath, err)
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer os.Remove(csvFilePath)
	defer csvFile.Close()
	{
		csvWriter := csv.NewWriter(csvFile)
		err = csvWriter.Write([]string{"id", "name", "category", "price", "created_at"})
		if err != nil {
			log.Printf("An error occurd on writeing to file %s", err)
			http.Error(response, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		for _, item := range items {
			var itemRow [5]string
			itemRow[0] = strconv.FormatInt(int64(item.Id), 10)
			itemRow[1] = item.Name
			itemRow[2] = item.Category
			itemRow[3] = fmt.Sprintf("%.2f", item.Price)
			itemRow[4] = item.CreateDate
			err = csvWriter.Write(itemRow[:])
			if err != nil {
				log.Printf("An error occurd on writeing to file %s", err)
				http.Error(response, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
		csvWriter.Flush()
		log.Printf("CSV file is ready - %s", csvFileName)
	}

	zipWriter := zip.NewWriter(response)
	defer zipWriter.Close()

	csvFileWriter, err := zipWriter.Create("data.csv")
	csvFile.Seek(0, 0)
	_, err = io.Copy(csvFileWriter, csvFile)
	if err != nil {
		log.Printf("An error occurd on writeing csv file to archive - %s", err)
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err != nil {
		log.Printf("An error occurd while writing zip to response - %s", err)
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	response.Header().Set("Content-Type", "application/zip")
	response.Header().Set("Content-Disposition", "attachment; filename=\"data.zip\"")
}
