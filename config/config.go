package config

import (
	"fmt"
	"time"
)

type ServerConfig struct {
	Host        string `env:"HOST"`
	Port        string `env:"PORT"`
	ExpireAfter int    `env:"TOKEN_EXPIRE"`
	SecretKey   string `env:"SECRET_KEY"`
}

func (sc *ServerConfig) TokenExpiration() time.Duration {
	return time.Hour * time.Duration(sc.ExpireAfter)
}

func (sc *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%s", sc.Host, sc.Port)
}

type ClientConfig struct {
	Host string `env:"HOST"`
	Port string `env:"PORT"`
}

func (c *ClientConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
