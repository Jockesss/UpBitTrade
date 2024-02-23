package rabbitmq

import (
	"common/config"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type Producer struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	cfg  *config.Config
}

func NewProducer(cfg *config.Config, connection RabbitConnection) (*Producer, error) {
	retries := 5
	conn, err := connection.ConnectWithRetries(cfg, retries)
	if err != nil {
		return nil, fmt.Errorf("producer failed to connect to Rabbit: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close() // Ensure connection is closed on channel creation failure
		return nil, fmt.Errorf("producer failed to open a channel: %s", err)
	}

	_, err = ch.QueueDeclare(
		"trade_queue",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		return nil, fmt.Errorf("producer failed to declare Trade Queue: %s", err)
	}

	_, err = ch.QueueDeclare(
		"ticker_queue",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		return nil, fmt.Errorf("producer failed to declare Ticker Queue: %s", err)
	}

	return &Producer{
		conn: conn,
		ch:   ch,
	}, nil
}

// SendMessage publishes messages to specific queue
func (p *Producer) SendMessage(queue, message string) error {
	err := p.ch.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to send message to queue %s: %v", queue, err)
	}
	return nil
}

func (p *Producer) Close() {
	if p.ch != nil {
		if err := p.ch.Close(); err != nil {
			log.Printf("Error closing AMQP channel: %v", err)
		}
	}
	if p.conn != nil {
		if err := p.conn.Close(); err != nil {
			log.Printf("Error closing AMQP connection: %v", err)
		}
	}
	log.Println("RabbitMQ producer closed successfully")
}
