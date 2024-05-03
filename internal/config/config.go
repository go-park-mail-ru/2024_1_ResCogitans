package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type Config struct {
	Env            string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath    string `yaml:"storage_path" env-required:"true"`
	HTTPServer     `yaml:"http_server"`
	Dsn            `yaml:"dsn"`
	Redis          `yaml:"redis"`
	FileUploadPath string `env:"FILE_UPLOAD_PATH"`
	Drive          `yaml:"token"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

type Dsn struct {
	Host     string `env:"DB_POSTGRES_HOST"`
	Port     int    `env:"DB_POSTGRES_PORT"`
	User     string `env:"DB_POSTGRES_USER"`
	Password string `env:"DB_POSTGRES_PASSWORD"`
	DBName   string `env:"DB_POSTGRES_NAME"`
}

type Redis struct {
	Host     string `yaml:"host" env:"DB_REDIS_HOST"`
	Port     int    `yaml:"port" env:"DB_REDIS_PORT"`
	DB       int    `yaml:"db" env:"DB_REDIS_DB"`
	Password string `yaml:"password" env:"DB_REDIS_PASSWORD"`
}

type Drive struct {
	Token string `yaml:"token" env:"YANDEX_TOKEN"`
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
