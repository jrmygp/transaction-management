package order

import (
	"github.com/jrmygp/transaction-management/models"
	"github.com/jrmygp/transaction-management/requests"
)

type OrderService interface {
	CreateOrder(orderForm requests.CreateOrderRequest) (models.OrderBook, error)
	BillPayment(orderID int) (models.OrderBook, error)
	MidtransWebhook(notification requests.MidtransWebhookRequest) error
	CheckPaymentStatus(midtransOrderID string) (string, error)
	FindByMidtransOrderID(midtransOrderID string) (models.OrderBook, error)
	RefundOrder(orderID int) (models.OrderBook, error)
}
