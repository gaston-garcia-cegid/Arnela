# 🚀 Quick Start - Redis & Docker

## Comandos Esenciales

### 1. Levantar Todo
```powershell
.\docker-setup.ps1 up
```

### 2. Verificar Estado
```powershell
.\docker-setup.ps1 status
```

### 3. Ver Logs
```powershell
.\docker-setup.ps1 logs
```

### 4. Detener Todo
```powershell
.\docker-setup.ps1 down
```

### 5. Reiniciar
```powershell
.\docker-setup.ps1 restart
```

---

## Health Check

```powershell
Invoke-RestMethod http://localhost:8080/health | ConvertTo-Json
```

**Esperado:**
```json
{
  "status": "healthy",
  "database": "connected",
  "redis": "connected",
  "workers": {
    "tasks_processed": 0,
    "tasks_failed": 0,
    "active_workers": 5
  }
}
```

---

## Tests

```bash
# Cache
go test ./pkg/cache/... -v

# Worker Pool
go test ./pkg/queue/... -v

# Todo el backend
go test ./... -v

# Build
go build ./cmd/api
```

---

## Redis CLI

```bash
# Conectar
docker-compose exec redis redis-cli -a arnela_redis_pass_2024

# Ver tareas
LRANGE arnela:tasks 0 -1

# Ver tareas fallidas
LRANGE arnela:tasks:failed 0 -1

# Ver todas las keys
KEYS *
```

---

## PostgreSQL

```bash
# Conectar
docker-compose exec postgres psql -U arnela_user -d arnela_db

# Ver tablas
\dt

# Ver migraciones
SELECT * FROM schema_migrations;
```

---

## Troubleshooting

### Redis no conecta
```bash
docker-compose restart redis
docker-compose logs redis
```

### Base de datos no conecta
```bash
docker-compose restart postgres
docker-compose logs postgres
```

### Workers no procesan tareas
```bash
docker-compose logs go-api | grep -i worker
```

### Limpiar todo
```powershell
.\docker-setup.ps1 clean
```

---

## Documentación

- **Técnica completa:** `REDIS_DOCKER_IMPLEMENTATION.md`
- **Resumen detallado:** `REDIS_DOCKER_COMPLETE.md`
- **Resumen ejecutivo:** `SUMMARY.md`
- **Quick start:** `QUICKSTART.md` (este archivo)

---

## URLs

- API: http://localhost:8080
- Health: http://localhost:8080/health
- Swagger: http://localhost:8080/swagger/index.html
- PostgreSQL: localhost:5432
- Redis: localhost:6379

---

## Estado Actual

✅ Redis integrado  
✅ Worker pool funcionando (5 workers)  
✅ Cache service listo  
✅ Docker completamente funcional  
✅ 58 tests pasando (42 anteriores + 16 nuevos)  
✅ Build exitoso  

---

_Para empezar: `.\docker-setup.ps1 up`_
