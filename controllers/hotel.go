package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

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
	image := strings.ReplaceAll(o.Image, "\\", "/")

	imageURL := ""
	if image != "" {
		imageURL = fmt.Sprintf("http://localhost:8081/%s", image)
	}

	return responses.HotelResponse{
		ID:    o.ID,
		Name:  o.Name,
		Price: o.Price,
		Image: imageURL,
	}
}

func (h *HotelController) GetAllHotels(c *gin.Context) {
	projects, err := h.service.GetAllHotels()
	if err != nil {
		webResponse := responses.Response{
			Code:   http.StatusBadRequest,
			Status: "ERROR",
			Data:   err,
		}
		c.JSON(http.StatusBadRequest, webResponse)
		return
	}

	var hotelResponse []responses.HotelResponse

	if len(projects) == 0 {
		webResponse := responses.Response{
			Code:   http.StatusOK,
			Status: "OK",
			Data:   []responses.HotelResponse{},
		}
		c.JSON(http.StatusOK, webResponse)
		return
	}

	for _, hotel := range projects {
		response := convertHotelResponse(hotel)

		hotelResponse = append(hotelResponse, response)
	}

	webResponse := responses.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   hotelResponse,
	}

	c.JSON(http.StatusOK, webResponse)
}

func (h *HotelController) CreateHotel(c *gin.Context) {
	var hotelForm requests.CreateHotelRequest

	err := c.ShouldBind(&hotelForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File upload failed",
		})
		return
	}

	destination := "public/hotel/"
	filePath := filepath.Join(destination, file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save file",
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
