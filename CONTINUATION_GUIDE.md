# 🎯 GUÍA COMPLETA DE CONTINUACIÓN - Proyecto Arnela

## 📅 Actualizado: 6 de diciembre de 2025

---

## ✅ LO QUE SE HA COMPLETADO

### **1. Infraestructura de Optimización** ✅ 100%

#### **A. Logger Centralizado**
**Archivo**: `frontend/src/lib/logger.ts`
- ✅ Métodos: info, warn, error, debug
- ✅ Solo logs en desarrollo
- ✅ Contexto estructurado
- ✅ Preparado para Sentry

#### **B. Validadores**
**Archivo**: `frontend/src/lib/validators.ts`
- ✅ 15 validadores implementados
- ✅ Validaciones españolas (DNI, CIF, teléfono, CP)
- ✅ Mensajes consistentes

#### **C. Error Handler Hook**
**Archivo**: `frontend/src/hooks/useErrorHandler.ts`
- ✅ Logging automático
- ✅ Toasts integrados
- ✅ Estado de error

#### **D. Sistema de Toasts**
- ✅ Sonner instalado con pnpm
- ✅ Configurado en layout
- ✅ Componente UI personalizado
- ✅ Guía completa: `TOASTS_GUIDE.md`

---

### **2. Archivos Migrados** ✅ 28% (7/25)

1. ✅ `hooks/useStats.ts`
2. ✅ `components/backoffice/CreateClientModal.tsx`
3. ✅ `components/backoffice/EditClientModal.tsx`
4. ✅ `components/backoffice/CreateEmployeeModal.tsx`
5. ✅ `components/backoffice/EditEmployeeModal.tsx`
6. ✅ `app/dashboard/backoffice/page.tsx`
7. ✅ `app/dashboard/backoffice/appointments/page.tsx`

---

## 🔄 LO QUE QUEDA POR HACER

### **PASO 1: Migrar 18 Archivos Restantes** (72% pendiente)

#### **Script para Migrar cada Archivo**:

```typescript
// 1. Añadir import
import { logError } from '@/lib/logger';

// 2. Buscar: console.error('...', err);
// 3. Reemplazar con: logError('descripción', err, { component: 'ComponentName' });
```

#### **Archivos Pendientes** (en orden de prioridad):

**Alta Prioridad** (Páginas Críticas):
1. ⏳ `app/dashboard/backoffice/employees/page.tsx`
   - console.error en línea ~64 y ~110
   
2. ⏳ `app/dashboard/backoffice/employees/[id]/page.tsx`
   - console.error en línea ~101, ~135, ~184
   
3. ⏳ `app/dashboard/backoffice/clients/page.tsx`
   - console.error en línea ~83, ~160

**Media Prioridad** (Facturación):
4. ⏳ `app/dashboard/backoffice/billing/page.tsx`
   - console.error en línea ~27

5. ⏳ `app/dashboard/backoffice/billing/invoices/page.tsx`
   - console.error en línea ~57, ~80

6. ⏳ `app/dashboard/backoffice/billing/expenses/page.tsx`
   - console.error en línea ~60, ~71

7. ⏳ `app/dashboard/backoffice/billing/expenses/new/page.tsx`
   - console.error en línea ~86
   - alert() en línea ~74, ~87

8. ⏳ `app/dashboard/backoffice/billing/invoices/new/page.tsx`
   - console.error en línea ~47
   - alert() en línea ~34, ~48

9. ⏳ `app/dashboard/backoffice/billing/categories/page.tsx`
   - console.error en línea ~51, ~66, ~78
   - alert() en línea ~67, ~79

---

### **PASO 2: Reemplazar alert() con toast**

**Archivos con alert()** (6 ocurrencias en 3 archivos):

1. `billing/expenses/new/page.tsx`:
```typescript
// ❌ Antes
alert("No estás autenticado");
alert("Error al crear el gasto");

// ✅ Después
import { toast } from 'sonner';
toast.error("No estás autenticado", { description: "Por favor, inicia sesión nuevamente" });
toast.error("Error al crear el gasto", { description: err.message });
```

2. `billing/invoices/new/page.tsx`:
```typescript
// ❌ Antes
alert("No estás autenticado");
alert("Error al crear la factura");

// ✅ Después
toast.error("No estás autenticado");
toast.error("Error al crear la factura", { description: err.message });
```

3. `billing/categories/page.tsx`:
```typescript
// ❌ Antes
alert("Error al crear la categoría");
alert("Error al eliminar la categoría");

// ✅ Después
toast.error("Error al crear la categoría");
toast.error("Error al eliminar la categoría");
```

---

### **PASO 3: Agregar Toasts de Confirmación en CRUD**

#### **Patrón para Success Toasts**:

**Create Operations**:
```typescript
// components/backoffice/CreateClientModal.tsx
try {
  const newClient = await api.clients.create(data, token);
  toast.success('Cliente creado', {
    description: `${newClient.firstName} ${newClient.lastName} agregado al sistema`
  });
  onSuccess(newClient);
} catch (err) {
  logError('Error creating client', err, { component: 'CreateClientModal' });
  toast.error('Error al crear cliente', { description: err.message });
}
```

**Update Operations**:
```typescript
// components/backoffice/EditClientModal.tsx
try {
  const updated = await api.clients.update(id, data, token);
  toast.success('Cliente actualizado', {
    description: 'Los cambios han sido guardados'
  });
  onSuccess(updated);
} catch (err) {
  logError('Error updating client', err, { component: 'EditClientModal' });
  toast.error('Error al actualizar', { description: err.message });
}
```

**Delete Operations**:
```typescript
// pages con delete
try {
  await api.clients.delete(id, token);
  toast.success('Cliente eliminado', {
    description: 'El cliente ha sido removido del sistema'
  });
  refreshList();
} catch (err) {
  logError('Error deleting client', err, { component: 'ClientsPage' });
  toast.error('Error al eliminar', { description: err.message });
}
```

#### **Archivos donde Agregar Toasts**:

1. ✅ `components/backoffice/CreateClientModal.tsx` - Ya tiene logError (agregar toast.success)
2. ✅ `components/backoffice/EditClientModal.tsx` - Ya tiene logError (agregar toast.success)
3. ✅ `components/backoffice/CreateEmployeeModal.tsx` - Ya tiene logError (agregar toast.success)
4. ✅ `components/backoffice/EditEmployeeModal.tsx` - Ya tiene logError (agregar toast.success)
5. ⏳ `app/dashboard/backoffice/employees/page.tsx` - Agregar toasts en create/delete
6. ⏳ `app/dashboard/backoffice/clients/page.tsx` - Agregar toasts en create/delete
7. ⏳ `components/appointments/*` - Agregar toasts en confirm/cancel

---

### **PASO 4: Tests para Validadores**

**Archivo a crear**: `frontend/src/lib/__tests__/validators.test.ts`

```typescript
import { describe, it, expect } from 'vitest'; // o jest
import {
  validateDNI,
  validateCIF,
  validateEmail,
  validatePhone,
  validatePostalCode,
} from '../validators';

describe('validators', () => {
  describe('validateDNI', () => {
    it('should validate correct DNI', () => {
      const result = validateDNI('12345678Z');
      expect(result.isValid).toBe(true);
    });

    it('should reject DNI with wrong letter', () => {
      const result = validateDNI('12345678A');
      expect(result.isValid).toBe(false);
      expect(result.error).toContain('letra');
    });

    it('should reject invalid format', () => {
      const result = validateDNI('1234567');
      expect(result.isValid).toBe(false);
    });

    it('should validate NIE', () => {
      const result = validateDNI('X1234567L');
      expect(result.isValid).toBe(true);
    });
  });

  describe('validateEmail', () => {
    it('should validate correct email', () => {
      const result = validateEmail('test@example.com');
      expect(result.isValid).toBe(true);
    });

    it('should reject invalid email', () => {
      const result = validateEmail('invalid-email');
      expect(result.isValid).toBe(false);
    });
  });

  describe('validatePhone', () => {
    it('should validate Spanish mobile', () => {
      const result = validatePhone('612345678');
      expect(result.isValid).toBe(true);
    });

    it('should validate Spanish landline', () => {
      const result = validatePhone('986123456');
      expect(result.isValid).toBe(true);
    });

    it('should reject invalid phone', () => {
      const result = validatePhone('123');
      expect(result.isValid).toBe(false);
    });
  });

  describe('validatePostalCode', () => {
    it('should validate correct postal code', () => {
      const result = validatePostalCode('28001');
      expect(result.isValid).toBe(true);
    });

    it('should reject invalid postal code', () => {
      const result = validatePostalCode('99999');
      expect(result.isValid).toBe(false);
    });
  });

  describe('validateCIF', () => {
    it('should validate correct CIF', () => {
      const result = validateCIF('A12345678');
      expect(result.isValid).toBe(true);
    });

    it('should reject invalid CIF', () => {
      const result = validateCIF('123456789');
      expect(result.isValid).toBe(false);
    });
  });
});
```

**Configurar Testing** (si no está configurado):

```bash
# Instalar dependencias
pnpm add -D vitest @testing-library/react @testing-library/jest-dom

# Crear vitest.config.ts
```

```typescript
// vitest.config.ts
import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react';
import path from 'path';

export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'jsdom',
    globals: true,
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
});
```

```json
// package.json
{
  "scripts": {
    "test": "vitest",
    "test:ui": "vitest --ui",
    "test:coverage": "vitest --coverage"
  }
}
```

**Ejecutar tests**:
```bash
pnpm test
```

---

## 📊 RESUMEN DE TAREAS

| Tarea | Archivos | Estado |
|-------|----------|--------|
| **1. Migrar console.error** | 18 | ⏳ 0/18 |
| **2. Reemplazar alert()** | 3 | ⏳ 0/3 |
| **3. Agregar toasts CRUD Success** | 7 | ⏳ 0/7 |
| **4. Tests validadores** | 1 | ⏳ 0/1 |

**Total**: 29 tareas pendientes

---

## 🚀 ORDEN RECOMENDADO DE EJECUCIÓN

### **Sesión 1** (1-2 horas):
1. Migrar 3 archivos críticos (employees, employees/[id], clients)
2. Reemplazar todos los alert()
3. Agregar toasts de success en modales ya migrados

### **Sesión 2** (1-2 horas):
4. Migrar archivos de facturación (6 archivos)
5. Agregar toasts en páginas de delete
6. Tests para validadores básicos

### **Sesión 3** (1 hora):
7. Migrar archivos restantes
8. Tests completos para todos los validadores
9. Verificación final

---

## 💡 TIPS PARA MIGRAR RÁPIDO

### **1. Usar Find & Replace en VSCode**:
```
Find: console\.error\('([^']+)', err\);
Replace: logError('$1', err, { component: 'COMPONENT_NAME' });
```

### **2. Añadir import automáticamente**:
```typescript
// Copiar esta línea al inicio de cada archivo
import { logError } from '@/lib/logger';
```

### **3. Template para toasts de success**:
```typescript
toast.success('ACCIÓN realizada', {
  description: 'DESCRIPCIÓN_ESPECÍFICA'
});
```

---

## 📚 DOCUMENTACIÓN DISPONIBLE

1. ✅ `OPTIMIZATION_PLAN.md` - Plan maestro
2. ✅ `OPTIMIZATION_GUIDE.md` - Guía de uso
3. ✅ `TOASTS_GUIDE.md` - Guía de Sonner
4. ✅ `OPTIMIZATION_FINAL.md` - Resumen final
5. ✅ `MIGRATION_STATUS.md` - Estado de migración
6. ✅ `CONTINUATION_GUIDE.md` - **ESTE DOCUMENTO**

---

## 🎯 CHECKLIST FINAL

- [ ] 18 archivos migrados con logError
- [ ] 6 alert() reemplazados con toast.error
- [ ] 7 operaciones CRUD con toast.success
- [ ] Tests para 5+ validadores
- [ ] Coverage de tests > 80%
- [ ] 0 console.error/log en producción
- [ ] Todas las operaciones con feedback visual

---

**Una vez completado**:
- ✅ Código production-ready
- ✅ UX mejorada con toasts
- ✅ Seguridad (no logs en producción)
- ✅ Mantenibilidad (validadores centralizados)
- ✅ Testing (validadores probados)

---

**Última actualización**: 6 de diciembre de 2025, 14:20 UTC
**Progreso actual**: 28% completado (7/25 archivos migrados)
