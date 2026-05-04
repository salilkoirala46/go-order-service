package repository

import (
	"github.com/example/go-microservices/order-service/internal/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *models.Order) error
	FindByID(id uint) (*models.Order, error)
	FindByUserID(userID uint) ([]models.Order, error)
	UpdateStatus(id uint, status models.OrderStatus) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) FindByID(id uint) (*models.Order, error) {
	var order models.Order
	if err := r.db.First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) UpdateStatus(id uint, status models.OrderStatus) error {
	return r.db.Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}
