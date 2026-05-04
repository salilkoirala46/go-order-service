package handlers

import (
	"net/http"

	"github.com/example/go-microservices/user-service/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.GetUint("user_id")
	user, err := h.svc.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}
