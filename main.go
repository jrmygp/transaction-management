package main

import (
	"github.com/jrmygp/transaction-management/config"
	"github.com/jrmygp/transaction-management/controllers"
	transactiongrpc "github.com/jrmygp/transaction-management/grpc"
	hotelRepo "github.com/jrmygp/transaction-management/repositories/hotel"
	orderRepo "github.com/jrmygp/transaction-management/repositories/order"
	hotelService "github.com/jrmygp/transaction-management/services/hotel"
	orderService "github.com/jrmygp/transaction-management/services/order"
)

func main() {
	db := config.DatabaseConnection()

	hotelRepository := hotelRepo.NewRepository(db)
	hotelService := hotelService.NewService(hotelRepository)
	hotelController := controllers.NewHotelController(hotelService)

	orderRepository := orderRepo.NewRepository(db)
	orderService := orderService.NewService(orderRepository, hotelRepository)
	orderController := controllers.NewOrderController(orderService)

	router := config.NewRouter(hotelController, orderController)

	go transactiongrpc.StartGRPCServer(orderService)

	router.Run(":8081")
}
