package infrastructure

import (
	"github.com/IBM/sarama"
	"time"
)

type KafkaProducerClient struct {
	Url string
}

func NewKafkaProducerClient(url string) *KafkaProducerClient {
	return &KafkaProducerClient{Url: url}
}

func (kpc *KafkaProducerClient) Connect() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Retry.Backoff = 500 * time.Millisecond

	conn, err := sarama.NewSyncProducer([]string{kpc.Url}, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
