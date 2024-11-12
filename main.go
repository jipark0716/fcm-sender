package main

import (
	"context"
	"github.com/jipark0716/fcm-sender/common"
	"github.com/jipark0716/fcm-sender/firebase"
	"github.com/jipark0716/fcm-sender/kafka"
	"os"
	"os/signal"
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
			Workers:       4,
		},
	}
}

func main() {
	killSign := make(chan os.Signal, 1)
	signal.Notify(killSign, os.Interrupt)

	service := initService(args)
	err := service.Run()
	if err != nil {
		panic(err)
	}

	<-killSign
	err = service.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
}

func initService(args *Args) common.Service {
	return InitFirebaseMessageService(args.KafkaConfig, args.FirebaseConfig)
}
