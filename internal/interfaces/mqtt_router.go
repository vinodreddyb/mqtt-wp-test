package interfaces

import (
	"strings"

	"mqtt-kafka-connector/internal/domain"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTRouter struct {
	telemetry chan<- domain.Message
	status    chan<- domain.Message
	command   chan<- domain.Message
	cmdexe    chan<- domain.Message
	boot      chan<- domain.Message
	groupCmd  chan<- domain.Message
}

func NewRouter(t, s, c, ce, b, gc chan<- domain.Message) *MQTTRouter {
	return &MQTTRouter{
		telemetry: t,
		status:    s,
		command:   c,
		cmdexe:    ce,
		boot:      b,
		groupCmd:  gc,
	}
}

func (r *MQTTRouter) Handler() mqtt.MessageHandler {
	return func(_ mqtt.Client, msg mqtt.Message) {
		m := domain.Message{Topic: msg.Topic(), Payload: msg.Payload()}
		switch {
		case msg.Topic() == "neevrfc/boot":

			r.boot <- m
		case strings.HasPrefix(msg.Topic(), "neevrfc/group"):
			r.groupCmd <- m
		case !strings.HasPrefix(msg.Topic(), "neevrfc/group") && strings.HasSuffix(msg.Topic(), "/cmd"):
			r.command <- m
		case strings.HasSuffix(msg.Topic(), "/telemetry"):
			r.telemetry <- m
		case strings.HasSuffix(msg.Topic(), "/status"):
			r.status <- m
		case strings.HasSuffix(msg.Topic(), "/cmdexe"):
			r.cmdexe <- m
		}
	}
}
