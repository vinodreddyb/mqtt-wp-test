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
		"localhost:9092",
	)
	defer kafkaProducer.Close()

	telemetryUC := usecase.NewTelemetryProcessor(kafkaProducer)
	statusUC := usecase.NewStatusProcessor(kafkaProducer)
	commandUC := usecase.NewCommandProcessor(kafkaProducer)
	telemetryPool := interfaces.NewWorkerPool(
		"telemetry", 5, 1000, telemetryUC,
	)
	statusPool := interfaces.NewWorkerPool(
		"status", 5, 500, statusUC,
	)
	commandPool := interfaces.NewWorkerPool(
		"command", 5, 100, commandUC,
	)
	commandExePool := interfaces.NewWorkerPool(
		"commandExe", 5, 100, commandUC,
	)
	groupPool := interfaces.NewWorkerPool(
		"group", 5, 100, commandUC,
	)
	bootPool := interfaces.NewWorkerPool(
		"boot", 5, 100, commandUC,
	)

	telemetryPool.Start(ctx, wg)
	statusPool.Start(ctx, wg)
	commandPool.Start(ctx, wg)
	commandExePool.Start(ctx, wg)
	groupPool.Start(ctx, wg)
	bootPool.Start(ctx, wg)

	router := interfaces.NewRouter(
		telemetryPool.Jobs(),
		statusPool.Jobs(),
		commandPool.Jobs(),
		commandExePool.Jobs(),
		bootPool.Jobs(),
		groupPool.Jobs(),
	)

	subs := []string{
		"$share/boot-group/neevrfc/boot",
		"$share/telemetry-group/neevrfc/+/+/cmdexe",
		"$share/telemetry-group/neevrfc/+/+/telemetry",
		"$share/status-group/neevrfc/+/+/status",
		"$share/cmd-group/neevrfc/+/+/cmd",
		"$share/cmd-group/neevrfc/group/+/cmd",
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
