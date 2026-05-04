package service

import (
	"encoding/json"
	"errors"

	"github.com/example/go-microservices/shared/pkg/events"
	"github.com/example/go-microservices/order-service/internal/models"
	"github.com/example/go-microservices/order-service/internal/repository"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type OrderService interface {
	Create(userID uint, product string, quantity int, price float64) (*models.Order, error)
	GetByID(id, userID uint) (*models.Order, error)
	ListByUser(userID uint) ([]models.Order, error)
	UpdateStatus(id, userID uint, status models.OrderStatus) (*models.Order, error)
}

type orderService struct {
	repo repository.OrderRepository
	nc   *nats.Conn
}

func NewOrderService(repo repository.OrderRepository, nc *nats.Conn) OrderService {
	return &orderService{repo: repo, nc: nc}
}

func (s *orderService) Create(userID uint, product string, quantity int, price float64) (*models.Order, error) {
	order := &models.Order{
		UserID:   userID,
		Product:  product,
		Quantity: quantity,
		Price:    price,
		Total:    float64(quantity) * price,
		Status:   models.StatusPending,
	}
	if err := s.repo.Create(order); err != nil {
		return nil, err
	}

	evt := events.OrderCreatedEvent{
		ID: order.ID, UserID: order.UserID,
		Product: order.Product, Quantity: order.Quantity, Total: order.Total,
	}
	data, _ := json.Marshal(evt)
	_ = s.nc.Publish(events.OrderCreated, data)

	return order, nil
}

func (s *orderService) GetByID(id, userID uint) (*models.Order, error) {
	order, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	if order.UserID != userID {
		return nil, errors.New("forbidden")
	}
	return order, nil
}

func (s *orderService) ListByUser(userID uint) ([]models.Order, error) {
	return s.repo.FindByUserID(userID)
}

func (s *orderService) UpdateStatus(id, userID uint, status models.OrderStatus) (*models.Order, error) {
	order, err := s.GetByID(id, userID)
	if err != nil {
		return nil, err
	}
	if err := s.repo.UpdateStatus(order.ID, status); err != nil {
		return nil, err
	}
	order.Status = status

	evt := events.OrderStatusUpdatedEvent{ID: order.ID, UserID: order.UserID, Status: string(status)}
	data, _ := json.Marshal(evt)
	_ = s.nc.Publish(events.OrderStatusUpdated, data)

	return order, nil
}
