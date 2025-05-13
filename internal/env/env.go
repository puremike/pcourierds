package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	PORT, DB_ADDR, ENV string
}

func GetEnv() *EnvConfig {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv(":PORT")
	if port == "" {
		port = ":5100"
	}

	db_addr := os.Getenv("DB_ADDR")
	if db_addr == "" {
		log.Fatal("DB_ADDR is not set")
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	return &EnvConfig{
		PORT:    port,
		DB_ADDR: db_addr,
		ENV:     env,
	}
}
