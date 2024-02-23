package rabbitmq

import (
	"common/config"
	"github.com/streadway/amqp"
)

type RabbitConnection interface {
	ConnectWithRetries(cfg *config.Config, retries int) (*amqp.Connection, error)
}
