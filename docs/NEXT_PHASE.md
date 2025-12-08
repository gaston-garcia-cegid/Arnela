# 🚀 Fase 1.2 - Siguiente Etapa

## ✅ Estado Actual (Fase 1.1 Completada)

- ✓ Infraestructura Docker (PostgreSQL, Redis)
- ✓ Backend Go con GIN (estructura Clean Architecture)
- ✓ Frontend Next.js 16 con TypeScript y Zustand
- ✓ Configuración corregida (sin errores de symlinks)
- ✓ Servidor de desarrollo funcionando

---

## 🎯 Fase 1.2 - Tareas Pendientes

### 1. Base de Datos y Migraciones
**Objetivo:** Configurar el sistema de migraciones para PostgreSQL

**Tareas:**
- [ ] Instalar `golang-migrate` o `goose` para migraciones
- [ ] Crear script de inicialización de BD
- [ ] Definir primera migración (tabla `users`)
- [ ] Integrar migraciones en el startup del backend

**Archivos a crear:**
- `backend/migrations/000001_create_users_table.up.sql`
- `backend/migrations/000001_create_users_table.down.sql`
- `backend/pkg/database/migrate.go`

---

### 2. Modelos de Dominio
**Objetivo:** Crear las entidades principales del sistema

**Tareas:**
- [ ] Modelo `User` (autenticación y roles)
- [ ] Modelo `Client` (clientes del sistema)
- [ ] Modelo `Employee` (empleados/profesionales)
- [ ] Modelo `Appointment` (citas)

**Archivos a crear:**
- `backend/internal/domain/user.go`
- `backend/internal/domain/client.go`
- `backend/internal/domain/employee.go`
- `backend/internal/domain/appointment.go`

**Estructura ejemplo (`user.go`):**
```go
package domain

import "time"

type User struct {
    ID        string    `json:"id" db:"id"`
    Email     string    `json:"email" db:"email"`
    Password  string    `json:"-" db:"password"` // No exponer en JSON
    FirstName string    `json:"firstName" db:"first_name"`
    LastName  string    `json:"lastName" db:"last_name"`
    Role      string    `json:"role" db:"role"`
    CreatedAt time.Time `json:"createdAt" db:"created_at"`
    UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
```

---

### 3. Autenticación JWT
**Objetivo:** Implementar sistema de autenticación con JWT

**Tareas:**
- [ ] Crear servicio de autenticación (`AuthService`)
- [ ] Implementar generación de JWT tokens
- [ ] Crear middleware de autenticación
- [ ] Endpoints de login/register

**Archivos a crear:**
- `backend/internal/service/auth_service.go`
- `backend/internal/middleware/auth_middleware.go`
- `backend/internal/handler/auth_handler.go`
- `backend/pkg/jwt/token.go`

**Endpoints a implementar:**
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `GET /api/v1/auth/me` (requiere autenticación)

---

### 4. Repositorios (Data Access Layer)
**Objetivo:** Implementar la capa de acceso a datos

**Tareas:**
- [ ] Interfaz `UserRepository`
- [ ] Implementación con PostgreSQL
- [ ] Operaciones CRUD básicas
- [ ] Manejo de errores y transacciones

**Archivos a crear:**
- `backend/internal/repository/user_repository.go`
- `backend/internal/repository/postgres/user_repository_impl.go`

---

### 5. Testing (TDD)
**Objetivo:** Configurar framework de testing

**Tareas:**
- [ ] Setup de `testify` para assertions
- [ ] Tests unitarios para `AuthService`
- [ ] Tests de integración para repositorios
- [ ] Mocks para dependencias

**Archivos a crear:**
- `backend/internal/service/auth_service_test.go`
- `backend/internal/repository/postgres/user_repository_test.go`
- `backend/test/helpers.go` (utilidades de testing)

---

### 6. Documentación API (Swagger)
**Objetivo:** Generar documentación automática con Swagger

**Tareas:**
- [ ] Instalar `swaggo/swag`
- [ ] Anotar handlers con comentarios Swagger
- [ ] Generar documentación automática
- [ ] Endpoint `/swagger/index.html`

**Ejemplo de anotación:**
```go
// Login godoc
// @Summary      Login de usuario
// @Description  Autentica un usuario y devuelve un JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Credenciales"
// @Success      200 {object} LoginResponse
// @Failure      401 {object} ErrorResponse
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
    // ...
}
```

---

## 📝 Orden Sugerido de Implementación

1. **Migraciones de BD** → Base para todo lo demás
2. **Modelos de Dominio** → Definir estructuras
3. **Repositorios** → Acceso a datos
4. **Servicio de Autenticación** → Lógica de negocio
5. **Handlers y Middleware** → Exposición de APIs
6. **Tests** → Validación (idealmente en paralelo con cada paso)
7. **Swagger** → Documentación (al final)

---

## 🔧 Dependencias a Instalar

```powershell
cd backend

# Migraciones
go get -u github.com/golang-migrate/migrate/v4

# Testing
go get -u github.com/stretchr/testify

# JWT
go get -u github.com/golang-jwt/jwt/v5

# Swagger
go install github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files

# Password hashing
go get -u golang.org/x/crypto/bcrypt

# UUID
go get -u github.com/google/uuid
```

---

## 📊 Criterios de Éxito

Al completar la Fase 1.2, deberías tener:

- ✅ Base de datos con tablas `users`, `clients`, `employees`
- ✅ Sistema de autenticación JWT funcional
- ✅ Endpoints protegidos con middleware
- ✅ Tests unitarios con >70% coverage
- ✅ Documentación Swagger accesible
- ✅ Capacidad de crear/login usuarios desde el frontend

---

## 🎯 Siguiente Prompt para Copilot

```
Iniciar Fase 1.2: Implementar sistema de migraciones de BD y modelo User.

1. Instalar dependencias necesarias (golang-migrate, uuid, bcrypt)
2. Crear estructura de migraciones en backend/migrations/
3. Crear primera migración: tabla users con campos (id, email, password_hash, first_name, last_name, role, created_at, updated_at)
4. Implementar dominio User en backend/internal/domain/user.go
5. Crear sistema de ejecución de migraciones en startup

Usar convenciones del Agent.md (camelCase en JSON, PascalCase en structs Go).
```

---

**Ready para continuar con la Fase 1.2!** 🚀
