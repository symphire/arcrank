package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RunPort      string
	MongoURL     string
	ElasticURL   string
	ElasticIndex string
	LogLevel     string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		RunPort:      get("RUN_PORT", "8080"),
		MongoURL:     get("MONGO_URL", "mongodb://localhost:27017"),
		ElasticURL:   get("ELASTIC_URL", "http://localhost:9200"),
		ElasticIndex: get("ELASTIC_INDEX", "players"),
		LogLevel:     get("LOG_LEVEL", "debug"),
	}
}

func get(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		log.Printf("[config]: read %s from .env: %s", key, value)
		return value
	}
	log.Printf("[config]: %s not set, using default: %s", key, fallback)
	return fallback
}
