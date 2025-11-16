# 🏥 Arnela - CRM/CMS para Oficina Profesional

Sistema empresarial personalizado para la gestión de clientes, empleados, citas y tareas. Desarrollado con arquitectura Modular Monolith (Backend) y Next.js (Frontend).

---

## 📦 Stack Tecnológico

| Componente | Tecnología |
|------------|------------|
| **Backend** | Go 1.23 + GIN Framework |
| **Frontend** | Next.js 16 + TypeScript + Zustand |
| **Database** | PostgreSQL 16 |
| **Cache** | Redis 7 |
| **Container** | Docker + Docker Compose |
| **Metodología** | TDD, Clean Architecture |

---

## 🚀 Inicio Rápido

### 1️⃣ Pre-requisitos

- **Docker** y **Docker Compose** instalados
- **Go 1.23+** (para desarrollo local del backend)
- **Node.js 22+** con **pnpm** (para desarrollo local del frontend)

### 2️⃣ Clonar y Configurar

```powershell
# Clonar el repositorio
git clone <tu-repo-url>
cd arnela

# Copiar variables de entorno (Backend)
cd backend
copy .env.example .env
cd ..

# Copiar variables de entorno (Frontend)
cd frontend
copy .env.example .env
cd ..
```

### 3️⃣ Iniciar con Docker

```powershell
# Iniciar todos los servicios (PostgreSQL, Redis, Go API)
docker-compose up -d

# Ver logs
docker-compose logs -f go-api

# Verificar que los servicios estén arriba
docker-compose ps
```

**URLs disponibles:**
- **Backend API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health
- **Swagger UI**: http://localhost:8080/swagger/index.html (📖 Documentación interactiva de la API)
- **API Auth**: 
  - POST `/api/v1/auth/register` - Registro
  - POST `/api/v1/auth/login` - Login
  - GET `/api/v1/auth/me` - Usuario actual (requiere token)
- **PostgreSQL**: localhost:5432 (usuario: `arnela_user`, password: `arnela_secure_pass_2024`)
- **Redis**: localhost:6379 (password: `arnela_redis_pass_2024`)

### 4️⃣ Desarrollo Local del Frontend

```powershell
cd frontend

# Instalar dependencias (solo primera vez)
pnpm install

# Iniciar servidor de desarrollo
pnpm dev
```

**Frontend disponible en**: http://localhost:3000

---

## 🏗️ Estructura del Proyecto

```
arnela/
├── backend/
│   ├── cmd/api/              # Punto de entrada de la aplicación
│   ├── internal/
│   │   ├── handler/          # HTTP handlers (Controllers)
│   │   ├── service/          # Lógica de negocio
│   │   ├── domain/           # Modelos de dominio
│   │   ├── repository/       # Acceso a datos (DB)
│   │   ├── middleware/       # Middlewares (Auth, CORS, etc.)
│   │   └── integration/      # Integraciones externas (Google Cal, SMS)
│   ├── pkg/
│   │   ├── database/         # Conexión PostgreSQL
│   │   └── cache/            # Conexión Redis
│   ├── config/               # Configuración centralizada
│   ├── go.mod                # Dependencias Go
│   └── Dockerfile            # Imagen Docker del backend
│
├── frontend/
│   ├── src/
│   │   ├── app/              # Next.js App Router (páginas/layouts)
│   │   ├── components/       # Componentes React
│   │   │   ├── ui/           # Componentes Shadcn UI
│   │   │   ├── common/       # Componentes compartidos
│   │   │   └── backoffice/   # Componentes del backoffice
│   │   ├── stores/           # Zustand stores (gestión de estado)
│   │   ├── hooks/            # Custom hooks
│   │   └── lib/              # Utilidades y cliente API
│   ├── package.json          # Dependencias frontend
│   └── tsconfig.json         # Configuración TypeScript
│
└── docker-compose.yml        # Orquestación de servicios
```

---

## 🛠️ Comandos Útiles

### Backend (Go)

```powershell
cd backend

# Ejecutar localmente (sin Docker)
go run cmd/api/main.go

# Ejecutar tests
go test ./...                    # Todos los tests
go test ./internal/... -v         # Tests con output detallado
go test ./internal/... -cover     # Con coverage

# Tests específicos
go test ./internal/service/... -v
go test ./internal/handler/... -v

# Actualizar dependencias
go mod tidy

# Compilar binario
go build -o main.exe cmd/api/main.go

# Generar documentación Swagger
swag init -g cmd/api/main.go -o docs
```

### Frontend (Next.js)

```powershell
cd frontend

# Desarrollo
pnpm dev

# Build para producción
pnpm build

# Iniciar en producción
pnpm start

# Linter
pnpm lint
```

### Docker

```powershell
# Iniciar servicios
docker-compose up -d

# Ver logs de un servicio específico
docker-compose logs -f go-api
docker-compose logs -f postgres

# Reiniciar servicios
docker-compose restart

# Detener servicios
docker-compose down

# Rebuild del backend
docker-compose up -d --build go-api

# Limpiar todo (incluye volúmenes)
docker-compose down -v
```

---

## 🧪 Testing

El proyecto implementa **TDD** (Test-Driven Development) con cobertura completa de tests unitarios.

### Ejecutar Tests

```powershell
cd backend

# Todos los tests
go test ./internal/... -v

# Ver resultados
# ✅ internal/handler - 10 tests passing
# ✅ internal/service - 8 tests passing
# Total: 18/18 tests passing
```

### Estructura de Tests

```
backend/internal/
├── handler/
│   └── auth_handler_test.go     # Tests de endpoints HTTP
├── service/
│   └── auth_service_test.go     # Tests de lógica de negocio
└── repository/mocks/
    └── user_repository_mock.go  # Mocks para testing
```

**Framework:** `stretchr/testify` para assertions y mocking

### Verificar Coverage

```powershell
go test ./internal/... -cover
```

---

## 📖 Documentación API (Swagger)

El proyecto usa **Swagger/OpenAPI 3.0** para documentación automática de la API.

### Acceder a Swagger UI

1. Iniciar el backend:
   ```powershell
   cd backend
   go run cmd/api/main.go
   ```

2. Abrir en el navegador:
   ```
   http://localhost:8080/swagger/index.html
   ```

### Regenerar Documentación

Después de cambiar endpoints o modelos:

```powershell
cd backend
swag init -g cmd/api/main.go -o docs
```

**Archivos generados:**
- `docs/docs.go` - Especificación en Go
- `docs/swagger.json` - Formato JSON
- `docs/swagger.yaml` - Formato YAML

---

## 📊 Logging Estructurado

El backend usa **zerolog** para logging de alto rendimiento con JSON estructurado.

### Desarrollo (Pretty Logs)

```powershell
$env:GO_ENV="development"
go run cmd/api/main.go
```

**Output:**
```
19:30:45 INF Starting Arnela API server port=8080
19:30:45 INF Database connected
19:30:45 INF Redis connected
```

### Producción (JSON Logs)

```powershell
go run cmd/api/main.go
```

**Output:**
```json
{"level":"info","time":"2024-11-15T19:30:45Z","message":"Starting Arnela API server","port":8080}
```

### Logs de Requests HTTP

Cada request HTTP es logueado automáticamente:

```json
{
  "level": "info",
  "method": "POST",
  "path": "/api/v1/auth/login",
  "status": 200,
  "duration": 45,
  "ip": "127.0.0.1",
  "message": "HTTP request completed"
}
```

---

## 🔑 Convenciones del Proyecto

### Backend (Go)

- **PascalCase**: Funciones/structs exportados (`GetUserByID`, `UserService`)
- **camelCase**: Variables/funciones privadas (`userName`, `validateInput`)
- **CONST_CASE**: Constantes públicas (`MaxRetries`, `DefaultTimeout`)
- **JSON Tags**: Siempre en `camelCase` para compatibilidad con frontend

```go
type CreateUserRequest struct {
    FirstName string `json:"firstName"`
    LastName  string `json:"lastName"`
    Email     string `json:"email"`
}
```

### Frontend (TypeScript)

- **PascalCase**: Componentes, interfaces, tipos (`UserList`, `UserProps`)
- **camelCase**: Props, variables, funciones (`firstName`, `handleClick`)
- **Zustand**: Para toda la gestión de estado global

```typescript
import { create } from 'zustand';

export const useAuthStore = create((set) => ({
  user: null,
  setUser: (user) => set({ user })
}));
```

---

## � Troubleshooting

### Error de Symlinks en Windows (Frontend)
Si ves errores relacionados con `EPERM: operation not permitted, symlink`:
- Ya está solucionado en la configuración (`next.config.js`)
- Ver `FRONTEND_FIX.md` para más detalles
- El modo `standalone` está deshabilitado para desarrollo local

### Problemas con pnpm
Si hay errores de autenticación con registros privados:
- El proyecto usa `.npmrc` local configurado con el registro público
- Eliminar archivos `.npmrc` globales si causan conflictos

---

## �📚 Referencias

- **Documentación detallada**: Ver `Agent.md` en la raíz del proyecto
- **Copilot Instructions**: `.github/copilot-instructions.md`
- **Fix Frontend**: Ver `FRONTEND_FIX.md` para detalles de correcciones

---

## 📞 Soporte

Para preguntas o issues, consulta la documentación técnica completa en `Agent.md`.

---

**✨ Happy Coding!**
