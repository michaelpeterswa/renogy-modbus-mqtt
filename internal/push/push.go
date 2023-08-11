package push

import (
	"encoding/json"

	gorenogymodbus "github.com/michaelpeterswa/go-renogy-modbus"
	"github.com/michaelpeterswa/renogy-modbus-mqtt/internal/mqtt"
)

type Pusher interface {
	Push() error
	Close() error
}

type MQTTPusher struct {
	client *mqtt.MQTTClient
	topic  string
	inChan chan gorenogymodbus.DynamicControllerInformation
}

func NewMQTTPusher(client *mqtt.MQTTClient, topic string, inChan chan gorenogymodbus.DynamicControllerInformation) *MQTTPusher {
	return &MQTTPusher{
		client: client,
		topic:  topic,
		inChan: inChan,
	}
}

func (m *MQTTPusher) Push() error {
	for dci := range m.inChan {
		// TODO: publish to mqtt
		dciJson, err := json.Marshal(dci)
		if err != nil {
			// log
			continue
		}

		token := m.client.Client.Publish(m.topic, 0, false, dciJson)
		go func() {
			<-token.Done()
			// if token.Error() != nil {
			// }
		}()
	}

	return nil
}

func (m *MQTTPusher) Close() error {
	return nil
}
