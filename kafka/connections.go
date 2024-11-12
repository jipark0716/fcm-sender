package kafka

import (
	"context"
	"errors"
	"github.com/IBM/sarama"
	"log"
	"sync"
	"time"
)

type Connection struct {
	*Config
	*Consumer
	Context context.Context
	Cancel  context.CancelFunc
	Wg      *sync.WaitGroup
}

func NewConnection(config *Config) *Connection {
	ctx, cancel := context.WithCancel(context.Background())

	return &Connection{
		config,
		&Consumer{
			make(chan string, config.Workers),
		},
		ctx,
		cancel,
		&sync.WaitGroup{},
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
		go func() {
			defer c.Wg.Done()
			for row := range c.Queue {
				log.Printf("dequeue %s", row)
				time.Sleep(time.Second * 2)
			}
		}()
	}

	return nil
}

func (c *Connection) Shutdown(context.Context) error {
	c.Cancel()
	c.Wg.Wait()
	return nil
}
