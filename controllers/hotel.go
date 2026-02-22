package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jrmygp/transaction-management/models"
	"github.com/jrmygp/transaction-management/requests"
	"github.com/jrmygp/transaction-management/responses"
	"github.com/jrmygp/transaction-management/services/hotel"
)

type HotelController struct {
	service hotel.HotelService
}

func NewHotelController(service hotel.HotelService) *HotelController {
	return &HotelController{service}
}

func convertHotelResponse(o models.Hotel) responses.HotelResponse {
	return responses.HotelResponse{
		ID:    o.ID,
		Name:  o.Name,
		Price: o.Price,
	}
}

func (h *HotelController) CreateHotel(c *gin.Context) {
	var hotelForm requests.CreateHotelRequest

	err := c.ShouldBindJSON(&hotelForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	hotel, err := h.service.CreateHotel(hotelForm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errors": err.Error(),
		})
		return
	}

	webResponse := responses.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   convertHotelResponse(hotel),
	}

	c.JSON(http.StatusOK, webResponse)
}
