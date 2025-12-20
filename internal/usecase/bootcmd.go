package usecase

import (
	"context"
	"mqtt-kafka-connector/internal/domain"
	"mqtt-kafka-connector/internal/infrastructure"
)

type BootCmdProcessor struct {
	kafka      *infrastructure.KafkaProducer
	kafkaTopic string
}

func NewBootCmdProcessor(k *infrastructure.KafkaProducer) *BootCmdProcessor {
	return &BootCmdProcessor{
		kafka:      k,
		kafkaTopic: "boot",
	}
}

func (b *BootCmdProcessor) Process(ctx context.Context, msg domain.Message) error {

	return b.kafka.Publish(ctx, b.kafkaTopic, msg.MsgId, msg.Payload)
}
