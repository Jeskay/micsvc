package auth

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	jwt "github.com/golang-jwt/jwt/v5"

	"github.com/Jeskay/micsvc/config"
	"github.com/Jeskay/micsvc/internal/dto"
)

type Service struct {
	storage   dto.UserRepository
	expiresAt time.Duration
	secretKey string
}

func NewAuthService(config *config.ServerConfig, storage dto.UserRepository) *Service {
	return &Service{
		storage:   storage,
		secretKey: config.SecretKey,
		expiresAt: config.TokenExpiration(),
	}
}

func (s *Service) Login(id int32, password string) (bool, string) {
	ok, user := s.storage.Get(id)
	if !ok {
		return false, ""
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false, ""
	}
	token, err := s.CreateToken(id)
	if err != nil {
		return false, ""
	}
	return true, token
}

func (s *Service) Register(user *dto.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return s.storage.Set(user.ID, user)
}

func (s *Service) CreateToken(userID int32) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiresAt)),
		},
		UserID: userID,
	})
	return claims.SignedString([]byte(s.secretKey))
}

func (s *Service) VerifyToken(tokenString string) (int32, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, errors.New("invalid token")
	}
	return claims.UserID, nil
}
