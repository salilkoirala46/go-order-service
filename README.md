# Go Order Services

A full-stack microservices application built with Go and Vue.js, featuring JWT authentication, NATS pub/sub messaging, MySQL persistence, a Vue 3 frontend, Docker Compose, and Kubernetes (Minikube) support.

---

## Services

| Service | Port (local) | NodePort (k8s) | Responsibility |
|---|---|---|---|
| **frontend** | 3000 | 30000 | Vue 3 SPA — UI for all services |
| **user-service** | 8080 | 30080 | User registration, login, JWT issuance |
| **order-service** | 8081 | 30081 | Order management (CRUD + status updates) |
| **notification-service** | 8082 | 30082 | Event-driven notifications via pub/sub |

---

## Architecture

```
                        Browser
                           │
                           ▼
                    ┌─────────────┐
                    │   Frontend  │  Vue 3 SPA (nginx)
                    │   :3000     │  NodePort :30000
                    └──────┬──────┘
                           │  /api/* proxied to backend
          ┌────────────────┼────────────────┐
          ▼                ▼                ▼
  ┌───────────────┐ ┌─────────────┐ ┌──────────────────────┐
  │ user-service  │ │order-service│ │ notification-service  │
  │ :8080         │ │ :8081       │ │ :8082                 │
  │               │ │             │ │                       │
  │ POST /auth/   │ │ POST /orders│ │ GET  /notifications   │
  │ GET  /users/  │ │ GET  /orders│ │ PATCH/:id/read        │
  └───────┬───────┘ └──────┬──────┘ └──────────┬────────────┘
          │                │                    │
          └────────────────┼────────────────────┘
                           ▼
                      ┌─────────┐
                      │  NATS   │  pub/sub broker
                      └─────────┘
                           │
                    ┌──────┴──────┐
                    ▼             ▼
               ┌────────┐   ┌────────┐
               │ MySQL  │   │ MySQL  │  (shared instance,
               │ userdb │   │orderdb │   3 databases)
               └────────┘   └────────┘
```

### NATS Event Flow

| Event | Published by | Consumed by | Result |
|---|---|---|---|
| `user.created` | user-service | notification-service | WELCOME notification |
| `order.created` | order-service | notification-service | ORDER_CREATED notification |
| `order.status.updated` | order-service | notification-service | ORDER_STATUS_CHANGE notification |

---

## Tech Stack

| Layer | Technology |
|---|---|
| **Frontend** | Vue 3, Vue Router 4, Pinia, Axios, Tailwind CSS, Vite |
| **Backend language** | Go 1.21 |
| **HTTP framework** | [Gin](https://github.com/gin-gonic/gin) |
| **ORM** | [GORM](https://gorm.io) |
| **Database** | MySQL 8.0 (one instance, 3 separate databases) |
| **Pub/Sub** | [NATS](https://nats.io) |
| **Authentication** | JWT (`golang-jwt/jwt/v5`) + bcrypt |
| **Module sharing** | Go workspace (`go.work`) |
| **Containers** | Docker + Docker Compose |
| **Orchestration** | Kubernetes / Minikube |

---

## Project Structure

```
go-order-services/
├── go.work                               # Go workspace — links all Go modules
├── docker-compose.yml                    # Full stack: frontend + 3 services + MySQL + NATS
├── Makefile                              # All dev, Docker and Kubernetes commands
├── scripts/
│   └── mysql-init.sql                    # Creates userdb, orderdb, notificationdb
│
├── frontend/                             # Vue 3 SPA
│   ├── Dockerfile                        # Multi-stage: node build → nginx serve
│   ├── nginx.conf                        # SPA fallback + /api/* proxy to backends
│   ├── vite.config.js                    # Dev proxy config (mirrors nginx.conf)
│   ├── tailwind.config.js
│   └── src/
│       ├── api/                          # Axios instance + per-service API modules
│       │   ├── axios.js                  # JWT interceptor, 401 auto-logout
│       │   ├── auth.js
│       │   ├── orders.js
│       │   └── notifications.js
│       ├── stores/                       # Pinia state management
│       │   ├── auth.js
│       │   ├── orders.js
│       │   └── notifications.js
│       ├── router/index.js               # Vue Router + auth/guest guards
│       ├── components/
│       │   ├── AppLayout.vue             # Shell: sidebar + navbar + <RouterView>
│       │   ├── AppSidebar.vue            # Nav links + unread badge + logout
│       │   ├── AppNavbar.vue             # Page title + notification bell
│       │   ├── StatusBadge.vue           # Coloured order status pill
│       │   ├── LoadingSpinner.vue
│       │   └── EmptyState.vue
│       └── views/
│           ├── LoginView.vue             # JWT login
│           ├── RegisterView.vue          # Account creation
│           ├── DashboardView.vue         # Stats + recent orders + notif preview
│           ├── OrdersView.vue            # Filterable orders table
│           ├── CreateOrderView.vue       # New order form with live total
│           ├── OrderDetailView.vue       # Detail + status timeline + updater
│           ├── NotificationsView.vue     # All notifications, mark-read
│           └── ProfileView.vue           # Account info + order summary
│
├── shared/                               # Shared Go library module
│   └── pkg/
│       ├── events/events.go              # NATS topic constants + event structs
│       └── middleware/jwt.go             # JWT generate & validate
│
├── services/
│   ├── user-service/
│   │   ├── cmd/server/main.go
│   │   ├── internal/
│   │   │   ├── config/                   # Env-based config
│   │   │   ├── handlers/                 # auth_handler, user_handler
│   │   │   ├── middleware/               # JWT Gin middleware
│   │   │   ├── models/                   # User (GORM)
│   │   │   ├── repository/               # DB layer
│   │   │   └── service/                  # Business logic
│   │   └── Dockerfile
│   ├── order-service/
│   │   ├── cmd/server/main.go
│   │   ├── internal/
│   │   │   ├── config/
│   │   │   ├── handlers/                 # order_handler (CRUD + status)
│   │   │   ├── middleware/
│   │   │   ├── models/                   # Order + OrderStatus enum
│   │   │   ├── repository/
│   │   │   └── service/
│   │   └── Dockerfile
│   └── notification-service/
│       ├── cmd/server/main.go
│       ├── internal/
│       │   ├── config/
│       │   ├── handlers/                 # notification_handler
│       │   ├── middleware/
│       │   ├── models/                   # Notification + NotificationType enum
│       │   ├── repository/
│       │   └── service/
│       └── Dockerfile
│
└── k8s/
    ├── namespace.yaml                    # namespace: microservices
    ├── secrets.yaml                      # JWT secret + MySQL password (base64)
    ├── configmap.yaml                    # MySQL init SQL
    ├── mysql.yaml                        # StatefulSet + headless Service
    ├── nats.yaml                         # Deployment + ClusterIP Service
    ├── frontend.yaml                     # Deployment + NodePort :30000
    ├── user-service.yaml                 # Deployment + NodePort :30080
    ├── order-service.yaml                # Deployment + NodePort :30081
    └── notification-service.yaml         # Deployment + NodePort :30082
```

---

## Prerequisites

| Tool | Purpose |
|---|---|
| Go 1.21+ | Build Go services locally |
| Node.js 18+ | Build / run Vue frontend locally |
| Docker & Docker Compose | Containerised local dev |
| `kubectl` | Apply Kubernetes manifests |
| `minikube` | Local Kubernetes cluster |

---

## Running with Docker Compose

The simplest way to run everything — one command starts MySQL, NATS, all three Go services and the Vue frontend.

```bash
# Build images and start the full stack
make docker-up

# Tail logs (all services including frontend)
make docker-logs

# Tear down containers and volumes
make docker-down
```

| Service | URL |
|---|---|
| **Frontend (Vue SPA)** | http://localhost:3000 |
| User Service (API) | http://localhost:8080 |
| Order Service (API) | http://localhost:8081 |
| Notification Service (API) | http://localhost:8082 |
| NATS monitoring | http://localhost:8222 |

The frontend's nginx container proxies all `/api/*` requests to the correct backend service — no CORS configuration required.

---

## Running on Minikube

The frontend and all backend services run inside the cluster. The nginx container resolves backend service names (`user-service`, `order-service`, `notification-service`) via Kubernetes DNS, so no extra configuration is needed.

```bash
# Step 1 — Start minikube and build all 4 images inside its Docker daemon
make k8s-setup

# Step 2 — Deploy everything to the cluster
make k8s-deploy

# Step 3 — Verify all pods are Running
make k8s-status
```

Expected output after `k8s-status`:

```
NAME                                        READY   STATUS    RESTARTS
pod/frontend-xxxx                           1/1     Running   0
pod/user-service-xxxx                       1/1     Running   0
pod/order-service-xxxx                      1/1     Running   0
pod/notification-service-xxxx               1/1     Running   0
pod/mysql-0                                 1/1     Running   0
pod/nats-xxxx                               1/1     Running   0

NAME                           TYPE        PORT(S)
service/frontend               NodePort    80:30000/TCP
service/user-service           NodePort    8080:30080/TCP
service/order-service          NodePort    8081:30081/TCP
service/notification-service   NodePort    8082:30082/TCP
service/mysql                  ClusterIP   None
service/nats                   ClusterIP   4222/TCP
```

```bash
# Get the frontend URL (opens the SPA in your browser)
make k8s-url-frontend

# Get backend service URLs (for direct API testing)
make k8s-url-user
make k8s-url-order
make k8s-url-notif
```

```bash
# Tail logs per service
make k8s-logs-frontend
make k8s-logs-user
make k8s-logs-order
make k8s-logs-notif

# Tear down everything
make k8s-teardown
```

> **How images work in Minikube:** `make k8s-setup` runs `eval $(minikube docker-env)` to point Docker at Minikube's internal daemon, then builds all images there. The manifests use `imagePullPolicy: Never` so Kubernetes uses these local images directly without needing a registry.

> **Rebuilding after code changes:** Run `make k8s-build` then `kubectl rollout restart deployment -n microservices` to redeploy with fresh images.

---

## Frontend — Vue 3 SPA

### Pages

| Route | Page | Description |
|---|---|---|
| `/login` | Login | Sign in with email + password, receive JWT |
| `/register` | Register | Create a new account |
| `/dashboard` | Dashboard | Stats overview, recent orders, notification preview |
| `/orders` | Orders | Filterable table of all orders (by status tab) |
| `/orders/new` | New Order | Form with live order total calculation |
| `/orders/:id` | Order Detail | Full detail, visual progress timeline, status updater |
| `/notifications` | Notifications | All notifications with type icons, mark-as-read |
| `/profile` | Profile | Account info, order summary stats, logout |

### Frontend local development

```bash
cd frontend
npm install
npm run dev     # http://localhost:5173
```

The Vite dev server proxies all `/api/*` calls to the running backend services on ports 8080–8082, so you can develop the frontend against live backends without any extra setup.

```bash
# Build the production bundle
npm run build

# Preview the production build locally
npm run preview
```

### How nginx proxying works

Both in Docker Compose and Kubernetes, the built frontend is served by nginx. The `nginx.conf` defines:

```
/api/auth/*           →  http://user-service:8080/auth/*
/api/users/*          →  http://user-service:8080/users/*
/api/orders/*         →  http://order-service:8081/orders/*
/api/notifications/*  →  http://notification-service:8082/notifications/*
/*                    →  index.html  (SPA fallback)
```

Service names resolve automatically:
- **Docker Compose** — via Docker's internal DNS on the `microservices` network
- **Kubernetes** — via CoreDNS within the `microservices` namespace

---

## API Reference

> All protected endpoints require the header: `Authorization: Bearer <token>`

### User Service — `:8080`

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/auth/register` | No | Register + receive JWT |
| `POST` | `/auth/login` | No | Login + receive JWT |
| `GET` | `/users/me` | Yes | Get current user profile |

#### Register
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com","password":"secret123"}'
```
```json
{ "token": "<jwt>", "user": { "id": 1, "name": "Alice", "email": "alice@example.com" } }
```

#### Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"secret123"}'
```

---

### Order Service — `:8081`

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/orders` | Yes | Create a new order |
| `GET` | `/orders` | Yes | List all orders for the current user |
| `GET` | `/orders/:id` | Yes | Get a single order |
| `PATCH` | `/orders/:id/status` | Yes | Update order status |

Valid statuses: `pending` → `confirmed` → `shipped` → `delivered` (or `cancelled`)

```bash
# Create order
curl -X POST http://localhost:8081/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"product":"Widget","quantity":3,"price":9.99}'

# Update status
curl -X PATCH http://localhost:8081/orders/1/status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status":"confirmed"}'
```

---

### Notification Service — `:8082`

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/notifications` | Yes | List all notifications for the current user |
| `GET` | `/notifications/:id` | Yes | Get a single notification |
| `PATCH` | `/notifications/:id/read` | Yes | Mark notification as read |

```bash
curl http://localhost:8082/notifications -H "Authorization: Bearer $TOKEN"
curl -X PATCH http://localhost:8082/notifications/1/read -H "Authorization: Bearer $TOKEN"
```

---

## End-to-End Example

```bash
# 1. Register — triggers user.created → WELCOME notification
TOKEN=$(curl -s -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com","password":"secret123"}' \
  | jq -r .token)

# 2. Place an order — triggers order.created → ORDER_CREATED notification
curl -s -X POST http://localhost:8081/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"product":"Widget","quantity":3,"price":9.99}' | jq .

# 3. Update order status — triggers order.status.updated → ORDER_STATUS_CHANGE notification
curl -s -X PATCH http://localhost:8081/orders/1/status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status":"shipped"}' | jq .

# 4. Read all notifications — should show 3 entries
curl -s http://localhost:8082/notifications \
  -H "Authorization: Bearer $TOKEN" | jq .
```

---

## Configuration

### Backend services (environment variables)

| Variable | Default | Description |
|---|---|---|
| `SERVER_PORT` | `8080 / 8081 / 8082` | HTTP listen port |
| `DATABASE_DSN` | local connection string | MySQL DSN |
| `NATS_URL` | `nats://localhost:4222` | NATS server address |
| `JWT_SECRET` | `super-secret-jwt-key` | HMAC signing key — **change in production** |
| `JWT_EXPIRY_HOURS` | `24` | Token TTL in hours (user-service only) |

### Kubernetes secrets

Secrets are stored base64-encoded in [k8s/secrets.yaml](k8s/secrets.yaml). To rotate them:

```bash
echo -n "your-new-secret" | base64
# paste the output into k8s/secrets.yaml, then:
kubectl apply -f k8s/secrets.yaml
kubectl rollout restart deployment -n microservices
```

---

## Development

```bash
# Update Go module dependencies
make tidy

# Build all Go binaries locally (no Docker)
make build

# Frontend dev server with hot reload
cd frontend && npm run dev
```

### Go workspace

The project uses a [Go workspace](https://go.dev/ref/mod#workspaces) (`go.work`) so all modules share a single dependency graph locally. Each service has its own `go.mod` with a `replace` directive for the shared module (`../../shared`), which also works inside Docker where `GOWORK=off` is set so only the `replace` directive applies.
