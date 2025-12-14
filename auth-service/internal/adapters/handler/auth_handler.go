package handler

import (
	"auth-service/internal/core/ports"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct { svc ports.AuthService }
func NewAuthHandler(svc ports.AuthService) *AuthHandler { return &AuthHandler{svc: svc} }

type LoginRequest struct { 
    Email string `json:"email" binding:"required"`; 
    Password string `json:"password" binding:"required"` 
}

type RegisterRequest struct { 
    Email string `json:"email" binding:"required"`; 
    Password string `json:"password" binding:"required"` 
}

func (h *AuthHandler) Login(c *gin.Context) (int, string, interface{}) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil { return 400, err.Error(), nil }
    token, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
    if err != nil { return 401, err.Error(), nil }
    return 200, "Login successful", gin.H{"token": token}
}

func (h *AuthHandler) Register(c *gin.Context) (int, string, interface{}) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil { return 400, err.Error(), nil }
    if err := h.svc.Register(c.Request.Context(), req.Email, req.Password); err != nil { return 500, err.Error(), nil }
    return 200, "User registered successfully", nil
}

func (h *AuthHandler) Me(c *gin.Context) (int, string, interface{}) {
    // Identity Propagation: Trust the header from Gateway
    uidStr := c.GetHeader("X-User-ID")
    if uidStr == "" { return 401, "Invalid User", nil }
    uid, _ := uuid.Parse(uidStr)

    log.Printf("Fetching user with ID: %s", uid.String())
    
    user, err := h.svc.GetUser(c.Request.Context(), uid)
    if err != nil { return 404, "User not found", nil }
    return 200, "User found", gin.H{"id": user.ID, "email": user.Email}
}


