package firebase

import "github.com/jipark0716/fcm-sender/kafka"

type MessageService struct {
	KafkaConfig *kafka.Config
}

func NewMessageService(config *kafka.Config) *MessageService {
	return &MessageService{
		config,
	}
}

func (MessageService) Run() {
}
