package db

import (
	"database/sql"
	"log"
	"project_sem/internal/models"
)

const (
	insertQuery                string = `INSERT INTO prices (id, name, category, price, create_date) VALUES ($1, $2, $3, $4, $5)`
	selectTotalPriceQuery      string = `SELECT coalesce(sum(price), 0) FROM prices`
	selectCategoriesCountQuery string = `SELECT count(distinct category) FROM prices`
	selectQuery                string = `SELECT id, name, category, price, create_date FROM prices`
)

type PricesRepoInterface interface {
	AddItems([]models.Item) (*models.TotalItemsResponse, error)
	GetItems() (*[]models.Item, error)
}

type PricesRepo struct {
	db *sql.DB
}

func (repo *PricesRepo) AddItems(items []models.Item) (*models.TotalItemsResponse, error) {
	tx, err := repo.db.Begin()
	defer tx.Rollback()
	if err != nil {
		log.Printf("An error occured on starting transaction %s", err.Error())
		return nil, err
	}

	totalItems := 0
	for _, item := range items {
		_, err := tx.Exec(insertQuery, item.Id, item.Name, item.Category, item.Price, item.CreateDate)
		if err != nil {
			log.Printf("An error %s occured on adding item %v", err, item)
		} else {
			totalItems += 1
		}
	}
	log.Printf("%d items were added", totalItems)

	var totalPrice float32
	err = tx.QueryRow(selectTotalPriceQuery).Scan(&totalPrice)
	if err != nil {
		log.Printf("An error occured on select total price - %s", err.Error())
		return nil, err
	}

	var totalCategories int
	err = tx.QueryRow(selectCategoriesCountQuery).Scan(&totalCategories)
	if err != nil {
		log.Printf("An error occured on select categories count - %s", err.Error())
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("An error occured on commit - %s", err.Error())
		return nil, err
	}

	responseItem := models.TotalItemsResponse{
		TotalItems:      totalItems,
		TotalCategories: totalCategories,
		TotalPrice:      totalPrice,
	}
	return &responseItem, nil
}

func (repo *PricesRepo) GetItems() ([]models.Item, error) {
	var items []models.Item

	rows, err := repo.db.Query(selectQuery)
	if err != nil {
		log.Printf("An error occured on select all rows %s", err.Error())
		return items, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Item
		err = rows.Scan(&item.Id, &item.Name, &item.Category, &item.Price, &item.CreateDate)
		if err != nil {
			log.Printf("An error occured on scan item - %s", err)
			return items, err
		}
		items = append(items, item)
	}
	return items, nil
}

func NewPricesRepo(db *sql.DB) *PricesRepo {
	return &PricesRepo{db: db}
}
