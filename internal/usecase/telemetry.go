package usecase

import (
	"context"
	"log/slog"
	"mqtt-kafka-connector/internal/domain"
	"mqtt-kafka-connector/internal/infrastructure"
	"strings"
)

type TelemetryProcessor struct {
	kafka *infrastructure.KafkaProducer
}

func NewTelemetryProcessor(k *infrastructure.KafkaProducer) *TelemetryProcessor {
	return &TelemetryProcessor{kafka: k}
}

func (t *TelemetryProcessor) Process(ctx context.Context, msg domain.Message) error {
	// Example key = deviceId
	parts := strings.Split(msg.Topic, "/")
	deviceID := parts[1]
	topic := parts[len(parts)-1]

	slog.Info("Received message on topic: ", topic, " message: ", string(msg.Payload))
	return t.kafka.Publish(ctx, topic, deviceID, msg.Payload)
}
