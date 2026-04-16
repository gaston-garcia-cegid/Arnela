# Deployment

Guía de despliegue de Arnela en producción.

## Arquitectura

```
Internet
  │
  ▼
Nginx (:80/443)
  ├── /api/*    → go-api (:8080)
  ├── /swagger/ → go-api (:8080)
  ├── /health   → go-api (:8080)
  └── /*        → frontend (:3000)
                      │
                go-api ──┬── PostgreSQL (:5432)
                         └── Redis (:6379)
```

Todos los servicios corren en contenedores Docker conectados a la red interna `arnela-network`. Solo Nginx expone puertos al host.

## Pre-requisitos

- Docker y Docker Compose instalados en el servidor
- Git para clonar el repositorio
- Certificado TLS (Let's Encrypt o similar) configurado a nivel de reverse proxy externo o directamente en Nginx

## Configuración

### 1. Crear archivo de variables

```bash
cp .env.prod.example .env.prod
```

Editar `.env.prod` con valores de producción:

```env
# Obligatorios
DB_PASSWORD=<contraseña segura>
REDIS_PASSWORD=<contraseña segura>
JWT_SECRET=<generar con: openssl rand -base64 32>
CORS_ORIGINS=https://tu-dominio.com
NEXT_PUBLIC_API_URL=https://tu-dominio.com/api/v1

# Opcionales
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=<email>
SMTP_PASSWORD=<app password>
SMTP_FROM=no-reply@arnela.es
GOOGLE_CALENDAR_CREDENTIALS=<JSON service account>
GOOGLE_CALENDAR_ID=<calendar id>
```

**Variables obligatorias**:

| Variable | Descripción |
|----------|-------------|
| `DB_PASSWORD` | Contraseña de PostgreSQL |
| `REDIS_PASSWORD` | Contraseña de Redis |
| `JWT_SECRET` | Clave para firmar tokens JWT |
| `CORS_ORIGINS` | URL pública del frontend (sin trailing slash) |
| `NEXT_PUBLIC_API_URL` | URL pública de la API vista desde el navegador del usuario |
| `NEXT_PUBLIC_SENTRY_DSN` | (Opcional) DSN de Sentry para el frontend; en Docker se pasa como **build arg** (ver `docker-compose.prod.yml`) |

### 2. Verificar rutas de datos

El docker-compose de producción monta volúmenes en rutas del host:

```yaml
postgres: /DATA/AppData/Arnela/data/postgres
redis:    /DATA/AppData/Arnela/data/redis
```

Crear estos directorios si no existen, o modificar las rutas en `docker-compose.prod.yml`.

## Deploy

### Build y arranque

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml --env-file .env.prod up -d --build
```

### Verificar estado

```bash
# Estado de contenedores
docker compose -f docker-compose.yml -f docker-compose.prod.yml ps

# Health check del backend
curl http://localhost:8080/health

# Readiness probe
curl http://localhost:8080/readiness

# Logs
docker compose -f docker-compose.yml -f docker-compose.prod.yml logs -f go-api
docker compose -f docker-compose.yml -f docker-compose.prod.yml logs -f frontend
docker compose -f docker-compose.yml -f docker-compose.prod.yml logs -f nginx
```

### Detener

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml --env-file .env.prod down
```

## Servicios

### Backend (go-api)

- **Imagen base**: `golang:1.23-alpine` (build) + `alpine:latest` (runtime)
- `GOTOOLCHAIN=auto` en el Dockerfile permite que Go descargue la versión que necesita `go.mod` (actualmente 1.25)
- Las migraciones SQL se ejecutan automáticamente al arrancar
- Healthcheck: `wget -qO- http://localhost:8080/readiness` cada 15s
- Modo producción: `gin.ReleaseMode`, JSON logs via zerolog

### Frontend

- **Build**: Multi-stage (deps → build → runner) con `node:22-alpine`
- `NEXT_PUBLIC_API_URL` se inyecta como build arg (baked en el JS del cliente)
- Output standalone: el runtime solo necesita `node server.js`
- Corre como usuario `nextjs` (UID 1001)

### Nginx

- Reverse proxy para frontend (/) y backend (/api/, /swagger/, /health)
- `client_max_body_size 10M`
- WebSocket headers para Next.js HMR
- Puerto expuesto al host: `8101:80` (modificar según necesidad)

### PostgreSQL

- `postgres:16-alpine`
- Healthcheck con `pg_isready`
- Datos persistidos en volumen del host

### Redis

- `redis:7-alpine` con password
- Healthcheck con `redis-cli ping`
- Datos persistidos en volumen del host

## HTTPS / TLS

El setup actual de Nginx solo escucha en puerto 80. Para HTTPS:

**Opción A: Reverse proxy externo (recomendado)**

Si tienes un proxy externo (Caddy, Traefik, o el propio Nginx del host) que maneja TLS, apunta al puerto 8101 del contenedor Nginx.

**Opción B: TLS en el Nginx del contenedor**

Montar certificados en el contenedor y modificar `nginx/nginx.conf`:

```nginx
server {
    listen 443 ssl;
    ssl_certificate /etc/nginx/certs/fullchain.pem;
    ssl_certificate_key /etc/nginx/certs/privkey.pem;
    # ... resto de la config
}

server {
    listen 80;
    return 301 https://$host$request_uri;
}
```

Agregar volumen en `docker-compose.prod.yml`:

```yaml
nginx:
  volumes:
    - ./nginx/nginx.conf:/etc/nginx/conf.d/default.conf:ro
    - /path/to/certs:/etc/nginx/certs:ro
```

## Actualización

```bash
# Obtener últimos cambios
git pull

# Rebuild y restart
docker compose -f docker-compose.yml -f docker-compose.prod.yml --env-file .env.prod up -d --build

# Verificar
docker compose -f docker-compose.yml -f docker-compose.prod.yml ps
curl http://localhost:8080/health
```

## Backups y recuperación (runbook)

Los nombres de contenedor asumen el `docker-compose.prod.yml` del repositorio (`arnela-postgres`, `arnela-redis`, `arnela-go-api`). Ajusta usuario y base si cambiaste `DB_USER` / `DB_NAME`.

### PostgreSQL

**Backup lógico (SQL plano):**

```bash
docker exec arnela-postgres pg_dump -U arnela_user arnela_db > backup_$(date +%Y%m%d).sql
```

**Backup en formato custom** (recomendado para restores controlados y compresión con `pg_restore`):

```bash
docker exec arnela-postgres pg_dump -U arnela_user -Fc arnela_db > backup_$(date +%Y%m%d).dump
```

**Restore desde SQL plano** (sobrescribe datos de la base destino; hacer en ventana de mantenimiento):

```bash
# Parar el API evita escrituras durante el restore
docker compose -f docker-compose.yml -f docker-compose.prod.yml stop go-api

cat backup.sql | docker exec -i arnela-postgres psql -U arnela_user arnela_db

docker compose -f docker-compose.yml -f docker-compose.prod.yml start go-api
```

**Restore desde `.dump`:**

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml stop go-api

docker exec -i arnela-postgres pg_restore -U arnela_user --clean --if-exists -d arnela_db < backup_YYYYMMDD.dump

docker compose -f docker-compose.yml -f docker-compose.prod.yml start go-api
```

**Comprobaciones post-restore:** `curl` a `/readiness`, login en el backoffice, revisar `schema_migrations` en Postgres si hubo dudas sobre el estado de migraciones.

**Retención:** guardar al menos el último backup diario fuera del mismo disco que el servidor (objeto S3, NAS u otra máquina).

### Redis

Redis persiste en el volumen del host (`appendonly` / RDB según imagen y comando). Para forzar un snapshot en caliente:

```bash
docker exec arnela-redis redis-cli -a "$REDIS_PASSWORD" BGSAVE
```

El fichero RDB suele estar bajo el directorio de datos del volumen. **Importante:** restaurar solo un volúmen de Redis coherente con el mismo “punto en el tiempo” que Postgres si usas colas como fuente de verdad; en Arnela la verdad de negocio está en Postgres — Redis puede reconstruirse vacío en un desastre total, a costa de reencolar tareas manualmente o aceptar pérdida de cola.

### Simulacro de incidente (criterio Fase 8)

Hacer al menos una vez al año (o tras cambios mayores de infra):

1. En un entorno de **staging** o máquina aislada, restaurar el último backup de Postgres siguiendo los pasos anteriores.
2. Arrancar `go-api` y comprobar `/health`, `/readiness` y un login real.
3. Documentar tiempo total de RTO observado y incidencias (permisos, versión de Postgres, espacio en disco).
4. Ajustar esta guía si los comandos o nombres de contenedor difieren de vuestro despliegue.

## Observabilidad (frontend)

El frontend puede enviar errores del **navegador** a **Sentry** de forma opcional usando `@sentry/browser` (no se instrumenta el runtime de servidor de Next, para evitar problemas de empaquetado con `standalone`).

| Variable | Descripción |
|----------|-------------|
| `NEXT_PUBLIC_SENTRY_DSN` | DSN del proyecto Sentry (público en el bundle). Si está vacío, no se envían eventos. |

El cliente HTTP envía **`X-Request-ID`** en cada petición a la API; el backend lo refleja y lo registra en logs, lo que ayuda a correlacionar un error de Sentry con una línea de log concreta.

## Rotación de secretos

| Secreto | Acción típica |
|---------|----------------|
| `JWT_SECRET` | Cambiar valor y **reiniciar** `go-api`. Todos los usuarios deberán volver a iniciar sesión. |
| `DB_PASSWORD` | Cambiar en Postgres y en `.env.prod`; reiniciar servicios en orden: `go-api` tras `postgres` sano. |
| `REDIS_PASSWORD` | Actualizar comando de Redis y variables de `go-api`; reinicio coordinado (cola vacía si es posible). |
| `GOOGLE_CALENDAR_CREDENTIALS` | Sustituir JSON de la cuenta de servicio en el entorno y reiniciar `go-api`. |
| Certificados TLS | Renovar en el proxy o volumen de Nginx y recargar Nginx. |

Tras cualquier rotación, comprobar login, una operación de escritura (p. ej. crear borrador de factura) y logs sin errores de autenticación.

## Endurecimiento y seguridad

- **Nginx** (`nginx/nginx.conf`): cabeceras `X-Frame-Options`, `X-Content-Type-Options`, `Referrer-Policy`, `Permissions-Policy` en la respuesta al cliente.
- **API (Gin):** las mismas cabeceras y **`X-Request-ID`** en cada respuesta; rate limit en `/api/v1/auth/login` y `/register` (ver `cmd/api/main.go`).
- **Roles:** rutas sensibles exigen `RequireRole` (empleados, facturación, etc.); revisar nuevos endpoints al añadir features.

## Troubleshooting

### `next build` falla en Windows con `EPERM` / `symlink` (output `standalone`)

Next.js intenta crear enlaces simbólicos al copiar trazas al directorio `standalone`. En Windows, activa **Modo de desarrollador** (Configuración → Privacidad y seguridad → Para desarrolladores) o ejecuta el build dentro de **Docker** / **WSL2** / CI Linux, donde el build de producción del `Dockerfile` no tiene esta limitación.

### El frontend no conecta con la API

1. Verificar `NEXT_PUBLIC_API_URL` en `.env.prod`: debe ser la URL que el navegador del usuario puede alcanzar
2. Verificar `CORS_ORIGINS`: debe coincidir exactamente con el origen del frontend (sin trailing slash)
3. Rebuild del frontend tras cambiar `NEXT_PUBLIC_API_URL` (es un build arg)

### Migraciones fallan

```bash
# Ver logs del backend
docker compose logs go-api | head -50

# Conectar a PostgreSQL
docker exec -it arnela-postgres psql -U arnela_user arnela_db

# Verificar versión de migración
SELECT * FROM schema_migrations;
```

### Contenedor se reinicia continuamente

```bash
# Ver logs de crash
docker compose logs --tail 50 <servicio>

# Verificar healthcheck
docker inspect arnela-go-api | jq '.[0].State.Health'
```

## CI/CD

GitHub Actions (`.github/workflows/ci.yml`) ejecuta en cada push/PR a `main`:

- **Backend**: vet → build → test (race + coverage)
- **Frontend**: lint → type-check → test → build

La versión de Go en CI se lee automáticamente de `backend/go.mod` (`go-version-file`).
