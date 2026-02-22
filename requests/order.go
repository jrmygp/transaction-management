package requests

type CreateOrderRequest struct {
	UserID  int `json:"userId" binding:"required"`
	HotelID int `json:"hotelId" binding:"required"`
	Nights  int `json:"nights" binding:"required"`
}

type MidtransWebhookRequest struct {
	OrderID           string `json:"order_id"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
	PaymentType       string `json:"payment_type"`
}
