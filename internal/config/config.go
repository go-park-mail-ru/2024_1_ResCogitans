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
	StoragePath    string `yaml:"storage_path" env-requiered:"true"`
	HTTPServer     `yaml:"http_server"`
	Dsn            `yaml:"dsn"`
	FileUploadPath string `env:"FILE_UPLOAD_PATH"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"user" env-requiered:"true"`
	Password    string        `yaml:"password" env-requiered:"true" env:"HTTP_SERVER_PASSWORD"`
}

type Dsn struct {
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	DBname   string `env:"DB_NAME"`
}

var DSN = Dsn{
	Host:     "localhost",
	Port:     5432,
	User:     "postgres",
	Password: "admin",
	DBname:   "res",
}

type redisData struct {
	Addr     string
	Username string
	Password string
	DB       int
}

var RedisData = redisData{
	Addr:     "redis-13041.c302.asia-northeast1-1.gce.cloud.redislabs.com:13041",
	Username: "default",
	Password: "Hwsuxke8YC8vT6E2jOKd7lTK6cPEvq5I",
	DB:       0,
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
