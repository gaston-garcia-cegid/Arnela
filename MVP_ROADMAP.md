# MVP y roadmap — Arnela CRM/CMS

> **Última actualización:** abril de 2026  
> **Propósito:** estado consolidado del MVP ya entregado y **fases siguientes** priorizadas para evolución del producto.

---

## 1. Resumen ejecutivo

Arnela es un **MVP operativo** para uso interno del gabinete: CRM/backoffice (clientes, empleados, citas, tareas), **facturación completa** (facturas, gastos, categorías, estadísticas), sitio público parcialmente contenido, **Docker** (dev/prod), **CI en GitHub Actions**, documentación en raíz y `docs/`.

Lo que queda no es “terminar el MVP core”, sino **pulir producto público**, **consistencia fiscal/datos**, **integraciones reales**, **observabilidad**, **tests E2E** y **funcionalidades de siguiente nivel** según prioridad del negocio.

---

## 2. Lo ya completado (referencia — no es backlog)

Estas líneas de trabajo **se consideran cerradas** en el código y documentación actual del repositorio:

| Bloque | Contenido principal |
|--------|---------------------|
| **Estabilización** | Auth coherente, CORS/env, rate limit en login/registro, recovery/logging Gin |
| **Core funcional** | Contacto / Sobre Arnela / Hero; facturación desde cita; cache en stats |
| **Operaciones y calidad** | CI (Go + Next), SEO en rutas públicas clave, tests billing backend, email SMTP en cola, health/readiness, healthchecks Docker |
| **Valor añadido** | PDF factura, dark mode, Google Calendar (base, opcional), fixes UI facturación (modal detalle) |

*Detalle técnico de API, env vars y estructura:* ver `README.md`, `backend/README.md`, `frontend/README.md`, `docs/DEPLOYMENT.md`.

---

## 3. Estado actual honesto (brechas conocidas)

| Área | Situación |
|------|-----------|
| **Sitio público** | **Formación**, **Intervención** y **Convenios** publicadas con textos e imágenes locales bajo `frontend/public/images/` alineados a [arnelagabinete.com](https://www.arnelagabinete.com/); revisar legal/cookies si se añaden widgets externos |
| **IVA / dominio factura** | **Unificado (Fase 6):** `vatRate` en puntos porcentuales (21 = 21%), `CalculateAmounts` coherente, migración `000016` corrige filas antiguas; tests en `domain` y preview en frontend (`invoiceFiscal`) |
| **Cola de trabajos** | TODOs en worker para SMS/WhatsApp y sync completo de Calendar según evolución del código |
| **Frontend observabilidad** | `logger` con TODO hacia servicio de monitoring (p. ej. Sentry) |
| **Página dedicada** `invoices/[id]` | No es obligatoria si el flujo es modal/listado; opcional para edición “pantalla completa” |

---

## 4. Roadmap por fases (por dónde avanzar)

### Fase 5 — Producto público y marca

**Objetivo:** alinear la web con la oferta real del gabinete.

- ✅ **Formación**, **Intervención** y **Convenios y colaboraciones**: contenido sustituido (textos oficiales + imágenes en `public/images/`).
- ✅ **Tests:** Vitest sobre `arnelaSiteAssets` (rutas locales, sin URLs remotas).
- Pendiente opcional: revisar CTAs, enlaces internos, metadatos SEO entre landing y esas páginas; valorar copiar logos EU/PRTR como assets locales si hace falta.

**Criterio de salida:** ninguna ruta enlazada desde el menú principal queda en *Under construction* — **cumplido** para las tres rutas anteriores.

---

### Fase 6 — Fiscalidad y consistencia de datos

**Objetivo:** facturación **auditables** (contabilidad / inspección).

- ✅ Unificar modelo de **IVA** (dominio `DefaultVATPercent` = 21, API JSON, UI nueva factura + detalle, PDF, migración `000016_fix_invoice_vat_percent` para datos legados; Vitest `invoiceFiscal` para el preview del formulario).
- Pendiente: revisar redondeos en casos límite, fechas de cobro, método de pago y coherencia listado ↔ detalle ↔ export en escenarios reales.

**Criterio de salida:** tests que documenten el contrato fiscal; flujo crear factura → PDF → export sin discrepancias en montos — **contrato IVA cubierto** con tests Go (`domain`) + Vitest (`invoiceFiscal`); integración `pkg/pdf` (PDF sin comprimir en test + misma terna base/IVA/total que CSV) y `invoiceFiscalExportSegment` alineado con export; seguir validando exportes con datos de producción tras migración.

---

### Fase 7 — Integraciones y canal al cliente

**Objetivo:** integraciones **verificables** en entorno de staging, no solo variables de entorno.

- Completar **Google Calendar** (crear / actualizar / borrar alineado con el ciclo de vida de la cita, según diseño acordado).
- **SMS / WhatsApp** o retirar de la cola y de la documentación si no hay roadmap claro.
- **Recordatorios** de cita (p. ej. 24h antes) reutilizando email y/o canal móvil si aplica.

**Criterio de salida:** checklist de pruebas con credenciales reales + evidencia (logs o tests de contrato donde sea viable).

---

### Fase 8 — Observabilidad, backups y endurecimiento

**Objetivo:** operación tranquila en servidor.

- Monitoring (p. ej. **Sentry** en frontend, correlación con errores API si aplica).
- Runbooks: backup/restore Postgres y Redis, rotación de secretos (extender `docs/DEPLOYMENT.md` si hace falta).
- Revisión **seguridad**: headers Nginx, permisos por rol en endpoints sensibles, límites de tasa donde tenga sentido.

**Criterio de salida:** simulacro de incidente recuperable siguiendo documentación interna.

---

### Fase 9 — Calidad automatizada “de usuario”

**Objetivo:** que refactors no rompan flujos críticos ya usados a mano.

- **E2E** (p. ej. Playwright): login → flujo mínimo de negocio (cliente / cita / factura / cobro o equivalente acordado).
- Opcional: validación de **contrato** OpenAPI frente a respuestas reales.

**Criterio de salida:** CI (o job nocturno) ejecuta al menos 1–3 flujos E2E verdes.

---

### Fase 10 — Producto “siguiente nivel gabinete”

**Objetivo:** solo cuando el día a día del staff esté estable. Elegir **1–2** hilos, no todos a la vez.

Ejemplos de hilos:

- **Portal cliente** ampliado (si encaja legal y de negocio).
- **Informes de sesión / subsidios / comunicación** si forman parte de la visión del backoffice tipo CMS.
- **Informes de negocio** (ocupación, previsión de cobro, por terapeuta).

**Criterio de salida:** feature piloto en producción con usuario interno y feedback recogido.

---

## 5. Orden recomendado

1. **Fase 5** si la web pública es compromiso inmediato con el gabinete.  
2. **Fase 6** en paralelo o justo después si ya hay facturación con datos reales.  
3. **Fase 7** cuando calendario y recordatorios dejen de ser “nice to have”.  
4. **Fase 8** antes o en paralelo al crecimiento de usuarios reales.  
5. **Fase 9** cuando suba el ritmo de cambios en UI/API.  
6. **Fase 10** cuando el MVP deje de generar incidencias habituales.

---

## 6. Nota sobre documentación histórica

Versiones anteriores de este archivo contenían sprints y porcentajes desalineados con el repo actual. **Esta versión es la fuente de verdad** del roadmap a partir de la fecha indicada en la cabecera. El detalle de implementación vive en `README.md`, `CONTRIBUTING.md` y `docs/`.
