package constants

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	AWS_SECRET_KEY string
	AWS_ACCESS_KEY string
	AWS_REGION     string
	AWS_BUCKET     string

	PORT string

	DB_HOST      string
	DB_USER      string
	DB_PASSWORD  string
	DB_PORT      string
	DB_NAME      string
	REDIS_SERVER string

	JWT_ACCESS_SECRET  string
	JWT_REFRESH_SECRET string

	FROM_EMAIL    string
	SMTP_HOST     string
	SMTP_PORT     string
	SMTP_USERNAME string
	SMTP_PASSWORD string

	RABBITMQ_SERVER string

	PAYSTACK_SECRET_KEY    string
	FLUTTERWAVE_SECRET_KEY string
	PAYMENT_CALLBACK_URL   string
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	} else {
		fmt.Println("Loaded .env file")
	}
}

func GetEnv() Env {
	return Env{
		AWS_SECRET_KEY:         os.Getenv("AWS_SECRET_KEY"),
		AWS_ACCESS_KEY:         os.Getenv("AWS_ACCESS_KEY"),
		AWS_REGION:             os.Getenv("AWS_REGION"),
		AWS_BUCKET:             os.Getenv("AWS_BUCKET"),
		PORT:                   os.Getenv("PORT"),
		DB_HOST:                os.Getenv("DB_HOST"),
		DB_USER:                os.Getenv("DB_USER"),
		DB_PASSWORD:            os.Getenv("DB_PASSWORD"),
		DB_PORT:                os.Getenv("DB_PORT"),
		DB_NAME:                os.Getenv("DB_NAME"),
		REDIS_SERVER:           os.Getenv("REDIS_SERVER"),
		JWT_ACCESS_SECRET:      os.Getenv("JWT_ACCESS_SECRET"),
		JWT_REFRESH_SECRET:     os.Getenv("JWT_REFRESH_SECRET"),
		FROM_EMAIL:             os.Getenv("FROM_EMAIL"),
		SMTP_HOST:              os.Getenv("SMTP_HOST"),
		SMTP_PORT:              os.Getenv("SMTP_PORT"),
		SMTP_USERNAME:          os.Getenv("SMTP_USERNAME"),
		SMTP_PASSWORD:          os.Getenv("SMTP_PASSWORD"),
		RABBITMQ_SERVER:        os.Getenv("RABBITMQ_SERVER"),
		PAYSTACK_SECRET_KEY:    os.Getenv("PAYSTACK_SECRET_KEY"),
		FLUTTERWAVE_SECRET_KEY: os.Getenv("FLUTTERWAVE_SECRET_KEY"),
		PAYMENT_CALLBACK_URL:   os.Getenv("PAYMENT_CALLBACK_URL"),
	}
}
