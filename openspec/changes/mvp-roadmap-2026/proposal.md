# Proposal: MVP estado (código) y roadmap Q2–Q3 2026

## Intent

Fijar una lectura **solo desde el código** (sin documentación de producto) del grado de completitud del MVP operativo y una **hoja de ruta faseada** con fechas/ventanas para cerrar brechas visibles y riesgos de entrega.

## Scope

### In Scope

- Inventario de capacidades **presentes** en API (`main.go`) y rutas **App Router** relevantes.
- Juicio de **madurez MVP** (usable en producción vs. huecos).
- Lista de **faltantes** priorizados y tabla **roadmap** (fases, fechas, duración relativa).

### Out of Scope

- Re-leer `MVP_ROADMAP.md` u otros docs largos (pedido explícito).
- Implementación de código en este change (solo propuesta).
- Estimación legal/fiscal fuera de técnico.

## Capabilities

### New Capabilities

- None — artefacto de planificación; no altera contratos de producto hasta que un change posterior cree specs.

### Modified Capabilities

- None

## Approach

- **Backend**: rutas bajo `/api/v1` — auth, clients, appointments, employees, **tasks**, stats, search, billing (invoices, expenses, categorías, stats agregados).
- **Frontend**: backoffice (dashboard, clientes, empleados, citas, facturación listas + altas, categorías); portal cliente (citas); marketing estático en raíz.
- **Brechas inferidas**: (1) **Tasks**: API existe; **no** hay páginas ni uso de `api.tasks` en `.tsx`. (2) **E2E**: `openspec/config.yaml` marca E2E *No*. (3) **Gasto detalle**: lista navega a `/billing/expenses/:id` pero **no** hay `expenses/[id]/page.tsx` en el árbol — riesgo 404. (4) Hardening: CORS/listas ya saneadas en commits recientes; mantener pruebas de regresión al tocar listados.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `openspec/changes/mvp-roadmap-2026/` | New | Esta propuesta |
| Futuro `frontend/…/tasks` | Planned | UI tareas si se aprueba fase 1 |
| Futuro `…/billing/expenses/[id]` | Planned | Detalle gasto si se aprueba fase 1 |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Roadmap desalineado con negocio sin leer MVP doc | Med | Reconciliar con `MVP_ROADMAP.md` en sesión corta stakeholder. |
| Detalle gasto roto en prod | Med | Priorizar ruta o quitar navegación hasta exista página. |

## Rollback Plan

Eliminar o archivar `openspec/changes/mvp-roadmap-2026/`; ningún runtime afectado.

## Dependencies

- Acceso a priorización negocio para ordenar fases 2+.

## Success Criteria

- [ ] Stakeholders validan fases/fechas o las ajustan en ≤1 semana.
- [ ] Cada fase posterior tiene change + spec cuando toque implementación.

---

## Evaluación MVP (código)

**Listo (alto nivel):** auth; CRM clientes; agenda (citas, slots, confirmación); empleados; facturación (facturas/gastos/categorías + PDF); búsqueda global; stats dashboard backoffice; despliegue Docker prod; observabilidad base (commits previos).

**Falta / débil:** UI **tareas**; **E2E**; **detalle gasto** (ruta ausente); tests integración front limitados; posibles más polish UX (no auditado en profundidad).

---

## Roadmap (fases, fechas 2026, esfuerzo orientativo)

*Ancla: semana del **21 abr** (ajustar al sprint real del equipo).*

| Fase | Qué | Inicio | Fin | Esfuerzo |
|------|-----|--------|-----|----------|
| **0** | Cierre brecha **gasto detalle** o deshabilitar navegación | 21 abr | 2 may | ~1 sem |
| **1** | **Tasks**: listado + detalle mínimo + sidebar | 5 may | 30 may | ~4 sem |
| **2** | **E2E** críticos (login, factura, cita) + CI | 2 jun | 27 jun | ~4 sem |
| **3** | Hardening: rate limits sensibles, auditoría permisos, métricas | 30 jun | 25 jul | ~4 sem |
| **4** | Buffer mejora UX / deuda menor post–MVP | 28 jul | 22 ago | ~4 sem |

**Hitos:** fin F0 = sin 404 en flujo gastos; fin F1 = tasks usables en backoffice; fin F2 = pipeline verde con smoke E2E; fin F3 = checklist seguridad/ops cerrado.
