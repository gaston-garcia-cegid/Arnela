# Migración a Tailwind CSS 4 - Arnela Frontend

## ✅ Cambios Realizados

### 1. **globals.css** - Migrado a sintaxis CSS-first de Tailwind 4
```css
@theme {
  /* Colores usando OKLCH (nuevo estándar de color) */
  --color-primary: oklch(65% 0.2 217);
  --color-background: oklch(97% 0.02 39);
  /* ... más colores */
}
```

**Ventajas:**
- Colores más precisos y perceptualmente uniformes
- Mejor soporte para dark mode
- Configuración más simple y mantenible
- No requiere tailwind.config.ts

### 2. **Eliminado tailwind.config.ts**
- Tailwind 4 usa CSS-first configuration
- Toda la configuración ahora está en `globals.css` con `@theme`
- Más fácil de mantener y compartir entre proyectos

### 3. **postcss.config.js** - Simplificado
```javascript
module.exports = {
  plugins: {
    '@tailwindcss/postcss': {},
  },
};
```
- Removido autoprefixer (incluido en Tailwind 4)
- Solo necesita @tailwindcss/postcss

### 4. **Actualizado sintaxis de clases**
| Tailwind 3 (Antiguo) | Tailwind 4 (Nuevo) |
|---------------------|-------------------|
| `bg-gradient-to-br` | `bg-linear-to-br` |
| `bg-gradient-to-r`  | `bg-linear-to-r`  |
| `supports-[backdrop-filter]` | `supports-backdrop-filter` |

**Archivos actualizados:**
- `components/landing/Hero.tsx`
- `components/landing/About.tsx`
- `components/landing/Testimonial.tsx`
- `components/common/Navbar.tsx`

## 🎨 Sistema de Diseño Estandarizado

### Paleta de Colores (Light Mode)
```
Background: oklch(97% 0.02 39)      - Beige muy claro
Foreground: oklch(25% 0.02 220)     - Gris oscuro cálido
Primary: oklch(65% 0.2 217)         - Azul profesional
Secondary: oklch(45% 0.15 142)      - Verde bosque (bienestar)
Accent: oklch(72% 0.15 27)          - Coral suave (energía)
Muted: oklch(95% 0.02 39)           - Beige claro
Destructive: oklch(55% 0.22 0)      - Rojo suave
```

### Uso de Colores
```tsx
// Backgrounds
className="bg-background"           // Fondo principal
className="bg-card"                 // Tarjetas
className="bg-muted"                // Fondos alternativos

// Text
className="text-foreground"         // Texto principal
className="text-muted-foreground"   // Texto secundario
className="text-primary"            // Texto destacado

// Borders
className="border-border"           // Bordes estándar
className="border-input"            // Bordes de inputs

// Interactive
className="ring-ring"               // Focus rings
className="hover:bg-accent"         // Hover states
```

### Border Radius
```
--radius-lg: 0.5rem    → rounded-lg
--radius-md: 0.4rem    → rounded-md  
--radius-sm: 0.3rem    → rounded-sm
```

## 📁 Estructura CSS Centralizada

**Único archivo CSS necesario:** `src/app/globals.css`

### Organización:
1. `@import "tailwindcss"` - Import principal
2. `@theme { }` - Configuración de colores y variables
3. `@media (prefers-color-scheme: dark)` - Dark mode
4. `@layer base` - Estilos base
5. `@layer utilities` - Utilidades custom

**NO crear más archivos .css** - Todo debe vivir en globals.css

## 🔧 Componentes Shadcn Compatible

Los componentes shadcn funcionan perfectamente con esta configuración:
- ✅ Button
- ✅ Dialog
- ✅ Form
- ✅ Input
- ✅ Label
- ✅ Select
- ✅ Alert

Todos usan las variables CSS definidas en `@theme`

## 📝 Próximos Pasos

1. **Reiniciar servidor de desarrollo**
   ```bash
   # Detener proceso existente primero
   pnpm dev
   ```

2. **Verificar en el navegador**
   - Los colores deben verse consistentes
   - Dark mode debe funcionar (si está habilitado)
   - Gradientes deben renderizarse correctamente

3. **Si hay problemas:**
   - Limpiar caché: `rm -rf .next`
   - Reinstalar deps: `pnpm install`
   - Verificar versión: Tailwind CSS v4.1.17

## 🎯 Beneficios de esta Configuración

✅ **Un solo archivo de configuración** - Todo en globals.css
✅ **Colores estandarizados** - Sistema de design tokens
✅ **Dark mode automático** - Usando prefers-color-scheme
✅ **Compatible con shadcn** - Variables CSS estándar
✅ **Sin conflictos** - No más archivos CSS duplicados
✅ **Mejor rendimiento** - CSS-first es más rápido
✅ **Fácil mantenimiento** - Cambios en un solo lugar

## 🐛 Debugging

Si los estilos no se aplican:
1. Verificar que no exista `tailwind.config.ts`
2. Verificar `@import "tailwindcss"` en globals.css
3. Verificar que layout.tsx importe `'./globals.css'`
4. Limpiar caché de Next.js: `rm -rf .next`
5. Reiniciar servidor de desarrollo
