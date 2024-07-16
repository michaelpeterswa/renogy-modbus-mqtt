package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	gorenogymodbus "github.com/michaelpeterswa/go-renogy-modbus"
	"github.com/michaelpeterswa/renogy-modbus-mqtt/internal/config"
	"github.com/michaelpeterswa/renogy-modbus-mqtt/internal/handlers"
	"github.com/michaelpeterswa/renogy-modbus-mqtt/internal/logging"
	"github.com/michaelpeterswa/renogy-modbus-mqtt/internal/mqtt"
	"github.com/michaelpeterswa/renogy-modbus-mqtt/internal/pull"
	"github.com/michaelpeterswa/renogy-modbus-mqtt/internal/push"
	"github.com/michaelpeterswa/renogy-modbus-mqtt/internal/redis"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func main() {
	var (
		mqttClient   *mqtt.MQTTClient
		modbusClient *gorenogymodbus.ModbusClient
		redisClient  *redis.RedisClient

		puller pull.Puller
		pusher push.Pusher
	)

	k, err := config.Get()
	if err != nil {
		log.Panicf("could not acquire config: %s", err.Error())
	}

	logger, err := logging.InitZap(k.String(config.LogLevel))
	if err != nil {
		log.Panicf("could not acquire zap logger: %s", err.Error())
	}
	logger.Info("renogy-modbus-mqtt init...")

	mqttConfig := mqtt.NewMQTTClientConfig(
		k.String(config.MQTTHost),
		k.String(config.MQTTClientID),
		k.String(config.MQTTUsername),
		k.String(config.MQTTPassword),
	)

	if k.String(config.PullMode) == "mqtt" || k.String(config.PushMode) == "mqtt" {
		mqttClient, err = mqtt.InitMQTT(mqttConfig)
		if err != nil {
			logger.Fatal("could not init mqtt client", zap.Error(err))
		}
	}

	if k.String(config.PullMode) == "redis" || k.String(config.PushMode) == "redis" {
		redisClient, err = redis.NewRedisClient(
			k.String(config.RedisHost),
			k.String(config.RedisPassword),
			k.Int(config.RedisDB),
		)
		if err != nil {
			logger.Fatal("could not init redis client", zap.Error(err))
		}
	}

	if k.String(config.PullMode) == "modbus" {
		modbusLogger, err := zap.NewStdLogAt(logger.Named("modbus"), zap.DebugLevel)
		if err != nil {
			logger.Fatal("could not init modbus logger", zap.Error(err))
		}
		modbusClient, err = gorenogymodbus.NewModbusClient(modbusLogger, k.String(config.ModBusHost), k.Duration(config.ModBusIdleTimeout))
		if err != nil {
			logger.Fatal("could not init modbus client", zap.Error(err))
		}
	}

	var dciChan = make(chan gorenogymodbus.DynamicControllerInformation, 100)

	switch k.String(config.PullMode) {
	case "modbus":
		puller = pull.NewModbusPuller(modbusClient, dciChan)
	case "redis":
		puller = pull.NewRedisPuller(redisClient, k.String(config.RedisQueueName), dciChan)
	default:
		logger.Fatal("invalid pull mode", zap.String("mode", k.String(config.PullMode)))
	}

	c := cron.New()
	_, err = c.AddFunc(k.String(config.PullCron), func() {
		err := puller.Pull(context.Background())
		if err != nil {
			logger.Error("could not pull", zap.Error(err))
		}
	})
	if err != nil {
		logger.Error("could not add cron job", zap.Error(err))
	}

	pusher = push.NewMQTTPusher(mqttClient, k.String(config.MQTTTopic), dciChan)

	go func() {
		err := pusher.Push()
		if err != nil {
			logger.Error("could not push", zap.Error(err))
		}
	}()

	c.Start()

	r := mux.NewRouter()
	r.HandleFunc("/healthcheck", handlers.HealthcheckHandler)
	r.Handle("/metrics", promhttp.Handler())
	http.Handle("/", r)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		logger.Fatal("could not start http server", zap.Error(err))
	}
}
