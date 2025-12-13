package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/google/uuid"

	"mqtt-kafka-connector/internal/infrastructure"
	"mqtt-kafka-connector/internal/interfaces"
	"mqtt-kafka-connector/internal/usecase"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	kafkaProducer, _ := infrastructure.NewKafkaProducer(
		"b-1.demokafkacluster.63zxci.c5.kafka.us-east-1.amazonaws.com:9096,b-2.demokafkacluster.63zxci.c5.kafka.us-east-1.amazonaws.com:9096",
	)
	defer kafkaProducer.Close()

	telemetryUC := usecase.NewTelemetryProcessor(kafkaProducer)
	statucUC := usecase.NewStatusProcessor(kafkaProducer)
	commandUC := usecase.NewCommandProcessor(kafkaProducer)
	telemetryPool := interfaces.NewWorkerPool(
		"telemetry", 5, 1000, telemetryUC,
	)
	statusPool := interfaces.NewWorkerPool(
		"status", 5, 500, statucUC,
	)
	commandPool := interfaces.NewWorkerPool(
		"command", 5, 100, commandUC,
	)

	telemetryPool.Start(ctx, wg)
	statusPool.Start(ctx, wg)
	commandPool.Start(ctx, wg)

	router := interfaces.NewRouter(
		telemetryPool.Jobs(),
		statusPool.Jobs(),
		commandPool.Jobs(),
	)

	subs := []string{
		"$share/telemetry-group/neevrfc/+/+/telemetry",
		"$share/status-group/neevrfc/+/+/status",
		"$share/cmd-group/neevrfc/+/+/cmd",
	}

	client := infrastructure.NewMQTTClient(
		"tcp://localhost:1883",
		"consumer-"+uuid.NewString(),
		subs,
		router.Handler(),
	)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	cancel() // stop workers
	client.Disconnect()
	wg.Wait()
}
