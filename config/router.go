package config

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jrmygp/transaction-management/controllers"
)

func NewRouter(hotelController *controllers.HotelController, orderController *controllers.OrderController) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	hotel := router.Group("/hotel")
	hotel.POST("/create-hotel", hotelController.CreateHotel)

	order := router.Group("/order")
	order.POST("/create-order", orderController.CreateOrder)
	order.POST("/bill-payment/:id", orderController.BillPayment)
	order.POST("/midtrans-webhook", orderController.MidtransWebhook)
	order.GET("/check-payment-status/:midtransOrderID", orderController.CheckPaymentStatus)
	order.GET("/refund/:id", orderController.RefundOrder)
	return router
}
