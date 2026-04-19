# Pruebas E2E (Playwright)

Smoke end-to-end del front contra una API real. Los specs viven en `frontend/e2e/`.

## Prerrequisitos

1. API Go y Postgres/Redis según `docker-compose.yml` o entorno local habitual.
2. Un usuario **admin** válido (misma base que usa el backoffice).
3. Node 22+ y pnpm en `frontend/`.

## Configuración

```powershell
cd frontend
Copy-Item e2e/.env.example e2e/.env
# Editar e2e/.env: E2E_ADMIN_EMAIL, E2E_ADMIN_PASSWORD, y opcionalmente E2E_BASE_URL / E2E_API_URL
```

`playwright.config.ts` carga automáticamente `e2e/.env` si existe (sin publicar secretos en git).

Alinear el front con la API:

- En `.env` del front, `NEXT_PUBLIC_API_URL` debe apuntar a la misma API que usará el navegador (p. ej. `http://127.0.0.1:8080/api/v1`).

## Instalar navegadores Playwright

```powershell
cd frontend
pnpm exec playwright install chromium
```

## Ejecutar E2E local

Con el **dev server** ya levantado (`pnpm dev` en otra terminal), evita que Playwright intente otro `pnpm dev`:

```powershell
cd frontend
$env:PLAYWRIGHT_SKIP_WEBSERVER="1"
pnpm test:e2e
```

Si preferís que Playwright arranque solo el front:

```powershell
cd frontend
pnpm test:e2e
```

## CI / GitHub Actions

- Usar **secrets** del repositorio (`E2E_ADMIN_EMAIL`, `E2E_ADMIN_PASSWORD`) nunca valores reales en el repo.
- Job recomendado: Postgres + Redis + migraciones + API en background + `pnpm build` + `pnpm start` + `pnpm exec playwright install --with-deps` + `pnpm test:e2e`. Plantilla pendiente en `.github/workflows/e2e.yml` (tarea 2.9–2.10 del change `mvp-roadmap-2026`).

## Alcance actual (smoke)

| Spec | Qué valida |
|------|-------------|
| `auth-backoffice.spec.ts` | Login desde home → URL bajo `/dashboard/backoffice`. |
| `billing-invoice.spec.ts` | Listado facturas, heading **Facturas**. |
| `appointment.spec.ts` | Agenda, heading **Citas** / **Mis Citas**. |

## Problemas frecuentes

- **Timeouts**: subir `timeout` en `playwright.config.ts` o comprobar red/API.
- **Login falla**: revisar credenciales y que el modal esté visible (viewport desktop para el botón *Iniciar sesión* del Navbar).
