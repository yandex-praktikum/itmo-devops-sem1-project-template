package models

type TotalItemsResponse struct {
	TotalItems      int     `json:"total_items"`
	TotalCategories int     `json:"total_categories"`
	TotalPrice      float32 `json:"total_price"`
}

type Item struct {
	Id         int
	Name       string
	Category   string
	Price      float64
	CreateDate string
}
