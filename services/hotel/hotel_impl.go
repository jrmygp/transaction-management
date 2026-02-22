package hotel

import (
	"fmt"

	"github.com/jrmygp/transaction-management/grpc"
	"github.com/jrmygp/transaction-management/models"
	"github.com/jrmygp/transaction-management/repositories/hotel"
	"github.com/jrmygp/transaction-management/requests"
)

type service struct {
	repository hotel.HotelRepository
}

func NewService(repository hotel.HotelRepository) *service {
	return &service{repository}
}

func (s *service) CreateHotel(hotelForm requests.CreateHotelRequest) (models.Hotel, error) {
	hotel := models.Hotel{
		Name:  hotelForm.Name,
		Price: hotelForm.Price,
	}

	userClient, err := grpc.NewUserClient()
	if err != nil {
		panic(err)
	}

	res, err := userClient.GetUserByID(1)
	if err != nil {
		panic(err)
	}
	fmt.Println("User from user-service:", res.Username)

	newHotel, err := s.repository.CreateHotel(hotel)
	return newHotel, err
}

func (s *service) GetHotelByID(id int) (models.Hotel, error) {
	hotel, err := s.repository.GetHotelByID(id)
	return hotel, err
}
