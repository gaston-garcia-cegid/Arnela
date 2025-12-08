# Bug Fix: Usuario no se reactiva al reactivar cliente eliminado

## 📋 Resumen del Bug

**Síntoma:** Al crear un cliente que ya fue eliminado (soft deleted), el sistema reactivaba el cliente correctamente pero el usuario asociado (tabla `users`) permanecía con `is_active = false`.

**Impacto:** El usuario no podía hacer login después de que el cliente fuera reactivado.

---

## 🔍 Análisis Detallado

### 1. ¿Qué estaba pasando exactamente?

**Flujo con el bug:**
```
1. Sistema encuentra cliente eliminado (deleted_at IS NOT NULL) ✅
2. Reactiva el cliente (deleted_at = NULL, is_active = true) ✅
3. Actualiza datos del cliente ✅
4. Intenta obtener usuario asociado con GetByID() ❌
5. GetByID() retorna error porque WHERE is_active = true
6. Como err != nil, no entra al bloque if err == nil
7. Usuario queda is_active = false ❌
```

**Código problemático (líneas 83-90 de `client_service.go`):**
```go
// ANTES (CON BUG)
if deletedClient.UserID != uuid.Nil {
    user, err := s.userRepo.GetByID(ctx, deletedClient.UserID)  // ❌ Filtra is_active = true
    if err == nil && !user.IsActive {  // ❌ Nunca entra aquí porque err != nil
        if reactivateErr := s.reactivateUser(ctx, deletedClient.UserID); reactivateErr != nil {
            log.Printf("[WARN] Failed to reactivate user %s: %v", deletedClient.UserID, reactivateErr)
        }
    }
}
```

### 2. ¿Por qué ocurría el bug?

**Causa Raíz:** Conflicto entre el propósito del método y su uso

- **`GetByID()` está diseñado para:** Obtener usuarios activos (caso de uso normal)
  ```sql
  SELECT * FROM users WHERE id = $1 AND is_active = true
  ```

- **Pero se necesitaba:** Obtener usuarios sin importar su estado (caso de reactivación)
  ```sql
  SELECT * FROM users WHERE id = $1  -- Sin filtro de is_active
  ```

**Cadena de fallos:**
1. Usuario inactivo existe en BD
2. `GetByID()` no lo encuentra por el filtro `is_active = true`
3. Retorna error "user not found"
4. Condición `if err == nil` es falsa
5. Bloque de reactivación nunca se ejecuta
6. Usuario permanece inactivo

### 3. ¿Cómo se solucionó?

**Solución:** Crear método específico para obtener usuarios sin filtrar por estado

#### A. Nuevo método en repositorio (`user_repository.go`)

```go
// GetByIDAll retrieves a user by their ID regardless of is_active status
// Used for reactivation flows where we need to check inactive users
func (r *userRepository) GetByIDAll(ctx context.Context, id uuid.UUID) (*domain.User, error) {
    query := `
        SELECT id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at
        FROM users
        WHERE id = $1  -- ✅ Sin filtro is_active
    `
    // ... implementación
}
```

#### B. Actualizar lógica de reactivación (`client_service.go`)

```go
// DESPUÉS (CORREGIDO)
if deletedClient.UserID != uuid.Nil {
    // ✅ Usar GetByIDAll para obtener usuario sin importar estado
    user, err := s.userRepo.GetByIDAll(ctx, deletedClient.UserID)
    if err != nil {
        log.Printf("[WARN] Failed to fetch user %s for reactivation: %v", deletedClient.UserID, err)
    } else if !user.IsActive {
        // ✅ Usuario existe pero está inactivo, reactivarlo
        log.Printf("[DEBUG] Reactivating user %s (email: %s)", user.ID, user.Email)
        if reactivateErr := s.reactivateUser(ctx, deletedClient.UserID); reactivateErr != nil {
            log.Printf("[ERROR] Failed to reactivate user %s: %v", deletedClient.UserID, reactivateErr)
        } else {
            log.Printf("[DEBUG] User %s reactivated successfully", user.ID)
        }
    } else {
        log.Printf("[DEBUG] User %s is already active, no reactivation needed", user.ID)
    }
}
```

### 4. ¿Cómo prevenir bugs similares en el futuro?

#### 🛡️ Estrategia 1: Naming Conventions Claras

**Problema:** `GetByID()` no deja claro que filtra por estado

**Solución:** Nomenclatura explícita
```go
// ❌ AMBIGUO
GetByID(id)  // ¿Incluye inactivos?

// ✅ CLARO
GetByID(id)              // Solo activos (caso común)
GetByIDAll(id)           // Todos los estados
GetByIDIncludingInactive(id)  // Alternativa descriptiva
```

**Aplicar a todas las entidades:**
```go
// Clientes
GetByID(id)        // WHERE deleted_at IS NULL
GetByIDAll(id)     // Sin filtro

// Empleados
GetByID(id)        // WHERE deleted_at IS NULL
GetByIDAll(id)     // Sin filtro

// Citas
GetByID(id)        // WHERE status != 'cancelled'
GetByIDAll(id)     // Todas las citas
```

#### 🛡️ Estrategia 2: Tests Específicos para Soft Delete

**Crear suite de tests para cada entidad con soft delete:**

```go
// Para cada entidad (Client, Employee, User), crear:

func TestEntity_SoftDelete_Lifecycle(t *testing.T) {
    // 1. Crear entidad
    // 2. Soft delete
    // 3. Verificar que GetByID no la encuentra
    // 4. Verificar que GetByIDAll sí la encuentra
    // 5. Reactivar
    // 6. Verificar que GetByID la encuentra nuevamente
}

func TestEntity_Reactivation_CascadesToRelated(t *testing.T) {
    // Verificar que reactivar una entidad reactiva sus dependencias
    // Ej: Client → User, Employee → User
}
```

**Tests implementados:**
- ✅ `TestClientService_CreateClient_ReactivatesInactiveUser` - Verifica reactivación de usuario
- ✅ `TestClientService_CreateClient_UserReactivation_WithGetByIDFailure` - Verifica manejo de errores

#### 🛡️ Estrategia 3: Documentación en Código

```go
// GetByID retrieves a user by their ID (only active users)
// 
// ⚠️  IMPORTANT: This method filters by is_active = true
// For reactivation flows, use GetByIDAll() instead
// 
// Use cases:
// - User login: Use GetByID (only active users can login)
// - User profile: Use GetByID (only show active users)
// - Reactivation: Use GetByIDAll (need to check inactive users)
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)

// GetByIDAll retrieves a user by their ID regardless of is_active status
// 
// ⚠️  IMPORTANT: This method does NOT filter by is_active
// Use this ONLY for:
// - Reactivation flows
// - Admin operations
// - Audit operations
// 
// For normal operations, use GetByID() instead
func (r *userRepository) GetByIDAll(ctx context.Context, id uuid.UUID) (*domain.User, error)
```

#### 🛡️ Estrategia 4: Patrón Repository con Opciones

**Implementar patrón de opciones funcionales:**

```go
type QueryOptions struct {
    IncludeInactive bool
    IncludeDeleted  bool
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID, opts ...QueryOption) (*domain.User, error) {
    options := &QueryOptions{
        IncludeInactive: false,  // Default: solo activos
        IncludeDeleted:  false,
    }
    
    for _, opt := range opts {
        opt(options)
    }
    
    query := "SELECT * FROM users WHERE id = $1"
    
    if !options.IncludeInactive {
        query += " AND is_active = true"
    }
    
    if !options.IncludeDeleted {
        query += " AND deleted_at IS NULL"
    }
    
    // ...
}

// Uso:
user, err := repo.GetByID(ctx, id)  // Solo activos
user, err := repo.GetByID(ctx, id, WithInactive())  // Incluye inactivos
user, err := repo.GetByID(ctx, id, WithDeleted())   // Incluye eliminados
```

#### 🛡️ Estrategia 5: Linting Rules Personalizadas

**Crear regla de linting para detectar uso incorrecto:**

```go
// .golangci.yml
linters-settings:
  gocritic:
    enabled-checks:
      - commentFormatting
    settings:
      commentFormatting:
        # Detectar uso de GetByID en contexto de reactivación
        patterns:
          - pattern: 'GetByID.*Reactivate'
            message: "Use GetByIDAll instead of GetByID in reactivation flows"
```

#### 🛡️ Estrategia 6: Code Review Checklist

**Agregar a checklist de PR:**

```markdown
## Soft Delete Review Checklist

Cuando trabajas con entidades que tienen soft delete:

- [ ] ¿Estás usando GetByID en flujo de reactivación?
      → Cambiar a GetByIDAll
      
- [ ] ¿Estás reactivando una entidad con relaciones?
      → Verificar que se reactivan las entidades relacionadas
      
- [ ] ¿Agregaste tests para el flujo de reactivación?
      → Incluir test que verifique reactivación en cascada
      
- [ ] ¿El método de repositorio filtra por estado?
      → Documentar claramente qué estados incluye/excluye
```

#### 🛡️ Estrategia 7: Integration Tests

**Crear tests de integración end-to-end:**

```go
func TestIntegration_ClientReactivation_ReactivatesUser(t *testing.T) {
    // Setup real database
    db := setupTestDB(t)
    defer cleanupDB(t, db)
    
    // 1. Create client + user
    client := createTestClient(t, db)
    user := getUserByClientID(t, db, client.ID)
    assert.True(t, user.IsActive)
    
    // 2. Soft delete client
    deleteClient(t, db, client.ID)
    user = getUserByID(t, db, user.ID)
    assert.False(t, user.IsActive)  // User is deactivated
    
    // 3. Re-create client (triggers reactivation)
    reactivatedClient := createClientWithEmail(t, db, client.Email)
    
    // 4. CRITICAL: Verify user was reactivated
    user = getUserByID(t, db, user.ID)
    assert.True(t, user.IsActive)  // ✅ User must be reactivated
    assert.Equal(t, client.ID, reactivatedClient.ID)  // Same client ID
}
```

---

## 📊 Archivos Modificados

### Backend Core
1. **`internal/repository/user_repository.go`**
   - Agregado método `GetByIDAll()` a la interface

2. **`internal/repository/postgres/user_repository.go`**
   - Implementado `GetByIDAll()` sin filtro `is_active`

3. **`internal/service/client_service.go`**
   - Cambiado `GetByID()` → `GetByIDAll()` en flujo de reactivación
   - Mejorado logging para debugging

4. **`internal/repository/mocks/user_repository_mock.go`**
   - Agregado mock para `GetByIDAll()`

### Tests
5. **`internal/service/client_service_reactivation_test.go`**
   - Actualizado para usar `GetByIDAll()`

6. **`internal/service/client_service_user_reactivation_test.go`** (NUEVO)
   - Test específico para verificar reactivación de usuario
   - Test para verificar comportamiento con GetByID fallando

---

## ✅ Verificación del Fix

```bash
# Ejecutar tests de reactivación
go test ./internal/service/... -v -run "Reactivat"

# Salida esperada:
# ✅ TestClientService_CreateClient_ReactivatesInactiveUser PASS
# ✅ TestClientService_CreateClient_UserReactivation_WithGetByIDFailure PASS

# Logs del test exitoso:
[DEBUG] Found deleted client with ID: xxx, reactivating...
[DEBUG] Reactivating user xxx (email: inactive.user@example.com)  # ✅ Usuario detectado
[DEBUG] User xxx reactivated successfully                          # ✅ Usuario reactivado
[DEBUG] Client reactivated successfully: ID=xxx
```

---

## 🎯 Conclusión

**Bug identificado:** `GetByID()` filtraba usuarios inactivos, impidiendo su reactivación

**Solución implementada:** Método `GetByIDAll()` que obtiene usuarios sin filtrar por estado

**Prevención futura:**
1. ✅ Naming conventions claras (`GetByID` vs `GetByIDAll`)
2. ✅ Tests específicos para soft delete
3. ✅ Documentación explícita en código
4. ✅ Code review checklist actualizada

**Lecciones aprendidas:**
- Los métodos de repositorio deben tener nombres que reflejen sus filtros
- Los flujos de reactivación necesitan métodos sin filtros de estado
- Siempre crear tests que cubran el ciclo completo: crear → eliminar → reactivar
- Documentar casos de uso específicos para cada método
