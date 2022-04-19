package config

import (
	"flag"
	"fmt"
	"sync"

	"github.com/caarlos0/env/v6"
)

var once sync.Once
var Config *config

func GetConfig() *config {
	once.Do(func() {
		Config = parseConfig()
	})

	return Config
}

type config struct {
	RunAddress  string `env:"RUN_ADDRESS" envDefault:":8080"`
	BaseURL     string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	DatabaseURI string `env:"DATABASE_URI"`
	SessionKey  string `env:"SESSION_KEY"`
}

func parseConfig() *config {
	cfg := &config{}
	if err := env.Parse(cfg); err != nil {
		fmt.Println("failed to parse config: %w", err)
	}

	flag.StringVar(&cfg.RunAddress, "a", cfg.RunAddress, "server address")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "base url")
	flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI, "database data source name")
	flag.StringVar(&cfg.SessionKey, "k", cfg.SessionKey, "session key")

	flag.Parse()

	return cfg
}
