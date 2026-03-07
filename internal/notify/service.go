package notify

import (
	"context"
	"fmt"

	"github.com/Jeskay/micsvc/internal/dto"
	"github.com/Jeskay/micsvc/internal/messager"
)

type EventListener interface {
	Listen(out chan<- dto.UserEvent)
	Close() error
}

type Service struct {
	msgSvc    *messager.Service
	eListener EventListener
	done      chan struct{}
}

func NewService(msgSvc *messager.Service, eListener EventListener) *Service {
	return &Service{msgSvc, eListener, make(chan struct{})}
}

func (s *Service) Run(id string) (err error) {
	out := make(chan dto.UserEvent, 5)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go s.eListener.Listen(out)
	defer func() {
		if closeErr := s.eListener.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()
	for {
		select {
		case <-s.done:
			return nil
		case e, ok := <-out:
			if !ok {
				return nil
			}
			data := fmt.Appendf([]byte{}, "event %s on %d", e.Event, e.UserID)
			if e.Error != "" {
				data = fmt.Appendf(data, " with error %s", e.Error)
			}
			msg := messager.Message{AuthorID: id, Binary: false, Data: data}
			if err := s.msgSvc.Message(ctx, msg); err != nil {
				return err
			}
		}
	}
}

func (s *Service) Close() {
	close(s.done)
}
