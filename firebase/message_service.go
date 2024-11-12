package firebase

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/jipark0716/fcm-sender/kafka"
	"log"
	"time"
)

type MessageService struct {
	FirebaseConfig *Config
	Consumer       *kafka.Connection
	Converter      kafka.Converter[Message]
}

func NewMessageService(firebaseConfig *Config, consumer *kafka.Connection) (m *MessageService) {
	m = &MessageService{
		FirebaseConfig: firebaseConfig,
		Consumer:       consumer,
		Converter:      kafka.NewJsonConverter[Message](),
	}
	consumer.Handler = m.Send

	return
}

func (m *MessageService) Send(message *sarama.ConsumerMessage, workerId int) {
	row, err := m.Converter.Convert(message)
	if err != nil {
		log.Printf("fail convert message %v", err)
		return
	}
	log.Printf("[%d] dequeue workerId:%d message:%s", time.Now().UnixNano()-StartTime, workerId, row)
}

var StartTime int64

func (m *MessageService) Run() error {
	StartTime = time.Now().UnixNano()
	err := m.Consumer.Run()
	if err != nil {
		log.Printf("consume fail %v", err)
		return err
	}

	return nil
}

func (m *MessageService) Shutdown(context context.Context) error {
	log.Printf("shutdown")

	err := m.Consumer.Shutdown(context)
	if err != nil {
		log.Printf("fail consumer shutdown with: %v", err)
		return err
	}

	return nil
}
