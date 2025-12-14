package infrastructure

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducer struct {
	producer *kafka.Producer
}

func NewKafkaProducer(brokers string) (*KafkaProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		/*"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "SCRAM-SHA-512",
		"sasl.username":     "kafka_user",
		"sasl.password":     "kafkaTest!123",*/
		"acks":    "all",
		"retries": 5,
	})
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		producer: p,
	}, nil
}

func (k *KafkaProducer) Publish(ctx context.Context, topic string, key string, value []byte) error {
	delivery := make(chan kafka.Event, 1)

	err := k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(key),
		Value: value,
	}, delivery)

	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case e := <-delivery:
		m := e.(*kafka.Message)
		return m.TopicPartition.Error
	}
}

func (k *KafkaProducer) Close() {
	k.producer.Flush(5000)
	k.producer.Close()
}
