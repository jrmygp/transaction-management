package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	channel *amqp091.Channel
	config  Config
}

type RefundRequestedMessage struct {
	OrderID         int       `json:"order_id"`
	UserID          int       `json:"user_id"`
	Amount          int       `json:"amount"`
	MidtransOrderID string    `json:"midtrans_order_id"`
	RequestedAt     time.Time `json:"requested_at"`
	RefundStatus    string    `json:"refund_status"`
	PaymentGateway  string    `json:"payment_gateway"`
}

func NewPublisher(channel *amqp091.Channel) (*Publisher, error) {
	config := loadConfig()
	if err := declareTopology(channel, config); err != nil {
		return nil, err
	}

	return &Publisher{
		channel: channel,
		config:  config,
	}, nil
}

func (p *Publisher) PublishRefundRequested(message RefundRequestedMessage) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal refund message: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return p.channel.PublishWithContext(
		ctx,
		p.config.ExchangeName,
		p.config.RefundRoutingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp091.Persistent,
			Timestamp:    time.Now(),
			Body:         body,
		},
	)
}
