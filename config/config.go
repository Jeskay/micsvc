package config

import (
	"fmt"
	"time"
)

type ServerConfig struct {
	Host         string `env:"HOST"`
	Port         string `env:"PORT"`
	ExpireAfter  int    `env:"TOKEN_EXPIRE"`
	SecretKey    string `env:"SECRET_KEY"`
	ConnTimeout  int    `env:"CONNECTION_TIMEOUT"`
	EventTopic   string `env:"EVENT_TOPIC"`
	KafkaAddress string `env:"KAFKA_ADDRESS"`
}

func (sc *ServerConfig) TokenExpiration() time.Duration {
	return time.Hour * time.Duration(sc.ExpireAfter)
}

func (sc *ServerConfig) ConnectionTimeout() time.Duration {
	return time.Second * time.Duration(sc.ConnTimeout)
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
