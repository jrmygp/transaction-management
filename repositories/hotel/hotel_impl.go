package hotel

import (
	"github.com/jrmygp/transaction-management/models"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetAllHotels() ([]models.Hotel, error) {
	var hotels []models.Hotel

	err := r.db.Find(&hotels).Error

	return hotels, err
}

func (r *repository) CreateHotel(hotel models.Hotel) (models.Hotel, error) {
	err := r.db.Create(&hotel).Error
	return hotel, err
}

func (r *repository) GetHotelByID(id int) (models.Hotel, error) {
	var hotel models.Hotel
	err := r.db.First(&hotel, id).Error
	return hotel, err
}
