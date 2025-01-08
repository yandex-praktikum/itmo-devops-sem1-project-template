package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type pgConfig struct {
	host     string
	port     string
	database string
	username string
	password string
}

func (conf pgConfig) getDNS() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		conf.username,
		conf.password,
		conf.host,
		conf.port,
		conf.database,
	)
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getConfig() pgConfig {
	return pgConfig{
		host:     GetEnv("DB_HOST", "127.0.0.1"),
		port:     GetEnv("DB_PORT", "5432"),
		database: GetEnv("DB_NAME", "project-sem-1"),
		username: GetEnv("DB_USER", "validator"),
		password: GetEnv("DB_PASSWORD", "val1dat0r"),
	}
}

func GetDB() *sql.DB {
	conf := getConfig()
	dns := conf.getDNS()
	log.Printf("DB dns %s", dns)
	db, err := sql.Open("postgres", dns)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(fmt.Sprintf("An error occured while ping db %s", err))
	}
	log.Print("Connection created...")
	return db
}
