package order

import "github.com/jrmygp/transaction-management/models"

type OrderRepository interface {
	CreateOrder(order models.OrderBook) (models.OrderBook, error)
	FindOrderByID(id int) (models.OrderBook, error)
	UpdateOrder(order models.OrderBook) (models.OrderBook, error)
	FindByMidtransOrderID(midtransOrderID string) (models.OrderBook, error)
}
