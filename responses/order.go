package responses

type CreateOrderResponse struct {
	ID         int    `json:"id"`
	UserID     int    `json:"userId"`
	HotelID    int    `json:"hotelId"`
	Nights     int    `json:"nights"`
	Bill       int    `json:"bill"`
	Status     string `json:"status"`
	PaymentURL string `json:"paymentUrl"`
}

type PaymentLinkResponse struct {
	PaymentURL string `json:"paymentUrl"`
}
