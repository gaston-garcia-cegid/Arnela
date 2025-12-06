# Backend Structure

```
backend/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── domain/                     # Core entities
│   │   ├── user.go                 # User entity (roles, auth)
│   │   └── client.go               # Client entity (Spanish fields)
│   ├── repository/                 # Data access interfaces
│   │   ├── user_repository.go      # User repository interface
│   │   ├── client_repository.go    # Client repository interface
│   │   ├── postgres/               # PostgreSQL implementations
│   │   │   ├── user_repository.go  # User repo with sqlx
│   │   │   └── client_repository.go # Client repo with sqlx
│   │   └── mocks/                  # Mock implementations for testing
│   │       ├── user_repository_mock.go
│   │       └── client_repository_mock.go
│   ├── service/                    # Business logic
│   │   ├── auth_service.go         # Registration, login, JWT
│   │   ├── user_service.go         # User CRUD operations
│   │   ├── client_service.go       # Client CRUD + Spanish validations
│   │   ├── auth_service_test.go    # Auth tests
│   │   ├── user_service_test.go    # User tests
│   │   └── client_service_test.go  # Client tests
│   ├── handler/                    # HTTP handlers (GIN)
│   │   ├── auth_handler.go         # Register, Login, Me
│   │   ├── user_handler.go         # User CRUD endpoints
│   │   └── client_handler.go       # Client CRUD endpoints
│   └── middleware/                 # HTTP middleware
│       └── auth_middleware.go      # JWT validation + role checks
├── pkg/
│   ├── database/                   # Database utilities
│   │   ├── postgres.go             # PostgreSQL connection (sqlx)
│   │   └── migrate.go              # Migration runner
│   └── utils/                      # Utilities
│       ├── jwt.go                  # JWT generation/validation
│       └── logger.go               # Logging utilities
├── migrations/                     # Database migrations
├── docs/                           # Swagger documentation
└── go.mod
```
