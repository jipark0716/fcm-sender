package firebase

import (
	"github.com/jipark0716/fcm-sender/kafka"
	"log"
)

type MessageService struct {
	FirebaseConfig *Config
	Consumer       *kafka.Connection
}

func NewMessageService(firebaseConfig *Config, consumer *kafka.Connection) *MessageService {
	return &MessageService{
		firebaseConfig,
		consumer,
	}
}

func (m *MessageService) Run() error {
	ch, err := m.Consumer.Consume()
	if err != nil {
		log.Printf("consume fail %v", err)
		return err
	}

	for row := range ch {
		println(row)
	}

	return nil
}
