package config

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

const (
	LogLevel     = "log.level"
	MQTTHost     = "mqtt.host"
	MQTTUsername = "mqtt.username"
	MQTTPassword = "mqtt.password"
	MQTTClientID = "mqtt.client.id"

	MQTTTopic         = "mqtt.topic"
	ModBusHost        = "modbus.host"
	ModBusIdleTimeout = "modbus.idle.timeout"

	PullMode = "pull.mode"
	PushMode = "push.mode"

	RedisHost     = "redis.host"
	RedisPassword = "redis.password"
	RedisDB       = "redis.db"

	PullCron = "pull.cron"

	RedisQueueName = "redis.queue.name"
)

func Get() (*koanf.Koanf, error) {
	k := koanf.New(".")

	err := k.Load(env.Provider("RMM_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "RMM_")), "_", ".", -1)
	}), nil)

	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return k, nil
}
