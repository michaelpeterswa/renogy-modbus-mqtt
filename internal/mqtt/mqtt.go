package mqtt

import (
	"crypto/tls"
	"errors"
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	ErrCouldNotConnectToMQTTServer = errors.New("could not connect to mqtt server")
)

type MQTTCLientConfig struct {
	Host     string
	ClientID string
	Username string
	Password string
}

type MQTTClient struct {
	Client MQTT.Client
}

func NewMQTTClientConfig(host, clientID, username, password string) *MQTTCLientConfig {
	return &MQTTCLientConfig{
		Host:     host,
		ClientID: clientID,
		Username: username,
		Password: password,
	}
}

func InitMQTT(config *MQTTCLientConfig) (*MQTTClient, error) {
	connOpts := MQTT.NewClientOptions().AddBroker(config.Host).SetClientID(config.ClientID).SetCleanSession(true)
	if config.Username != "" && config.Password != "" {
		connOpts.SetUsername(config.Username)
		connOpts.SetPassword(config.Password)
	}

	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("%w: %s", ErrCouldNotConnectToMQTTServer, token.Error())
	}

	return &MQTTClient{Client: client}, nil
}

func (m *MQTTClient) Noop() {}
