package handler

import (
	"net/http"
	"order-service/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc ports.OrderService
}

func NewHandler(svc ports.OrderService) *Handler {
	return &Handler{svc: svc}
}

type CreateOrderRequest struct {
	ProductID int `json:"product_id" binding:"required"`
}

// Handler returns values for the Response Wrapper
func (h *Handler) CreateOrder(c *gin.Context) (int, string, interface{}) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return http.StatusBadRequest, "Invalid Request Body", nil
	}

	userID := c.GetHeader("X-User-ID") // Trusted header from Gateway
	
	order, err := h.svc.PlaceOrder(c.Request.Context(), userID, req.ProductID)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}

	return http.StatusCreated, "Order Placed Successfully", order
}
