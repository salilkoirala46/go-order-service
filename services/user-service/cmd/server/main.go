package main

import (
	"log"

	"github.com/example/go-microservices/user-service/internal/config"
	"github.com/example/go-microservices/user-service/internal/handlers"
	"github.com/example/go-microservices/user-service/internal/middleware"
	"github.com/example/go-microservices/user-service/internal/models"
	"github.com/example/go-microservices/user-service/internal/repository"
	"github.com/example/go-microservices/user-service/internal/service"
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
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("db migrate: %v", err)
	}

	nc, err := nats.Connect(cfg.NATSUrl)
	if err != nil {
		log.Fatalf("nats connect: %v", err)
	}
	defer nc.Drain()

	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo, nc, cfg)

	authH := handlers.NewAuthHandler(svc)
	userH := handlers.NewUserHandler(svc)

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok", "service": "user-service"}) })

	auth := r.Group("/auth")
	{
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
	}

	users := r.Group("/users")
	users.Use(middleware.Auth(cfg.JWTSecret))
	{
		users.GET("/me", userH.GetMe)
	}

	log.Printf("user-service listening on :%s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server: %v", err)
	}
}
