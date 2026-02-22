package order

import (
	"fmt"
	"os"
	"time"

	"github.com/jrmygp/transaction-management/grpcclient"
	"github.com/jrmygp/transaction-management/helper"
	"github.com/jrmygp/transaction-management/models"
	"github.com/jrmygp/transaction-management/repositories/hotel"
	"github.com/jrmygp/transaction-management/repositories/order"
	"github.com/jrmygp/transaction-management/requests"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type service struct {
	repository      order.OrderRepository
	hotelRepository hotel.HotelRepository
}

func NewService(repository order.OrderRepository, hotelRepository hotel.HotelRepository) *service {
	return &service{repository, hotelRepository}
}

func (s *service) CreateOrder(orderForm requests.CreateOrderRequest) (models.OrderBook, error) {

	hotel, err := s.hotelRepository.GetHotelByID(orderForm.HotelID)
	if err != nil {
		return models.OrderBook{}, err
	}

	order := models.OrderBook{
		UserID:  orderForm.UserID,
		HotelID: orderForm.HotelID,
		Nights:  orderForm.Nights,
		Bill:    hotel.Price * orderForm.Nights,
		Status:  "pending",
	}

	newOrder, err := s.repository.CreateOrder(order)
	return newOrder, err
}

func (s *service) BillPayment(orderID int) (models.OrderBook, error) {
	order, err := s.repository.FindOrderByID(orderID)
	if err != nil {
		return order, err
	}

	snapClient := helper.NewSnapClient()

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  fmt.Sprintf("ORDER-%d", order.ID),
			GrossAmt: int64(order.Bill),
		},
		EnabledPayments: []snap.SnapPaymentType{
			snap.PaymentTypeCreditCard,
			snap.PaymentTypeBankTransfer,
			snap.PaymentTypeGopay,
		},
	}

	snapResp, err := snapClient.CreateTransaction(req)

	if err != nil && (snapResp == nil || snapResp.RedirectURL == "") {
		return order, err
	}

	order.MidtransOrderID = req.TransactionDetails.OrderID
	order.PaymentURL = snapResp.RedirectURL
	order.Status = "pending"

	return s.repository.UpdateOrder(order)
}

func (s *service) CheckPaymentStatus(midtransOrderID string) (string, error) {
	var c coreapi.Client
	c.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	resp, err := c.CheckTransaction(midtransOrderID)
	if err != nil {
		return "", err
	}

	return resp.TransactionStatus, nil
}

func (s *service) MidtransWebhook(notification requests.MidtransWebhookRequest) error {
	order, err := s.repository.FindByMidtransOrderID(notification.OrderID)
	if err != nil {
		return err
	}

	if order.Status == "paid" {
		return nil
	}

	switch notification.TransactionStatus {

	case "settlement":
		order.Status = "paid"
		now := time.Now()
		order.PaidAt = &now

	case "capture":
		if notification.FraudStatus == "accept" {
			order.Status = "paid"
			now := time.Now()
			order.PaidAt = &now
		}

	case "expire":
		order.Status = "expired"

	case "cancel":
		order.Status = "cancelled"

	case "deny":
		order.Status = "failed"
	}

	order.PaymentMethod = notification.PaymentType

	_, err = s.repository.UpdateOrder(order)
	return err
}

func (s *service) FindByMidtransOrderID(midtransOrderID string) (models.OrderBook, error) {
	return s.repository.FindByMidtransOrderID(midtransOrderID)
}

func (s *service) RefundOrder(orderID int) (models.OrderBook, error) {
	order, err := s.repository.FindOrderByID(orderID)
	if err != nil {
		return order, err
	}

	if order.Status == "refunded" {
		return order, fmt.Errorf("order already refunded")
	}

	if order.Status != "paid" && order.Status != "refund_pending" {
		return order, fmt.Errorf("order cannot be refunded")
	}

	order.Status = "refund_pending"
	if _, err := s.repository.UpdateOrder(order); err != nil {
		return order, err
	}

	userClient, conn, err := grpcclient.NewUserClient()
	if err != nil {
		return order, err
	}
	defer conn.Close()

	_, err = userClient.RefundBalance(
		int32(order.UserID),
		int32(order.Bill),
	)
	if err != nil {
		return order, err
	}

	order.Status = "refunded"
	return s.repository.UpdateOrder(order)
}
