package user

import (
	"errors"
	"log/slog"

	"golang.org/x/crypto/bcrypt"

	dto "github.com/Jeskay/micsvc/internal/dto"
)

type EventSender interface {
	Send(event dto.UserEvent) error
}

type Service struct {
	notifySvc EventSender
	storage   dto.UserRepository
	logger    *slog.Logger
}

func NewUserService(logger *slog.Logger, storage dto.UserRepository, notifySvc EventSender) *Service {
	return &Service{storage: storage, notifySvc: notifySvc, logger: logger}
}

func (s *Service) Add(user *dto.User) (err error) {
	defer func() {
		if notifyErr := s.notify(user.ID, dto.Add, err); notifyErr != nil {
			err = errors.Join(err, notifyErr)
		}
	}()
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
	defer func() {
		if notifyErr := s.notify(id, dto.GetAll, err); notifyErr != nil {
			err = errors.Join(err, notifyErr)
		}
	}()
	ok, user := s.storage.Get(id)
	if !ok {
		return nil, errors.New("user not found")
	}
	return
}

func (s *Service) GetAll() (users []*dto.User, err error) {
	defer func() {
		if notifyErr := s.notify(-1, dto.GetAll, err); notifyErr != nil {
			err = errors.Join(err, notifyErr)
		}
	}()
	return s.storage.GetAll(), nil
}

func (s *Service) Update(id int32, user *dto.User) (err error) {
	defer func() {
		if notifyErr := s.notify(id, dto.Update, err); notifyErr != nil {
			err = errors.Join(err, notifyErr)
		}
	}()
	ok, _ := s.storage.Get(id)
	if !ok {
		return errors.New("user not found")
	}
	return s.storage.Set(id, user)
}

func (s *Service) Delete(id int32) (err error) {
	defer func() {
		if notifyErr := s.notify(id, dto.Delete, err); notifyErr != nil {
			err = errors.Join(err, notifyErr)
		}
	}()
	return s.storage.Remove(id)
}

func (s *Service) notify(id int32, eventType dto.EventType, err error) error {
	event := dto.UserEvent{UserID: id, Event: eventType}
	if err != nil {
		s.logger.Error("Operation failed", slog.String("operation", string(eventType)), slog.Int("UserID", int(id)), slog.Any("error", err))
		event.Error = err.Error()
	}
	s.logger.Info("Operation completed", slog.String("operation", string(eventType)), slog.Int("UserID", int(id)))
	return s.notifySvc.Send(event)
}
