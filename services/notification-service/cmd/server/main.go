package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/example/go-microservices/shared/pkg/events"
	"github.com/example/go-microservices/notification-service/internal/config"
	"github.com/example/go-microservices/notification-service/internal/handlers"
	"github.com/example/go-microservices/notification-service/internal/middleware"
	"github.com/example/go-microservices/notification-service/internal/models"
	"github.com/example/go-microservices/notification-service/internal/repository"
	"github.com/example/go-microservices/notification-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open(mysql.Open(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	if err := db.AutoMigrate(&models.Notification{}); err != nil {
		log.Fatalf("db migrate: %v", err)
	}

	nc, err := nats.Connect(cfg.NATSUrl)
	if err != nil {
		log.Fatalf("nats connect: %v", err)
	}
	defer nc.Drain()

	repo := repository.NewNotificationRepository(db)
	svc := service.NewNotificationService(repo)

	// user.created → welcome notification
	nc.Subscribe(events.UserCreated, func(msg *nats.Msg) {
		var evt events.UserCreatedEvent
		if err := json.Unmarshal(msg.Data, &evt); err != nil {
			log.Printf("user.created decode: %v", err)
			return
		}
		message := fmt.Sprintf("Welcome to the platform, %s! Your account has been created.", evt.Name)
		if _, err := svc.Create(evt.ID, models.TypeWelcome, message); err != nil {
			log.Printf("create welcome notification: %v", err)
		}
		log.Printf("notification-service: welcome notification created for user %d", evt.ID)
	})

	// order.created → order confirmation notification
	nc.Subscribe(events.OrderCreated, func(msg *nats.Msg) {
		var evt events.OrderCreatedEvent
		if err := json.Unmarshal(msg.Data, &evt); err != nil {
			log.Printf("order.created decode: %v", err)
			return
		}
		message := fmt.Sprintf("Your order #%d for %dx %s (total: $%.2f) has been placed successfully.", evt.ID, evt.Quantity, evt.Product, evt.Total)
		if _, err := svc.Create(evt.UserID, models.TypeOrderCreated, message); err != nil {
			log.Printf("create order notification: %v", err)
		}
		log.Printf("notification-service: order notification created for user %d, order %d", evt.UserID, evt.ID)
	})

	// order.status.updated → status change notification
	nc.Subscribe(events.OrderStatusUpdated, func(msg *nats.Msg) {
		var evt events.OrderStatusUpdatedEvent
		if err := json.Unmarshal(msg.Data, &evt); err != nil {
			log.Printf("order.status.updated decode: %v", err)
			return
		}
		message := fmt.Sprintf("Your order #%d status has been updated to: %s.", evt.ID, evt.Status)
		if _, err := svc.Create(evt.UserID, models.TypeOrderStatusChange, message); err != nil {
			log.Printf("create status notification: %v", err)
		}
		log.Printf("notification-service: status notification created for user %d, order %d → %s", evt.UserID, evt.ID, evt.Status)
	})

	h := handlers.NewNotificationHandler(svc)

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok", "service": "notification-service"}) })

	notif := r.Group("/notifications")
	notif.Use(middleware.Auth(cfg.JWTSecret))
	{
		notif.GET("", h.List)
		notif.GET("/:id", h.Get)
		notif.PATCH("/:id/read", h.MarkRead)
	}

	log.Printf("notification-service listening on :%s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server: %v", err)
	}
}
