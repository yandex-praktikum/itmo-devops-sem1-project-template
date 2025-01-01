package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	APIHost      string
	DBHost       string
	DBPort       string
	DBName       string
	DBUser       string
	DBPassword   string
	GraceTimeout time.Duration
}

func NewConfig() (*Config, error) {
	err := godotenv.Load("devops_project.env")
	if err != nil {
		return nil, fmt.Errorf("unable to load config from env %w", err)
	}

	config := &Config{}

	apiHost, err := url.Parse(os.Getenv("API_HOST"))
	if err != nil {
		return nil, fmt.Errorf("unable to parse api host url %w", err)
	}

	config.APIHost = apiHost.Host

	config.DBHost = os.Getenv("DB_HOST")
	config.DBPort = os.Getenv("DB_PORT")
	config.DBName = os.Getenv("DB_NAME")
	config.DBUser = os.Getenv("DB_USER")
	config.DBPassword = os.Getenv("DB_PASSWORD")

	graceTimeout, err := strconv.Atoi(os.Getenv("GRACE_TIMEOUT"))
	if err != nil {
		graceTimeout = 5
	}

	config.GraceTimeout = time.Duration(graceTimeout) * time.Second

	return config, nil
}
