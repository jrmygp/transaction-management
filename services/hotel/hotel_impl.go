package hotel

import (
	"mime/multipart"
	"path/filepath"

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

func convertFileToPath(file *multipart.FileHeader) string {
	baseDirectory := "public/hotel/"

	filePath := filepath.Join(baseDirectory, file.Filename)

	return filePath
}

func (s *service) GetAllHotels() ([]models.Hotel, error) {
	hotels, err := s.repository.GetAllHotels()
	return hotels, err
}

func (s *service) CreateHotel(hotelForm requests.CreateHotelRequest) (models.Hotel, error) {
	hotel := models.Hotel{
		Name:  hotelForm.Name,
		Price: hotelForm.Price,
		Image: convertFileToPath(hotelForm.Image),
	}

	newHotel, err := s.repository.CreateHotel(hotel)
	return newHotel, err
}

func (s *service) GetHotelByID(id int) (models.Hotel, error) {
	hotel, err := s.repository.GetHotelByID(id)
	return hotel, err
}
