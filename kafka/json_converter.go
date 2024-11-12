package kafka

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

type JsonConverter[T any] struct{}

func NewJsonConverter[T any]() *JsonConverter[T] {
	return &JsonConverter[T]{}
}

func (j *JsonConverter[T]) Convert(message *sarama.ConsumerMessage) (T, error) {
	var result T
	err := json.Unmarshal(message.Value, &result)
	return result, err
}
