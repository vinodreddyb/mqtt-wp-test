package interfaces

import (
	"strconv"
	"strings"

	"mqtt-kafka-connector/internal/domain"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTRouter struct {
	telemetry chan<- domain.Message
	status    chan<- domain.Message
	cmdexe    chan<- domain.Message
	boot      chan<- domain.Message
}

func NewRouter(t, s, ce, b chan<- domain.Message) *MQTTRouter {
	return &MQTTRouter{
		telemetry: t,
		status:    s,

		cmdexe: ce,
		boot:   b,
	}
}

func (r *MQTTRouter) Handler() mqtt.MessageHandler {
	return func(_ mqtt.Client, msg mqtt.Message) {
		m := domain.Message{Topic: msg.Topic(), Payload: msg.Payload(), MsgId: strconv.Itoa(int(msg.MessageID()))}
		switch {
		case msg.Topic() == "neevrfc/boot":
			r.boot <- m
		case strings.HasSuffix(msg.Topic(), "/telemetry"):
			r.telemetry <- m
		case strings.HasSuffix(msg.Topic(), "/status"):
			r.status <- m
		case strings.HasSuffix(msg.Topic(), "/cmdexe"):
			r.cmdexe <- m
		}
	}
}
