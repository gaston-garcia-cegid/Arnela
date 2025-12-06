# 💻 Workflow Commands

## 13.1. Local Setup

```bash
# Start Docker services
docker-compose up -d

# Run API
go run cmd/api/main.go
```

## 13.2. Testing

```bash
# Run all tests
go test ./...

# Run specific package
go test ./internal/service/... -v
```

## 13.3. Swagger Regeneration

```bash
swag init -g cmd/api/main.go -o docs
```

## 13.4. Database Migrations

```bash
# Create migration
migrate create -ext sql -dir migrations -seq <name>

# Run migrations
migrate -path migrations -database "postgres://..." up
```
