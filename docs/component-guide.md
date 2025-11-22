# Guía de Componentes Reutilizables - Arnela

## 🎨 Sistema de Diseño

Esta guía define los patrones de diseño y componentes reutilizables para mantener consistencia en toda la aplicación.

## 📦 Componentes Base (Shadcn)

Todos los componentes base están en `src/components/ui/` y siguen el patrón de Shadcn.

### Button
```tsx
import { Button } from '@/components/ui/button';

// Variantes
<Button variant="default">Primary</Button>
<Button variant="secondary">Secondary</Button>
<Button variant="destructive">Delete</Button>
<Button variant="outline">Outline</Button>
<Button variant="ghost">Ghost</Button>
<Button variant="link">Link</Button>

// Tamaños
<Button size="sm">Small</Button>
<Button size="default">Default</Button>
<Button size="lg">Large</Button>
```

### Card
```tsx
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui/card';

<Card>
  <CardHeader>
    <CardTitle>Título</CardTitle>
    <CardDescription>Descripción</CardDescription>
  </CardHeader>
  <CardContent>
    Contenido de la tarjeta
  </CardContent>
  <CardFooter>
    <Button>Acción</Button>
  </CardFooter>
</Card>
```

### Form
```tsx
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { useForm } from 'react-hook-form';

const form = useForm();

<Form {...form}>
  <form onSubmit={form.handleSubmit(onSubmit)}>
    <FormField
      control={form.control}
      name="email"
      render={({ field }) => (
        <FormItem>
          <FormLabel>Email</FormLabel>
          <FormControl>
            <Input type="email" {...field} />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  </form>
</Form>
```

### Alert
```tsx
import { Alert, AlertDescription } from '@/components/ui/alert';
import { AlertCircle } from 'lucide-react';

<Alert variant="default">
  <AlertCircle className="h-4 w-4" />
  <AlertDescription>Mensaje informativo</AlertDescription>
</Alert>

<Alert variant="destructive">
  <AlertCircle className="h-4 w-4" />
  <AlertDescription>Error message</AlertDescription>
</Alert>
```

## 🎭 Patrones de Composición

### Layout Container
```tsx
// Para contenido con ancho máximo
<div className="container mx-auto px-4">
  {/* Contenido centrado con padding */}
</div>

// Para secciones full-width con contenido centrado
<section className="w-full bg-background">
  <div className="container mx-auto px-4 py-16">
    {/* Contenido */}
  </div>
</section>
```

### Grid Responsive
```tsx
// Grid de 2-4 columnas según viewport
<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
  <Card>...</Card>
  <Card>...</Card>
  <Card>...</Card>
</div>
```

### Stack Vertical
```tsx
// Elementos apilados con espacio
<div className="space-y-4">
  <Card>...</Card>
  <Card>...</Card>
</div>

// Con más espacio
<div className="space-y-6 md:space-y-8">
  <Section />
  <Section />
</div>
```

### Flex Layouts
```tsx
// Centrado horizontal y vertical
<div className="flex items-center justify-center min-h-screen">
  <LoginModal />
</div>

// Header con logo y acciones
<header className="flex items-center justify-between p-4">
  <Logo />
  <nav className="flex items-center gap-4">
    <Button>Login</Button>
  </nav>
</header>

// Lista con iconos
<div className="flex items-center gap-3">
  <Icon className="h-5 w-5 text-muted-foreground" />
  <span>Texto</span>
</div>
```

## 🎨 Clases Utilitarias Comunes

### Spacing (usa múltiplos de 4px)
```
p-1   = 4px      gap-1  = 4px
p-2   = 8px      gap-2  = 8px
p-3   = 12px     gap-3  = 12px
p-4   = 16px  ⭐ gap-4  = 16px  ⭐
p-6   = 24px     gap-6  = 24px
p-8   = 32px     gap-8  = 32px
p-12  = 48px     gap-12 = 48px
p-16  = 64px     gap-16 = 64px
```

### Typography
```tsx
// Headings
<h1 className="text-4xl font-bold tracking-tight">
<h2 className="text-3xl font-semibold">
<h3 className="text-2xl font-semibold">
<h4 className="text-xl font-medium">

// Body
<p className="text-base text-foreground">        // Normal
<p className="text-sm text-muted-foreground">   // Small
<p className="text-lg">                         // Large

// Labels
<label className="text-sm font-medium">
```

### Borders & Shadows
```tsx
// Borders
className="border border-border"              // Borde simple
className="border-b border-border"            // Solo abajo
className="divide-y divide-border"            // Entre elementos

// Shadows
className="shadow-sm"                         // Sutil
className="shadow-md"                         // Medio
className="shadow-lg"                         // Grande
className="shadow-xl"                         // Extra grande

// Rings (focus)
className="focus:ring-2 focus:ring-ring focus:ring-offset-2"
```

### Rounded Corners
```tsx
className="rounded-sm"     // 0.3rem - Botones pequeños
className="rounded-md"     // 0.4rem - Inputs
className="rounded-lg"     // 0.5rem - Cards
className="rounded-xl"     // 0.75rem - Modales
className="rounded-full"   // Circular - Avatars
```

### Interactive States
```tsx
// Hover
className="hover:bg-accent hover:text-accent-foreground"
className="hover:scale-105 transition-transform"

// Focus
className="focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"

// Disabled
className="disabled:pointer-events-none disabled:opacity-50"

// Active
className="active:scale-95"
```

## 📐 Breakpoints Responsive

```tsx
// Mobile-first approach
sm:   // @media (min-width: 640px)
md:   // @media (min-width: 768px)
lg:   // @media (min-width: 1024px)
xl:   // @media (min-width: 1280px)
2xl:  // @media (min-width: 1536px)

// Ejemplo
<div className="
  text-sm      /* Mobile */
  md:text-base /* Tablet */
  lg:text-lg   /* Desktop */
">
```

## 🎯 Patrones de Componentes Específicos

### Dashboard Card con Estadística
```tsx
<Card>
  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
    <CardTitle className="text-sm font-medium">
      Total Clientes
    </CardTitle>
    <Users className="h-4 w-4 text-muted-foreground" />
  </CardHeader>
  <CardContent>
    <div className="text-2xl font-bold">254</div>
    <p className="text-xs text-muted-foreground">
      +12% desde el mes pasado
    </p>
  </CardContent>
</Card>
```

### Lista con Acciones
```tsx
<div className="rounded-lg border border-border">
  <table className="w-full">
    <thead className="border-b border-border bg-muted/50">
      <tr>
        <th className="px-4 py-3 text-left text-sm font-medium">Nombre</th>
        <th className="px-4 py-3 text-left text-sm font-medium">Email</th>
        <th className="px-4 py-3 text-right text-sm font-medium">Acciones</th>
      </tr>
    </thead>
    <tbody className="divide-y divide-border">
      {items.map(item => (
        <tr key={item.id} className="hover:bg-muted/50">
          <td className="px-4 py-3">{item.name}</td>
          <td className="px-4 py-3">{item.email}</td>
          <td className="px-4 py-3 text-right">
            <Button variant="ghost" size="sm">Ver</Button>
          </td>
        </tr>
      ))}
    </tbody>
  </table>
</div>
```

### Modal / Dialog Pattern
```tsx
<Dialog open={isOpen} onOpenChange={setIsOpen}>
  <DialogContent className="sm:max-w-[425px]">
    <DialogHeader>
      <DialogTitle>Título del Modal</DialogTitle>
      <DialogDescription>
        Descripción breve del contenido
      </DialogDescription>
    </DialogHeader>
    
    {/* Contenido */}
    <div className="space-y-4 py-4">
      <FormField />
    </div>
    
    <DialogFooter>
      <Button variant="outline" onClick={() => setIsOpen(false)}>
        Cancelar
      </Button>
      <Button onClick={handleSubmit}>
        Guardar
      </Button>
    </DialogFooter>
  </DialogContent>
</Dialog>
```

### Empty State
```tsx
<div className="flex flex-col items-center justify-center py-12">
  <div className="rounded-full bg-muted p-3 mb-4">
    <InboxIcon className="h-6 w-6 text-muted-foreground" />
  </div>
  <h3 className="text-lg font-semibold mb-2">No hay clientes</h3>
  <p className="text-sm text-muted-foreground mb-4">
    Comienza agregando tu primer cliente
  </p>
  <Button>
    <Plus className="mr-2 h-4 w-4" />
    Agregar Cliente
  </Button>
</div>
```

## ⚠️ Anti-patrones (NO HACER)

❌ **NO crear archivos CSS separados**
```tsx
// ❌ MAL
import './MyComponent.css'

// ✅ BIEN - Usar Tailwind classes
<div className="bg-card p-4 rounded-lg">
```

❌ **NO usar estilos inline**
```tsx
// ❌ MAL
<div style={{ padding: '16px', backgroundColor: '#fff' }}>

// ✅ BIEN
<div className="p-4 bg-card">
```

❌ **NO hardcodear colores**
```tsx
// ❌ MAL
<div className="bg-[#3b82f6]">

// ✅ BIEN - Usar variables del tema
<div className="bg-primary">
```

❌ **NO repetir clases comunes**
```tsx
// ❌ MAL - Repetir mismo pattern
<button className="px-4 py-2 bg-primary text-primary-foreground rounded-md">
<button className="px-4 py-2 bg-primary text-primary-foreground rounded-md">

// ✅ BIEN - Usar componente Button
<Button>Acción 1</Button>
<Button>Acción 2</Button>
```

## 📏 Convenciones de Código

### Orden de clases Tailwind
```tsx
// Orden recomendado:
// 1. Layout (flex, grid, block)
// 2. Positioning (relative, absolute)
// 3. Size (w-, h-, min-, max-)
// 4. Spacing (p-, m-, gap-)
// 5. Typography (text-, font-)
// 6. Visual (bg-, border-, shadow-)
// 7. Interactive (hover:, focus:, active:)

<div className="
  flex items-center justify-between
  relative
  w-full h-12
  px-4 py-2 gap-2
  text-sm font-medium
  bg-card border border-border rounded-lg shadow-sm
  hover:bg-accent focus:ring-2
">
```

### Responsive Design
```tsx
// Mobile-first: base styles para mobile, breakpoints para desktop
<div className="
  grid grid-cols-1    /* Mobile: 1 columna */
  md:grid-cols-2      /* Tablet: 2 columnas */
  lg:grid-cols-3      /* Desktop: 3 columnas */
  xl:grid-cols-4      /* Desktop XL: 4 columnas */
  gap-4
">
```

## 🔄 Actualizaciones

Este documento se actualiza cuando:
- Se agregan nuevos componentes base
- Se establecen nuevos patrones de diseño
- Se identifican anti-patrones comunes
- Se modifican las convenciones de estilo

**Última actualización:** 22 Nov 2025
