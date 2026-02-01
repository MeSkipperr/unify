package http

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "unify-backend/internal/services"
)

type AuthHandler struct {
    AuthService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
    return &AuthHandler{AuthService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    token, err := h.AuthService.Login(req.Username, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}
