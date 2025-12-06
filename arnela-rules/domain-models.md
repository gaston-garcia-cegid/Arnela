# 📦 Domain Models

## 5.1. User Domain

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

## 5.2. Client Domain

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
    City        string     `json:"city,omitempty" db:"city"`
    IsActive    bool       `json:"isActive" db:"is_active"`
    // ... complete struct in codebase
}
```
