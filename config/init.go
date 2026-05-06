package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresHost   string
	PostgresPort   string
	PostgresUser   string
	PostgresPass   string
	PostgresDBName string

	TelegramBotToken string

	OpenrouterToken string
	ModelMain       string
	ModelSmall      string
}

func Init() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		PostgresHost:   os.Getenv("PostgresHost"),
		PostgresPort:   os.Getenv("PostgresPort"),
		PostgresUser:   os.Getenv("PostgresUser"),
		PostgresPass:   os.Getenv("PostgresPass"),
		PostgresDBName: os.Getenv("PostgresDBName"),

		TelegramBotToken: os.Getenv("TelegramBotToken"),

		OpenrouterToken: os.Getenv("OpenrouterToken"),
		ModelMain:       os.Getenv("ModelMain"),
		ModelSmall:      os.Getenv("ModelSmall"),
	}, nil
}
