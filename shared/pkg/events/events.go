package events

// Subject constants for NATS topics
const (
	UserCreated        = "user.created"
	OrderCreated       = "order.created"
	OrderStatusUpdated = "order.status.updated"
)

type UserCreatedEvent struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type OrderCreatedEvent struct {
	ID       uint    `json:"id"`
	UserID   uint    `json:"user_id"`
	Product  string  `json:"product"`
	Quantity int     `json:"quantity"`
	Total    float64 `json:"total"`
}

type OrderStatusUpdatedEvent struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"user_id"`
	Status string `json:"status"`
}
