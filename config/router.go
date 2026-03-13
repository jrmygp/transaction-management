package config

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jrmygp/transaction-management/controllers"
)

func NewRouter(hotelController *controllers.HotelController, orderController *controllers.OrderController) *gin.Engine {
	router := gin.Default()

	router.Static("/public", "./public")
	router.MaxMultipartMemory = 8 << 20

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Authorization", "Content-Type"},
	}))

	hotel := router.Group("/hotel")
	hotel.GET("/", hotelController.GetAllHotels)
	hotel.POST("/create-hotel", hotelController.CreateHotel)

	order := router.Group("/order")
	order.POST("/create-order", orderController.CreateOrder)
	order.POST("/bill-payment/:id", orderController.BillPayment)
	order.POST("/midtrans-webhook", orderController.MidtransWebhook)
	order.POST("/mark-refunded/:id", orderController.MarkOrderRefunded)
	order.GET("/check-payment-status/:midtransOrderID", orderController.CheckPaymentStatus)
	order.GET("/refund/:id", orderController.RefundOrder)
	return router
}
