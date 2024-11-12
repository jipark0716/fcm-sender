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

var args *Args

func init() {
	args = &Args{
		KafkaConfig: &kafka.Config{
			Endpoint:      "127.0.0.1:9092",
			Topic:         "example-topic",
			ConsumerGroup: "example-consumer",
		},
	}
}

func main() {
	err := initService(args).Run()
	if err != nil {
		panic(err)
	}
}

func initService(args *Args) common.Service {
	return InitFirebaseMessageService(args.KafkaConfig, args.FirebaseConfig)
}
