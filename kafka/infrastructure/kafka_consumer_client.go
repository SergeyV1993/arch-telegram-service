package infrastructure

import (
	"github.com/IBM/sarama"
	"time"
)

type KafkaConsumerClient struct {
	Url string
}

func NewKafkaConsumerClient(url string) *KafkaConsumerClient {
	return &KafkaConsumerClient{Url: url}
}

func (kcc *KafkaConsumerClient) Connect() (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.MaxWaitTime = 100 * time.Millisecond
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	backOff := time.Duration(250)
	config.Metadata.Retry.Backoff = backOff * time.Millisecond
	waitInMinutes := 60
	config.Metadata.Retry.Max = waitInMinutes * 60 * 1000 / int(backOff)

	conn, err := sarama.NewConsumer([]string{kcc.Url}, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
