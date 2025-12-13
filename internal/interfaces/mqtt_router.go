package interfaces

import (
	"strings"

	"mqtt-wp-test/internal/domain"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTRouter struct {
	telemetry chan<- domain.Message
	status    chan<- domain.Message
	command   chan<- domain.Message
}

func NewRouter(t, s, c chan<- domain.Message) *MQTTRouter {
	return &MQTTRouter{
		telemetry: t,
		status:    s,
		command:   c,
	}
}

func (r *MQTTRouter) Handler() mqtt.MessageHandler {
	return func(_ mqtt.Client, msg mqtt.Message) {
		m := domain.Message{Topic: msg.Topic(), Payload: msg.Payload()}

		switch {
		case strings.HasSuffix(msg.Topic(), "/telemetry"):
			r.telemetry <- m
		case strings.HasSuffix(msg.Topic(), "/status"):
			r.status <- m
		case strings.HasSuffix(msg.Topic(), "/cmd"):
			r.command <- m
		}
	}
}
