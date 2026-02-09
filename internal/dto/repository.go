package dto

type UserRepository interface {
	Get(key int32) (bool, *User)
	Set(key int32, value *User) error
	Remove(key int32) error
	GetAll() []*User
}
