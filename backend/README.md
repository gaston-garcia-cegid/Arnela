# 🏥 Arnela Backend - CRM/CMS API

Backend API for Arnela Professional Office CRM/CMS system built with Go, GIN Framework, PostgreSQL, and Redis.

---

## 📋 Table of Contents

- [Tech Stack](#-tech-stack)
- [Project Structure](#-project-structure)
- [Getting Started](#-getting-started)
- [API Endpoints](#-api-endpoints)
- [Authentication](#-authentication)
- [Testing](#-testing)
- [Database](#-database)
- [Development](#-development)

---

## 🚀 Tech Stack

- **Go 1.23** - Programming language
- **GIN** - Web framework
- **sqlx** - SQL extensions for Go
- **PostgreSQL 16** - Primary database
- **Redis** - Caching & session management
- **JWT** - Authentication
- **Swagger/OpenAPI 3.0** - API documentation
- **testify** - Testing framework
- **Docker** - Containerization

---

## 📁 Project Structure

```
backend/
├── cmd/api/                    # Application entry point
├── internal/
│   ├── domain/                 # Core entities (User, Client, etc.)
│   ├── repository/             # Data access layer
│   │   ├── postgres/           # PostgreSQL implementations
│   │   └── mocks/              # Mock implementations for tests
│   ├── service/                # Business logic layer
│   ├── handler/                # HTTP handlers (controllers)
│   └── middleware/             # HTTP middleware (auth, logging)
├── pkg/
│   ├── database/               # Database utilities
│   └── utils/                  # Helper utilities (JWT, logger)
├── migrations/                 # Database migrations
├── docs/                       # Swagger documentation
├── go.mod
└── .env.example
```

---

## 🏁 Getting Started

### Prerequisites

- **Go 1.23+**
- **Docker & Docker Compose**
- **PostgreSQL 16** (or use Docker)
- **Redis** (or use Docker)

### Installation

1. **Clone the repository:**
```bash
git clone <repo-url>
cd arnela/backend
```

2. **Install dependencies:**
```bash
go mod download
```

3. **Setup environment:**
```bash
cp .env.example .env
```

Edit `.env` with your configuration:
```env
# Server
PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=arnela
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
```

4. **Start Docker services:**
```bash
docker-compose up -d
```

5. **Run migrations (automatic on startup):**
```bash
go run cmd/api/main.go
```

6. **Access the API:**
- **Backend:** http://localhost:8080
- **Swagger UI:** http://localhost:8080/swagger/index.html

---

## 🌐 API Endpoints

### Authentication

| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| POST | `/api/v1/auth/register` | Public | Register new user |
| POST | `/api/v1/auth/login` | Public | Login with credentials |
| GET | `/api/v1/auth/me` | Authenticated | Get current user |

**Example - Register:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@arnela.com",
    "password": "Admin123!",
    "firstName": "Admin",
    "lastName": "User",
    "role": "admin"
  }'
```

**Example - Login:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@arnela.com",
    "password": "Admin123!"
  }'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid",
    "email": "admin@arnela.com",
    "firstName": "Admin",
    "lastName": "User",
    "role": "admin",
    "isActive": true
  }
}
```

### User Management (Admin Only)

| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| POST | `/api/v1/users` | Admin | Create new user |
| GET | `/api/v1/users` | Admin | List all users |
| GET | `/api/v1/users/:id` | Admin | Get user by ID |
| PUT | `/api/v1/users/:id` | Admin | Update user |
| DELETE | `/api/v1/users/:id` | Admin | Delete user |

**Example - Create User:**
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "employee@arnela.com",
    "password": "Emp123!",
    "firstName": "John",
    "lastName": "Doe",
    "role": "employee"
  }'
```

### Client Management

| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| POST | `/api/v1/clients` | Admin, Employee | Create new client |
| GET | `/api/v1/clients` | Admin, Employee | List clients (with filters) |
| GET | `/api/v1/clients/:id` | Admin, Employee | Get client by ID |
| PUT | `/api/v1/clients/:id` | Admin, Employee | Update client |
| DELETE | `/api/v1/clients/:id` | Admin | Soft delete client |
| GET | `/api/v1/clients/me` | Client | Get own profile |

**Example - Create Client:**
```bash
curl -X POST http://localhost:8080/api/v1/clients \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Juan",
    "lastName": "Pérez",
    "email": "juan@example.com",
    "phone": "612345678",
    "dni": "12345678Z",
    "dateOfBirth": "1990-01-15",
    "address": "Calle Mayor 123",
    "city": "Madrid",
    "postalCode": "28001",
    "province": "Madrid",
    "notes": "Cliente preferente"
  }'
```

**Example - List Clients with Filters:**
```bash
# Get all clients (page 1, 20 per page)
curl "http://localhost:8080/api/v1/clients?page=1&pageSize=20" \
  -H "Authorization: Bearer $TOKEN"

# Search for "Juan"
curl "http://localhost:8080/api/v1/clients?search=Juan" \
  -H "Authorization: Bearer $TOKEN"

# Filter by city and active status
curl "http://localhost:8080/api/v1/clients?city=Madrid&isActive=true" \
  -H "Authorization: Bearer $TOKEN"

# Combine filters
curl "http://localhost:8080/api/v1/clients?search=Juan&city=Madrid&page=1&pageSize=10" \
  -H "Authorization: Bearer $TOKEN"
```

**Query Parameters:**
- `page` (default: 1) - Page number
- `pageSize` (default: 20, max: 100) - Items per page
- `search` - Search across firstName, lastName, email, dni, phone
- `isActive` - Filter by active status (true/false)
- `city` - Filter by city
- `province` - Filter by province

---

## 🔐 Authentication

### JWT Token

All protected endpoints require a JWT token in the `Authorization` header:

```bash
Authorization: Bearer <your-jwt-token>
```

**Token includes:**
- User ID (UUID)
- Email
- Role (admin, employee, client)
- Expiration (24 hours)

### Role-Based Access Control

| Role | Permissions |
| :--- | :--- |
| **admin** | Full access to all endpoints, can delete resources |
| **employee** | Can manage clients, appointments, tasks (no delete) |
| **client** | Can only access own profile and appointments |

---

## 🧪 Testing

### Run All Tests

```bash
go test ./...
```

### Run Specific Package Tests

```bash
# Test all services
go test ./internal/service/... -v

# Test specific service
go test ./internal/service/... -v -run TestClientService
```

### Run with Coverage

```bash
go test -cover ./...
```

### Test Results

```bash
$ go test ./internal/service/... -v

=== RUN   TestAuthService_Register
    ✓ successful registration
    ✓ email already exists
    ✓ weak password
    ... (5/5 passing)

=== RUN   TestAuthService_Login
    ✓ successful login
    ✓ invalid credentials
    ... (3/3 passing)

=== RUN   TestUserService_CreateUser
    ... (4/4 passing)

=== RUN   TestClientService_CreateClient
    ✓ successful creation
    ✓ email already exists
    ✓ DNI already exists
    ✓ invalid email format
    ✓ invalid phone format
    ✓ invalid DNI format
    ... (6/6 passing)

Total: 28/28 tests passing ✅
```

---

## 🗄️ Database

### Migrations

Migrations run automatically on application startup. To manage manually:

```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path migrations -database "postgres://user:pass@localhost:5432/arnela?sslmode=disable" up

# Rollback last migration
migrate -path migrations -database "postgres://user:pass@localhost:5432/arnela?sslmode=disable" down 1

# Create new migration
migrate create -ext sql -dir migrations -seq <migration_name>
```

### Current Schema

**Users Table:**
- Authentication and user management
- Roles: admin, employee, client
- Password hashing with bcrypt

**Clients Table:**
- Client information with Spanish-specific fields
- DNI/NIE validation
- Soft delete support
- Comprehensive indexes for performance

---

## 🛠️ Development

### Build

```bash
go build -o arnela-api.exe cmd/api/main.go
```

### Run

```bash
# Development mode (with hot reload)
air

# Production mode
go run cmd/api/main.go
```

### Swagger Documentation

After modifying handler annotations, regenerate Swagger docs:

```bash
# Install swag (once)
go install github.com/swaggo/swag/cmd/swag@latest

# Regenerate docs
swag init -g cmd/api/main.go -o docs
```

Access Swagger UI at: http://localhost:8080/swagger/index.html

### Docker

```bash
# Start services
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f

# Rebuild
docker-compose up -d --build
```

---

## 📝 Spanish-Specific Validations

### DNI/NIE Validation

**Valid DNI:** `12345678Z` (8 digits + check letter)  
**Valid NIE:** `X1234567L` (X/Y/Z + 7 digits + check letter)

The system validates the check digit using the modulo 23 algorithm.

### Phone Normalization

Automatically normalizes Spanish phone numbers to `+34XXXXXXXXX` format:

- `612345678` → `+34612345678`
- `+34612345678` → `+34612345678`
- `34612345678` → `+34612345678`

---

## 🔍 URLs Disponibles

### Local Development

- **Backend API:** http://localhost:8080
- **Swagger UI:** http://localhost:8080/swagger/index.html
- **Swagger JSON:** http://localhost:8080/swagger/doc.json
- **PostgreSQL:** localhost:5432
- **Redis:** localhost:6379

### API Base URL

```
http://localhost:8080/api/v1
```

---

## 📚 Documentation

- **Agent.md** - Technical project definition and architecture
- **PHASE_1.3_COMPLETE.md** - User management implementation details
- **PHASE_1.4_COMPLETE.md** - Client management implementation details
- **Swagger UI** - Interactive API documentation

---

## 🐛 Troubleshooting

### Database Connection Issues

```bash
# Check if PostgreSQL is running
docker-compose ps

# View PostgreSQL logs
docker-compose logs postgres

# Restart PostgreSQL
docker-compose restart postgres
```

### Migration Errors

```bash
# Force migration version
migrate -path migrations -database "postgres://..." force <version>

# Check current version
migrate -path migrations -database "postgres://..." version
```

### JWT Token Expired

If you get "token expired" errors, login again to get a new token. Tokens are valid for 24 hours.

---

## 📈 Performance

### Database Indexes

All performance-critical fields are indexed:
- Email (unique)
- DNI (unique, excluding soft-deleted)
- Phone
- City, Province
- Last name
- Active status

### Pagination

List endpoints support pagination to handle large datasets efficiently:
- Default page size: 20
- Maximum page size: 100

---

## 🚦 Status

**Current Phase:** 1.4 - Client Management ✅  
**Tests:** 28/28 passing ✅  
**API Endpoints:** 15 endpoints implemented  
**Next Phase:** Employee Management or Appointments System

---

## 📄 License

Proprietary - Arnela Professional Office

---

## 🤝 Contributing

This is a private project. For internal team members:

1. Create a feature branch
2. Write tests for new features
3. Ensure all tests pass
4. Update documentation
5. Submit pull request

---

**Built with ❤️ using Go, GIN, and PostgreSQL**
