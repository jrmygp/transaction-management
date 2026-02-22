package order

import (
	"github.com/jrmygp/transaction-management/models"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateOrder(order models.OrderBook) (models.OrderBook, error) {
	err := r.db.Create(&order).Error
	return order, err
}

func (r *repository) FindOrderByID(id int) (models.OrderBook, error) {
	var order models.OrderBook
	err := r.db.First(&order, id).Error
	return order, err
}

func (r *repository) UpdateOrder(order models.OrderBook) (models.OrderBook, error) {
	err := r.db.Save(&order).Error
	return order, err
}

func (r *repository) FindByMidtransOrderID(midtransOrderID string) (models.OrderBook, error) {
	var order models.OrderBook
	err := r.db.Where("midtrans_order_id = ?", midtransOrderID).First(&order).Error
	return order, err
}
