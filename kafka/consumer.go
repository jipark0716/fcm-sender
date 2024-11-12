package kafka

import (
	"github.com/IBM/sarama"
	"log"
)

type Consumer struct {
	Queue chan string
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
			session.MarkMessage(message, "")
			c.Queue <- string(message.Value)
		case <-session.Context().Done():
			return nil
		}
	}
}
