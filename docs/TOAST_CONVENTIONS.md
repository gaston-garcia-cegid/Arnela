# 🎯 Convenciones de Toast Notifications

> **Documentado:** Diciembre 8, 2025  
> **Propósito:** Prevenir toasts duplicados y mantener consistencia en la UX

---

## 🐛 Problema Resuelto

**Bug identificado:** Toasts duplicados en operaciones CRUD (crear/actualizar clientes, empleados, etc.)

**Causa raíz:** Toast se mostraba tanto en el **modal** como en el **callback del padre**, resultando en notificaciones duplicadas al usuario.

---

## ✅ Convención Establecida

### **Regla Principal:**

> **Los toasts de éxito se muestran ÚNICAMENTE en el componente que ejecuta la operación (modal/hook), NUNCA en el callback del padre.**

---

## 📐 Patrones Correctos

### **Patrón 1: Modal con CRUD**

```typescript
// ✅ CORRECTO: Modal muestra toast
// components/backoffice/CreateClientModal.tsx

const CreateClientModal = ({ onSuccess }) => {
  const handleSubmit = async (data) => {
    try {
      const newClient = await api.clients.create(data, token);
      
      // ✅ Toast se muestra AQUÍ
      toast.success('Cliente creado', {
        description: `${newClient.firstName} ${newClient.lastName} agregado al sistema`
      });
      
      reset();
      onOpenChange(false);
      onSuccess(newClient); // Solo notifica al padre, SIN toast
      
    } catch (err) {
      logError('Error creating client', err, { component: 'CreateClientModal' });
      toast.error('Error al crear cliente', {
        description: err.message
      });
    }
  };
};
```

```typescript
// ✅ CORRECTO: Página padre NO muestra toast
// app/dashboard/backoffice/clients/page.tsx

const ClientsPage = () => {
  const handleCreateSuccess = (client: Client) => {
    loadClients(); // ✅ Solo actualiza estado
    setIsCreateModalOpen(false);
    // ❌ NO AGREGAR toast.success() AQUÍ
    // El toast ya fue mostrado por CreateClientModal
  };

  return (
    <CreateClientModal 
      open={isCreateModalOpen}
      onSuccess={handleCreateSuccess}
    />
  );
};
```

---

### **Patrón 2: Hook personalizado con operación async**

```typescript
// ✅ CORRECTO: Hook muestra toast
// hooks/useAppointments.ts

export function useAppointments() {
  const confirmAppointment = useCallback(async (id: string, notes?: string) => {
    try {
      await api.appointments.confirm(id, { notes }, token);
      
      // ✅ Toast se muestra AQUÍ
      toast.success('Cita confirmada', {
        description: 'La cita ha sido confirmada correctamente'
      });
      
      return true;
    } catch (err) {
      toast.error('Error al confirmar cita');
      return false;
    }
  }, [token]);

  return { confirmAppointment };
}
```

```typescript
// ✅ CORRECTO: Componente NO muestra toast
// components/appointments/AppointmentList.tsx

const AppointmentList = () => {
  const { confirmAppointment } = useAppointments();

  const handleConfirm = async (id: string) => {
    const success = await confirmAppointment(id);
    if (success) {
      refreshList(); // ✅ Solo actualiza estado
      // ❌ NO AGREGAR toast.success() AQUÍ
      // El toast ya fue mostrado por el hook
    }
  };
};
```

---

## 🎨 Tipos de Toasts y Colores

El sistema soporta 4 tipos de toasts con colores personalizados que respetan la paleta cálida:

### **1. Success (Verde) - `toast.success()`**
**Uso:** Operaciones completadas exitosamente
**Color:** Verde suave y cálido (`#f0fdf4` background, `#86efac` border)

```typescript
toast.success('Cliente creado', {
  description: 'Juan Pérez ha sido agregado al sistema'
});
```

**Ejemplos:**
- Cliente/empleado creado o actualizado
- Cita confirmada o cancelada
- Factura marcada como pagada
- Operación guardada correctamente

---

### **2. Error (Rojo) - `toast.error()`**
**Uso:** Errores, fallos en operaciones
**Color:** Rojo cálido (`#fef2f2` background, `#fca5a5` border)

```typescript
toast.error('Error al crear cliente', {
  description: 'El email ya está registrado'
});
```

**Ejemplos:**
- Validación fallida
- Error de conexión
- Operación no permitida
- Credenciales incorrectas

---

### **3. Warning (Amarillo) - `toast.warning()`**
**Uso:** Advertencias, acciones que requieren atención
**Color:** Amarillo cálido (`#fffbeb` background, `#fcd34d` border)

```typescript
toast.warning('Cambios sin guardar', {
  description: 'Tienes cambios pendientes que se perderán',
  action: {
    label: 'Guardar',
    onClick: () => handleSave()
  }
});
```

**Ejemplos:**
- Cambios no guardados
- Sesión próxima a expirar
- Límite de recursos alcanzado
- Acción que requiere confirmación

---

### **4. Info (Beige) - `toast()` o `toast.info()`**
**Uso:** Información general, notificaciones neutras
**Color:** Beige cálido - paleta actual (`#fdfaf7` background, `#e8ddd0` border)

```typescript
toast('Sincronizando datos', {
  description: 'La información se está actualizando'
});

// o explícitamente
toast.info('Nueva funcionalidad disponible', {
  description: 'Ahora puedes exportar facturas a PDF'
});
```

**Ejemplos:**
- Sincronización en progreso
- Nueva funcionalidad disponible
- Información del sistema
- Recordatorios generales

---

## ❌ Anti-Patrones (NO HACER)

### **Anti-Patrón 1: Toast duplicado**

```typescript
// ❌ INCORRECTO: Toast en modal Y en callback

// Modal
const newClient = await api.clients.create(data, token);
toast.success('Cliente creado'); // 🔴 PRIMER toast
onSuccess(newClient);

// Página padre
const handleCreateSuccess = (client: Client) => {
  loadClients();
  toast.success('Cliente creado exitosamente'); // 🔴 SEGUNDO toast (DUPLICADO)
};

// Resultado: Usuario ve DOS toasts idénticos ❌
```

---

### **Anti-Patrón 2: Toast solo en el padre**

```typescript
// ❌ INCORRECTO: No mostrar toast en el modal

// Modal
const newClient = await api.clients.create(data, token);
// ❌ Falta feedback inmediato al usuario
onSuccess(newClient);

// Página padre
const handleCreateSuccess = (client: Client) => {
  toast.success('Cliente creado'); // ⚠️ Feedback retrasado
  loadClients();
};

// Problema: El usuario no recibe feedback inmediato después de submit
```

---

## 🎯 Excepciones a la Regla

### **Operaciones de eliminación (Delete)**

Las operaciones de **eliminación** suelen ejecutarse directamente en la página padre (no en modal), por lo tanto el toast SÍ debe mostrarse allí:

```typescript
// ✅ CORRECTO: Delete en página padre
const handleDeleteConfirm = async () => {
  try {
    await api.clients.delete(clientId, token);
    
    // ✅ Toast se muestra AQUÍ porque la operación se ejecuta aquí
    toast.success('Cliente eliminado', {
      description: `${clientName} ha sido eliminado del sistema`
    });
    
    setClients(clients.filter((c) => c.id !== clientId));
    setClientToDelete(null);
    
  } catch (err) {
    logError('Error deleting client', err, { component: 'ClientsPage' });
    toast.error('Error al eliminar cliente');
  }
};
```

---

## 🔍 Checklist para Code Review

Al revisar PRs con operaciones CRUD, verificar:

- [ ] **¿El toast de éxito se muestra solo UNA vez?**
- [ ] **¿El toast está en el componente que EJECUTA la operación?**
  - Si es modal → toast en modal
  - Si es hook → toast en hook
  - Si es página → toast en página
- [ ] **¿El callback `onSuccess` NO muestra toasts duplicados?**
- [ ] **¿Los errores se manejan con `logError` + `toast.error`?**
- [ ] **¿Los mensajes de toast son claros y descriptivos?**

---

## 🧪 Testing

### **Test para prevenir duplicación**

```typescript
// __tests__/CreateClientModal.test.tsx
import { render, waitFor } from '@testing-library/react';
import { vi } from 'vitest';
import { toast } from 'sonner';

it('should show toast only once on success', async () => {
  const mockOnSuccess = vi.fn();
  const mockToast = vi.spyOn(toast, 'success');
  
  render(<CreateClientModal onSuccess={mockOnSuccess} />);
  
  // Simular submit
  await userEvent.type(screen.getByLabelText('Nombre'), 'Juan');
  await userEvent.click(screen.getByText('Crear'));
  
  await waitFor(() => {
    expect(mockToast).toHaveBeenCalledTimes(1); // ✅ Solo UNA vez
    expect(mockOnSuccess).toHaveBeenCalledTimes(1);
  });
});
```

---

## 📚 Archivos de Referencia

### **Implementaciones correctas:**

1. ✅ `frontend/src/components/backoffice/CreateClientModal.tsx` (línea 95)
2. ✅ `frontend/src/components/backoffice/EditClientModal.tsx` (línea 141)
3. ✅ `frontend/src/components/backoffice/CreateEmployeeModal.tsx` (línea 90)
4. ✅ `frontend/src/components/backoffice/EditEmployeeModal.tsx` (línea 121)
5. ✅ `frontend/src/components/appointments/ConfirmAppointmentModal.tsx` (línea 44)
6. ✅ `frontend/src/components/appointments/AppointmentDetailsModal.tsx` (línea 69)

### **Páginas corregidas:**

1. ✅ `frontend/src/app/dashboard/backoffice/clients/page.tsx` (callbacks sin toasts)
2. ✅ `frontend/src/app/dashboard/backoffice/employees/page.tsx` (callbacks sin toasts)

---

## 🎓 Recursos Adicionales

- **Guía de Toasts:** `docs/TOASTS_GUIDE.md`
- **Error Handling:** `docs/OPTIMIZATION_GUIDE.md`
- **Convenciones de Código:** `docs/DEVELOPMENT_GUIDE.md`

---

## 📝 Historial de Cambios

| Fecha | Cambio | Responsable |
|-------|--------|-------------|
| 2025-12-08 | Documentación inicial + corrección de toasts duplicados | AI Development Team |

---

**Recuerda:** El objetivo es proporcionar feedback **claro, inmediato y no repetitivo** al usuario. Un solo toast en el momento correcto es mejor que múltiples toasts confusos.
