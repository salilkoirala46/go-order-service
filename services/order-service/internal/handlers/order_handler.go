package handlers

import (
	"net/http"
	"strconv"

	"github.com/example/go-microservices/order-service/internal/models"
	"github.com/example/go-microservices/order-service/internal/service"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	svc service.OrderService
}

func NewOrderHandler(svc service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

type createRequest struct {
	Product  string  `json:"product"  binding:"required"`
	Quantity int     `json:"quantity" binding:"required,min=1"`
	Price    float64 `json:"price"    binding:"required,gt=0"`
}

type updateStatusRequest struct {
	Status models.OrderStatus `json:"status" binding:"required"`
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order, err := h.svc.Create(c.GetUint("user_id"), req.Product, req.Quantity, req.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) List(c *gin.Context) {
	orders, err := h.svc.ListByUser(c.GetUint("user_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	order, err := h.svc.GetByID(uint(id), c.GetUint("user_id"))
	if err != nil {
		status := http.StatusNotFound
		if err.Error() == "forbidden" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req updateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order, err := h.svc.UpdateStatus(uint(id), c.GetUint("user_id"), req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}
