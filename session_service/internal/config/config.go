package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type Config struct {
	Env    string `yaml:"env" env:"ENV" env-required:"true"`
	Redis  `yaml:"redis"`
	Server `yaml:"server"`
}

type Redis struct {
	Host     string `yaml:"host" env:"DB_SESSION_HOST"`
	Port     int    `yaml:"port" env:"DB_SESSION_PORT"`
	DB       int    `yaml:"db" env:"DB_SESSION_DB"`
	Password string `yaml:"password" env:"DB_SESSION_PASSWORD"`
}

type Server struct {
	Protocol string `yaml:"protocol" env:"SERVICE_PROTOCOL"`
	Port     string `yaml:"port" env:"SERVICE_PORT"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, errors.Wrap(err, "error loading .env file")
	}

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		return nil, errors.New("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, errors.Errorf("config file %s is not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, errors.Wrap(err, "cannot read config")
	}

	return &cfg, nil
}
