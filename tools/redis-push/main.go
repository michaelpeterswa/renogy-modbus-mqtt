package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/go-redis/redis/v8"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", "localhost", 6379),
		Password: "",
		DB:       0,
	})

	b, err := load("data/2023-08-10-1uwviv.bin")
	if err != nil {
		panic(err)
	}

	_, err = client.LPush(context.Background(), "data", string(b)).Result()
	if err != nil {
		panic(err)
	}
}

func load(fileName string) ([]byte, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return b, nil
}
