FROM golang:1.21-alpine AS builder

ENV GOWORK=off
WORKDIR /workspace

COPY shared/ ./shared/
COPY services/user-service/ ./services/user-service/

WORKDIR /workspace/services/user-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /app/server ./cmd/server

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
ENTRYPOINT ["./server"]
