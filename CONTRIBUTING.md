# Guía de Contribución

Guía para el desarrollo de Arnela.

## Setup del entorno

### Requisitos

- Go 1.25+
- Node.js 22+ con pnpm
- Docker y Docker Compose
- Git

### Primer setup

```powershell
# Clonar
git clone <repo-url>
cd arnela

# Backend env
copy backend\.env.example backend\.env

# Frontend env
copy frontend\.env.example frontend\.env

# Levantar PostgreSQL + Redis
docker compose up -d

# Backend
cd backend
go mod download
go run cmd/api/main.go

# Frontend (otra terminal)
cd frontend
pnpm install
pnpm dev
```

## Workflow de desarrollo

### Branches

- `main`: branch principal, protegido por CI
- Feature branches: `feat/nombre-feature`
- Fix branches: `fix/descripcion-del-bug`

### Commits

Usar Conventional Commits:

```
feat: add invoice PDF export
fix: correct VAT calculation for zero amounts
chore: update CI Go version
docs: update backend README with new endpoints
test: add billing stats service tests
refactor: extract email templates to separate package
```

### Pull Requests

1. Crear branch desde `main`
2. Implementar cambios con tests
3. Verificar que CI pasa localmente
4. Crear PR con descripción clara
5. Esperar review y CI verde

## Estructura del código

### Backend - Clean Architecture

```
Flujo de una request:

  HTTP Request
      │
      ▼
  middleware/     (auth, logging, rate limit, CORS)
      │
      ▼
  handler/        (parse request, call service, format response)
      │
      ▼
  service/        (business logic, validation, orchestration)
      │
      ▼
  repository/     (data access, SQL queries)
      │
      ▼
  PostgreSQL / Redis
```

**Reglas**:
- Los handlers nunca acceden directamente a la base de datos
- Los services reciben interfaces de repository (dependency injection)
- El dominio no tiene dependencias externas
- Los errores se propagan con `pkg/errors` y se manejan en los handlers

### Frontend - Componentes y Estado

```
Flujo de datos:

  Component
      │
      ├── useHook()         → fetch data via lib/api.ts
      ├── useAuthStore()    → Zustand global state
      └── react-hook-form   → form state + zod validation
```

**Reglas**:
- Componentes en `PascalCase`, un componente por archivo
- Estado global solo via Zustand stores
- Formularios siempre con react-hook-form + zod
- API calls centralizadas en `lib/api.ts`
- Errores manejados con `useErrorHandler` + sonner toasts

## Deploy remoto asistido (Cursor)

Para un flujo guiado por el agente (preguntas → SSH → `git pull` → Docker Compose en el servidor), usar la skill **`.cursor/skills/remote-deploy/SKILL.md`** (triggers: *deploy remoto*, *ssh deploy*, *remote deploy*, …). La guía operativa sigue siendo **`docs/DEPLOYMENT.md`**.

## Testing

### Backend

```powershell
cd backend

# Todos los tests
go test ./... -count=1

# Con race detector (recomendado)
go test ./... -count=1 -race

# CI ejecuta -race en Ubuntu (con CGO). En Windows, sin CGO habilitado,
# `go test -race` puede fallar o no ejecutarse; usa WSL/Docker o
# `set CGO_ENABLED=1` con un compilador C para reproducir lo mismo que GitHub Actions.

# Con cobertura
go test ./... -cover

# Un paquete específico
go test ./internal/service/... -v -run TestInvoiceService
```

**Patrón**: tests unitarios con mocks de repository (testify/mock). Redis tests con miniredis.

### Frontend

```powershell
cd frontend

# Watch mode
pnpm test

# Single run
pnpm test -- --run

# Un archivo
pnpm test -- --run src/lib/__tests__/api.test.ts

# Coverage
pnpm test:coverage
```

**Patrón**: Vitest + Testing Library + jsdom. Mock de fetch y stores.

### Verificar antes de push

```powershell
# Backend
cd backend
go vet ./...
go build ./...
go test ./... -count=1 -race

# Frontend
cd frontend
pnpm lint
pnpm exec tsc --noEmit
pnpm test -- --run
pnpm build
```

Estos son los mismos checks que ejecuta el CI.

## Agregar un nuevo endpoint

1. **Dominio**: Definir/actualizar entidad en `internal/domain/`
2. **Repository**: Agregar interfaz en `internal/repository/` e implementación en `repository/postgres/`
3. **Mock**: Agregar mock en `repository/mocks/`
4. **Service**: Implementar lógica en `internal/service/`
5. **Handler**: Crear handler en `internal/handler/`
6. **Routes**: Registrar en `cmd/api/main.go`
7. **Tests**: Escribir tests para service y handler
8. **Swagger**: Agregar anotaciones swag y regenerar docs

## Agregar una migración

```powershell
cd backend
migrate create -ext sql -dir migrations -seq nombre_descriptivo
```

Editar los archivos `.up.sql` y `.down.sql` generados. Las migraciones se ejecutan automáticamente al arrancar el backend.

## Agregar un componente frontend

1. Si es un primitivo UI (Button, Dialog, etc.): `components/ui/`
2. Si es compartido (Navbar, Footer, ThemeToggle): `components/common/`
3. Si es específico de feature: `components/<feature>/`
4. Si necesita datos: crear custom hook en `hooks/`
5. Si tiene estado global: agregar al store en `stores/`

## Variables de entorno

- Backend: ver `backend/.env.example` y `backend/README.md`
- Frontend: ver `frontend/.env.example` y `frontend/README.md`
- Producción: ver `.env.prod.example` y `docs/DEPLOYMENT.md`

Nunca commitear archivos `.env` con credenciales reales.
