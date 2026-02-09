package db

import (
	"sync"

	dto "github.com/Jeskay/micsvc/internal/dto"
)

type UserStorage struct {
	storage map[int32]*dto.User
	sync.Mutex
}

func NewUserStorage() *UserStorage {
	return &UserStorage{storage: make(map[int32]*dto.User)}
}

func (s *UserStorage) Set(key int32, value *dto.User) error {
	s.Lock()
	s.storage[key] = value
	s.Unlock()
	return nil
}

func (s *UserStorage) Get(key int32) (ok bool, value *dto.User) {
	s.Lock()
	value, ok = s.storage[key]
	s.Unlock()
	return
}

func (s *UserStorage) GetAll() (users []*dto.User) {
	users = make([]*dto.User, len(s.storage))
	i := 0
	s.Lock()
	for _, v := range s.storage {
		users[i] = v
		i++
	}
	s.Unlock()
	return
}

func (s *UserStorage) Remove(key int32) error {
	delete(s.storage, key)
	return nil
}
