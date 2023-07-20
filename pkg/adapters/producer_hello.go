package adapters

import (
	"fmt"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

type ProducerHello struct {
	driver     *kafka.Writer
	BrokerUrls string `example:"localhost:9092,localhost:9093"`
	Topic      string
	ClientID   string
}

func (p *ProducerHello) Open() (*kafka.Writer, error) {
	if p.driver == nil {
		return nil, fmt.Errorf("driver was failed to connected")
	}
	return p.driver, nil
}

func (p *ProducerHello) Connect() error {
	brokerUrls := strings.Split(p.BrokerUrls, ",")

	if len(brokerUrls) == 0 {
		return fmt.Errorf("broker urls is empty")
	}

	if p.Topic == "" {
		return fmt.Errorf("topic is empty")
	}

	if p.ClientID == "" {
		return fmt.Errorf("client id is empty")
	}

	dialer := &kafka.Dialer{
		Timeout:  time.Second * 10,
		ClientID: p.ClientID,
	}

	p.driver = kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokerUrls,
		Topic:   p.Topic,
		Dialer:  dialer,
	})

	return nil
}

func (k *ProducerHello) Disconnect() error {
	return k.driver.Close()
}

func WithProducerHello(driver Driver[*kafka.Writer]) Option {
	return func(a *Adapter) {
		if err := driver.Connect(); err != nil {
			panic(err)
		}

		open, err := driver.Open()
		if err != nil {
			panic(err)
		}

		a.ProducerHello = open
	}
}
