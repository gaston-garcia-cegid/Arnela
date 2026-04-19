# Tasks — mvp-roadmap-2026

## Fase 0: Detalle de gasto (hueco 404)

- [x] 0.1 Tests Vitest + RTL para vista de detalle (`ExpenseDetailView`): proveedor/importe/categoría, categoría sin nombre, estado de carga.
- [x] 0.2 Ruta App Router `frontend/src/app/dashboard/backoffice/billing/expenses/[id]/page.tsx`: `getById`, token/id inválidos, `NotFoundError`, enlace volver al listado.

## Fase 1: Tasks backoffice (listado + detalle mínimo + sidebar)

- [x] 1.1 Backend: `GET /api/v1/tasks/:id` en `backend/internal/handler/task_handler.go` y registro en `backend/cmd/api/main.go` **después** de `GET /tasks/me` (evitar capturar `me` como id); usar `TaskService.GetTask`; 404 si no existe; roles `admin`+`employee` como en list/update.
- [x] 1.2 Backend: test del handler o servicio para `GET /tasks/:id` (éxito + 404).
- [x] 1.3 Frontend: `frontend/src/types/task.ts` con `Task`, `TaskStatus`, `TaskPriority` alineados a `domain.Task` JSON (camelCase).
- [x] 1.4 Frontend: en `frontend/src/lib/api.ts` añadir `api.tasks` — `list` (`GET /tasks` + query), `mine` (`GET /tasks/me`), `getById` (`GET /tasks/:id` tras 1.1), `create`, `update`, `delete`; normalizar `data` a array vacío si el backend devuelve `null`.
- [ ] 1.5 `frontend/src/components/backoffice/BackofficeSidebar.tsx`: ítem **Tareas** con enlace `/dashboard/backoffice/tasks` (icono tipo checklist); visible para admin y employee.
- [ ] 1.6 `frontend/src/app/dashboard/backoffice/tasks/page.tsx`: tabla o cards; **admin** usa `api.tasks.list` (paginación/filtros `assigneeId`, `status`); **employee** usa `api.tasks.mine`; acciones: ver detalle (`router.push` a `[id]`), enlace a nueva tarea.
- [ ] 1.7 `frontend/src/app/dashboard/backoffice/tasks/[id]/page.tsx`: detalle lectura (título, descripción, estado, prioridad, fechas, ids creator/assignee); botón volver; opcional cambio de estado vía `update` si encaja en una sesión.
- [ ] 1.8 `frontend/src/app/dashboard/backoffice/tasks/new/page.tsx`: formulario POST `create`; **admin**: selector de asignatario (empleados activos vía API existente, p. ej. `api.employees.list`); **employee**: obtener `assigneeId` con `GET /employees/me` y fijarlo sin selector.
- [ ] 1.9 Strict TDD: tests Vitest/RTL **antes** de cerrar UI (p. ej. componente presentacional de fila o detalle, o helper de etiquetas estado/prioridad) — ciclo RED→GREEN→REFACTOR documentado en `apply-progress.md`.
- [ ] 1.10 Cierre: `pnpm test -- --run` y `pnpm exec tsc --noEmit` en `frontend/`; smoke manual admin vs employee (listado, crear, detalle).

## Fase 2: E2E críticos + CI (propuesta)

- [ ] 2.1 Desglosar en change/tooling: Playwright (o acordado), login + flujo factura + flujo cita; job CI; umbral mínimo documentado.

## Fase 3: Hardening (rate limits, permisos, métricas)

- [ ] 3.1 Tareas técnicas según `proposal.md` cuando se abra el change correspondiente.

## Fase 4: Buffer UX / deuda menor

- [ ] 4.1 Backlog liviano post-MVP según priorización negocio.
