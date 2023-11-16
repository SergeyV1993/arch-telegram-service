package main

import (
	"arch-telegram-service/config"
	"arch-telegram-service/kafka"
	kafkaClient "arch-telegram-service/kafka/infrastructure"
	"arch-telegram-service/message_consumer"
	"arch-telegram-service/telegram"
	telegramClient "arch-telegram-service/telegram/infrastructure"
	"arch-telegram-service/update_producer"
	"context"
	"log"
)

func main() {
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	ctx := context.Background()

	//envs
	env, err := config.InitEnvs()
	if err != nil {
		log.Fatal("get envs failed", err)
	}

	//telegram service
	tc := telegramClient.NewTelegramClient(env.TelegramToken, env.PaginationOffset)
	telegramService := telegram.NewTelegramService(
		tc,
		telegram.ServiceConfig{
			Offset:                  config.DefaultTelegramOffset,
			Limit:                   config.DefaultTelegramLimit,
			DefaultPaginationOffset: env.PaginationOffset,
		},
	)

	//kafkaProducer
	kafkaProducerClient := kafkaClient.NewKafkaProducerClient(env.KafkaUrl)
	kafkaProducerConnect, kpErr := kafkaProducerClient.Connect()
	if kpErr != nil {
		log.Fatalf("kafka producer connection is failed: %v", kpErr)
	}
	defer func() {
		err := kafkaProducerConnect.Close()
		if err != nil {
			log.Fatalf("kafka producer connection close is failed: %v", err)
		}
	}()
	kafkaProducer := kafka.NewProducer(kafkaProducerConnect)
	log.Print("kafka producer started")

	//kafkaConsumer
	kafkaConsumerClient := kafkaClient.NewKafkaConsumerClient(env.KafkaUrl)
	kafkaConsumerConnect, kcErr := kafkaConsumerClient.Connect()
	if kcErr != nil {
		log.Fatal("kafka consumer connection is failed", kcErr)
	}
	defer func() {
		err := kafkaConsumerConnect.Close()
		if err != nil {
			log.Fatal("kafka consumer connection close is failed", err)
		}
	}()
	kafkaConsumer := kafka.NewConsumer(kafkaConsumerConnect)
	log.Print("kafka consumer started")

	//consume messages and send to telegram
	go func() {
		eventConsumer := message_consumer.NewMessageConsumer(telegramService, kafkaConsumer, env.KafkaTelegramTopic)
		if consErr := eventConsumer.Start(ctx); consErr != nil {
			log.Fatal("telegram event consumer is stopped", consErr)
		}
		log.Print("telegram event consumer started")
	}()

	//update processor
	log.Print("update processor started")
	updateProcessor := update_producer.NewUpdateProcessor(
		kafkaProducer,
		telegramService,
		update_producer.Config{
			Topic: env.KafkaTopic,
		},
	)
	if eventErr := updateProcessor.Start(ctx); eventErr != nil {
		log.Fatal("update processor is stopped", eventErr)
	}
}
