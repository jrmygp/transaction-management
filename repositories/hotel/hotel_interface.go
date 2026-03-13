package hotel

import "github.com/jrmygp/transaction-management/models"

type HotelRepository interface {
	GetAllHotels() ([]models.Hotel, error)
	CreateHotel(hotel models.Hotel) (models.Hotel, error)
	GetHotelByID(id int) (models.Hotel, error)
}
