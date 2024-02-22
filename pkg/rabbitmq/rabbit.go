package rabbitmq

import (
	"github.com/streadway/amqp"
	"upbit/internal/config"
)

type RabbitConnection interface {
	ConnectWithRetries(cfg *config.Config, retries int) (*amqp.Connection, error)
}
