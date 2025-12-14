package handler

import (
	"order-service/internal/core/ports"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct { 
	svc ports.OrderService 
}

func NewOrderHandler(svc ports.OrderService) *OrderHandler { 
	return &OrderHandler{svc: svc} 
}

type BookRequest struct { 
	UserID string `json:"user_id" binding:"required"`; 
	EventID int `json:"event_id" binding:"required"`; 
	SeatID string `json:"seat_id" binding:"required"` 
}

func (h *OrderHandler) BookSeat(c *gin.Context) (int, string, interface{}) {
    var req BookRequest
    if err := c.ShouldBindJSON(&req); err != nil { return 400, err.Error(), nil }
    uid, _ := uuid.Parse(req.UserID)
    if err := h.svc.BookSeat(c.Request.Context(), uid, req.EventID, req.SeatID); 
		err != nil { 
			return 500, err.Error(), nil 
		}
    return 200, "Order processing", nil
}
