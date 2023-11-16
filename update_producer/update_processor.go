package update_producer

import (
	"arch-telegram-service/models"
	"context"
	"log"
	"time"
)

type KafkaProducerInterface interface {
	PushMsgToQueue(ctx context.Context, topic string, msgs []models.Update) error
}

type UpdateProcessor struct {
	kafkaProducer KafkaProducerInterface

	telegramService TelegramInterface

	config Config
}

type TelegramInterface interface {
	GetUpdatesFromTelegram() ([]models.Update, error)
}

type Config struct {
	Topic string
}

func NewUpdateProcessor(
	kafkaProducer KafkaProducerInterface,
	telegramService TelegramInterface,

	config Config,
) *UpdateProcessor {
	return &UpdateProcessor{
		kafkaProducer:   kafkaProducer,
		telegramService: telegramService,
		config:          config,
	}
}

func (c *UpdateProcessor) Start(ctx context.Context) error {
	for {
		gotEvents, err := c.telegramService.GetUpdatesFromTelegram()
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		c.HandleUpdates(ctx, gotEvents)
	}
}

func (c *UpdateProcessor) HandleUpdates(ctx context.Context, updates []models.Update) {
	if err := c.kafkaProducer.PushMsgToQueue(ctx, c.config.Topic, updates); err != nil {
		log.Printf("can't send event to kafka: %s", err.Error())
	}
}
