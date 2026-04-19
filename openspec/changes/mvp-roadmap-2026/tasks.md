# Tasks — mvp-roadmap-2026

## Fase 0: Detalle de gasto (hueco 404)

- [x] 0.1 Tests Vitest + RTL para vista de detalle (`ExpenseDetailView`): proveedor/importe/categoría, categoría sin nombre, estado de carga.
- [x] 0.2 Ruta App Router `frontend/src/app/dashboard/backoffice/billing/expenses/[id]/page.tsx`: `getById`, token/id inválidos, `NotFoundError`, enlace volver al listado.

## Fase 1: Tasks backoffice (listado + detalle mínimo + sidebar)

- [x] 1.1 Backend: `GET /api/v1/tasks/:id` en `backend/internal/handler/task_handler.go` y registro en `backend/cmd/api/main.go` **después** de `GET /tasks/me`; `TaskService.GetTask`; 404; roles `admin`+`employee`.
- [x] 1.2 Backend: tests servicio + handler para `GET /tasks/:id`.
- [x] 1.3 `frontend/src/types/task.ts` (JSON camelCase alineado a `domain.Task`).
- [x] 1.4 `frontend/src/lib/api.ts` → `api.tasks` + `taskList.ts` / tests de normalización.
- [x] 1.5 `frontend/src/components/backoffice/BackofficeSidebar.tsx`: entrada **Tareas** → `/dashboard/backoffice/tasks` (p. ej. `ListTodo`); visible admin y employee (no `employeeHidden`).
- [x] 1.6 `frontend/src/app/dashboard/backoffice/tasks/page.tsx`: **admin** `api.tasks.list` con `page`/`pageSize` y filtros opcionales `assigneeId`/`status`; **employee** `api.tasks.mine`; tabla shadcn; loading/error/vacío; `router.push` a `./tasks/[id]` y enlace a `./tasks/new`.
- [x] 1.7 `frontend/src/app/dashboard/backoffice/tasks/[id]/page.tsx`: detalle lectura (título, descripción, estado, prioridad, fechas, creator/assignee); volver al listado; opcional selector estado + `api.tasks.update` + toast.
- [x] 1.8 `frontend/src/app/dashboard/backoffice/tasks/new/page.tsx`: POST `api.tasks.create`; **admin**: `Select` asignatario desde `api.employees.list` (activos); **employee**: `api.employees.getMyProfile` y `assigneeId` fijo sin selector; validación mínima (título, asignatario).
- [x] 1.9 Strict TDD: `frontend/src/lib/taskLabels.ts` (etiquetas ES `status`/`priority`) + `src/lib/__tests__/taskLabels.test.ts` **antes** de usar en UI; componente presentacional p. ej. `TaskDetailView.tsx` + tests RTL en `src/components/tasks/__tests__/` (RED→GREEN→REFACTOR en `apply-progress.md`).
- [x] 1.10 Cierre F1: `pnpm test -- --run` y `pnpm exec tsc --noEmit`; smoke manual admin (list+filtros+alta) y employee (mis tareas+alta+detalle).

## Fase 2: E2E críticos + CI (2 jun–27 jun, propuesta)

- [ ] 2.1 `frontend/package.json`: devDependency `@playwright/test`, script `test:e2e` (p. ej. `playwright test`) y opcional `test:e2e:ui`.
- [ ] 2.2 `frontend/playwright.config.ts`: `baseURL` (Next dev o preview), `testDir: e2e`, timeouts razonables, `forbidOnly` en CI.
- [ ] 2.3 `frontend/e2e/.env.example` (o doc): `E2E_BASE_URL`, credenciales/usuario seed **sin** secretos reales; alinear con `docker-compose` / API local.
- [ ] 2.4 `frontend/e2e/auth-backoffice.spec.ts`: login backoffice (reutilizar selectores de `LoginModal` / rutas reales) → expect URL bajo `/dashboard/backoffice`.
- [ ] 2.5 `frontend/e2e/billing-invoice.spec.ts`: smoke facturación (p. ej. abrir listado facturas o flujo corto acordado con datos seed).
- [ ] 2.6 `frontend/e2e/appointment.spec.ts`: smoke citas (listado o crear mínimo según seed estable).
- [ ] 2.7 `.github/workflows/ci.yml` (o nuevo `e2e.yml`): job **e2e** con servicios (Postgres/Redis si aplica), migraciones + API + `pnpm build && pnpm start` o `webServer` de Playwright; `pnpm exec playwright install --with-deps`; subir artefacto (trace/screenshot) en fallo.
- [ ] 2.8 `frontend/README.md` o `docs/E2E.md`: cómo ejecutar E2E local, variables obligatorias, alcance smoke vs regresión completa.

## Fase 3: Hardening (30 jun–25 jul)

- [ ] 3.1 Inventario rate limits sensibles (login/registro ya existentes) + gaps en `backend/`; tareas de implementación en change dedicado.
- [ ] 3.2 Matriz permisos admin/employee/client vs rutas Gin y páginas App Router; acciones correctivas documentadas.
- [ ] 3.3 Métricas/observabilidad (Prometheus/OpenTelemetry o mínimo logs estructurados) según acuerdo ops.

## Fase 4: Buffer UX (28 jul–22 ago)

- [ ] 4.1 Backlog deuda menor post-MVP priorizado con negocio (referencia `MVP_ROADMAP.md`).
