# Apply progress — mvp-roadmap-2026

**Change**: mvp-roadmap-2026  
**Mode**: Strict TDD  

## TDD cycle evidence (acumulado)

| Task | RED | GREEN | REFACTOR |
|------|-----|-------|----------|
| 0.1 `ExpenseDetailView` | Import inexistente tras tests RTL. | Componente + 3 tests. | `Readonly` props; `aria-live` en carga. |
| 0.2 `expenses/[id]` | — | Página con `getById` y errores. | `id` sin ternario anidado. |
| 1.1 `GET /tasks/:id` | Tests handler/servicio. | Handler + ruta tras `/me`; servicio `ErrTaskNotFound`. | — |
| 1.2 Tests backend GET | Suite nueva. | `task_handler_test` + `task_service_test` GetTask. | — |
| 1.3 Tipos `Task` | **Triangulation skipped** (solo tipos). | `types/task.ts`. | — |
| 1.4 `api.tasks` | `taskList.test` / `api.test` sin impl. | `taskList.ts`, `api.tasks`. | Cast `buildQueryParams`. |
| 1.5–1.10 Tasks UI | Varios (ver commits previos). | Sidebar, páginas, `taskLabels`, `TaskDetailView`, tests. | Sonar / formato fechas. |
| 2.1–2.3 Playwright tooling | **Triangulation skipped** (deps + config + `.env.example`). | `@playwright/test`, scripts, `playwright.config.ts`, `e2e/.env.example`, carga opcional `e2e/.env`. | `PLAYWRIGHT_SKIP_WEBSERVER`; `vitest` exclude `e2e/`; eslint ignore `e2e/`. |
| 2.4 `loginBackoffice` | — | Helper con Navbar + `LoginModal` (labels Email/Contraseña, botón Ingresar). | — |
| 2.5–2.7 E2E specs | Sin credenciales → `test.skip` (3 skipped, exit 0). | Specs auth, facturas, citas con `loginBackoffice` + headings reales. | — |
| 2.8 + 2.11 + 2.12 Docs | — | `docs/E2E.md`, enlace en `frontend/README.md`. | — |
| 2.9–2.10 CI workflow | Pendiente implementación. | — | — |

## Archivos tocados (batch Fase 2 parcial 2.1–2.8, 2.11–2.12)

- `frontend/package.json`, `frontend/pnpm-lock.yaml`
- `frontend/playwright.config.ts`
- `frontend/e2e/.env.example`, `frontend/e2e/helpers/loginBackoffice.ts`
- `frontend/e2e/auth-backoffice.spec.ts`, `frontend/e2e/billing-invoice.spec.ts`, `frontend/e2e/appointment.spec.ts`
- `frontend/vitest.config.ts`, `frontend/eslint.config.mjs`, `frontend/.gitignore`
- `docs/E2E.md`, `frontend/README.md`
- `openspec/changes/mvp-roadmap-2026/tasks.md`

## Desviaciones

Specs usan `test.skip` si faltan `E2E_ADMIN_*` para que CI/Vitest no fallen sin secrets; ejecución real requiere `e2e/.env`. Workflow `e2e.yml` (2.9–2.10) no incluido en este batch.

## Issues

Ninguno bloqueante.

## Estado

**Fase 1** completa. **Fase 2**: 2.1–2.8 y 2.11–2.12 listas; **pendientes 2.9–2.10** (GitHub Actions full stack + artefactos).
