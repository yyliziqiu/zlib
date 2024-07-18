package zkafka

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func NewProducer(config Config) (*kafka.Producer, error) {
	producer, err := kafka.NewProducer(config.Map())
	if err != nil {
		return nil, fmt.Errorf("create producer failed [%v]", err)
	}

	go func(deliveredCallback func(partition kafka.TopicPartition), deliverFailedCallback func(partition kafka.TopicPartition)) {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error == nil {
					if deliveredCallback != nil {
						deliveredCallback(ev.TopicPartition)
					}
				} else {
					if deliverFailedCallback != nil {
						deliverFailedCallback(ev.TopicPartition)
					}
				}
			}
		}
	}(config.DeliveredCallback, config.DeliverFailedCallback)

	return producer, nil
}

func Produce(producer *kafka.Producer, topic string, message []byte) error {
	return producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, nil)
}

func ProduceObject(producer *kafka.Producer, topic string, object interface{}) error {
	message, err := json.Marshal(object)
	if err != nil {
		return err
	}
	return Produce(producer, topic, message)
}
