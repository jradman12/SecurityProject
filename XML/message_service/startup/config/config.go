package config

import "os"

type Config struct {
	Port               string
	MessageDBHost      string
	MessageDBPort      string
	PublicKey          string
	UserCommandSubject string
	UserReplySubject   string
	NatsHost           string
	NatsUser           string
	NatsPort           string
	NatsPass           string
}

func NewConfig() *Config {

	return &Config{
		Port:               os.Getenv("MESSAGE_SERVICE_PORT"),
		MessageDBHost:      os.Getenv("MESSAGE_DB_HOST"),
		MessageDBPort:      os.Getenv("MESSAGE_DB_PORT"),
		PublicKey:          "PostService",
		UserCommandSubject: os.Getenv("USER_COMMAND_SUBJECT"),
		UserReplySubject:   os.Getenv("USER_REPLY_SUBJECT"),
		NatsPort:           os.Getenv("NATS_PORT"),
		NatsHost:           os.Getenv("NATS_HOST"),
		NatsPass:           os.Getenv("NATS_PASS"),
		NatsUser:           os.Getenv("NATS_USER"),
	}
}
