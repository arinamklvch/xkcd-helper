package config

import (
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ServerTimeout  int    `yaml:"server_shutdown_timeout"`
	DbTimeout      int    `yaml:"db_connect_timeout"`
	Port           int    `yaml:"port"`
	MaxWorkers     int    `yaml:"max_workers"`
	RateLimit      int    `yaml:"rate_limit"`
	RateBurst      int    `yaml:"rate_burst"`
	TokenTTL       int    `yaml:"token_TTL"`
	MaxFoundComics int    `yaml:"max_found_comics"`
	LoggerLevel    string `yaml:"logger_level"`
	JWTsecretKey   string
	DatabaseURL    string
}

func New() (Config, error) {
	data, err := os.ReadFile("config.yml")
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	err = godotenv.Load()
	if err != nil {
		return Config{}, err
	}

	config.JWTsecretKey = os.Getenv("JWT_SECRET_KEY")
	config.DatabaseURL = os.Getenv("DATABASE_URL")

	return config, nil
}
