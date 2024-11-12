//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/jipark0716/fcm-sender/firebase"
	"github.com/jipark0716/fcm-sender/kafka"
)

func InitFirebaseMessageService(k *kafka.Config, f *firebase.Config) *firebase.MessageService {
	wire.Build(firebase.NewMessageService, kafka.NewConsumer, kafka.NewConnection)
	return &firebase.MessageService{}
}
