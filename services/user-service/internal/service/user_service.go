package service

import (
	"encoding/json"
	"errors"

	"github.com/example/go-microservices/shared/pkg/events"
	sharedmw "github.com/example/go-microservices/shared/pkg/middleware"
	"github.com/example/go-microservices/user-service/internal/config"
	"github.com/example/go-microservices/user-service/internal/models"
	"github.com/example/go-microservices/user-service/internal/repository"
	"github.com/nats-io/nats.go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Register(name, email, password string) (*models.User, string, error)
	Login(email, password string) (*models.User, string, error)
	GetByID(id uint) (*models.User, error)
}

type userService struct {
	repo repository.UserRepository
	nc   *nats.Conn
	cfg  *config.Config
}

func NewUserService(repo repository.UserRepository, nc *nats.Conn, cfg *config.Config) UserService {
	return &userService{repo: repo, nc: nc, cfg: cfg}
}

func (s *userService) Register(name, email, password string) (*models.User, string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user := &models.User{Name: name, Email: email, Password: string(hashed)}
	if err := s.repo.Create(user); err != nil {
		return nil, "", errors.New("email already registered")
	}

	token, err := sharedmw.GenerateToken(user.ID, user.Email, s.cfg.JWTSecret, s.cfg.JWTExpiry)
	if err != nil {
		return nil, "", err
	}

	evt := events.UserCreatedEvent{ID: user.ID, Email: user.Email, Name: user.Name}
	data, _ := json.Marshal(evt)
	_ = s.nc.Publish(events.UserCreated, data)

	return user, token, nil
}

func (s *userService) Login(email, password string) (*models.User, string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("invalid credentials")
		}
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := sharedmw.GenerateToken(user.ID, user.Email, s.cfg.JWTSecret, s.cfg.JWTExpiry)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *userService) GetByID(id uint) (*models.User, error) {
	return s.repo.FindByID(id)
}
