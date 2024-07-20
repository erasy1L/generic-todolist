package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Configs struct {
	APP app
	DB  DB
}

type DB struct {
	Username string
	Password string
	Host     string
	Port     string
	Name     string
}

type app struct {
	Port string
	Path string
}

func New() (cfg Configs, err error) {
	root, err := os.Getwd()
	if err != nil {
		return
	}

	err = godotenv.Load(filepath.Join(root, ".env"))
	if err != nil {
		return
	}

	if err = envconfig.Process("DB", &cfg.DB); err != nil {
		return
	}

	if err = envconfig.Process("APP", &cfg.APP); err != nil {
		return
	}

	return
}
