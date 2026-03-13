package requests

import "mime/multipart"

type CreateHotelRequest struct {
	Name  string                `form:"name" binding:"required"`
	Price int                   `form:"price" binding:"required"`
	Image *multipart.FileHeader `form:"image" binding:"required"`
}
