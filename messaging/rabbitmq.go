package messaging

import (
	"fmt"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

const (
	defaultRabbitMQURL      = "amqp://guest:guest@localhost:5672/"
	defaultExchangeName     = "hotel.events"
	defaultExchangeType     = "topic"
	defaultRefundQueue      = "user.refund.queue"
	defaultRefundRoutingKey = "refund.balance"
)

type Config struct {
	URL              string
	ExchangeName     string
	ExchangeType     string
	RefundQueue      string
	RefundRoutingKey string
}

func loadConfig() Config {
	return Config{
		URL:              getEnv("RABBITMQ_URL", defaultRabbitMQURL),
		ExchangeName:     getEnv("RABBITMQ_EXCHANGE", defaultExchangeName),
		ExchangeType:     getEnv("RABBITMQ_EXCHANGE_TYPE", defaultExchangeType),
		RefundQueue:      getEnv("RABBITMQ_REFUND_QUEUE", defaultRefundQueue),
		RefundRoutingKey: getEnv("RABBITMQ_REFUND_ROUTING_KEY", defaultRefundRoutingKey),
	}
}

func NewConnection() (*amqp091.Connection, *amqp091.Channel, error) {
	config := loadConfig()

	conn, err := amqp091.Dial(config.URL)
	if err != nil {
		return nil, nil, fmt.Errorf("connect rabbitmq: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("create rabbitmq channel: %w", err)
	}

	if err := declareTopology(channel, config); err != nil {
		channel.Close()
		conn.Close()
		return nil, nil, err
	}

	return conn, channel, nil
}

func declareTopology(channel *amqp091.Channel, config Config) error {
	if err := channel.ExchangeDeclare(
		config.ExchangeName,
		config.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("declare rabbitmq exchange: %w", err)
	}

	queue, err := channel.QueueDeclare(
		config.RefundQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("declare rabbitmq queue: %w", err)
	}

	if err := channel.QueueBind(
		queue.Name,
		config.RefundRoutingKey,
		config.ExchangeName,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("bind rabbitmq queue: %w", err)
	}

	return nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
