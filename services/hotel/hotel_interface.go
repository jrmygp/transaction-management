package hotel

import (
	"github.com/jrmygp/transaction-management/models"
	"github.com/jrmygp/transaction-management/requests"
)

type HotelService interface {
	GetAllHotels() ([]models.Hotel, error)
	CreateHotel(hotel requests.CreateHotelRequest) (models.Hotel, error)
	GetHotelByID(id int) (models.Hotel, error)
}
