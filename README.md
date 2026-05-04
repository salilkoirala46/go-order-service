# Go Microservices

A production-ready microservices project built with Go, featuring JWT authentication, NATS pub/sub messaging, MySQL persistence, Docker, and Kubernetes (Minikube) support.

---

## Services

| Service | Port | Responsibility |
|---|---|---|
| **user-service** | 8080 | User registration, login, JWT issuance |
| **order-service** | 8081 | Order management (CRUD + status updates) |
| **notification-service** | 8082 | Event-driven notifications via pub/sub |

---

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌────────────────────────┐
│  User Service   │    │  Order Service  │    │  Notification Service  │
│  :8080          │    │  :8081          │    │  :8082                 │
│                 │    │                 │    │                        │
│ POST /auth/reg  │    │ POST /orders    │    │ GET  /notifications    │
│ POST /auth/login│    │ GET  /orders    │    │ GET  /notifications/:id│
│ GET  /users/me  │    │ GET  /orders/:id│    │ PATCH/:id/read         │
└────────┬────────┘    └───────┬─────────┘    └───────────┬────────────┘
         │                    │                           │
         └──────────┬─────────┘              subscribes  │
                    ▼                                     │
               ┌─────────┐ ◄───────────────────────────── ┘
               │  NATS   │   pub/sub broker
               └─────────┘
```

### NATS Event Flow

| Event | Published by | Consumed by | Notification created |
|---|---|---|---|
| `user.created` | user-service | notification-service | WELCOME message |
| `order.created` | order-service | notification-service | ORDER_CREATED message |
| `order.status.updated` | order-service | notification-service | ORDER_STATUS_CHANGE message |

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.21 |
| HTTP Framework | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io) |
| Database | MySQL 8.0 (one instance, 3 separate databases) |
| Pub/Sub | [NATS](https://nats.io) |
| Authentication | JWT (`golang-jwt/jwt/v5`) + bcrypt |
| Module sharing | Go workspace (`go.work`) |
| Containers | Docker + Docker Compose |
| Orchestration | Kubernetes / Minikube |

---

## Project Structure

```
go-microservices/
├── go.work                          # Go workspace (links all modules)
├── docker-compose.yml
├── Makefile
├── scripts/
│   └── mysql-init.sql               # Creates userdb, orderdb, notificationdb
├── shared/                          # Shared library module
│   └── pkg/
│       ├── events/events.go         # NATS topic constants + event structs
│       └── middleware/jwt.go        # JWT generate & validate helpers
├── services/
│   ├── user-service/
│   │   ├── cmd/server/main.go
│   │   ├── internal/
│   │   │   ├── config/              # Env-based configuration
│   │   │   ├── handlers/            # HTTP handlers (auth, user)
│   │   │   ├── middleware/          # JWT auth middleware (Gin)
│   │   │   ├── models/              # GORM model: User
│   │   │   ├── repository/          # DB access layer
│   │   │   └── service/             # Business logic
│   │   └── Dockerfile
│   ├── order-service/
│   │   ├── cmd/server/main.go
│   │   ├── internal/
│   │   │   ├── config/
│   │   │   ├── handlers/            # HTTP handlers (order CRUD)
│   │   │   ├── middleware/
│   │   │   ├── models/              # GORM model: Order (+ OrderStatus enum)
│   │   │   ├── repository/
│   │   │   └── service/
│   │   └── Dockerfile
│   └── notification-service/
│       ├── cmd/server/main.go
│       ├── internal/
│       │   ├── config/
│       │   ├── handlers/            # HTTP handlers (list, get, mark-read)
│       │   ├── middleware/
│       │   ├── models/              # GORM model: Notification (+ Type enum)
│       │   ├── repository/
│       │   └── service/
│       └── Dockerfile
└── k8s/
    ├── namespace.yaml
    ├── secrets.yaml                 # JWT secret + MySQL password (base64)
    ├── configmap.yaml               # MySQL init SQL
    ├── mysql.yaml                   # StatefulSet + headless Service
    ├── nats.yaml                    # Deployment + Service
    ├── user-service.yaml            # Deployment + NodePort :30080
    ├── order-service.yaml           # Deployment + NodePort :30081
    └── notification-service.yaml    # Deployment + NodePort :30082
```

---

## Prerequisites

- Go 1.21+
- Docker & Docker Compose
- `kubectl`
- `minikube` (for Kubernetes testing)

---

## Running with Docker Compose

```bash
# Build images and start all services
make docker-up

# Tail service logs
make docker-logs

# Tear down (removes volumes too)
make docker-down
```

Services will be available at:
- User Service: `http://localhost:8080`
- Order Service: `http://localhost:8081`
- Notification Service: `http://localhost:8082`
- NATS monitoring: `http://localhost:8222`

---

## Running on Minikube

```bash
# 1. Start minikube and build images inside its Docker daemon
make k8s-setup

# 2. Deploy all Kubernetes resources
make k8s-deploy

# 3. Check pod and service status
make k8s-status

# 4. Get service URLs
make k8s-url-user   # NodePort 30080
make k8s-url-order  # NodePort 30081
make k8s-url-notif  # NodePort 30082

# Tail logs for a service
make k8s-logs-user
make k8s-logs-order
make k8s-logs-notif

# Tear down all resources
make k8s-teardown
```

> **Note:** Images are built directly into minikube's Docker daemon (`eval $(minikube docker-env)`), so `imagePullPolicy: Never` works without a registry.

---

## API Reference

### Health checks (no auth required)

```
GET /health   →  each service exposes this
```

---

### User Service — `http://localhost:8080`

#### Register
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com","password":"secret123"}'
```

Response:
```json
{
  "token": "<jwt>",
  "user": { "id": 1, "name": "Alice", "email": "alice@example.com", "created_at": "..." }
}
```

#### Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"secret123"}'
```

#### Get current user (protected)
```bash
curl http://localhost:8080/users/me \
  -H "Authorization: Bearer <token>"
```

---

### Order Service — `http://localhost:8081`

All routes require `Authorization: Bearer <token>`.

#### Create order
```bash
curl -X POST http://localhost:8081/orders \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"product":"Widget","quantity":3,"price":9.99}'
```

#### List my orders
```bash
curl http://localhost:8081/orders \
  -H "Authorization: Bearer <token>"
```

#### Get order by ID
```bash
curl http://localhost:8081/orders/1 \
  -H "Authorization: Bearer <token>"
```

#### Update order status
```bash
curl -X PATCH http://localhost:8081/orders/1/status \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"status":"confirmed"}'
```

Valid statuses: `pending` `confirmed` `shipped` `delivered` `cancelled`

---

### Notification Service — `http://localhost:8082`

All routes require `Authorization: Bearer <token>`.

#### List my notifications
```bash
curl http://localhost:8082/notifications \
  -H "Authorization: Bearer <token>"
```

#### Get notification by ID
```bash
curl http://localhost:8082/notifications/1 \
  -H "Authorization: Bearer <token>"
```

#### Mark as read
```bash
curl -X PATCH http://localhost:8082/notifications/1/read \
  -H "Authorization: Bearer <token>"
```

---

## End-to-End Example

```bash
# 1. Register (triggers user.created → WELCOME notification)
TOKEN=$(curl -s -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com","password":"secret123"}' \
  | jq -r .token)

# 2. Create an order (triggers order.created → ORDER_CREATED notification)
curl -s -X POST http://localhost:8081/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"product":"Widget","quantity":3,"price":9.99}' | jq .

# 3. Update order status (triggers order.status.updated → ORDER_STATUS_CHANGE notification)
curl -s -X PATCH http://localhost:8081/orders/1/status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status":"shipped"}' | jq .

# 4. Read all notifications (should see 3: WELCOME, ORDER_CREATED, ORDER_STATUS_CHANGE)
curl -s http://localhost:8082/notifications \
  -H "Authorization: Bearer $TOKEN" | jq .
```

---

## Configuration

All services are configured via environment variables:

| Variable | Default | Description |
|---|---|---|
| `SERVER_PORT` | `8080/8081/8082` | HTTP listen port |
| `DATABASE_DSN` | local DSN | MySQL connection string |
| `NATS_URL` | `nats://localhost:4222` | NATS server address |
| `JWT_SECRET` | `super-secret-jwt-key` | HMAC signing secret — **change in production** |
| `JWT_EXPIRY_HOURS` | `24` | Token lifetime in hours (user-service only) |

---

## Development

```bash
# Install/update dependencies for all modules
make tidy

# Build binaries locally (no Docker)
make build
```

### Module structure

The project uses a [Go workspace](https://go.dev/ref/mod#workspaces) (`go.work`) so all modules share a single dependency graph during local development. Each service has its own `go.mod` with a `replace` directive pointing to `../../shared`, which also works inside Docker (`GOWORK=off` is set in Dockerfiles so only the `replace` directive is used).
