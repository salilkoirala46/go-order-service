package main

import (
	"encoding/json"
	"log"

	"github.com/example/go-microservices/shared/pkg/events"
	"github.com/example/go-microservices/order-service/internal/config"
	"github.com/example/go-microservices/order-service/internal/handlers"
	"github.com/example/go-microservices/order-service/internal/middleware"
	"github.com/example/go-microservices/order-service/internal/models"
	"github.com/example/go-microservices/order-service/internal/repository"
	"github.com/example/go-microservices/order-service/internal/service"
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
	if err := db.AutoMigrate(&models.Order{}); err != nil {
		log.Fatalf("db migrate: %v", err)
	}

	nc, err := nats.Connect(cfg.NATSUrl)
	if err != nil {
		log.Fatalf("nats connect: %v", err)
	}
	defer nc.Drain()

	// Subscribe to user.created — log new users known to this service
	if _, err := nc.Subscribe(events.UserCreated, func(msg *nats.Msg) {
		var evt events.UserCreatedEvent
		if err := json.Unmarshal(msg.Data, &evt); err != nil {
			log.Printf("user.created decode error: %v", err)
			return
		}
		log.Printf("order-service: new user registered id=%d email=%s", evt.ID, evt.Email)
	}); err != nil {
		log.Fatalf("nats subscribe: %v", err)
	}

	repo := repository.NewOrderRepository(db)
	svc := service.NewOrderService(repo, nc)
	h := handlers.NewOrderHandler(svc)

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok", "service": "order-service"}) })

	orders := r.Group("/orders")
	orders.Use(middleware.Auth(cfg.JWTSecret))
	{
		orders.POST("", h.Create)
		orders.GET("", h.List)
		orders.GET("/:id", h.Get)
		orders.PATCH("/:id/status", h.UpdateStatus)
	}

	log.Printf("order-service listening on :%s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server: %v", err)
	}
}
