package messager

import (
	"context"
	"errors"
	"fmt"
)

type Service struct {
	logf func(string)
	subs map[string]*subscriber
}

func NewMessagerSvc(logf func(string)) *Service {
	return &Service{logf: logf, subs: make(map[string]*subscriber)}
}

func (s *Service) Subscribe(ctx context.Context, id string) (chan Message, error) {
	if _, exists := s.subs[id]; exists {
		return nil, errors.New("already exists")
	}
	sub := &subscriber{msg: make(chan Message)}
	s.subs[id] = sub
	return sub.msg, nil
}

func (s *Service) Unsubscribe(ctx context.Context, id string) error {
	_, exists := s.subs[id]
	if !exists {
		return errors.New("not found")
	}
	delete(s.subs, id)
	return nil
}

func (s *Service) Message(ctx context.Context, msg Message) error {
	for _, sub := range s.subs {
		sub.msg <- msg
		if err := s.logMsg(msg); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) logMsg(msg Message) error {
	if !msg.Binary {
		s.logf(fmt.Sprintf("%s: %v", msg.AuthorID, string(msg.Data)))
	} else {
		s.logf(fmt.Sprintf("%s: binary file of size %d", msg.AuthorID, len(msg.Data)))
	}

	return nil
}
