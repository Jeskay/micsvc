package messager

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Service struct {
	logf      func(string)
	subs      map[string]*subscriber
	opTimeout time.Duration
	sync.Mutex
}

func NewMessagerSvc(logf func(string), opTimeout time.Duration) *Service {
	return &Service{logf: logf, subs: make(map[string]*subscriber), opTimeout: opTimeout}
}

func (s *Service) Subscribe(ctx context.Context, id string) (chan Message, error) {
	s.Lock()
	defer s.Unlock()
	if _, exists := s.subs[id]; exists {
		return nil, errors.New("already exists")
	}
	sub := &subscriber{msg: make(chan Message)}
	s.subs[id] = sub
	return sub.msg, nil
}

func (s *Service) Unsubscribe(ctx context.Context, id string) error {
	s.Lock()
	defer s.Unlock()
	sub, exists := s.subs[id]
	if !exists {
		return errors.New("not found")
	}
	close(sub.msg)
	delete(s.subs, id)
	return nil
}

func (s *Service) Message(ctx context.Context, msg Message) error {
	s.Lock()
	defer s.Unlock()
	timer := time.NewTimer(s.opTimeout)
	defer timer.Stop()
	for id, sub := range s.subs {
		timer.Reset(s.opTimeout)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			close(sub.msg)
			delete(s.subs, id)
			continue
		case sub.msg <- msg:
			if err := s.logMsg(msg); err != nil {
				return err
			}
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
