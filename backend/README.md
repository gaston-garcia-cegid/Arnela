# Arnela Backend

API REST para el sistema CRM/CMS de Arnela Gabinete. Go + Gin + PostgreSQL + Redis.

## Requisitos

- Go 1.25+
- PostgreSQL 16
- Redis 7

O usar Docker Compose desde la raíz del proyecto.

## Inicio rápido

```powershell
# Copiar configuración
copy .env.example .env

# Levantar dependencias (PostgreSQL + Redis)
docker compose up -d   # desde la raíz del proyecto

# Ejecutar
go run cmd/api/main.go
```

Las migraciones se ejecutan automáticamente al arrancar. El servidor escucha en `http://localhost:8080`.

## Estructura

```
backend/
├── cmd/api/main.go              # Entrypoint, wiring, routes, graceful shutdown
├── config/config.go             # Lectura de env vars con defaults
├── internal/
│   ├── domain/                  # Entidades: User, Client, Employee, Appointment,
│   │                            #   Task, Invoice, Expense, ExpenseCategory
│   ├── repository/              # Interfaces de repositorio
│   │   ├── postgres/            # Implementaciones PostgreSQL (sqlx)
│   │   └── mocks/               # Mocks para tests
│   ├── service/                 # Lógica de negocio (un service por entidad)
│   ├── handler/                 # HTTP handlers Gin (un handler por entidad + health)
│   └── middleware/              # Auth (JWT), logging, rate limiting, CORS
├── pkg/
│   ├── cache/                   # CacheService (Redis, get/set/invalidate con TTL)
│   ├── database/                # NewPostgresDB, RunMigrations
│   ├── email/                   # SMTP Mailer, HTML templates, queue handler
│   ├── errors/                  # AppError, RespondWithError, error codes
│   ├── gcal/                    # Google Calendar service (service account)
│   ├── jwt/                     # TokenManager (generate, validate, claims)
│   ├── logger/                  # Zerolog wrapper (dev=pretty, prod=JSON)
│   ├── pdf/                     # GenerateInvoicePDF (go-pdf/fpdf)
│   └── queue/                   # WorkerPool, TaskHandler, Redis-backed task queue
├── migrations/                  # 000001..000015 SQL up/down
├── docs/                        # Swagger auto-generado (swag init)
└── go.mod
```

## Variables de entorno

| Variable | Default | Descripción |
|----------|---------|-------------|
| `ENVIRONMENT` | `development` | `development` / `production` |
| `SERVER_PORT` | `8080` | Puerto del servidor |
| `CORS_ORIGINS` | `http://localhost:3000` | Orígenes CORS (comma-separated) |
| `DB_HOST` | `localhost` | Host PostgreSQL |
| `DB_PORT` | `5432` | Puerto PostgreSQL |
| `DB_USER` | `arnela_user` | Usuario DB |
| `DB_PASSWORD` | | Contraseña DB |
| `DB_NAME` | `arnela_db` | Nombre DB |
| `DB_SSLMODE` | `disable` | SSL mode |
| `REDIS_HOST` | `localhost` | Host Redis |
| `REDIS_PORT` | `6379` | Puerto Redis |
| `REDIS_PASSWORD` | | Contraseña Redis |
| `REDIS_DB` | `0` | Base de datos Redis |
| `JWT_SECRET` | | Clave secreta JWT (obligatorio en producción) |
| `JWT_EXPIRY_HOURS` | `168` | Expiración del token en horas |
| `SMTP_HOST` | | Host SMTP (opcional, sin configurar = logs only) |
| `SMTP_PORT` | `587` | Puerto SMTP |
| `SMTP_USERNAME` | | Usuario SMTP |
| `SMTP_PASSWORD` | | Contraseña SMTP |
| `SMTP_FROM` | `no-reply@arnela.es` | Dirección from |
| `GOOGLE_CALENDAR_CREDENTIALS` | | JSON credentials service account (opcional) |
| `GOOGLE_CALENDAR_ID` | | ID del calendario (opcional) |

## API Endpoints

### Health

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/health` | Estado detallado (DB, Redis, Go runtime, uptime) |
| GET | `/readiness` | Probe liviano para healthchecks |
| GET | `/swagger/*` | Swagger UI |

### Auth (`/api/v1/auth`)

| Método | Ruta | Auth | Descripción |
|--------|------|------|-------------|
| POST | `/register` | No (rate limited) | Registrar usuario |
| POST | `/login` | No (rate limited) | Login, devuelve JWT |
| GET | `/me` | JWT | Usuario actual |

### Clients (`/api/v1/clients`)

| Método | Ruta | Rol | Descripción |
|--------|------|-----|-------------|
| GET | `/me` | Any | Perfil propio del cliente |
| POST | `/` | Admin, Employee | Crear cliente |
| GET | `/` | Admin, Employee | Listar (paginado, filtros) |
| GET | `/:id` | Admin, Employee | Detalle |
| PUT | `/:id` | Admin, Employee | Actualizar |
| DELETE | `/:id` | Admin | Soft delete |

**Filtros de listado**: `?search=`, `?isActive=`, `?city=`, `?province=`, `?page=`, `?pageSize=`

### Employees (`/api/v1/employees`)

| Método | Ruta | Rol | Descripción |
|--------|------|-----|-------------|
| GET | `/` | Any | Listar empleados |
| GET | `/me` | Employee | Mi perfil de empleado |
| GET | `/:id` | Any | Detalle |
| GET | `/specialty/:specialty` | Any | Por especialidad |
| POST | `/` | Admin | Crear |
| PUT | `/:id` | Admin | Actualizar |
| DELETE | `/:id` | Admin | Eliminar |

### Appointments (`/api/v1/appointments`)

| Método | Ruta | Rol | Descripción |
|--------|------|-----|-------------|
| GET | `/therapists` | Any | Terapeutas disponibles |
| GET | `/available-slots` | Any | Horarios disponibles |
| POST | `/` | Any | Crear cita |
| GET | `/:id` | Any | Detalle |
| PUT | `/:id` | Admin, Employee | Actualizar |
| POST | `/:id/cancel` | Any | Cancelar |
| GET | `/me` | Client | Mis citas |
| GET | `/` | Admin, Employee | Listar todas |
| POST | `/:id/confirm` | Admin, Employee | Confirmar (envía email + calendar sync) |

### Tasks (`/api/v1/tasks`)

| Método | Ruta | Rol | Descripción |
|--------|------|-----|-------------|
| POST | `/` | Admin, Employee | Crear tarea |
| GET | `/` | Admin, Employee | Listar (filtros) |
| GET | `/me` | Employee | Mis tareas |
| PUT | `/:id` | Admin, Employee | Actualizar |
| DELETE | `/:id` | Admin | Eliminar |

### Billing - Invoices (`/api/v1/billing/invoices`)

| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/` | Crear factura |
| GET | `/` | Listar facturas |
| GET | `/:id` | Detalle |
| GET | `/number/:number` | Por número |
| PUT | `/:id` | Actualizar |
| DELETE | `/:id` | Eliminar (no pagadas) |
| POST | `/:id/mark-paid` | Marcar como pagada |
| GET | `/:id/pdf` | Descargar PDF |
| GET | `/client/:clientId` | Facturas de un cliente |
| GET | `/unpaid` | Facturas pendientes |

### Billing - Expenses (`/api/v1/billing/expenses`)

| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/` | Crear gasto |
| GET | `/` | Listar gastos |
| GET | `/:id` | Detalle |
| PUT | `/:id` | Actualizar |
| DELETE | `/:id` | Eliminar |
| GET | `/category/:categoryId` | Por categoría |
| GET | `/supplier/:supplier` | Por proveedor |

### Billing - Categories (`/api/v1/billing/expense-categories`)

| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/` | Crear categoría |
| GET | `/` | Listar |
| GET | `/tree` | Árbol jerárquico |
| GET | `/parents` | Solo padres |
| GET | `/:id` | Detalle |
| PUT | `/:id` | Actualizar |
| DELETE | `/:id` | Eliminar |
| GET | `/:id/subcategories` | Subcategorías |

### Billing - Stats (`/api/v1/billing`)

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/dashboard` | Resumen financiero |
| GET | `/revenue-by-month` | Ingresos mensuales |
| GET | `/expenses-by-category` | Gastos por categoría |
| GET | `/balance` | Balance general |

### Stats (`/api/v1/stats`)

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/dashboard` | Estadísticas CRM (clientes, citas, empleados) |

### Search

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/api/v1/search?q=` | Búsqueda global |

## Testing

```powershell
go test ./... -count=1                # Todos los tests
go test ./... -race                    # Con race detector
go test ./... -cover                   # Con cobertura
go test ./internal/service/... -v      # Services verbose
go test ./internal/handler/... -v      # Handlers verbose
go test ./pkg/... -v                   # Packages verbose
```

17 archivos de test cubriendo:
- **Services**: auth, client (incl. reactivación, isActive bug), employee, appointment, task, invoice, billing stats, search
- **Handlers**: auth, search
- **Packages**: cache (miniredis), errors, queue (worker pool)

Framework: `testify` (assert + mock) + `miniredis` para tests de Redis.

## Migraciones

Las migraciones se ejecutan automáticamente al arrancar. Para gestión manual:

```powershell
# Instalar CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Ejecutar
migrate -path migrations -database "postgres://user:pass@localhost:5432/arnela_db?sslmode=disable" up

# Rollback
migrate -path migrations -database "postgres://..." down 1

# Crear nueva migración
migrate create -ext sql -dir migrations -seq nombre_migracion
```

### Migraciones actuales (15)

| # | Nombre |
|---|--------|
| 1 | create_users_table |
| 2 | create_clients_table |
| 3 | add_nif_field |
| 4 | create_appointments |
| 5 | create_employees_table |
| 6 | update_appointments_employee_fk |
| 7 | add_room_to_appointments |
| 8 | consolidate_dni_cif |
| 9 | create_invoices_table |
| 10 | create_expense_categories_table |
| 11 | create_expenses_table |
| 12 | seed_expense_categories |
| 13 | add_soft_delete_billing |
| 14 | add_due_date_to_invoices |
| 15 | create_tasks_table |

## Swagger

Regenerar después de cambiar anotaciones en handlers:

```powershell
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/api/main.go -o docs
```

## Integraciones opcionales

### Email (SMTP)

Configurar `SMTP_*` en `.env`. Sin configurar, los envíos se loguean sin enviar. Se usa para:
- Confirmación de cita
- Cancelación de cita

Los emails se procesan asincrónicamente a través del worker pool de Redis.

### Google Calendar

Configurar `GOOGLE_CALENDAR_CREDENTIALS` (ruta al JSON de service account) y `GOOGLE_CALENDAR_ID` (ID del calendario destino). La cuenta de servicio debe tener permiso de escritura sobre ese calendario.

Sin configurar, las tareas `sync_calendar` se descartan con un log claro (no hay llamada a Google).

Con configurar: al **confirmar**, **actualizar** o **cancelar** una cita, la API encola un trabajo que crea, actualiza o borra el evento y persiste `google_calendar_event_id` en la tabla `appointments`. Los eventos incluyen recordatorio por email a **24 h** y popup a **30 min** (Google Calendar).

**Comprobación manual:** confirmar una cita de prueba → ver el evento en Google Calendar → cambiar hora desde el backoffice → cancelar; revisar logs del API (`Google Calendar event created/updated/deleted`).
