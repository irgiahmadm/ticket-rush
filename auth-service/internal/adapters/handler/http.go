package handler

import (
	"auth-service/internal/core/ports"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc ports.AuthService
}

func NewHandler(svc ports.AuthService) *Handler {
	return &Handler{svc: svc}
}

type LoginRequest struct {
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) (int, string, interface{}) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return http.StatusBadRequest, "Invalid Body", nil
	}

	token, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		return http.StatusUnauthorized, err.Error(), nil
	}

	return http.StatusOK, "Success", gin.H{"token": token}
}

func (h *Handler) Register(c *gin.Context) (int, string, interface{}) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        return http.StatusBadRequest, "Invalid Body", nil
    }

    if err := h.svc.Register(c.Request.Context(), req.Email, req.Password); err != nil {
        return http.StatusInternalServerError, err.Error(), nil
    }

    return http.StatusOK, "Success", true
}
