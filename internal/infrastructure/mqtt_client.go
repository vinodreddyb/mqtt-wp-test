package infrastructure

import mqtt "github.com/eclipse/paho.mqtt.golang"

type MQTTClient struct {
	client mqtt.Client
}

func NewMQTTClient(broker, clientID string, subs []string, handler mqtt.MessageHandler) *MQTTClient {

	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
		SetAutoReconnect(true).
		SetOnConnectHandler(func(c mqtt.Client) {
			for _, s := range subs {
				c.Subscribe(s, 1, handler)
			}
		})

	client := mqtt.NewClient(opts)
	client.Connect().Wait()

	return &MQTTClient{client: client}
}

func (m *MQTTClient) Disconnect() {
	m.client.Disconnect(250)
}
