package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

const DefaultTelegramLimit = 50
const DefaultTelegramOffset = 0

type Env struct {
	TelegramToken string

	KafkaUrl   string
	KafkaTopic string

	KafkaTelegramTopic string

	PaginationOffset int
}

func NewEnv(token, kafkaUrl, kafkaTopic, kafkaTelegramTopic string, paginationOffset int) *Env {
	return &Env{
		TelegramToken:      token,
		KafkaUrl:           kafkaUrl,
		KafkaTopic:         kafkaTopic,
		KafkaTelegramTopic: kafkaTelegramTopic,
		PaginationOffset:   paginationOffset,
	}
}

func InitEnvs() (*Env, error) {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

	offset, err := strconv.Atoi(os.Getenv("PAGINATION_OFFSET"))
	if err != nil {
		return nil, fmt.Errorf("convert pagination offset:%w", err)
	}

	e := NewEnv(
		os.Getenv("TELEGRAMM_TOKEN"),
		os.Getenv("KAFKA_URL"),
		os.Getenv("KAFKA_TOPIC"),
		os.Getenv("KAFKA_TELEGRAM_TOPIC"),
		offset,
	)

	return e, nil
}
