package csv

import (
	"bytes"
	"encoding/csv"
	"io"
	"project_sem/src/db"
	"strconv"
)

func recordToPrice(rec []string) (db.Price, error) {
	id, err := strconv.Atoi(rec[0])

	if err != nil {
		return db.Price{}, err
	}

	price, err := strconv.ParseFloat(rec[3], 64)

	if err != nil {
		return db.Price{}, err
	}

	return db.Price{
		Id:         id,
		Name:       rec[1],
		Category:   rec[2],
		Price:      price,
		CreateDate: rec[4],
	}, nil
}

func priceToRecord(price db.Price) []string {
	return []string{
		strconv.Itoa(price.Id),
		price.Name,
		price.Category,
		strconv.FormatFloat(price.Price, 'f', 2, 64),
		price.CreateDate,
	}
}

func PricesToCsv(prices []db.Price) ([]byte, error) {
	var buf bytes.Buffer

	w := csv.NewWriter(&buf)

	headers := []string{"id", "name", "category", "price", "create_date"}
	err := w.Write(headers)

	if err != nil {
		return nil, err
	}

	for _, price := range prices {
		err = w.Write(priceToRecord(price))

		if err != nil {
			return nil, err
		}
	}

	w.Flush()

	return buf.Bytes(), err
}

func GetPricesFromCsv(r io.Reader) ([]db.Price, error) {
	csvR := csv.NewReader(r)

	csvR.FieldsPerRecord = 5
	_, err := csvR.Read()

	if err != nil {
		return nil, err
	}

	var prices []db.Price
	for {
		rec, e := csvR.Read()

		if e != nil {
			if e == io.EOF {
				break
			}
			return nil, e
		}

		price, err := recordToPrice(rec)

		if err != nil {
			continue
		}

		prices = append(prices, price)
	}

	return prices, nil
}
