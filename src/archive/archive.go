package archive

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"project_sem/src/csv"
	"project_sem/src/db"
)

const (
	CsvFileName = "data.csv"
)

func GetPricesFromZip(r io.ReaderAt, size int64) ([]db.Price, error) {
	archive, err := zip.NewReader(r, size)

	if err != nil {
		return nil, err
	}

	file, err := archive.File[0].Open()

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()

	return csv.GetPricesFromCsv(file)
}

func WritePricesIntoZip(prices []db.Price) ([]byte, error) {

	b, err := csv.PricesToCsv(prices)

	if err != nil {
		return nil, err
	}

	var buff bytes.Buffer

	w := zip.NewWriter(&buff)

	f, err := w.Create(CsvFileName)
	if err != nil {
		return nil, err
	}

	_, err = f.Write(b)

	if err != nil {
		return nil, err
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func GetPricesFromTar(r io.Reader) ([]db.Price, error) {
	gzr, err := gzip.NewReader(r)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := gzr.Close(); err != nil {
			log.Println(err)
		}
	}()

	archive := tar.NewReader(gzr)

	_, err = archive.Next()

	if err != nil {
		return nil, err
	}

	return csv.GetPricesFromCsv(archive)
}
