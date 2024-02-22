package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
	"upbit/internal/config"
	"upbit/pkg/log"
)

type Connection struct {
	instance *amqp.Connection
	cfg      *config.Config
	//once     sync.Once
}

func NewConnectWithRetries(cfg *config.Config) *Connection {
	return &Connection{cfg: cfg}
}

func (c *Connection) ConnectWithRetries(cfg *config.Config, retries int) (*amqp.Connection, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	URL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.Rabbit.Username,
		cfg.Rabbit.Password,
		cfg.Rabbit.Host,
		cfg.Rabbit.Port,
	)

	var conn *amqp.Connection
	var err error
	for i := 1; i <= retries; i++ {
		conn, err = amqp.Dial(URL)
		if err == nil {
			log.Logger.Info("Connected to RabbitMQ!")
			return conn, nil
		}

		log.Logger.Info(fmt.Sprintf("Failed to connect to RabbitMQ (attempt %d/%d). Retrying in 5 seconds...\n", i, retries))
		time.Sleep(5 * time.Second)
	}
	return nil, fmt.Errorf("unable to establish connection after %d retries: %v", retries, err)
}
