package broker

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/IBM/sarama"

	"github.com/Jeskay/micsvc/internal/dto"
)

type EventConsumer struct {
	consumer  sarama.Consumer
	topic     string
	pConsumer sarama.PartitionConsumer
	done      chan struct{}
}

func NewEventConsumer(consumer sarama.Consumer, topic string, partition int32, offset int64) (*EventConsumer, error) {
	partitionConsumer, err := consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		return nil, err
	}
	return &EventConsumer{
		consumer:  consumer,
		pConsumer: partitionConsumer,
		topic:     topic,
		done:      make(chan struct{}),
	}, nil
}

func (c *EventConsumer) Listen(out chan<- dto.UserEvent) {
	in := c.pConsumer.Messages()
	for {
		select {
		case msg, ok := <-in:
			log.Printf("Received message: %v", msg)
			if !ok {
				return
			}
			var event dto.UserEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				continue
			}
			out <- event
		case <-c.done:
			return
		}
	}
}

func (c *EventConsumer) Close() error {
	return errors.Join(c.pConsumer.Close(), c.consumer.Close())
}
