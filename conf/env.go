package conf

import (
	"fmt"
	"github.com/caarlos0/env"
)

type Config struct {
	MongoURL string `env:"MONGO_URL" envDefault:"mongodb://localhost:27017"`
	Port     string `env:"PORT" envDefault:"8002"`
	DBName   string `env:"DB_NAME" envDefault:"crawl"`
}

var cfg Config

func SetEnv() {
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("Failed to read environment variables: %v", err)
		return
	}
}

func GetEnv() Config {
	return cfg
}
