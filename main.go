package main

import (
	"log"

	"github.com/jrmygp/transaction-management/config"
	"github.com/jrmygp/transaction-management/controllers"
	transactiongrpc "github.com/jrmygp/transaction-management/grpc"
	"github.com/jrmygp/transaction-management/messaging"
	hotelRepo "github.com/jrmygp/transaction-management/repositories/hotel"
	orderRepo "github.com/jrmygp/transaction-management/repositories/order"
	hotelService "github.com/jrmygp/transaction-management/services/hotel"
	orderService "github.com/jrmygp/transaction-management/services/order"
)

func main() {
	db := config.DatabaseConnection()

	rabbitConn, rabbitChannel, err := messaging.NewConnection()
	if err != nil {
		log.Printf("rabbitmq disabled in transaction-engine: %v", err)
	}
	if rabbitChannel != nil {
		defer rabbitChannel.Close()
	}
	if rabbitConn != nil {
		defer rabbitConn.Close()
	}

	var publisher *messaging.Publisher
	if rabbitChannel != nil {
		publisher, err = messaging.NewPublisher(rabbitChannel)
		if err != nil {
			log.Printf("rabbitmq publisher unavailable: %v", err)
		}
	}

	hotelRepository := hotelRepo.NewRepository(db)
	hotelService := hotelService.NewService(hotelRepository)
	hotelController := controllers.NewHotelController(hotelService)

	orderRepository := orderRepo.NewRepository(db)
	orderService := orderService.NewService(orderRepository, hotelRepository, publisher)
	orderController := controllers.NewOrderController(orderService)

	router := config.NewRouter(hotelController, orderController)

	go transactiongrpc.StartGRPCServer(orderService)

	router.Run(":8081")
}
