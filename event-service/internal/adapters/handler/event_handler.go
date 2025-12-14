package handler

import (
	"event-service/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type EventHandler struct{ 
	svc ports.EventService 
}

func NewEventHandler(s ports.EventService) *EventHandler { 
	return &EventHandler{svc: s} 
}

func (h *EventHandler) GetAll(c *gin.Context) (int, string, interface{}) {
    events, err := h.svc.GetAll()
    if err != nil { return 500, "Failed to retrieve events", nil }
    return 200, "Events retrieved", events
}