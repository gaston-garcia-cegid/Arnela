# Redis & Docker Implementation - Arnela

## ✅ Implementación Completa

### Componentes Implementados

#### 1. **Redis Client** (`pkg/cache/redis.go`)
- Conexión a Redis con health check
- Configuración desde variables de entorno
- Pool de conexiones optimizado

#### 2. **Worker Pool** (`pkg/queue/worker.go`)
- Sistema de cola de tareas asíncronas
- 5 workers concurrentes
- Retry con exponential backoff
- Dead Letter Queue (DLQ) para tareas fallidas
- 4 tipos de tareas soportados:
  - `send_email` - Envío de emails
  - `send_sms` - Envío de SMS
  - `send_whatsapp` - Envío de WhatsApp
  - `sync_calendar` - Sincronización con Google Calendar

#### 3. **Cache Service** (`pkg/cache/service.go`)
- Capa de cache para queries frecuentes
- Cache-Aside pattern con `GetOrSet`
- Generadores de keys predefinidos
- TTLs configurables (Short/Medium/Long/Day)
- Invalidación por patrón

#### 4. **Docker Setup**
- PostgreSQL 16 con health checks
- Redis 7 con autenticación
- Go API con hot-reload
- Volúmenes persistentes
- Red privada entre servicios

---

## 🚀 Quick Start

### 1. Levantar servicios con Docker

```powershell
# Construir y levantar todos los servicios
.\docker-setup.ps1 up

# O manualmente con docker-compose
docker-compose up -d
```

### 2. Verificar estado de servicios

```powershell
# Usando el script
.\docker-setup.ps1 status

# Health check manual
Invoke-RestMethod -Uri "http://localhost:8080/health" | ConvertTo-Json
```

### 3. Ver logs

```powershell
# Logs de todos los servicios
.\docker-setup.ps1 logs

# Logs de un servicio específico
docker-compose logs -f go-api
docker-compose logs -f redis
docker-compose logs -f postgres
```

---

## 📋 Comandos Docker

### Script PowerShell (`docker-setup.ps1`)

```powershell
# Construir imágenes
.\docker-setup.ps1 build

# Iniciar servicios
.\docker-setup.ps1 up

# Detener servicios
.\docker-setup.ps1 down

# Reiniciar servicios
.\docker-setup.ps1 restart

# Ver logs en tiempo real
.\docker-setup.ps1 logs

# Ver estado
.\docker-setup.ps1 status

# Limpiar todo (⚠️ elimina volúmenes)
.\docker-setup.ps1 clean
```

### Docker Compose Manual

```bash
# Construir y levantar
docker-compose up -d --build

# Detener
docker-compose down

# Detener y eliminar volúmenes
docker-compose down -v

# Ver logs
docker-compose logs -f

# Ejecutar comando en contenedor
docker-compose exec go-api sh
docker-compose exec postgres psql -U arnela_user -d arnela_db
docker-compose exec redis redis-cli -a arnela_redis_pass_2024
```

---

## 🔧 Configuración

### Variables de Entorno (`.env`)

```bash
# Application
APP_ENV=development
APP_PORT=8080

# Database
DB_HOST=localhost          # 'postgres' en Docker
DB_PORT=5432
DB_USER=arnela_user
DB_PASSWORD=arnela_secure_pass_2024
DB_NAME=arnela_db
DB_SSLMODE=disable

# Redis
REDIS_HOST=localhost       # 'redis' en Docker
REDIS_PORT=6379
REDIS_PASSWORD=arnela_redis_pass_2024
REDIS_DB=0

# JWT
JWT_SECRET=your_jwt_secret_key_change_in_production
```

**Nota:** En Docker, usar nombres de servicio (`postgres`, `redis`) en lugar de `localhost`.

---

## 📊 Health Check

El endpoint `/health` retorna el estado de todos los servicios:

```json
{
  "status": "healthy",
  "database": "connected",
  "redis": "connected",
  "workers": {
    "tasks_processed": 42,
    "tasks_failed": 3,
    "active_workers": 5
  }
}
```

---

## 🔄 Sistema de Tareas Asíncronas

### Encolar Tareas

```go
// En cualquier handler o servicio
workerPool.EnqueueTask(queue.TaskTypeSendEmail, map[string]interface{}{
    "to":      "client@example.com",
    "subject": "Confirmación de Cita",
    "body":    "Su cita ha sido confirmada...",
})
```

### Tipos de Tareas Disponibles

1. **TaskTypeSendEmail** - Envío de emails
   ```go
   payload: {
       "to": "email@example.com",
       "subject": "Asunto",
       "body": "Contenido HTML"
   }
   ```

2. **TaskTypeSendSMS** - Envío de SMS via Twilio
   ```go
   payload: {
       "phone": "+34123456789",
       "message": "Texto del SMS"
   }
   ```

3. **TaskTypeSendWhatsApp** - WhatsApp via Twilio
   ```go
   payload: {
       "phone": "+34123456789",
       "message": "Mensaje de WhatsApp"
   }
   ```

4. **TaskTypeSyncCalendar** - Sincronizar con Google Calendar
   ```go
   payload: {
       "event_id": "123",
       "action": "create|update|delete",
       "appointment_data": {...}
   }
   ```

### Características del Worker Pool

- **Retry automático:** Hasta 3 intentos con backoff exponencial
- **Dead Letter Queue:** Tareas fallidas se mueven a `arnela:tasks:failed`
- **Graceful shutdown:** Espera a que terminen tareas en progreso
- **Métricas:** Tracking de tareas procesadas y fallidas

---

## 💾 Sistema de Cache

### Uso Básico

```go
// Inicializar cache service
cacheService := cache.NewCacheService(redisClient.Client)

// Set
err := cacheService.Set(ctx, "client:123", clientData, cache.CacheTTLMedium)

// Get
var client domain.Client
err := cacheService.Get(ctx, "client:123", &client)

// Delete
err := cacheService.Delete(ctx, "client:123")

// Delete por patrón
err := cacheService.DeletePattern(ctx, "client:*")
```

### Cache-Aside Pattern (GetOrSet)

```go
var clients []domain.Client

err := cacheService.GetOrSet(
    ctx,
    cache.ClientListCacheKey(),
    &clients,
    cache.CacheTTLMedium,
    func() (interface{}, error) {
        // Loader: se llama solo si no está en cache
        return clientRepo.List(ctx)
    },
)
```

### Keys Predefinidos

```go
cache.ClientCacheKey(id)                     // "client:123"
cache.ClientListCacheKey()                   // "clients:list"
cache.EmployeeCacheKey(id)                   // "employee:456"
cache.EmployeeListCacheKey()                 // "employees:list"
cache.EmployeesBySpecialtyCacheKey("physio") // "employees:specialty:physio"
cache.AppointmentCacheKey(id)                // "appointment:789"
cache.ClientAppointmentsCacheKey(id)         // "appointments:client:123"
cache.EmployeeAppointmentsCacheKey(id, date) // "appointments:employee:456:date:2024-01-15"
cache.DashboardStatsCacheKey()               // "stats:dashboard"
```

### TTL Preconfigurados

```go
cache.CacheTTLShort   // 5 minutos  - datos que cambian frecuentemente
cache.CacheTTLMedium  // 15 minutos - datos moderadamente dinámicos
cache.CacheTTLLong    // 1 hora     - datos relativamente estables
cache.CacheTTLDay     // 24 horas   - datos muy estables
```

---

## 🧪 Testing

### Tests de Integración

```bash
# Tests de cache
cd backend/pkg/cache
go test -v

# Tests de worker pool
cd backend/pkg/queue
go test -v

# Todos los tests
cd backend
go test ./... -v
```

### Cobertura de Tests

- ✅ Cache: Set/Get/Delete/DeletePattern/GetOrSet/Expiration
- ✅ Worker Pool: Enqueue/Process/Retry/DLQ/GracefulShutdown
- ✅ Mock Redis con miniredis (no requiere Redis real)

---

## 📁 Estructura de Archivos

```
backend/
├── cmd/api/
│   └── main.go              # ✅ Redis + Worker Pool integrados
├── pkg/
│   ├── cache/
│   │   ├── redis.go         # ✅ Redis client wrapper
│   │   ├── service.go       # ✅ Cache service con GetOrSet
│   │   └── service_test.go  # ✅ Tests
│   └── queue/
│       ├── worker.go        # ✅ Worker pool + task queue
│       └── worker_test.go   # ✅ Tests
├── docker-compose.yml       # ✅ Postgres + Redis + Go API
├── Dockerfile               # ✅ Multi-stage build
└── .env                     # ✅ Variables de entorno

root/
└── docker-setup.ps1         # ✅ Script de gestión Docker
```

---

## 🔍 Debugging

### Redis CLI

```bash
# Conectar a Redis
docker-compose exec redis redis-cli -a arnela_redis_pass_2024

# Ver todas las keys
KEYS *

# Ver contenido de la cola
LRANGE arnela:tasks 0 -1
LRANGE arnela:tasks:processing 0 -1
LRANGE arnela:tasks:failed 0 -1

# Ver una key específica
GET client:123

# Monitorear comandos en tiempo real
MONITOR

# Ver info del servidor
INFO
```

### PostgreSQL

```bash
# Conectar a PostgreSQL
docker-compose exec postgres psql -U arnela_user -d arnela_db

# Ver tablas
\dt

# Ver migraciones aplicadas
SELECT * FROM schema_migrations;

# Query de prueba
SELECT * FROM users LIMIT 5;
```

---

## 🚨 Troubleshooting

### Error: "Redis connection refused"

```bash
# Verificar que Redis está corriendo
docker-compose ps

# Ver logs de Redis
docker-compose logs redis

# Reiniciar Redis
docker-compose restart redis
```

### Error: "Database connection failed"

```bash
# Verificar health de PostgreSQL
docker-compose exec postgres pg_isready -U arnela_user

# Ver logs
docker-compose logs postgres

# Verificar variables de entorno
docker-compose exec go-api env | grep DB_
```

### Worker Pool no procesa tareas

```bash
# Verificar logs del API
docker-compose logs go-api | grep -i worker

# Ver cola en Redis
docker-compose exec redis redis-cli -a arnela_redis_pass_2024 LLEN arnela:tasks

# Health check
curl http://localhost:8080/health | jq
```

---

## 📈 Próximos Pasos

### Pendiente de Implementar

1. **Email Service** (Sprint 2)
   - Integración con SendGrid
   - Templates HTML
   - Handlers en worker pool

2. **SMS/WhatsApp** (Sprint 1)
   - Integración con Twilio
   - Templates de notificación
   - Confirmaciones de citas

3. **Google Calendar** (Sprint 1)
   - OAuth 2.0 setup
   - Sync bidireccional
   - Handler en worker pool

4. **Cache en Repositories** (Sprint 5)
   - Envolver queries frecuentes con cache
   - Invalidación automática en updates
   - Métricas de hit rate

---

## 📚 Referencias

- [Redis Go Client](https://github.com/go-redis/redis)
- [Docker Compose Reference](https://docs.docker.com/compose/)
- [Worker Pool Pattern](https://gobyexample.com/worker-pools)
- [Cache-Aside Pattern](https://docs.microsoft.com/en-us/azure/architecture/patterns/cache-aside)

---

**Estado:** ✅ Implementación completa y funcional  
**Última actualización:** 2024-12-01  
**Próxima fase:** Integraciones externas (Google Calendar, Twilio)
