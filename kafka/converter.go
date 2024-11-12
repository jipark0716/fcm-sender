package kafka

import "github.com/IBM/sarama"

type Converter[T any] interface {
	Convert(*sarama.ConsumerMessage) (T, error)
}
