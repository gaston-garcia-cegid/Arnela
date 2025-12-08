# 🐛 Bug de Navegación: Dashboard Siempre Seleccionado

> **Identificado:** Diciembre 8, 2025  
> **Resuelto:** Diciembre 8, 2025  
> **Componente:** `BackofficeSidebar.tsx`

---

## 📋 Resumen del Bug

**Síntoma:** Al navegar a cualquier página del backoffice (Empleados, Clientes, etc.), la opción "Dashboard" permanecía visualmente seleccionada junto con la opción correcta.

**Impacto:** UX confusa - dos opciones del menú activas simultáneamente.

---

## 🔍 Análisis Técnico

### **Causa Raíz:**

La lógica de detección de ruta activa usaba `startsWith()` sin considerar que todas las rutas del backoffice comparten el mismo prefijo:

```typescript
// ❌ CÓDIGO INCORRECTO (antes)
const isActive =
  pathname === item.href || pathname?.startsWith(item.href + "/");

// Problema:
// Dashboard href: "/dashboard/backoffice"
// Empleados pathname: "/dashboard/backoffice/employees"
// 
// Evaluación:
// "/dashboard/backoffice/employees".startsWith("/dashboard/backoffice/")
// => true ❌ (Dashboard marcado como activo incorrectamente)
```

### **Rutas Afectadas:**

| Ruta | Dashboard Activo (Bug) | Debería Ser |
|------|------------------------|-------------|
| `/dashboard/backoffice` | ✅ Correcto | ✅ Activo |
| `/dashboard/backoffice/employees` | ❌ Incorrecto | ⬜ No activo |
| `/dashboard/backoffice/clients` | ❌ Incorrecto | ⬜ No activo |
| `/dashboard/backoffice/appointments` | ❌ Incorrecto | ⬜ No activo |
| `/dashboard/backoffice/billing` | ❌ Incorrecto | ⬜ No activo |

---

## ✅ Solución Implementada

### **Código Corregido:**

```typescript
// ✅ CÓDIGO CORRECTO (después)
const isActive = item.href === "/dashboard/backoffice"
  ? pathname === item.href // Coincidencia EXACTA solo para Dashboard
  : pathname === item.href || pathname?.startsWith(item.href + "/");
```

### **Lógica:**

1. **Dashboard (`/dashboard/backoffice`):**
   - Solo se activa si la ruta es **exactamente** `"/dashboard/backoffice"`
   - NO se activa para subrutas como `/dashboard/backoffice/employees`

2. **Otros items del menú:**
   - Se activan con coincidencia exacta O subrutas
   - Ejemplo: "Empleados" activo en `/dashboard/backoffice/employees` y `/dashboard/backoffice/employees/123`

### **Resultado:**

| Ruta | Dashboard | Empleados | Comportamiento |
|------|-----------|-----------|----------------|
| `/dashboard/backoffice` | ✅ Activo | ⬜ No | ✅ Correcto |
| `/dashboard/backoffice/employees` | ⬜ No | ✅ Activo | ✅ Correcto |
| `/dashboard/backoffice/clients` | ⬜ No | ⬜ No | ✅ Correcto |

---

## 🚫 Cómo Prevenir Bugs Similares

### **1. Principio de Coincidencia de Rutas**

**Regla General:**
> "Las rutas padre con hijos deben usar coincidencia exacta; las rutas sin hijos pueden usar `startsWith()`"

**Aplicar en:**
- Menús de navegación
- Breadcrumbs
- Tabs de navegación
- Cualquier sistema de routing visual

---

### **2. Patrón de Detección de Ruta Activa**

#### **Opción A: Ruta exacta vs. ruta con hijos**

```typescript
// Para rutas "parent" que tienen hijos (no deben activarse en subrutas)
const isParentRoute = href === "/dashboard/backoffice";

const isActive = isParentRoute
  ? pathname === href // Exacta
  : pathname === href || pathname?.startsWith(href + "/"); // Exacta o hijos
```

#### **Opción B: Lista de rutas exactas**

```typescript
const exactMatchRoutes = ["/dashboard/backoffice", "/dashboard"];

const isActive = exactMatchRoutes.includes(item.href)
  ? pathname === item.href // Exacta solo para rutas en lista
  : pathname === item.href || pathname?.startsWith(item.href + "/");
```

#### **Opción C: Flag en configuración de menú**

```typescript
const menuItems = [
  {
    title: "Dashboard",
    href: "/dashboard/backoffice",
    icon: LayoutDashboard,
    exactMatch: true, // 🎯 Flag para forzar coincidencia exacta
  },
  {
    title: "Empleados",
    href: "/dashboard/backoffice/employees",
    icon: UserCircle,
    // exactMatch: false por defecto
  },
];

// En el componente:
const isActive = item.exactMatch
  ? pathname === item.href
  : pathname === item.href || pathname?.startsWith(item.href + "/");
```

---

### **3. Checklist para Code Review**

Al revisar código de navegación, verificar:

- [ ] **¿La ruta padre tiene hijos?**
  - Si SÍ → Usar coincidencia exacta
  - Si NO → Puede usar `startsWith()`

- [ ] **¿Todas las rutas comparten prefijo común?**
  - Ejemplo: `/dashboard/backoffice/*`
  - Requiere lógica especial para la ruta raíz

- [ ] **¿Se probó la navegación manualmente?**
  - Navegar a cada ruta
  - Verificar que solo UN item esté activo

- [ ] **¿Hay rutas anidadas de 3+ niveles?**
  - Ejemplo: `/dashboard/backoffice/billing/invoices/123`
  - Verificar que items intermedios NO se activen incorrectamente

---

### **4. Testing Automatizado**

#### **Test para prevenir el bug:**

```typescript
// __tests__/BackofficeSidebar.test.tsx
import { render, screen } from '@testing-library/react';
import { BackofficeSidebar } from '@/components/backoffice/BackofficeSidebar';

// Mock usePathname
vi.mock('next/navigation', () => ({
  usePathname: vi.fn(),
}));

describe('BackofficeSidebar - Active Route Detection', () => {
  it('should only activate Dashboard when on exact /dashboard/backoffice route', () => {
    vi.mocked(usePathname).mockReturnValue('/dashboard/backoffice');
    
    render(<BackofficeSidebar />);
    
    const dashboardLink = screen.getByText('Dashboard').closest('a');
    const employeesLink = screen.getByText('Empleados').closest('a');
    
    expect(dashboardLink).toHaveClass('bg-primary'); // ✅ Activo
    expect(employeesLink).not.toHaveClass('bg-primary'); // ⬜ No activo
  });

  it('should NOT activate Dashboard when on child route /dashboard/backoffice/employees', () => {
    vi.mocked(usePathname).mockReturnValue('/dashboard/backoffice/employees');
    
    render(<BackofficeSidebar />);
    
    const dashboardLink = screen.getByText('Dashboard').closest('a');
    const employeesLink = screen.getByText('Empleados').closest('a');
    
    expect(dashboardLink).not.toHaveClass('bg-primary'); // ⬜ No activo
    expect(employeesLink).toHaveClass('bg-primary'); // ✅ Activo
  });

  it('should activate Empleados for nested routes like /dashboard/backoffice/employees/123', () => {
    vi.mocked(usePathname).mockReturnValue('/dashboard/backoffice/employees/123');
    
    render(<BackofficeSidebar />);
    
    const employeesLink = screen.getByText('Empleados').closest('a');
    
    expect(employeesLink).toHaveClass('bg-primary'); // ✅ Activo en subruta
  });
});
```

---

### **5. Patrones Comunes de Errores**

#### **Error 1: Usar solo `startsWith()` sin validación**

```typescript
// ❌ INCORRECTO
const isActive = pathname?.startsWith(item.href);

// Problema: "/dashboard" activo en "/dashboard/settings"
```

#### **Error 2: No considerar el slash final**

```typescript
// ❌ INCORRECTO
const isActive = pathname?.startsWith(item.href);

// Problema:
// href: "/dashboard"
// pathname: "/dashboard-settings"
// => startsWith retorna true! ❌
```

**Solución:**
```typescript
// ✅ CORRECTO
const isActive = pathname === item.href || pathname?.startsWith(item.href + "/");
// Asegura que después del prefijo haya un "/"
```

#### **Error 3: No normalizar rutas**

```typescript
// ❌ INCORRECTO
const isActive = pathname === item.href;

// Problema:
// href: "/dashboard"
// pathname: "/dashboard/" (con slash final)
// => false negativo ❌
```

**Solución:**
```typescript
// ✅ CORRECTO
const normalizedPathname = pathname?.replace(/\/$/, ""); // Remover slash final
const normalizedHref = item.href.replace(/\/$/, "");
const isActive = normalizedPathname === normalizedHref;
```

---

### **6. Documentación en Código**

Agregar comentarios claros en la lógica de routing:

```typescript
// ✅ BIEN DOCUMENTADO
const isActive = item.href === "/dashboard/backoffice"
  ? pathname === item.href // Dashboard: solo ruta exacta (tiene hijos)
  : pathname === item.href || pathname?.startsWith(item.href + "/"); // Otros: ruta o subrutas
```

---

### **7. Guía de Estilo para Rutas**

Agregar a `DEVELOPMENT_GUIDE.md`:

```markdown
## Navegación y Rutas Activas

### Reglas para detectar rutas activas:

1. **Rutas raíz con hijos:** Usar coincidencia EXACTA
   ```typescript
   const isActive = pathname === "/dashboard/backoffice";
   ```

2. **Rutas sin hijos:** Permitir coincidencia con subrutas
   ```typescript
   const isActive = 
     pathname === href || pathname?.startsWith(href + "/");
   ```

3. **Siempre agregar "/" después del prefijo en startsWith**
   ```typescript
   // ✅ Correcto
   pathname?.startsWith(href + "/")
   
   // ❌ Incorrecto
   pathname?.startsWith(href)
   ```

4. **Considerar normalización de slashes finales**
   ```typescript
   const normalized = pathname?.replace(/\/$/, "");
   ```
```

---

## 📚 Referencias

- **Archivo corregido:** `frontend/src/components/backoffice/BackofficeSidebar.tsx` (líneas 200-203)
- **Documentación de Next.js:** [usePathname](https://nextjs.org/docs/app/api-reference/functions/use-pathname)
- **Patrón similar en:** Breadcrumbs, Tabs, Nested Navigation

---

## 📊 Métricas

| Métrica | Antes | Después |
|---------|-------|---------|
| Rutas con selección correcta | 1/5 (20%) | 5/5 (100%) |
| Items simultáneamente activos | 2 | 1 |
| Confusión de usuario | Alta | Ninguna |

---

## ✅ Validación

Para verificar que el bug está resuelto:

1. **Iniciar la aplicación:**
   ```bash
   cd frontend
   pnpm dev
   ```

2. **Navegar a cada ruta:**
   - `/dashboard/backoffice` → Solo Dashboard activo ✅
   - `/dashboard/backoffice/employees` → Solo Empleados activo ✅
   - `/dashboard/backoffice/clients` → Solo Clientes activo ✅
   - `/dashboard/backoffice/appointments` → Solo Citas activo ✅
   - `/dashboard/backoffice/billing` → Solo Facturación activo ✅

3. **Verificar visualmente:**
   - Un solo item del menú con fondo de color
   - Otros items en color gris/muted

---

**Estado:** ✅ Bug corregido y documentado  
**Fecha:** Diciembre 8, 2025  
**Responsable:** AI Development Team
