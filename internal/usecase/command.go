package usecase

import (
	"context"
	"mqtt-wp-test/internal/infrastructure"
	"strings"

	"mqtt-wp-test/internal/domain"
)

type CommandProcessor struct {
	kafka *infrastructure.KafkaProducer
}

func NewCommandProcessor(k *infrastructure.KafkaProducer) *CommandProcessor {
	return &CommandProcessor{
		kafka: k,
	}
}

func (c *CommandProcessor) Process(ctx context.Context, msg domain.Message) error {
	// Example key = deviceId
	parts := strings.Split(msg.Topic, "/")
	deviceID := parts[1]
	topic := parts[len(parts)-1]

	return c.kafka.Publish(ctx, topic, deviceID, msg.Payload)
}
