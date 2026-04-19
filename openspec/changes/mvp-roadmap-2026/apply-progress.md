# Apply progress — mvp-roadmap-2026

**Change**: mvp-roadmap-2026  
**Mode**: Strict TDD  
**Batch**: Fase 0 (detalle gasto)

## TDD cycle evidence

| Task | RED | GREEN | REFACTOR |
|------|-----|-------|----------|
| 0.1 `ExpenseDetailView` | Suite falla: import inexistente `ExpenseDetailView.tsx` tras añadir tests que describen proveedor, categoría sin nombre y carga. | Implementado `ExpenseDetailView.tsx`; 3 tests en verde. | Props `Readonly<…>`; mensaje de carga con `aria-live="polite"` en lugar de `role="status"`. |
| 0.2 Ruta `[id]/page` | Cubierto por 0.1 (UI); página cliente delgada que compone la vista y llama a `api.billing.expenses.getById`. Sin test E2E (no hay runner en proyecto). | Página creada con `useParams`, guards token/id, manejo `NotFoundError`, toast y botón “Volver al listado”. | Extracción de `id` sin ternario anidado (legibilidad / Sonar). |

## Archivos tocados

- `frontend/src/components/billing/ExpenseDetailView.tsx` (nuevo)
- `frontend/src/components/billing/__tests__/ExpenseDetailView.test.tsx` (nuevo)
- `frontend/src/app/dashboard/backoffice/billing/expenses/[id]/page.tsx` (nuevo)

## Desviaciones del diseño

No hay `design.md` en este change; alineado con `proposal.md` (cerrar 404 en detalle de gasto).

## Issues

Ninguno.

## Estado

Fase 0: **completa** (tareas 0.1–0.2). Listo para verify u otro batch.
