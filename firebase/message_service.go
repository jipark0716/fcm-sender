package firebase

import (
	"context"
	"github.com/jipark0716/fcm-sender/kafka"
	"log"
	"os"
)

type MessageService struct {
	FirebaseConfig *Config
	Consumer       *kafka.Connection
	KillSign       chan os.Signal
}

func NewMessageService(firebaseConfig *Config, consumer *kafka.Connection) *MessageService {
	return &MessageService{
		firebaseConfig,
		consumer,
		make(chan os.Signal),
	}
}

func (m *MessageService) Run() error {
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
