# Arnela - CRM/CMS para Gabinete Profesional

Sistema de gestión para Arnela Gabinete (Vigo). Gestión de clientes, empleados, citas, facturación, gastos y tareas internas.

## Stack

| Capa | Tecnología |
|------|------------|
| Backend | Go 1.25, Gin, Clean Architecture |
| Frontend | Next.js 16, TypeScript, Zustand, Tailwind v4, Shadcn/Radix |
| Base de datos | PostgreSQL 16, 15 migraciones SQL (golang-migrate) |
| Cache / Queue | Redis 7 (cache-aside + worker pool async) |
| Infraestructura | Docker Compose, Nginx reverse proxy, GitHub Actions CI |
| Auth | JWT con roles (admin, employee, client), rate limiting |
| Docs | Swagger/OpenAPI auto-generado |

## Inicio rápido

### Pre-requisitos

- Docker y Docker Compose
- Go 1.25+ (backend local)
- Node.js 22+ con pnpm (frontend local)

### Levantar servicios

```powershell
# Copiar env
copy backend\.env.example backend\.env

# Levantar PostgreSQL + Redis
docker-compose up -d

# Backend (en otra terminal)
cd backend
go run cmd/api/main.go

# Frontend (en otra terminal)
cd frontend
pnpm install
pnpm dev
```

### URLs

| Servicio | URL |
|----------|-----|
| Frontend | http://localhost:3000 |
| API | http://localhost:8080/api/v1 |
| Swagger | http://localhost:8080/swagger/index.html |
| Health | http://localhost:8080/health |
| Readiness | http://localhost:8080/readiness |

## Estructura del proyecto

```
arnela/
├── backend/
│   ├── cmd/api/                 # Entrypoint
│   ├── config/                  # Configuración desde env vars
│   ├── internal/
│   │   ├── domain/              # Entidades de dominio
│   │   ├── repository/          # Interfaces + implementación Postgres
│   │   ├── service/             # Lógica de negocio
│   │   ├── handler/             # HTTP handlers (Gin)
│   │   └── middleware/          # Auth, logging, rate limiting
│   ├── pkg/
│   │   ├── cache/               # Redis cache service
│   │   ├── database/            # PostgreSQL connection + health
│   │   ├── email/               # SMTP mailer + templates + notifications
│   │   ├── errors/              # Structured error responses
│   │   ├── gcal/                # Google Calendar integration
│   │   ├── jwt/                 # Token manager
│   │   ├── logger/              # Zerolog structured logging
│   │   ├── pdf/                 # Invoice PDF generation
│   │   └── queue/               # Redis task queue + worker pool
│   ├── migrations/              # 15 SQL migrations
│   └── docs/                    # Swagger auto-generated
│
├── frontend/
│   ├── src/
│   │   ├── app/                 # Next.js App Router (páginas + layouts)
│   │   ├── components/          # UI, common, landing, auth, backoffice, billing
│   │   ├── hooks/               # Custom hooks
│   │   ├── lib/                 # API client, errors, validators, utils
│   │   ├── stores/              # Zustand stores
│   │   └── types/               # TypeScript types
│   └── Dockerfile
│
├── nginx/                       # Nginx reverse proxy config
├── docker-compose.yml           # Dev (Postgres + Redis)
├── docker-compose.prod.yml      # Producción (todos los servicios)
└── .github/workflows/ci.yml     # CI pipeline
```

## API

### Endpoints principales

| Grupo | Endpoints | Acceso |
|-------|-----------|--------|
| Auth | register, login, me | Público (rate limited) |
| Clients | CRUD, me, search | Admin/Employee |
| Employees | CRUD, specialties | Admin (write), All (read) |
| Appointments | CRUD, confirm, cancel, available-slots | Authenticated |
| Tasks | CRUD, my-tasks | Employee |
| Billing > Invoices | CRUD, mark-paid, PDF download, unpaid | Admin/Employee |
| Billing > Expenses | CRUD, by-category, by-supplier | Admin/Employee |
| Billing > Categories | CRUD, tree, subcategories | Admin/Employee |
| Billing > Stats | dashboard, revenue-by-month, balance | Admin/Employee |
| Stats | dashboard (CRM) | Admin/Employee |
| Search | global search | Authenticated |

Documentación interactiva completa en Swagger: `http://localhost:8080/swagger/index.html`

### Autenticación

JWT Bearer token en header `Authorization: Bearer <token>`.

Roles: `admin` (acceso total), `employee` (gestión sin borrar), `client` (solo perfil propio y citas).

## Testing

### Backend (17 test files, Go)

```powershell
cd backend
go test ./... -count=1 -race    # Todos con race detector
go test ./... -cover             # Con cobertura
```

Cobertura: services (auth, client, employee, appointment, task, invoice, billing stats, search), handlers (auth, search), packages (cache, errors, queue).

### Frontend (14 test files, Vitest)

```powershell
cd frontend
pnpm test -- --run              # Todos
pnpm test:coverage              # Con cobertura
```

Cobertura: hooks, API client, utils, validators, LoginModal, CreateClientModal, EditClientModal.

### CI/CD

GitHub Actions ejecuta en cada push/PR a `main`:
- **Backend**: `go vet` → `go build` → `go test` (race + coverage)
- **Frontend**: `pnpm lint` → `tsc --noEmit` → `vitest` → `next build`

## Features implementadas

### Sitio público
- Landing page con Hero, About, Services, Testimonials, Reviews
- Páginas: Sobre Arnela, Intervención, Formación, Convenios, Contacto
- SEO metadata por página (OpenGraph incluido)
- Dark mode con toggle (light/dark, next-themes)

### Backoffice (dashboard)
- Dashboard con estadísticas (clientes, citas, empleados)
- CRUD completo de clientes con filtros y búsqueda
- Gestión de empleados y especialidades
- Sistema de citas: crear, confirmar, cancelar, slots disponibles
- Facturación: facturas con IVA auto, gastos, categorías jerárquicas
- Exportación de facturas a PDF
- Estadísticas de billing: revenue mensual, gastos por categoría, balance
- Búsqueda global (Ctrl+K)
- Sidebar colapsable con menú por rol
- Dark mode integrado

### Backend
- Clean Architecture (domain, repository, service, handler)
- Rate limiting en endpoints de auth
- Structured logging con Zerolog
- Cache Redis en dashboard stats
- Worker pool async para tareas (email, SMS, calendar)
- Notificaciones email: confirmación y cancelación de citas (SMTP)
- Google Calendar sync para citas
- Health check detallado (/health) + readiness probe (/readiness)
- Migraciones automáticas al arrancar

## Deployment

### Producción con Docker Compose

```powershell
# Copiar y editar variables
copy .env.prod.example .env.prod

# Deploy
docker compose -f docker-compose.yml -f docker-compose.prod.yml --env-file .env.prod up -d --build
```

Requiere en `.env.prod`: `DB_PASSWORD`, `REDIS_PASSWORD`, `JWT_SECRET`, `CORS_ORIGINS`, `NEXT_PUBLIC_API_URL`.

Opcionales: `SMTP_*` (email), `GOOGLE_CALENDAR_*` (calendar sync).

### Variables de entorno

Ver `backend/.env.example` y `.env.prod.example` para la lista completa.

## Convenciones

- **Go**: PascalCase exports, camelCase private, JSON tags en camelCase
- **TypeScript**: PascalCase componentes, camelCase props/variables
- **Git**: Conventional Commits (`feat:`, `fix:`, `chore:`)
- **Arquitectura**: Clean Architecture, repository pattern, dependency injection
- **Errores**: Structured error responses con `pkg/errors`
- **Estado**: Zustand stores, no prop drilling
