package usecase

import (
	"context"
	"mqtt-kafka-connector/internal/infrastructure"
	"strings"

	"mqtt-kafka-connector/internal/domain"
)

type CommandProcessor struct {
	kafka      *infrastructure.KafkaProducer
	kafkaTopic string
}

func NewCommandProcessor(k *infrastructure.KafkaProducer) *CommandProcessor {
	return &CommandProcessor{
		kafka:      k,
		kafkaTopic: "command",
	}
}

func (c *CommandProcessor) Process(ctx context.Context, msg domain.Message) error {
	// Example key = deviceId
	parts := strings.Split(msg.Topic, "/")
	deviceID := parts[1]

	return c.kafka.Publish(ctx, c.kafkaTopic, deviceID, msg.Payload)
}
