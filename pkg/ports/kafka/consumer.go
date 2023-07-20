package kafkaconsumer

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

var (
	ErrConsumerAlreadyStarted = errors.New("consumer already started")
)

type Consumer struct {
	errCh chan error

	reader *kafka.Reader
}

func NewConsumer(brokerUrls, topic string, opts ...Option) *Consumer {
	c := Consumer{}

	config := kafka.ReaderConfig{
		Brokers:  strings.Split(brokerUrls, ","),
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		MaxWait:  10 * time.Second,
	}

	for _, opt := range opts {
		err := opt(&config)
		if err != nil {
			panic(err)
		}
	}

	c.reader = kafka.NewReader(config)

	return &c
}

func (c *Consumer) Listen() error {
	ctx := context.Background()

	go func() {
		for {
			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				c.errCh <- err
				return
			}

			timestamp := m.Time.UTC().Format(time.RFC3339)

			log.Info().Msgf("message at offset %d: %s\n", m.Offset, timestamp)
			log.Info().Msgf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		}
	}()

	return nil
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

func (c *Consumer) Error() chan error {
	return c.errCh
}
