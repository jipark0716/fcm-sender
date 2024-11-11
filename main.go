package main

import (
	"github.com/jipark0716/fcm-sender/common"
	"github.com/jipark0716/fcm-sender/firebase"
	"github.com/jipark0716/fcm-sender/kafka"
)

type Args struct {
	KafkaConfig    *kafka.Config
	FirebaseConfig *firebase.Config
}

func main() {
	args := &Args{}
	service := initService(args)
	service.Run()
}

func initService(args *Args) common.Service {
	return InitFirebaseMessageService(args.KafkaConfig, args.FirebaseConfig)
}
