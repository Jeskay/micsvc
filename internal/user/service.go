package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	dto "github.com/Jeskay/micsvc/internal/dto"
)

type EventSender interface {
	Send(event dto.UserEvent) error
}

type Service struct {
	notifySvc EventSender
	storage   dto.UserRepository
}

func NewUserService(storage dto.UserRepository, notifySvc EventSender) *Service {
	return &Service{storage: storage, notifySvc: notifySvc}
}

func (s *Service) Add(user *dto.User) (err error) {
	defer s.notify(user.ID, dto.Add, err)
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	if ok, _ := s.storage.Get(user.ID); !ok {
		return s.storage.Set(user.ID, user)
	}
	return errors.New("user already exists")
}

func (s *Service) Get(id int32) (user *dto.User, err error) {
	defer s.notify(id, dto.GetAll, err)
	ok, user := s.storage.Get(id)
	if !ok {
		return nil, errors.New("user not found")
	}
	return
}

func (s *Service) GetAll() (users []*dto.User, err error) {
	defer s.notify(-1, dto.GetAll, err)
	return s.storage.GetAll(), nil
}

func (s *Service) Update(id int32, user *dto.User) (err error) {
	defer s.notify(id, dto.Update, err)
	ok, _ := s.storage.Get(id)
	if !ok {
		return errors.New("user not found")
	}
	return s.storage.Set(id, user)
}

func (s *Service) Delete(id int32) (err error) {
	defer s.notify(id, dto.Delete, err)
	return s.storage.Remove(id)
}

func (s *Service) notify(id int32, eventType dto.EventType, err error) error {
	event := dto.UserEvent{UserID: id, Event: eventType}
	if err != nil {
		event.Error = err.Error()
	}
	return s.notifySvc.Send(event)
}
