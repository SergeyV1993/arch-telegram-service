package kafka

import (
	"arch-telegram-service/models"
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
)

type Consumer struct {
	connection sarama.Consumer
}

func NewConsumer(connection sarama.Consumer) *Consumer {
	return &Consumer{connection: connection}
}

func (c *Consumer) ConsumeMsg(ctx context.Context, topic string) <-chan models.OutgoingMessage {
	consumer, err := c.connection.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return nil
	}

	dataCh := make(chan models.OutgoingMessage)
	go func() {
		defer close(dataCh)
		for {
			select {
			case err := <-consumer.Errors():
				log.Printf("consumer error message: %s", err.Error())
			case msg, ok := <-consumer.Messages():
				if !ok {
					return
				}
				var outMsg models.OutgoingMessage
				err := json.Unmarshal(msg.Value, &outMsg)
				if err != nil {
					log.Printf("unmarshal consumed msg: %s", err.Error())
					continue
				}

				dataCh <- outMsg
			case <-ctx.Done():
				return
			}
		}
	}()

	return dataCh
}
