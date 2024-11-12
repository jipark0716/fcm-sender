package kafka

import (
	"context"
	"errors"
	"github.com/IBM/sarama"
	"log"
	"sync"
)

type Connection struct {
	*Config
	*Consumer
	Context context.Context
	Cancel  context.CancelFunc
	Handler func(*sarama.ConsumerMessage, int)
	Wg      *sync.WaitGroup
}

func NewConnection(config *Config, consumer *Consumer) *Connection {
	ctx, cancel := context.WithCancel(context.Background())
	return &Connection{
		Cancel:   cancel,
		Context:  ctx,
		Config:   config,
		Wg:       &sync.WaitGroup{},
		Consumer: consumer,
	}
}

func (c *Connection) Run() error {
	client, err := sarama.NewConsumerGroup(
		[]string{c.Config.Endpoint},
		c.Config.ConsumerGroup,
		sarama.NewConfig(),
	)

	if err != nil {
		log.Printf("fail connect kafak with: %v", err)
		return err
	}

	go func() {
		defer close(c.Queue)
		for {
			if err := client.Consume(c.Context, []string{c.Config.Topic}, c.Consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					log.Printf("kafka closed %v", err)
					return
				}
				log.Panicf("Error from consumer: %v", err)
			}
			if c.Context.Err() != nil {
				log.Printf("kafka closed")
				return
			}
		}
	}()

	c.Wg.Add(c.Workers)
	for i := 0; i < c.Workers; i++ {
		go func(i int) {
			defer c.Wg.Done()
			for row := range c.Queue {
				c.Handler(row, i)
			}
		}(i)
	}

	return nil
}

func (c *Connection) Shutdown(context.Context) error {
	c.Cancel()
	c.Wg.Wait()
	return nil
}
