package pull

import (
	"context"
	"fmt"

	gorenogymodbus "github.com/michaelpeterswa/go-renogy-modbus"
	"github.com/michaelpeterswa/renogy-modbus-mqtt/internal/redis"
)

type Puller interface {
	Pull(ctx context.Context) error
	Close() error
}

type RedisPuller struct {
	client    *redis.RedisClient
	queueName string
	outChan   chan gorenogymodbus.DynamicControllerInformation
}

func NewRedisPuller(client *redis.RedisClient, queueName string, outChan chan gorenogymodbus.DynamicControllerInformation) *RedisPuller {
	return &RedisPuller{
		client:    client,
		queueName: queueName,
		outChan:   outChan,
	}
}

func (r *RedisPuller) Pull(ctx context.Context) error {
	res, err := r.client.RPop(ctx, r.queueName)
	if err != nil {
		return fmt.Errorf("failed to rpop: %w", err)
	}

	dci, err := gorenogymodbus.Parse([]byte(res))
	if err != nil {
		return fmt.Errorf("failed to parse data: %w", err)
	}

	r.outChan <- *dci
	return nil
}

func (r *RedisPuller) Close() error {
	err := r.client.Close()
	if err != nil {
		return fmt.Errorf("failed to close redis client: %w", err)
	}
	return nil
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

type ModbusPuller struct {
	client  *gorenogymodbus.ModbusClient
	outChan chan gorenogymodbus.DynamicControllerInformation
}

func NewModbusPuller(client *gorenogymodbus.ModbusClient, outChan chan gorenogymodbus.DynamicControllerInformation) *ModbusPuller {
	return &ModbusPuller{
		client:  client,
		outChan: outChan,
	}
}

func (r *ModbusPuller) Pull(ctx context.Context) error {
	res, err := r.client.ReadData()
	if err != nil {
		return fmt.Errorf("failed to read data: %w", err)
	}

	dci, err := gorenogymodbus.Parse(res)
	if err != nil {
		return fmt.Errorf("failed to parse data: %w", err)
	}

	r.outChan <- *dci

	return nil
}

func (r *ModbusPuller) Close() error {
	// currently a no-op
	return nil
}
