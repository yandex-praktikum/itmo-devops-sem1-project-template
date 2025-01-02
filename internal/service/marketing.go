package service

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"project_sem/internal/model"

	"github.com/shopspring/decimal"
)

var ErrCSVData = errors.New("failed to parse csv data")

const CSVColumnsNumber = 5

type MarketingRepository interface {
	UploadProducts(ctx context.Context, products []model.Product) error
	LoadProducts(ctx context.Context) ([]model.Product, error)
}

type MarketingService struct {
	repository MarketingRepository
}

func NewMarketingService(repository MarketingRepository) *MarketingService {
	return &MarketingService{repository: repository}
}

func (s *MarketingService) SaveProducts(ctx context.Context, file *multipart.FileHeader) (*model.LoadResult, error) {
	openedFile, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer openedFile.Close()

	zipReader, err := zip.NewReader(openedFile, file.Size)
	if err != nil {
		return nil, fmt.Errorf("failed to read zip archive: %w", err)
	}

	var products []model.Product

	for _, zipFile := range zipReader.File {
		// Пропускаем папки
		if zipFile.FileInfo().IsDir() {
			continue
		}

		// Обрабатываем только CSV-файлы
		if !strings.HasSuffix(zipFile.Name, ".csv") {
			continue
		}

		productsFile, err := processProductCSV(zipFile)
		if err != nil {
			return nil, fmt.Errorf("process zip file %w", err)
		}

		products = append(products, productsFile...)
	}

	err = s.repository.UploadProducts(ctx, products)
	if err != nil {
		return nil, fmt.Errorf("save products in db %w", err)
	}

	return formLoadResult(products), nil
}

func (s *MarketingService) LoadProducts(ctx context.Context) ([]byte, error) {
	products, err := s.repository.LoadProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("load products from db %w", err)
	}

	// Создаём буфер для CSV-данных
	var csvBuffer bytes.Buffer
	csvWriter := csv.NewWriter(&csvBuffer)

	// Записываем заголовок
	if err := csvWriter.Write([]string{"id", "name", "category", "price", "create_date"}); err != nil {
		return nil, fmt.Errorf("write CSV header: %w", err)
	}

	// Записываем продукты
	for _, product := range products {
		record := []string{
			fmt.Sprintf("%d", product.ID),
			product.Name,
			product.Category,
			product.Price.String(),
			product.CreateDate.Format(time.RFC3339),
		}
		if err := csvWriter.Write(record); err != nil {
			return nil, fmt.Errorf("write CSV record: %w", err)
		}
	}

	// Очищаем буфер CSV-писателя
	csvWriter.Flush()

	if err := csvWriter.Error(); err != nil {
		return nil, fmt.Errorf("flush CSV writer: %w", err)
	}

	// Создаём буфер для ZIP-архива
	var zipBuffer bytes.Buffer
	zipWriter := zip.NewWriter(&zipBuffer)

	// Добавляем CSV-файл в архив
	csvFileInZip, err := zipWriter.Create("data.csv")
	if err != nil {
		return nil, fmt.Errorf("create file in ZIP archive: %w", err)
	}

	// Записываем CSV-данные в архив
	if _, err := csvFileInZip.Write(csvBuffer.Bytes()); err != nil {
		return nil, fmt.Errorf("write CSV data to ZIP archive: %w", err)
	}

	// Закрываем ZIP-архив
	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("close ZIP archive: %w", err)
	}

	// Возвращаем ZIP-архив как []byte
	return zipBuffer.Bytes(), nil
}

func processProductCSV(zipFile *zip.File) ([]model.Product, error) {
	zippedFile, err := zipFile.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open zipped file: %w", err)
	}
	defer zippedFile.Close()

	reader := csv.NewReader(zippedFile)

	// Пропускаем заголовок (первую строку)
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	var products []model.Product
	// Парсим строки CSV
	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("failed to read CSV record: %w", err)
		}

		// Конвертация строки в модель Product
		product, err := parseCSVProduct(record)
		if err != nil {
			return nil, fmt.Errorf("failed to parse product: %w", err)
		}

		products = append(products, product)
	}

	return products, nil
}

func parseCSVProduct(record []string) (model.Product, error) {
	if len(record) != CSVColumnsNumber {
		return model.Product{}, fmt.Errorf("invalid record length: %w", ErrCSVData)
	}

	id, err := strconv.Atoi(record[0])
	if err != nil {
		return model.Product{}, fmt.Errorf("invalid id: %w", ErrCSVData)
	}

	price, err := decimal.NewFromString(record[3])
	if err != nil {
		return model.Product{}, fmt.Errorf("invalid price: %w", ErrCSVData)
	}

	date, err := time.Parse("2006-01-02", record[4])
	if err != nil {
		return model.Product{}, fmt.Errorf("invalid create_date: %w", ErrCSVData)
	}

	return model.Product{
		ID:         id,
		Name:       record[1],
		Category:   record[2],
		Price:      price,
		CreateDate: date,
	}, nil
}

func formLoadResult(products []model.Product) *model.LoadResult {
	res := model.LoadResult{}
	cat := make(map[string]struct{})

	for i := 0; i < len(products); i++ {
		res.TotalQuantity++
		cat[products[i].Category] = struct{}{}
		res.TotalPrice = res.TotalPrice.Add(products[i].Price)
	}

	res.TotalCategories = len(cat)

	return &res
}
