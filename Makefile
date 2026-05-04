.PHONY: help tidy build docker-build docker-up docker-down docker-logs \
        k8s-setup k8s-build k8s-deploy k8s-teardown k8s-status \
        k8s-logs-user k8s-logs-order k8s-logs-notif

# ── Local dev ────────────────────────────────────────────────────────────────

help:
	@echo ""
	@echo "  make tidy           Run go mod tidy for all modules"
	@echo "  make build          Build all service binaries locally"
	@echo ""
	@echo "  make docker-build   Build Docker images"
	@echo "  make docker-up      Start everything with Docker Compose"
	@echo "  make docker-down    Stop Docker Compose"
	@echo "  make docker-logs    Tail all service logs"
	@echo ""
	@echo "  make k8s-setup      Start minikube + build images inside it"
	@echo "  make k8s-deploy     Apply all manifests to minikube"
	@echo "  make k8s-teardown   Delete all resources from minikube"
	@echo "  make k8s-status     Show pod/service status"
	@echo "  make k8s-url-user   Print user-service NodePort URL"
	@echo "  make k8s-url-order  Print order-service NodePort URL"
	@echo "  make k8s-url-notif  Print notification-service NodePort URL"
	@echo ""

tidy:
	cd shared && go mod tidy
	cd services/user-service && go mod tidy
	cd services/order-service && go mod tidy
	cd services/notification-service && go mod tidy
	go work sync

build:
	cd services/user-service && go build ./cmd/server
	cd services/order-service && go build ./cmd/server
	cd services/notification-service && go build ./cmd/server

# ── Docker ───────────────────────────────────────────────────────────────────

docker-build:
	docker build -t user-service:latest         -f services/user-service/Dockerfile .
	docker build -t order-service:latest        -f services/order-service/Dockerfile .
	docker build -t notification-service:latest -f services/notification-service/Dockerfile .

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down -v

docker-logs:
	docker compose logs -f user-service order-service notification-service

# ── Minikube / Kubernetes ─────────────────────────────────────────────────────

k8s-setup:
	minikube start --driver=docker
	# Build images directly inside minikube's docker daemon so imagePullPolicy: Never works
	eval $$(minikube docker-env) && \
		docker build -t user-service:latest         -f services/user-service/Dockerfile . && \
		docker build -t order-service:latest        -f services/order-service/Dockerfile . && \
		docker build -t notification-service:latest -f services/notification-service/Dockerfile .

k8s-build:
	eval $$(minikube docker-env) && \
		docker build -t user-service:latest         -f services/user-service/Dockerfile . && \
		docker build -t order-service:latest        -f services/order-service/Dockerfile . && \
		docker build -t notification-service:latest -f services/notification-service/Dockerfile .

k8s-deploy:
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/secrets.yaml
	kubectl apply -f k8s/configmap.yaml
	kubectl apply -f k8s/nats.yaml
	kubectl apply -f k8s/mysql.yaml
	kubectl apply -f k8s/user-service.yaml
	kubectl apply -f k8s/order-service.yaml
	kubectl apply -f k8s/notification-service.yaml

k8s-teardown:
	kubectl delete -f k8s/ --ignore-not-found

k8s-status:
	kubectl get pods,svc -n microservices

k8s-url-user:
	minikube service user-service -n microservices --url

k8s-url-order:
	minikube service order-service -n microservices --url

k8s-url-notif:
	minikube service notification-service -n microservices --url

k8s-logs-user:
	kubectl logs -f -l app=user-service -n microservices

k8s-logs-order:
	kubectl logs -f -l app=order-service -n microservices

k8s-logs-notif:
	kubectl logs -f -l app=notification-service -n microservices
