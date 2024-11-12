package kafka

import (
	"context"
	"errors"
	"github.com/IBM/sarama"
	"log"
)

type Connection struct {
	*Config
	*Consumer
}

type Consumer struct {
}

func (c Consumer) Setup(sarama.ConsumerGroupSession) error {
	log.Printf("kafka connect")
	return nil
}

func (c Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	log.Printf("kafka cleanup")
	return nil
}

func (c Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

func NewConnection(config *Config) *Connection {
	return &Connection{
		config,
		&Consumer{},
	}
}

func (c *Connection) Consume() (chan string, error) {
	client, err := sarama.NewConsumerGroup(
		[]string{c.Config.Endpoint},
		c.Config.ConsumerGroup,
		sarama.NewConfig(),
	)

	if err != nil {
		log.Printf("fail connect kafak with: %v", err)
		return nil, err
	}

	for {
		if err := client.Consume(context.Background(), []string{c.Config.Topic}, c.Consumer); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				log.Printf("fail consume %v", err)
			}
			log.Panicf("Error from consumer: %v", err)
		}
	}

	return make(chan string), nil
}
