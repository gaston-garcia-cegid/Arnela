# Apply progress — mvp-roadmap-2026

**Change**: mvp-roadmap-2026  
**Mode**: Strict TDD  

## TDD cycle evidence (acumulado)

| Task | RED | GREEN | REFACTOR |
|------|-----|-------|----------|
| 0.1 `ExpenseDetailView` | Suite falla: import inexistente tras tests RTL. | `ExpenseDetailView.tsx` + 3 tests en verde. | Props `Readonly<…>`; `aria-live="polite"` en carga. |
| 0.2 Ruta `[id]/page` | Cubierto por vista; página con `getById` y errores. | `expenses/[id]/page.tsx`. | Extracción de `id` sin ternario anidado. |
| 1.1 `GET /tasks/:id` | Tests handler/servicio exigen 404/200 antes de estabilizar contrato HTTP. | `GetTask` en handler; ruta en `main.go` tras `/me`; `GetTask` en servicio mapea `sql.ErrNoRows` → `ErrTaskNotFound`. | — |
| 1.2 Tests backend GET | `TestTaskService_GetTask_*` + `TestTaskHandler_GetTask_*` (invalid, 404, success). | `go test` handler+service en verde. | — |
| 1.3 Tipos `Task` | Sin tests dedicados (solo tipos); cubierto por uso en 1.4. **Triangulation skipped:** artefacto estructural sin ramas. | `frontend/src/types/task.ts`. | — |
| 1.4 `api.tasks` + normalización | `taskList.test.ts` falla sin `normalizeTaskListResponse`; `api.test.ts` describe `api.tasks` sin implementación previa. | `taskList.ts`, `api.tasks.*`, tests `taskList` + `api.tasks`. | `buildQueryParams` cast a `Record<string, any>` alineado al helper existente. |

## Archivos tocados (este batch 1.1–1.4)

- `backend/internal/service/task_service.go` — `GetTask` con `ErrTaskNotFound`.
- `backend/internal/service/task_service_test.go` — `GetTask` success / not found.
- `backend/internal/handler/task_handler.go` — `GetTask`.
- `backend/internal/handler/task_handler_test.go` — httptest invalid / 404 / OK.
- `backend/cmd/api/main.go` — `GET /:id` tras `/me`.
- `frontend/src/types/task.ts` (nuevo).
- `frontend/src/lib/taskList.ts` (nuevo).
- `frontend/src/lib/__tests__/taskList.test.ts` (nuevo).
- `frontend/src/lib/api.ts` — `api.tasks`.
- `frontend/src/lib/__tests__/api.test.ts` — `api.tasks`.

## Desviaciones

No hay `design.md`; 1.3 sin tests unitarios propios (tipos); evidencia en 1.4.

## Issues

Ninguno.

## Estado

Fase 0 completa. Fase 1: **1.1–1.4 completas**; pendientes **1.5–1.10** (sidebar, páginas, formulario, TDD UI, cierre).
