package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jrmygp/transaction-management/models"
	"github.com/jrmygp/transaction-management/requests"
	"github.com/jrmygp/transaction-management/responses"
	"github.com/jrmygp/transaction-management/services/order"
)

type OrderController struct {
	service order.OrderService
}

func NewOrderController(service order.OrderService) *OrderController {
	return &OrderController{service}
}

func convertOrderResponse(o models.OrderBook) responses.CreateOrderResponse {
	return responses.CreateOrderResponse{
		ID:         o.ID,
		HotelID:    o.HotelID,
		UserID:     o.UserID,
		Nights:     o.Nights,
		Bill:       o.Bill,
		Status:     o.Status,
		PaymentURL: o.PaymentURL,
	}
}

func convertPaymentLinkResponse(o models.OrderBook) responses.PaymentLinkResponse {
	return responses.PaymentLinkResponse{
		PaymentURL: o.PaymentURL,
	}
}

func (h *OrderController) CreateOrder(c *gin.Context) {
	var orderForm requests.CreateOrderRequest
	err := c.ShouldBindJSON(&orderForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	order, err := h.service.CreateOrder(orderForm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errors": err.Error(),
		})
		return
	}

	webResponse := responses.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   convertOrderResponse(order),
	}

	c.JSON(http.StatusOK, webResponse)
}

func (h *OrderController) BillPayment(c *gin.Context) {
	orderID := c.Param("id")
	id, err := strconv.Atoi(orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}

	order, err := h.service.BillPayment(id)
	if err != nil {
		log.Printf("BillPayment error: %#v\n", err)

		c.JSON(500, gin.H{
			"errors": fmt.Sprintf("%v", err),
		})
		return
	}

	webResponse := responses.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   convertPaymentLinkResponse(order),
	}

	c.JSON(http.StatusOK, webResponse)
}

func (h *OrderController) MidtransWebhook(c *gin.Context) {
	var notif requests.MidtransWebhookRequest

	if err := c.ShouldBindJSON(&notif); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.service.MidtransWebhook(notif)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errors": err.Error(),
		})
		return
	}

	webResponse := responses.Response{
		Code:   http.StatusOK,
		Status: "OK",
	}

	c.JSON(http.StatusOK, webResponse)
}

func (h *OrderController) CheckPaymentStatus(c *gin.Context) {
	midtransOrderID := c.Param("midtransOrderID")
	status, err := h.service.CheckPaymentStatus(midtransOrderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errors": err.Error(),
		})
		return
	}

	webResponse := responses.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   status,
	}

	c.JSON(http.StatusOK, webResponse)
}
