package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-requiered:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address       string        `yaml:"address" env-default:"localhost:8080"`
	Timeout       time.Duration `yaml:"timeout" env-default:"4s"`
	Idle_timeoute time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User          string        `yaml:"user" env-requiered:"true"`
	Password      string        `yaml:"password" env-requiered:"true" env:"HTTP_SERVER_PASSWORD"`
}

func LoadConfig() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file %s is not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("cannot read config, error: %s", err)
	}

	return &cfg
}
