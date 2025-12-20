package usecase

import (
	"context"
	"mqtt-kafka-connector/internal/domain"
	"mqtt-kafka-connector/internal/infrastructure"
	"strings"
)

type CommandExecutorProcess struct {
	kafka      *infrastructure.KafkaProducer
	kafkaTopic string
}

func NewCommandExecutorProcess(k *infrastructure.KafkaProducer) *CommandExecutorProcess {
	return &CommandExecutorProcess{
		kafka:      k,
		kafkaTopic: "commandExe",
	}
}

func (s *CommandExecutorProcess) Process(ctx context.Context, msg domain.Message) error {
	// Example key = deviceId
	parts := strings.Split(msg.Topic, "/")
	deviceID := parts[1]
	return s.kafka.Publish(ctx, s.kafkaTopic, deviceID, msg.Payload)
}
