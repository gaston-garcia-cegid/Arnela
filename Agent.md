# 🤖 Agent.md - Technical Project Definition# 🤖 Agent.md - Technical Project Definition



## 1. 📝 Project Overview (CRM/CMS for Professional Office)## 1. 📝 Project Overview (CRM/CMS for Professional Office)



The project involves the development of a custom enterprise web application for a professional service office/cabinet. The system will replace manual processes (Excel files, external appointment management) and provide a unified platform for managing clients, employees, appointments, and tasks.The project involves the development of a custom enterprise web application for a professional service office/cabinet. The system will replace manual processes (Excel files, external appointment management) and provide a unified platform for managing clients, employees, appointments, and tasks.



The system consists of three main interfaces:The system consists of three main interfaces:



1.  **Landing Page:** A replica of the current website with a login modal added.1.  **Landing Page:** A replica of the current website with a login modal added.

2.  **Client Area:** Allows clients to request, consult, and manage their own appointments.2.  **Client Area:** Allows clients to request, consult, and manage their own appointments.

3.  **Backoffice (CRM + CMS):** An internal panel for the administration of clients, employees, tasks, session reports, subsidies, and social communication/marketing management.3.  **Backoffice (CRM + CMS):** An internal panel for the administration of clients, employees, tasks, session reports, subsidies, and social communication/marketing management.



The system will be implemented using a **Modular Monolith** architecture for the Backend (Golang) and a **Monolithic/BFF** application for the Frontend (Next.js), optimizing development for a single, part-time developer.The system will be implemented using a **Modular Monolith** architecture for the Backend (Golang) and a **Monolithic/BFF** application for the Frontend (Next.js), optimizing development for a single, part-time developer.



------



## 2. ⚙️ Technology Stack## 2. ⚙️ Technology Stack



| Component | Technology | Standards || Component | Technology | Standards |

| :--- | :--- | :--- || :--- | :--- | :--- |

| **Backend** | **Go 1.23** + **GIN** Framework | Modular Monolith, GIN Middleware, **sqlx** for DB access || **Backend** | **Go** (Latest Version) + **GIN** Framework | Modular Monolith, GIN Middleware. |

| **Frontend** | **Next.js 16** (TypeScript) | App Router, React Components, **Zustand** for state management || **Frontend** | **Next.js 16** (TypeScript) | App Router, React Components, **Zustand** for state management. |

| **Styling** | React CSS, **Shadcn UI** | Accessible components, utility/classes || **Styling** | React CSS, **Shadcn UI** | Accessible components, utility/classes. |

| **Database** | **PostgreSQL 16** | Transactional source of truth, **golang-migrate** for migrations || **Database** | **PostgreSQL** (Latest Version) | Transactional source of truth. |

| **Cache/Broker** | **Redis** | Session caching, read caching, **Asynchronous Task Queue** || **Cache/Broker** | **Redis** | Session caching, read caching, **Asynchronous Task Queue**. |

| **API Documentation** | **Swagger/OpenAPI 3.0** | Automated generation via **swaggo** (Go comments) || **API Documentation** | **Swagger/OpenAPI 3.0** | Automated generation (via Go comments). |

| **Development** | **Docker** | Replicable local environment (Go, PG, Redis) || **Development** | **Docker** | Replicable local environment (Go, PG, Redis). |

| **Testing** | **testify/mock, testify/assert** | TDD focus on backend business logic || **Methodology** | **TDD** (Test-Driven Development) | Primary focus on Backend development and business logic. |

| **Methodology** | **TDD** (Test-Driven Development) | Primary focus on Backend development and business logic |

---

---

## 3. 🏛️ Architecture Patterns

## 3. 🏛️ Architecture Patterns

### 3.1. Backend Structure (Clean Architecture / Modular Monolith)

### 3.1. Backend Structure (Clean Architecture / Modular Monolith)

The Golang Backend (GIN) will follow the principles of Clean Architecture to separate concerns from infrastructure to business logic.

The Golang Backend (GIN) follows the principles of Clean Architecture to separate concerns from infrastructure to business logic.

| Layer | Responsibility | Go Implementation |

| Layer | Responsibility | Go Implementation || :--- | :--- | :--- |

| :--- | :--- | :--- || **Transport / Adapters** | Entry/exit points: APIs, middlewares, GIN configuration. | `handler/`, `cmd/`, `middleware/` |

| **Transport / Adapters** | Entry/exit points: APIs, middlewares, GIN configuration | `handler/`, `cmd/`, `middleware/` || **Services / Business Logic**| Specific business rules (e.g., Schedule Appointment, Assign Task). | `service/` (Interfaces and Core Logic) |

| **Services / Business Logic**| Specific business rules (e.g., Schedule Appointment, Assign Task) | `service/` (Interfaces and Core Logic) || **Domain / Entities** | Core data models (Client, Employee, Appointment, Task). | `domain/` (Core Structs) |

| **Domain / Entities** | Core data models (User, Client, Employee, Appointment, Task) | `domain/` (Core Structs) || **Persistence / Repository** | Abstraction of database access (PostgreSQL, Redis). | `repository/` (DB Access Interfaces) |

| **Persistence / Repository** | Abstraction of database access (PostgreSQL, Redis) | `repository/` (DB Access Interfaces) || **External Integrations** | Communication with GCal, Notification APIs (WhatsApp/SMS). | `integration/` |

| **External Integrations** | Communication with GCal, Notification APIs (WhatsApp/SMS) | `integration/` |

### 3.2. Frontend Structure (Next.js + TypeScript + Zustand)

### 3.2. Frontend Structure (Next.js + TypeScript + Zustand)

The Frontend will use Next.js **App Router**. Global state management is delegated to **Zustand** for simplicity and performance, especially for handling user session, notifications, and UI states.

The Frontend will use Next.js **App Router**. Global state management is delegated to **Zustand** for simplicity and performance, especially for handling user session, notifications, and UI states.

| Folder | Responsibility | Details |

| Folder | Responsibility | Details || :--- | :--- | :--- |

| :--- | :--- | :--- || `app/` | Routes and Layout | Layouts, pages, route templates. |

| `app/` | Routes and Layout | Layouts, pages, route templates || `components/` | Reusable Components | `ui/` (Shadcn), `common/` (custom), `backoffice/` |

| `components/` | Reusable Components | `ui/` (Shadcn), `common/` (custom), `backoffice/` || `stores/` | Global State Management | **Zustand** files (e.g., `useAuthStore`, `useTaskStore`). |

| `stores/` | Global State Management | **Zustand** files (e.g., `useAuthStore`, `useTaskStore`) || `hooks/` | Interface Logic | `useDebounce`, custom hooks for accessing stores. |

| `hooks/` | Interface Logic | `useDebounce`, custom hooks for accessing stores || `lib/` | Utilities and Configuration | Formatting functions, constants, API client. |

| `lib/` | Utilities and Configuration | Formatting functions, constants, API client |

---

---

## 4. 🔠 Naming Conventions

## 4. 🔠 Naming Conventions

### 4.1. Go (Backend Services)

### 4.1. Go (Backend Services)

| Convention | Context | Correct Example |

| Convention | Context | Correct Example || :--- | :--- | :--- |

| :--- | :--- | :--- || **PascalCase** | Exported (public) Structs, Exported Functions. | `type UserService struct {}`, `func GetUserByID() {}` |

| **PascalCase** | Exported (public) Structs, Exported Functions | `type UserService struct {}`, `func GetUserByID() {}` || **camelCase** | Unexported (private) Structs, Variables, Parameters. | `type carService struct {}`, `var userName string` |

| **camelCase** | Unexported (private) Structs, Variables, Parameters | `type carService struct {}`, `var userName string` || **CONST_CASE** | Public constants (acronyms in uppercase). | `const MaxRetries = 3` |

| **CONST_CASE** | Public constants (acronyms in uppercase) | `const MaxRetries = 3` |

### 4.2. TypeScript (Frontend with Zustand)

### 4.2. TypeScript (Frontend with Zustand)

| Convention | Context | Correct Example |

| Convention | Context | Correct Example || :--- | :--- | :--- |

| :--- | :--- | :--- || **PascalCase** | Interfaces, Types, React Components. | `interface UserProps`, `const UserList` |

| **PascalCase** | Interfaces, Types, React Components | `interface UserProps`, `const UserList` || **camelCase** | Interface properties, variables, functions. | `firstName: string`, `const fetchUsers = () => {}` |

| **camelCase** | Interface properties, variables, functions | `firstName: string`, `const fetchUsers = () => {}` |

### 4.3. JSON API Conventions

### 4.3. JSON API Conventions

All keys used in API request and response bodies must use **camelCase**.

All keys used in API request and response bodies must use **camelCase**.

* **Go Structs:** The `json:"camelCase"` tag must be used to ensure Go serialization respects the TypeScript/JavaScript convention.

* **Go Structs:** The `json:"camelCase"` tag must be used to ensure Go serialization respects the TypeScript/JavaScript convention.

```go

```go// ✅ Correct - camelCase in JSON tag

// ✅ Correct - camelCase in JSON tagtype CreateUserRequest struct {

type CreateUserRequest struct {    FirstName string `json:"firstName"` 

    FirstName string `json:"firstName"`     LastName  string `json:"lastName"` 

    LastName  string `json:"lastName"`     Email     string `json:"email"`

    Email     string `json:"email"`}
}
```

---

## 5. 📦 Domain Models (Implemented)

### 5.1. User Domain (`internal/domain/user.go`)

Core user entity with authentication and role-based access control.

```go
type User struct {
    ID        uuid.UUID `json:"id" db:"id"`
    Email     string    `json:"email" db:"email"`
    Password  string    `json:"-" db:"password"`
    FirstName string    `json:"firstName" db:"first_name"`
    LastName  string    `json:"lastName" db:"last_name"`
    Role      UserRole  `json:"role" db:"role"`
    IsActive  bool      `json:"isActive" db:"is_active"`
    CreatedAt time.Time `json:"createdAt" db:"created_at"`
    UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type UserRole string

const (
    RoleAdmin    UserRole = "admin"
    RoleEmployee UserRole = "employee"
    RoleClient   UserRole = "client"
)
```

**Key Features:**
- Password hashing with bcrypt (cost 10)
- Role-based access control (admin, employee, client)
- Helper methods: `FullName()`, `CanAccessBackoffice()`, `HasRole(role)`

### 5.2. Client Domain (`internal/domain/client.go`)

Client management with Spanish-specific fields (DNI, phone).

```go
type Client struct {
    ID          uuid.UUID  `json:"id" db:"id"`
    UserID      *uuid.UUID `json:"userId,omitempty" db:"user_id"`
    FirstName   string     `json:"firstName" db:"first_name"`
    LastName    string     `json:"lastName" db:"last_name"`
    Email       string     `json:"email" db:"email"`
    Phone       string     `json:"phone" db:"phone"`
    DNI         string     `json:"dni" db:"dni"`
    DateOfBirth *time.Time `json:"dateOfBirth,omitempty" db:"date_of_birth"`
    Address     string     `json:"address,omitempty" db:"address"`
    City        string     `json:"city,omitempty" db:"city"`
    PostalCode  string     `json:"postalCode,omitempty" db:"postal_code"`
    Province    string     `json:"province,omitempty" db:"province"`
    IsActive    bool       `json:"isActive" db:"is_active"`
    LastVisit   *time.Time `json:"lastVisit,omitempty" db:"last_visit"`
    Notes       string     `json:"notes,omitempty" db:"notes"`
    CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
    UpdatedAt   time.Time  `json:"updatedAt" db:"updated_at"`
    DeletedAt   *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
}
```

**Key Features:**
- Soft delete support (deleted_at timestamp)
- Spanish ID validation (DNI/NIE/NIF with check digit algorithm)
- Phone normalization (+34 format)
- Helper methods: `FullName()`, `IsDeleted()`, `Age()`

---

## 6. 🔐 Authentication & Authorization

### 6.1. JWT Authentication

- **Token expiration:** 24 hours
- **Secret:** Environment variable `JWT_SECRET`
- **Claims:** UserID (UUID), Email, Role, Expiration
- **Algorithm:** HS256

### 6.2. Role-Based Access Control

| Role | Permissions | Access |
| :--- | :--- | :--- |
| **admin** | Full system access | All endpoints + user/client/employee management + delete operations |
| **employee** | Operational access | Client management, appointments, tasks (no delete) |
| **client** | Self-service access | Own profile (`/clients/me`), own appointments |

### 6.3. Middleware Implementation

**Files:**
- `internal/middleware/auth_middleware.go`

**Functions:**
- `AuthMiddleware()` - JWT validation, user context injection (userID as string)
- `RequireRole(...roles)` - Role-based endpoint protection

**Example:**
```go
// Admin only
router.DELETE("/clients/:id", middleware.AuthMiddleware(), middleware.RequireRole("admin"), handler.DeleteClient)

// Admin or Employee
router.POST("/clients", middleware.AuthMiddleware(), middleware.RequireRole("admin", "employee"), handler.CreateClient)

// Client only
router.GET("/clients/me", middleware.AuthMiddleware(), middleware.RequireRole("client"), handler.GetMyClient)
```

---

## 7. 🗄️ Database Schema (PostgreSQL)

### 7.1. Users Table (Migration 000001)

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'employee', 'client')),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_is_active ON users(is_active);
```

### 7.2. Clients Table (Migration 000002)

```sql
CREATE TABLE clients (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    dni VARCHAR(20) NOT NULL,
    date_of_birth DATE,
    address TEXT,
    city VARCHAR(100),
    postal_code VARCHAR(10),
    province VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    last_visit TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Indexes for performance and uniqueness
CREATE UNIQUE INDEX idx_clients_email_unique ON clients(email) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_clients_dni_unique ON clients(dni) WHERE deleted_at IS NULL;
CREATE INDEX idx_clients_phone ON clients(phone) WHERE deleted_at IS NULL;
CREATE INDEX idx_clients_city ON clients(city) WHERE deleted_at IS NULL;
CREATE INDEX idx_clients_province ON clients(province) WHERE deleted_at IS NULL;
CREATE INDEX idx_clients_user_id ON clients(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_clients_last_name ON clients(last_name) WHERE deleted_at IS NULL;
CREATE INDEX idx_clients_is_active ON clients(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_clients_deleted_at ON clients(deleted_at);

-- Auto-update trigger
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_clients_updated_at BEFORE UPDATE ON clients
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

---

## 8. 🌐 API Endpoints (Implemented)

### 8.1. Authentication Endpoints

| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| POST | `/api/v1/auth/register` | Public | Register new user (returns JWT token) |
| POST | `/api/v1/auth/login` | Public | Login and get JWT token |
| GET | `/api/v1/auth/me` | Authenticated | Get current user profile |

**Example Register/Login:**
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

### 8.2. User Management Endpoints

| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| POST | `/api/v1/users` | Admin | Create new user |
| GET | `/api/v1/users` | Admin | List all users |
| GET | `/api/v1/users/:id` | Admin | Get user by ID |
| PUT | `/api/v1/users/:id` | Admin | Update user |
| DELETE | `/api/v1/users/:id` | Admin | Delete user |

### 8.3. Client Management Endpoints

| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| POST | `/api/v1/clients` | Admin, Employee | Create new client |
| GET | `/api/v1/clients` | Admin, Employee | List clients (with filters & pagination) |
| GET | `/api/v1/clients/:id` | Admin, Employee | Get client by ID |
| PUT | `/api/v1/clients/:id` | Admin, Employee | Update client |
| DELETE | `/api/v1/clients/:id` | Admin | Soft delete client |
| GET | `/api/v1/clients/me` | Client | Get own client profile |

**List Clients Query Parameters:**
- `page` (default: 1)
- `pageSize` (default: 20, max: 100)
- `search` - Search across firstName, lastName, email, dni, phone
- `isActive` - Filter by active status (true/false)
- `city` - Filter by city
- `province` - Filter by province

**Example:**
```bash
# Create client
curl -X POST http://localhost:8080/api/v1/clients \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Juan",
    "lastName": "Pérez",
    "email": "juan@example.com",
    "phone": "612345678",
    "dni": "12345678Z",
    "city": "Madrid"
  }'

# List with filters
curl "http://localhost:8080/api/v1/clients?search=Juan&page=1&pageSize=10" \
  -H "Authorization: Bearer $TOKEN"
```

---

## 9. ✅ Validation Rules

### 9.1. Spanish DNI/NIE Validation

**DNI Format:** `12345678Z` (8 digits + check letter)  
**NIE Format:** `X1234567L` (X/Y/Z + 7 digits + check letter)

**Validation Algorithm:**
```go
func isValidSpanishDNI(dni string) bool {
    dni = strings.ToUpper(strings.TrimSpace(dni))
    
    // DNI: 8 digits + 1 letter
    if len(dni) == 9 && unicode.IsDigit(rune(dni[0])) {
        return validateDNILetter(dni)
    }
    
    // NIE: X/Y/Z + 7 digits + 1 letter
    if len(dni) == 9 && (dni[0] == 'X' || dni[0] == 'Y' || dni[0] == 'Z') {
        return validateNIELetter(dni)
    }
    
    return false
}

func validateDNILetter(dni string) bool {
    letters := "TRWAGMYFPDXBNJZSQVHLCKE"
    numStr := dni[:8]
    expectedLetter := dni[8]
    
    num, err := strconv.Atoi(numStr)
    if err != nil {
        return false
    }
    
    return letters[num%23] == expectedLetter
}
```

**Valid Examples:**
- DNI: `12345678Z`, `87654321X`
- NIE: `X1234567L`, `Y7654321K`, `Z9876543A`

### 9.2. Spanish Phone Validation

**Accepted formats:**
- `612345678` → Normalized to `+34612345678`
- `+34612345678` → Already normalized
- `34612345678` → Normalized to `+34612345678`

**Rules:**
- Must start with 6, 7, 8, or 9 (after removing prefix)
- Must be 9 digits after prefix
- Automatically adds `+34` if missing

```go
func normalizePhone(phone string) string {
    phone = strings.ReplaceAll(phone, " ", "")
    phone = strings.ReplaceAll(phone, "-", "")
    
    if strings.HasPrefix(phone, "+34") {
        return phone
    }
    if strings.HasPrefix(phone, "34") {
        return "+" + phone
    }
    return "+34" + phone
}
```

### 9.3. Email Validation

Standard email format: `name@domain.com`

```go
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
```

---

## 10. 🧪 Testing Strategy

### 10.1. Unit Tests

- **Framework:** `testify/assert` and `testify/mock`
- **Coverage:** Service layer (business logic)
- **Pattern:** Table-driven tests

**Example Test Structure:**
```go
func TestClientService_CreateClient(t *testing.T) {
    tests := []struct {
        name          string
        request       *CreateClientRequest
        mockSetup     func(*mocks.ClientRepository)
        expectedError string
    }{
        {
            name: "successful creation",
            request: &CreateClientRequest{
                FirstName: "Juan",
                Email: "juan@example.com",
                DNI: "12345678Z",
                Phone: "612345678",
            },
            mockSetup: func(repo *mocks.ClientRepository) {
                repo.On("EmailExists", mock.Anything, "juan@example.com", uuid.Nil).Return(false, nil)
                repo.On("DNIExists", mock.Anything, "12345678Z", uuid.Nil).Return(false, nil)
                repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Client")).Return(nil)
            },
            expectedError: "",
        },
        // ... more test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := new(mocks.ClientRepository)
            tt.mockSetup(mockRepo)
            
            service := NewClientService(mockRepo)
            _, err := service.CreateClient(context.Background(), tt.request)
            
            if tt.expectedError != "" {
                assert.ErrorContains(t, err, tt.expectedError)
            } else {
                assert.NoError(t, err)
            }
            
            mockRepo.AssertExpectations(t)
        })
    }
}
```

### 10.2. Test Results

**Phase 1.3 - User Service Tests:**
- `TestAuthService_Register`: 5/5 passing
- `TestAuthService_Login`: 3/3 passing
- `TestUserService_CreateUser`: 4/4 passing
- `TestUserService_GetUser`: 2/2 passing
- `TestUserService_UpdateUser`: 2/2 passing
- `TestUserService_DeleteUser`: 2/2 passing
- **Total: 18/18 tests passing ✅**

**Phase 1.4 - Client Service Tests:**
- `TestClientService_CreateClient`: 6/6 passing
- `TestClientService_GetClient`: 2/2 passing
- `TestClientService_DeleteClient`: 2/2 passing
- **Total: 10/10 tests passing ✅**

**Overall: 28/28 tests passing ✅**

---

## 11. 📚 API Documentation (Swagger)

### 11.1. Generation

```bash
# Install swag (once)
go install github.com/swaggo/swag/cmd/swag@latest

# Regenerate docs
swag init -g cmd/api/main.go -o docs
```

### 11.2. Access

- **Swagger UI:** `http://localhost:8080/swagger/index.html`
- **OpenAPI JSON:** `http://localhost:8080/swagger/doc.json`

### 11.3. Swagger Annotations

All handlers include comprehensive Swagger annotations:

```go
// CreateClient godoc
// @Summary Create a new client
// @Description Create a new client (admin or employee only)
// @Tags clients
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body service.CreateClientRequest true "Client creation request"
// @Success 201 {object} domain.Client
// @Failure 400 {object} handler.ErrorResponse "Invalid request or validation error"
// @Failure 409 {object} handler.ErrorResponse "Email or DNI already exists"
// @Router /api/v1/clients [post]
func (h *ClientHandler) CreateClient(c *gin.Context) { ... }
```

---

## 12. 🏗️ Project Structure (Backend)

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
│   │   ├── auth_service_test.go    # Auth tests (8 tests)
│   │   ├── user_service_test.go    # User tests (10 tests)
│   │   └── client_service_test.go  # Client tests (10 tests)
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
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_create_clients_table.up.sql
│   └── 000002_create_clients_table.down.sql
├── docs/                           # Swagger documentation
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── .env.example
├── README.md
├── PHASE_1.3_COMPLETE.md           # Phase 1.3 documentation
└── PHASE_1.4_COMPLETE.md           # Phase 1.4 documentation
```

---

## 13. 🚀 Development Workflow

### 13.1. Local Setup

```bash
# 1. Clone repository
git clone <repo-url>
cd arnela/backend

# 2. Install dependencies
go mod download

# 3. Setup environment
cp .env.example .env
# Edit .env with your configuration:
# - DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
# - JWT_SECRET
# - REDIS_HOST, REDIS_PORT

# 4. Start Docker services (PostgreSQL, Redis)
docker-compose up -d

# 5. Run migrations (automatic on startup)
go run cmd/api/main.go

# 6. Access API
# Backend: http://localhost:8080
# Swagger: http://localhost:8080/swagger/index.html
```

### 13.2. Testing

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/service/... -v

# Run with coverage
go test -cover ./...

# Run specific test
go test ./internal/service/... -v -run TestClientService_CreateClient
```

### 13.3. Swagger Regeneration

```bash
# After adding/modifying handler annotations
swag init -g cmd/api/main.go -o docs
```

### 13.4. Database Migrations

```bash
# Create new migration
migrate create -ext sql -dir migrations -seq <migration_name>

# Run migrations (automatic on app start, or manually)
migrate -path migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" up

# Rollback last migration
migrate -path migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" down 1
```

---

## 14. 📝 Implementation Status

### ✅ Phase 1.1 - Initial Setup
- [x] Project structure (Clean Architecture)
- [x] Docker configuration (PostgreSQL, Redis)
- [x] Database connection with sqlx
- [x] Basic GIN middleware setup
- [x] Environment configuration

### ✅ Phase 1.2 - Authentication
- [x] User registration with JWT
- [x] User login with JWT
- [x] JWT generation and validation
- [x] Role-based access control (admin, employee, client)
- [x] Auth middleware (JWT validation)

### ✅ Phase 1.3 - User Management
- [x] User domain model
- [x] User CRUD operations
- [x] User repository (PostgreSQL with sqlx)
- [x] User service with tests (18/18 passing)
- [x] Swagger documentation
- [x] Logging system

### ✅ Phase 1.4 - Client Management
- [x] Client domain model (Spanish fields)
- [x] Client CRUD operations
- [x] Client repository (PostgreSQL with sqlx)
- [x] Client service with Spanish validations (DNI/NIE, phone)
- [x] Client service tests (10/10 passing)
- [x] Role-based endpoints (admin/employee/client)
- [x] Soft delete implementation
- [x] List with filters and pagination
- [x] Swagger documentation

### 🔄 Phase 1.5 - Employee Management (Pending)
- [ ] Employee domain model
- [ ] Employee CRUD operations
- [ ] Employee repository
- [ ] Employee service
- [ ] Employee tests
- [ ] Swagger documentation

### 🔄 Phase 1.6 - Appointments System (Pending)
- [ ] Appointment domain model
- [ ] Appointment CRUD operations
- [ ] Appointment scheduling logic
- [ ] Google Calendar integration
- [ ] WhatsApp/SMS notifications

### 🔄 Phase 1.7 - Tasks Management (Pending)
- [ ] Task domain model
- [ ] Task assignment logic
- [ ] Task tracking
- [ ] Task notifications

---

## 15. 🔗 Integration Points (Future)

### 15.1. External Services

| Service | Purpose | Implementation |
| :--- | :--- | :--- |
| **Google Calendar API** | Appointment synchronization | `integration/gcal/` |
| **WhatsApp/SMS API** | Client notifications | `integration/notifications/` |
| **Email Service** | Appointment reminders, reports | `integration/email/` |

### 15.2. Redis Usage

| Use Case | Implementation |
| :--- | :--- |
| **Session Cache** | JWT token cache for fast validation |
| **Async Tasks** | Queue for email/SMS notifications |
| **Read Cache** | Cache frequently accessed client/appointment data |

---

## 16. 🐛 Known Issues & Limitations

### Current Limitations
- No integration tests (only unit tests)
- No Redis cache implementation (infrastructure ready)
- No external integrations (Google Calendar, WhatsApp)
- No file upload functionality
- No audit log system

### Future Enhancements
- Integration tests with test database
- Redis session caching
- Appointment reminders via email/SMS
- Bulk client import/export
- Advanced search with Elasticsearch
- Real-time notifications with WebSockets

---

## 17. 📖 References

- **Go Documentation:** https://go.dev/doc/
- **GIN Framework:** https://gin-gonic.com/docs/
- **sqlx:** https://jmoiron.github.io/sqlx/
- **Next.js 16:** https://nextjs.org/docs
- **Zustand:** https://zustand-demo.pmnd.rs/
- **Shadcn UI:** https://ui.shadcn.com/
- **Swagger/OpenAPI:** https://swagger.io/specification/
- **Clean Architecture:** https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
- **testify:** https://github.com/stretchr/testify

---

**Last Updated:** November 15, 2024  
**Phase:** 1.4 Complete (Client Management)  
**Next Phase:** 1.5 (Employee Management) or 1.6 (Appointments System)  
**Backend Tests:** 28/28 passing ✅  
**API Endpoints:** 15 endpoints implemented  
**Documentation:** Swagger UI available at `/swagger/index.html`
