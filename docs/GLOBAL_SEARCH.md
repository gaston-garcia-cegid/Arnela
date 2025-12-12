# Global Search - Búsqueda Global

## 📋 Overview

Sistema de búsqueda global que permite buscar rápidamente en clientes, empleados, citas y facturas desde cualquier lugar del backoffice mediante un modal accesible con `Ctrl+K`.

**Fecha de implementación:** 12 de diciembre de 2025  
**Sprint:** 2.2 - Optimistic UI & Export  
**Horas:** 12h

---

## 🎯 Objetivos Cumplidos

- ✅ Búsqueda en tiempo real con debounce (500ms)
- ✅ Backend con búsqueda paralela en 4 entidades
- ✅ Tests TDD completos (12 tests: 7 handler + 5 service)
- ✅ Navegación por teclado (↑↓ Enter Esc)
- ✅ Atajo global `Ctrl+K` / `Cmd+K`
- ✅ Resultados agrupados por tipo de entidad
- ✅ Máximo 5 resultados por categoría
- ✅ Case-insensitive search
- ✅ Integración en navbar del backoffice

---

## 🏗️ Arquitectura

### Backend (Go + PostgreSQL)

#### 1. Domain Layer (`internal/domain/search.go`)
```go
type SearchResults struct {
    Clients      []SearchClient      `json:"clients"`
    Employees    []SearchEmployee    `json:"employees"`
    Appointments []SearchAppointment `json:"appointments"`
    Invoices     []SearchInvoice     `json:"invoices"`
    TotalResults int                 `json:"totalResults"`
}

type SearchService interface {
    GlobalSearch(ctx context.Context, query string, limit int) (*SearchResults, error)
}
```

**Entidades de búsqueda:**
- `SearchClient`: ID, FirstName, LastName, Email, Phone, DNICIF
- `SearchEmployee`: ID, Name, Email, Phone, Specialties[], AvatarColor
- `SearchAppointment`: ID, Title, StartTime, EndTime, Status, ClientName, EmployeeName
- `SearchInvoice`: ID, InvoiceNumber, ClientName, TotalAmount, Status, IssueDate

#### 2. Handler Layer (`internal/handler/search_handler.go`)

**Endpoint:** `GET /api/v1/search?q=query`

**Validaciones:**
- Query parameter `q` requerido
- Mínimo 2 caracteres

**Responses:**
- `200 OK`: Resultados exitosos (incluso si es vacío)
- `400 Bad Request`: Query vacío o < 2 caracteres
- `500 Internal Server Error`: Error en servicio

**Swagger:**
```go
// @Summary Global search across entities
// @Description Search clients, employees, appointments, and invoices
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Search query (min 2 characters)"
// @Success 200 {object} domain.SearchResults
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /search [get]
// @Security BearerAuth
```

#### 3. Service Layer (`internal/service/search_service.go`)

**Características:**
- **Búsqueda paralela**: Usa goroutines para buscar en 4 entidades simultáneamente
- **Coordinación**: `sync.WaitGroup` para esperar todas las goroutines
- **Thread-safe**: `sync.Mutex` para agregar resultados de forma segura
- **Graceful degradation**: Continúa con resultados parciales si una búsqueda falla
- **Case-insensitive**: Normalización a lowercase

```go
func (s *SearchService) GlobalSearch(ctx context.Context, query string, limit int) (*SearchResults, error) {
    var wg sync.WaitGroup
    var mu sync.Mutex
    
    results := &domain.SearchResults{
        Clients:      []domain.SearchClient{},
        Employees:    []domain.SearchEmployee{},
        Appointments: []domain.SearchAppointment{},
        Invoices:     []domain.SearchInvoice{},
    }
    
    queryLower := strings.ToLower(query)
    
    // 4 goroutines paralelas
    wg.Add(4)
    
    go func() { /* SearchClients */ }()
    go func() { /* SearchEmployees */ }()
    go func() { /* SearchAppointments */ }()
    go func() { /* SearchInvoices */ }()
    
    wg.Wait()
    return results, nil
}
```

#### 4. Repository Layer (`internal/repository/postgres/search_repository.go`)

**Queries PostgreSQL:**

**Clientes:**
```sql
SELECT id, first_name, last_name, email, COALESCE(phone, ''), COALESCE(dni_cif, '')
FROM clients
WHERE deleted_at IS NULL AND is_active = true
  AND (LOWER(first_name) LIKE LOWER($1)
    OR LOWER(last_name) LIKE LOWER($1)
    OR LOWER(email) LIKE LOWER($1)
    OR LOWER(dni_cif) LIKE LOWER($1)
    OR LOWER(phone) LIKE LOWER($1))
ORDER BY 
  CASE 
    WHEN LOWER(first_name) = LOWER($2) THEN 1
    WHEN LOWER(last_name) = LOWER($2) THEN 2
    WHEN LOWER(email) = LOWER($2) THEN 3
    ELSE 4
  END,
  first_name, last_name
LIMIT $3
```

**Empleados:**
```sql
SELECT e.id, CONCAT(u.first_name, ' ', u.last_name), u.email, 
       COALESCE(e.phone, ''), COALESCE(e.specialties, '{}'), 
       COALESCE(e.avatar_color, '#000000')
FROM employees e
INNER JOIN users u ON e.user_id = u.id
WHERE e.deleted_at IS NULL AND e.is_active = true
  AND (LOWER(u.first_name) LIKE LOWER($1)
    OR LOWER(u.last_name) LIKE LOWER($1)
    OR LOWER(u.email) LIKE LOWER($1)
    OR LOWER(e.phone) LIKE LOWER($1)
    OR EXISTS (SELECT 1 FROM unnest(e.specialties) s WHERE LOWER(s) LIKE LOWER($1)))
LIMIT $2
```

**Citas:**
```sql
SELECT a.id, a.title, a.start_time, a.end_time, a.status,
       CONCAT(c.first_name, ' ', c.last_name),
       CONCAT(u.first_name, ' ', u.last_name)
FROM appointments a
INNER JOIN clients c ON a.client_id = c.id
INNER JOIN employees e ON a.employee_id = e.id
INNER JOIN users u ON e.user_id = u.id
WHERE a.deleted_at IS NULL
  AND (LOWER(a.title) LIKE LOWER($1)
    OR LOWER(c.first_name) LIKE LOWER($1)
    OR LOWER(c.last_name) LIKE LOWER($1)
    OR LOWER(u.first_name) LIKE LOWER($1)
    OR LOWER(u.last_name) LIKE LOWER($1))
ORDER BY a.start_time DESC
LIMIT $2
```

**Facturas:**
```sql
SELECT i.id, i.invoice_number, CONCAT(c.first_name, ' ', c.last_name),
       i.total_amount, i.status, i.issue_date
FROM invoices i
INNER JOIN clients c ON i.client_id = c.id
WHERE i.deleted_at IS NULL
  AND (LOWER(i.invoice_number) LIKE LOWER($1)
    OR LOWER(c.first_name) LIKE LOWER($1)
    OR LOWER(c.last_name) LIKE LOWER($1))
ORDER BY 
  CASE 
    WHEN LOWER(i.invoice_number) = LOWER($2) THEN 1
    ELSE 2
  END,
  i.issue_date DESC
LIMIT $3
```

---

### Frontend (Next.js 16 + TypeScript)

#### 1. GlobalSearch Component (`components/search/GlobalSearch.tsx`)

**Props:**
```typescript
interface GlobalSearchProps {
  isOpen: boolean;
  onClose: () => void;
}
```

**Estados:**
- `query`: String de búsqueda
- `results`: SearchResults del backend
- `isLoading`: Indicador de carga
- `selectedIndex`: Índice del item seleccionado (navegación teclado)

**Features:**
- **Debounce**: 500ms usando `useDebounce` hook
- **Auto-focus**: Input se enfoca al abrir el modal
- **Reset**: Limpia estado al cerrar
- **Overlay**: Fondo semitransparente `bg-black/50`
- **Máx altura**: `max-h-[500px]` con scroll

**Navegación por teclado:**
```typescript
const handleKeyDown = (e: React.KeyboardEvent) => {
  switch (e.key) {
    case 'ArrowDown': // Siguiente resultado
    case 'ArrowUp':   // Resultado anterior
    case 'Enter':     // Navegar al resultado seleccionado
    case 'Escape':    // Cerrar modal
  }
}
```

**Navegación a páginas:**
- Clientes → `/dashboard/backoffice/clients?id={id}`
- Empleados → `/dashboard/backoffice/team?id={id}`
- Citas → `/dashboard/backoffice/appointments?id={id}`
- Facturas → `/dashboard/backoffice/billing/invoices?id={id}`

#### 2. useKeyboardShortcut Hook (`hooks/useKeyboardShortcut.ts`)

```typescript
export function useKeyboardShortcut(
  key: string,
  callback: () => void,
  options: {
    ctrl?: boolean;
    alt?: boolean;
    shift?: boolean;
    meta?: boolean;
  } = {}
)
```

**Uso:**
```typescript
useKeyboardShortcut('k', () => setSearchOpen(true), { ctrl: true });
```

- Soporta `Ctrl+K` en Windows/Linux
- Soporta `Cmd+K` en macOS (`metaKey`)
- Previene evento por defecto
- Ignora repeticiones (`!event.repeat`)

#### 3. Integración en Layout (`app/dashboard/backoffice/layout.tsx`)

**Botón en Header:**
```tsx
<Button variant="outline" size="sm" onClick={() => setSearchOpen(true)}>
  <Search className="h-4 w-4" />
  <span className="hidden md:inline">Buscar...</span>
  <kbd className="hidden lg:inline-flex">⌘K</kbd>
</Button>
```

**Modal al final:**
```tsx
<GlobalSearch isOpen={searchOpen} onClose={() => setSearchOpen(false)} />
```

---

## 🧪 Testing

### Handler Tests (`internal/handler/search_handler_test.go`)

**7 tests:**
1. ✅ `TestGlobalSearch_Success` - Búsqueda exitosa con 4 tipos de entidades
2. ✅ `TestGlobalSearch_EmptyQuery` - Validación de query vacío (400)
3. ✅ `TestGlobalSearch_QueryTooShort` - Query < 2 caracteres (400)
4. ✅ `TestGlobalSearch_NoResults` - Búsqueda sin resultados (200, arrays vacíos)
5. ✅ `TestGlobalSearch_ServiceError` - Error en servicio (500)
6. ✅ `TestGlobalSearch_CaseInsensitive` - Query uppercase funciona
7. ✅ `TestGlobalSearch_MaxResultsPerType` - Máximo 5 por tipo

**Mock Service:**
```go
type MockSearchService struct {
    mock.Mock
}

func (m *MockSearchService) GlobalSearch(ctx context.Context, query string, limit int) (*domain.SearchResults, error) {
    args := m.Called(ctx, query, limit)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.SearchResults), args.Error(1)
}
```

### Service Tests (`internal/service/search_service_test.go`)

**5 tests:**
1. ✅ `TestSearchService_GlobalSearch_Success` - Búsqueda paralela exitosa
2. ✅ `TestSearchService_GlobalSearch_NoResults` - Sin resultados (arrays vacíos)
3. ✅ `TestSearchService_GlobalSearch_PartialError` - Degradación elegante (continúa con resultados parciales)
4. ✅ `TestSearchService_GlobalSearch_CaseInsensitive` - Normalización a lowercase
5. ✅ `TestSearchService_GlobalSearch_LimitPerType` - Verifica límite de 5

**Mock Repository:**
```go
type MockSearchRepository struct {
    mock.Mock
}

func (m *MockSearchRepository) SearchClients(ctx context.Context, query string, limit int) ([]domain.SearchClient, error)
func (m *MockSearchRepository) SearchEmployees(ctx context.Context, query string, limit int) ([]domain.SearchEmployee, error)
func (m *MockSearchRepository) SearchAppointments(ctx context.Context, query string, limit int) ([]domain.SearchAppointment, error)
func (m *MockSearchRepository) SearchInvoices(ctx context.Context, query string, limit int) ([]domain.SearchInvoice, error)
```

**Ejecución:**
```bash
# Handler tests
go test -v ./internal/handler -run TestGlobalSearch

# Service tests
go test -v ./internal/service -run TestSearchService

# Todos los tests
go test -v ./internal/handler ./internal/service
```

**Resultado esperado:**
```
=== RUN   TestGlobalSearch_Success
--- PASS: TestGlobalSearch_Success (0.00s)
=== RUN   TestGlobalSearch_EmptyQuery
--- PASS: TestGlobalSearch_EmptyQuery (0.00s)
...
PASS
ok      github.com/gaston-garcia-cegid/arnela/backend/internal/handler  0.494s

=== RUN   TestSearchService_GlobalSearch_Success
--- PASS: TestSearchService_GlobalSearch_Success (0.00s)
...
PASS
ok      github.com/gaston-garcia-cegid/arnela/backend/internal/service  0.539s
```

---

## 🚀 Uso

### Para Usuarios

1. **Abrir búsqueda:**
   - Presiona `Ctrl+K` (Windows/Linux) o `Cmd+K` (macOS)
   - O clic en botón "Buscar..." en navbar

2. **Buscar:**
   - Escribe al menos 2 caracteres
   - Resultados aparecen después de 500ms (debounce)

3. **Navegar:**
   - `↑` / `↓`: Navegar entre resultados
   - `Enter`: Ir al resultado seleccionado
   - `Esc`: Cerrar modal
   - Clic directo en cualquier resultado

4. **Resultados agrupados:**
   - **Clientes** 👤: Nombre, email, DNI/CIF o teléfono
   - **Empleados** 👥: Nombre, email, especialidades
   - **Citas** 📅: Título, cliente-empleado, fecha/hora
   - **Facturas** 📄: Número, cliente, monto, estado

### Para Desarrolladores

**Integrar en nuevo layout:**
```tsx
import { GlobalSearch } from '@/components/search/GlobalSearch';
import { useKeyboardShortcut } from '@/hooks/useKeyboardShortcut';

function MyLayout() {
  const [searchOpen, setSearchOpen] = useState(false);
  
  useKeyboardShortcut('k', () => setSearchOpen(true), { ctrl: true });
  
  return (
    <>
      {/* Your layout */}
      <GlobalSearch isOpen={searchOpen} onClose={() => setSearchOpen(false)} />
    </>
  );
}
```

**Extender búsqueda a nuevas entidades:**

1. **Backend:**
   - Agregar tipo en `domain/search.go`
   - Agregar método en `SearchService` interface
   - Implementar query en `search_repository.go`
   - Agregar goroutine en `search_service.go`
   - Actualizar tests

2. **Frontend:**
   - Agregar interface en `GlobalSearch.tsx`
   - Agregar en `SearchResults` type
   - Crear `<ResultSection>` con icono apropiado
   - Agregar navegación en `navigateToItem()`

---

## ⚡ Performance

### Backend

**Búsqueda Paralela:**
- 4 queries SQL ejecutan simultáneamente
- Tiempo total = MAX(query1, query2, query3, query4)
- Ejemplo: Si cada query toma 50ms, total es ~50ms (no 200ms)

**Límites:**
- Máximo 5 resultados por tipo (20 total)
- Reduce transferencia de datos
- Mantiene UI manejable

**Índices recomendados:**
```sql
-- Clientes
CREATE INDEX idx_clients_search ON clients(LOWER(first_name), LOWER(last_name), LOWER(email));

-- Empleados
CREATE INDEX idx_employees_search ON employees USING gin(specialties);

-- Citas
CREATE INDEX idx_appointments_search ON appointments(LOWER(title), start_time);

-- Facturas
CREATE INDEX idx_invoices_search ON invoices(LOWER(invoice_number), issue_date);
```

### Frontend

**Debounce:**
- 500ms delay antes de hacer request
- Evita requests innecesarios mientras el usuario escribe
- Ejemplo: "juan" → espera 500ms → hace 1 request (no 4)

**Optimizaciones:**
- Componente `ResultSection` reutilizable
- `useCallback` para funciones estables
- `useEffect` con deps correctas
- Cleanup de listeners en unmount

---

## 📊 Métricas

### Cobertura de Tests
- **Handler:** 7/7 tests ✅ (100%)
- **Service:** 5/5 tests ✅ (100%)
- **Total:** 12 tests, 100% pasan

### Búsqueda Cobertura
- ✅ Clientes: first_name, last_name, email, dni_cif, phone
- ✅ Empleados: name (first + last), email, phone, specialties
- ✅ Citas: title, client_name, employee_name
- ✅ Facturas: invoice_number, client_name

### UX
- **Debounce:** 500ms (balance entre rapidez y performance)
- **Min chars:** 2 (evita resultados excesivos)
- **Max results:** 5 por tipo (UI manejable)
- **Keyboard shortcuts:** ✅ Ctrl+K, ↑↓, Enter, Esc

---

## 🔐 Seguridad

**Autenticación:**
- Requiere `BearerAuth` (JWT token)
- Middleware `RequireAuth()` en ruta

**Autorización:**
- Solo usuarios autenticados (admin/employee)
- Clients pueden buscar en su contexto (futuro)

**SQL Injection:**
- ✅ Queries parametrizadas (`$1`, `$2`)
- ✅ No concatenación de strings
- ✅ `LOWER()` aplica a parámetros binding

**Rate Limiting:**
- Debounce frontend (500ms)
- TODO: Rate limit backend (futuro)

---

## 🐛 Troubleshooting

### "Query parameter 'q' is required"
**Causa:** Request sin query parameter  
**Solución:** Asegurar `?q=valor` en URL

### "Query must be at least 2 characters"
**Causa:** Query muy corto  
**Solución:** Escribir ≥2 caracteres

### "No se encontraron resultados"
**Causa:** No hay coincidencias en BD  
**Solución:** Normal, probar con otro término

### Búsqueda lenta
**Causa:** Sin índices en BD  
**Solución:** Aplicar índices recomendados (ver Performance)

### Ctrl+K no funciona
**Causa:** Focus en input u otro elemento  
**Solución:** Click fuera del input primero

---

## 🔄 Roadmap Futuro

### Sprint 2.3
- [ ] Búsqueda en descripciones de citas
- [ ] Búsqueda por rangos de fechas
- [ ] Historial de búsquedas recientes

### Sprint 2.4
- [ ] Highlighting de términos coincidentes
- [ ] Filtros avanzados (fecha, estado, monto)
- [ ] Export de resultados a CSV

### Sprint 2.5
- [ ] Búsqueda por voz (Web Speech API)
- [ ] Sugerencias automáticas (autocomplete)
- [ ] Rate limiting backend

---

## 📚 Referencias

**Backend:**
- `internal/domain/search.go` - Domain types
- `internal/handler/search_handler.go` - HTTP handler
- `internal/service/search_service.go` - Business logic
- `internal/repository/postgres/search_repository.go` - Database queries

**Frontend:**
- `components/search/GlobalSearch.tsx` - Main component
- `hooks/useKeyboardShortcut.ts` - Keyboard hook
- `hooks/useDebounce.ts` - Debounce hook
- `app/dashboard/backoffice/layout.tsx` - Integration

**Tests:**
- `internal/handler/search_handler_test.go` - 7 handler tests
- `internal/service/search_service_test.go` - 5 service tests

**Docs:**
- `OPTIMISTIC_UI_UPDATES.md` - Feature previa relacionada
- `MVP_ROADMAP.md` - Roadmap general del proyecto

---

**✨ Feature completada el 12 de diciembre de 2025**
