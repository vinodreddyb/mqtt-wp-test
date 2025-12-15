package usecase

import (
	"context"
	"mqtt-kafka-connector/internal/infrastructure"
	"strings"

	"mqtt-kafka-connector/internal/domain"
)

type StatusProcessor struct {
	kafka      *infrastructure.KafkaProducer
	kafkaTopic string
}

func NewStatusProcessor(k *infrastructure.KafkaProducer) *StatusProcessor {
	return &StatusProcessor{
		kafka:      k,
		kafkaTopic: "status",
	}
}

func (s *StatusProcessor) Process(ctx context.Context, msg domain.Message) error {
	// Example key = deviceId
	parts := strings.Split(msg.Topic, "/")
	deviceID := parts[1]
	return s.kafka.Publish(ctx, s.kafkaTopic, deviceID, msg.Payload)
}
