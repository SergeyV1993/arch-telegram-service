package message_consumer

import (
	"arch-telegram-service/models"
	"context"
	"log"
)

type KafkaConsumerInterface interface {
	ConsumeMsg(ctx context.Context, topic string) <-chan models.OutgoingMessage
}

type MessageConsumer struct {
	telegramService TelegramService
	messageConsumer KafkaConsumerInterface
	topic           string
}

type TelegramService interface {
	SendMessage(ctx context.Context, message models.OutgoingMessage) error
}

func NewMessageConsumer(telegramService TelegramService, messageConsumer KafkaConsumerInterface, topic string) *MessageConsumer {
	return &MessageConsumer{
		telegramService: telegramService,
		messageConsumer: messageConsumer,
		topic:           topic,
	}
}

func (tc *MessageConsumer) Start(ctx context.Context) error {
	msgCh := tc.messageConsumer.ConsumeMsg(ctx, tc.topic)

	for msg := range msgCh {
		select {
		case <-ctx.Done():
			log.Printf("ctx is done")
			return nil
		default:
			if err := tc.telegramService.SendMessage(ctx, msg); err != nil {
				log.Printf("[ERR] ошибка обработки: %s", err.Error())
				continue
			}
		}
	}

	return nil
}
