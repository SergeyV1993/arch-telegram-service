package kafka

import (
	"arch-telegram-service/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
)

type Producer struct {
	connection sarama.SyncProducer
}

func NewProducer(connection sarama.SyncProducer) *Producer {
	return &Producer{connection: connection}
}

func (p *Producer) PushMsgToQueue(ctx context.Context, topic string, msgs []models.Update) error {
	kafkaMsgs, err := p.createMsgs(ctx, topic, msgs)
	if err != nil {
		return fmt.Errorf("failed to create msgs: %w", err)
	}

	sErr := p.connection.SendMessages(kafkaMsgs)
	if sErr != nil {
		for _, err := range sErr.(sarama.ProducerErrors) {
			log.Println("Write to kafka failed:", err)
		}
		return fmt.Errorf("failed to send message: %w", sErr)
	}

	return nil
}

func (p *Producer) createMsgs(ctx context.Context, topic string, msgs []models.Update) ([]*sarama.ProducerMessage, error) {
	var kafkaMgs []*sarama.ProducerMessage
	for _, v := range msgs {
		encoded, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to encode message: %w", err)
		}

		msg := sarama.ByteEncoder(encoded)
		pm := &sarama.ProducerMessage{
			Topic: topic,
			Value: msg,
		}

		kafkaMgs = append(kafkaMgs, pm)
	}

	return kafkaMgs, nil
}
