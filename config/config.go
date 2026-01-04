package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Env string

const (
	EnvDev  Env = "dev"
	EnvTest Env = "test"
)

type Config struct {
	ApiserverPort    string `env:"APISERVER_PORT"`
	ApiserverHost    string `env:"APISERVER_HOST"`
	DatabaseName     string `env:"DB_NAME"`
	DatabaseHost     string `env:"DB_HOST"`
	DatabasePort     string `env:"DB_PORT"`
	DatabasePortTest string `env:"DB_PORT_TEST"`
	DatabaseUser     string `env:"DB_USER"`
	DatabasePassword string `env:"DB_PASSWORD"`
	Env              Env    `env:"ENV,default=dev"`
	ProjectRoot      string `env:"PROJECT_ROOT"`
	JWTSecret        string `env:"JWT_SECRET"`
}

func (c *Config) DatabaseUrl() string {
	port := c.DatabasePort
	if c.Env == EnvTest {
		port = c.DatabasePortTest
	}

	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		c.DatabaseUser,
		c.DatabasePassword,
		c.DatabaseHost,
		port,
		c.DatabaseName)
}

func New() (*Config, error) {
	envFilePath, err := filepath.Abs("./dev.env")
	if err != nil {
		return nil, fmt.Errorf("failed to get env file path: %w", err)
	}
	err = godotenv.Load(envFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load env: %w", err)
	}
	config := &Config{
		ApiserverPort:    os.Getenv("APISERVER_PORT"),
		ApiserverHost:    os.Getenv("APISERVER_HOST"),
		DatabaseName:     os.Getenv("DB_NAME"),
		DatabaseHost:     os.Getenv("DB_HOST"),
		DatabasePort:     os.Getenv("DB_PORT"),
		DatabasePortTest: os.Getenv("DB_PORT_TEST"),
		DatabaseUser:     os.Getenv("DB_USER"),
		DatabasePassword: os.Getenv("DB_PASSWORD"),
		Env:              Env(os.Getenv("ENV")),
		ProjectRoot:      os.Getenv("PROJECT_ROOT"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
	}
	return config, nil
}
