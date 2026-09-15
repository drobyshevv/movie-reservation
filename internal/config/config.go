package config

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env    string     `yaml:"env" env:"ENV"`
	Server HttpServer `yaml:"http-server"`
	DB     DataBase   `yaml:"db"`
}

type HttpServer struct {
	Host string `yaml:"host" env:"SERVER_HOST" env-default:"localhost"`
	Port int    `yaml:"port" env:"SERVER_PORT" env-default:"8081"`
}

type DataBase struct {
	User     string `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Password string `yaml:"" env:"DB_PASSWORD"`
	Host     string `yaml:"" env:"DB_HOST" env-default:"localhost"`
	Port     int    `yaml:"" env:"DB_PORT" env-default:"5432"`
	Name     string `yaml:"" env:"DB_NAME" env-default:"movie-reservation"`
	SslMode  string `yaml:"" env:"DB_SSLMODE" env-default:"disable"`
}

// postgres://user:password@host:5432/mydb?sslmode=require

func BuildDSN(db DataBase) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
		db.SslMode,
	)
}

func MustLoadConfig() *Config {
	var cfg Config

	path := fetchConfigPath()
	if path == "" {
		panic("failed to read config path")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(fmt.Sprintf("config file does not exist: %s", path))
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(fmt.Sprintf("failed to read config: %s", err))
	}

	if err := validator(cfg); err != nil {
		panic(fmt.Sprintf("incorrect data config: %s", err))
	}

	return &cfg
}

func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "config-path", "", "path to config")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}

func validator(cfg Config) error {
	if cfg.Server.Port < 0 {
		return errors.New("server port should be positive")
	}
	if cfg.DB.Port < 0 {
		return errors.New("db port should be positive")
	}
	if cfg.DB.Password == "" {
		return errors.New("db password must not be empty")
	}
	return nil
}
