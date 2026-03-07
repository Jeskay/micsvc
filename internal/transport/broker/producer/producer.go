package broker

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/IBM/sarama"

	"github.com/Jeskay/micsvc/internal/dto"
)

type EventProducer struct {
	topic string
	prod  sarama.AsyncProducer
	in    chan dto.UserEvent
	done  chan struct{}
	timeout time.Duration
}

func NewEventProducer(producer sarama.AsyncProducer, topic string, sendTimeout time.Duration) *EventProducer {
	return &EventProducer{
		prod:  producer,
		topic: topic,
		in:    make(chan dto.UserEvent, 5),
		done:  make(chan struct{}),
		timeout: sendTimeout,
	}
}

func (p *EventProducer) Run() error {
	for {
		select {
		case d, ok := <-p.in:
			if !ok {
				return nil
			}
			data, err := json.Marshal(d)
			if err != nil {
				continue
			}
			msg := &sarama.ProducerMessage{
				Topic: p.topic,
				Key:   sarama.ByteEncoder(strconv.Itoa(int(d.UserID))),
				Value: sarama.ByteEncoder(data),
			}
			p.prod.Input() <- msg
		case <-p.done:
			return nil
		case err := <-p.prod.Errors():
			return err
		}
	}
}

func (p *EventProducer) Shutdown() {
	close(p.done)
	close(p.in)
}

func (p *EventProducer) Send(event dto.UserEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	select {
	case p.in <- event:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
