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
| 1.5 Sidebar Tareas | — | `ListTodo` + ruta en `BackofficeSidebar`. | — |
| 1.6 `tasks/page` | — | Listado admin/employee, filtros, paginación, navegación. | `useEffect` auth sin `return` redundante; `tableRows` sin ternario anidado; empleados activos en filtro. |
| 1.7 `tasks/[id]` | — | `TaskDetailView` + `getById` + `update` con payload completo + selector estado. | Auth `useEffect` simplificado. |
| 1.8 `tasks/new` | — | Formulario create; admin select / employee `getMyProfile`. | Auth `useEffect` simplificado. |
| 1.9 `taskLabels` + `TaskDetailView` | Imports rotos (`taskLabels`, `TaskDetailView`). | Implementación + tests (incl. fallback raw). | `formatDate` sin `dateStyle`+`timeStyle` (jsdom). |
| 1.10 Cierre F1 | — | `pnpm test -- --run` + `tsc --noEmit` verdes. | Smoke manual: pendiente verificación humana admin/employee. |

## Archivos tocados (último batch 1.5–1.10)

- `frontend/src/lib/taskLabels.ts`, `frontend/src/lib/__tests__/taskLabels.test.ts`
- `frontend/src/components/tasks/TaskDetailView.tsx`, `frontend/src/components/tasks/__tests__/TaskDetailView.test.tsx`
- `frontend/src/components/backoffice/BackofficeSidebar.tsx`
- `frontend/src/app/dashboard/backoffice/tasks/page.tsx`
- `frontend/src/app/dashboard/backoffice/tasks/[id]/page.tsx`
- `frontend/src/app/dashboard/backoffice/tasks/new/page.tsx`

## Desviaciones

Ninguna respecto a `tasks.md`. Sin `design.md` en el change.

## Issues

**Update task**: el backend hace merge en servicio; el front envía siempre campos completos en el PUT de cambio de estado para no pisar título vacío.

## Estado

**Fase 1 completa** (1.1–1.10). Siguiente: **Fase 2** (Playwright + CI) u otro change.
