package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	dto "github.com/Jeskay/micsvc/internal/dto"
)

type Service struct {
	storage dto.UserRepository
}

func NewUserService(storage dto.UserRepository) *Service {
	return &Service{storage: storage}
}

func (s *Service) Add(user *dto.User) error {
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

func (s *Service) Get(id int32) (*dto.User, error) {
	ok, user := s.storage.Get(id)
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *Service) GetAll() ([]*dto.User, error) {
	return s.storage.GetAll(), nil
}

func (s *Service) Update(id int32, user *dto.User) error {
	ok, _ := s.storage.Get(id)
	if !ok {
		return errors.New("user not found")
	}
	return s.storage.Set(id, user)
}

func (s *Service) Delete(id int32) error {
	return s.storage.Remove(id)
}
