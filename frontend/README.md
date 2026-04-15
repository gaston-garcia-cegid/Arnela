# Arnela Frontend

Aplicación web para el CRM/CMS de Arnela Gabinete. Next.js 16, React 19, TypeScript, Tailwind v4, Zustand.

## Requisitos

- Node.js 22+
- pnpm

## Inicio rápido

```powershell
copy .env.example .env
pnpm install
pnpm dev
```

Abre `http://localhost:3000`. Requiere que el backend esté corriendo en `http://localhost:8080`.

## Variables de entorno

| Variable | Default | Descripción |
|----------|---------|-------------|
| `NEXT_PUBLIC_API_URL` | `http://localhost:8080/api/v1` | URL base de la API |

## Estructura

```
frontend/src/
├── app/                         # Next.js App Router
│   ├── layout.tsx               # Root layout (ThemeProvider, Toaster)
│   ├── page.tsx                 # Landing page
│   ├── globals.css              # Tailwind v4 + custom theme (light/dark)
│   ├── contacto/                # Página de contacto
│   ├── sobre-arnela/            # Sobre nosotros
│   ├── intervencion/            # Servicios de intervención
│   ├── formacion/               # Servicios de formación
│   ├── convenios-y-colaboraciones/
│   └── dashboard/
│       ├── backoffice/          # Dashboard admin/employee
│       │   ├── layout.tsx       # Sidebar + header + theme toggle
│       │   ├── page.tsx         # Dashboard principal con stats
│       │   ├── clients/         # CRUD de clientes
│       │   ├── employees/       # Gestión de empleados
│       │   ├── appointments/    # Gestión de citas
│       │   └── billing/         # Facturación
│       │       ├── invoices/    # Facturas (CRUD + nueva)
│       │       ├── expenses/    # Gastos (CRUD + nuevo)
│       │       └── categories/  # Categorías de gasto
│       └── client/              # Dashboard del cliente
│           └── appointments/    # Mis citas
│
├── components/
│   ├── ui/                      # Shadcn/Radix primitivos (Button, Dialog, etc.)
│   ├── common/                  # Navbar, Footer, ThemeToggle, ThemeProvider
│   ├── landing/                 # Hero, About, Services, Testimonials, Reviews
│   ├── auth/                    # LoginModal
│   ├── appointments/            # AppointmentForm, AppointmentList, Calendar
│   ├── backoffice/              # ClientModals, EmployeeForms, DataTables
│   ├── billing/                 # InvoiceForm, ExpenseForm, StatsCards
│   ├── dashboard/               # StatsCards, QuickActions
│   ├── search/                  # GlobalSearch (Ctrl+K)
│   └── examples/                # Componentes de ejemplo
│
├── hooks/                       # Custom hooks
│   ├── useAppointments.ts       # CRUD citas + optimistic updates
│   ├── useDebounce.ts           # Debounce para búsquedas
│   ├── useErrorHandler.ts       # Error handling centralizado
│   └── useStats.ts              # Dashboard stats
│
├── lib/
│   ├── api.ts                   # Cliente API (fetch wrapper + auth headers)
│   ├── errors.ts                # Error types + handling
│   ├── validators.ts            # Validaciones (DNI, teléfono, email)
│   ├── utils.ts                 # cn(), formatters
│   ├── appointmentUtils.ts      # Helpers de citas
│   ├── exportUtils.ts           # Exportación Excel (xlsx)
│   └── performanceUtils.ts      # Debounce, throttle, memo, virtualScroll
│
├── stores/
│   ├── authStore.ts             # Auth state (login, logout, token, user)
│   └── ...                      # Otros stores Zustand
│
├── types/                       # TypeScript interfaces y types
│
└── test/
    └── setup.ts                 # Vitest setup (testing-library, jsdom)
```

## Scripts

| Comando | Descripción |
|---------|-------------|
| `pnpm dev` | Servidor de desarrollo con hot reload |
| `pnpm build` | Build de producción (standalone) |
| `pnpm start` | Iniciar build de producción |
| `pnpm lint` | ESLint |
| `pnpm test` | Vitest en modo watch |
| `pnpm test -- --run` | Vitest single run |
| `pnpm test:ui` | Vitest con UI web |
| `pnpm test:coverage` | Coverage report |

## Testing

14 archivos de test con Vitest + Testing Library + jsdom:

**Hooks** (`hooks/__tests__/`):
- `useAppointments.test.ts` - CRUD y paginación de citas
- `useDebounce.test.ts` - Debounce timing
- `useErrorHandler.test.ts` - Error handling
- `useStats.test.ts` - Dashboard statistics

**Lib** (`lib/__tests__/`):
- `api.test.ts` - API client, auth headers, error handling
- `appointmentUtils.test.ts` - Utilidades de citas
- `errors.test.ts` - Error types
- `exportUtils.test.ts` - Excel export
- `performanceUtils.test.ts` - Performance utilities
- `utils.test.ts` - General utilities
- `validators.test.ts` - Validación DNI, teléfono, email

**Componentes** (`components/*/__tests__/`):
- `LoginModal.test.tsx` - Login flow, validación, errores backend
- `CreateClientModal.test.tsx` - Creación de clientes
- `EditClientModal.test.tsx` - Edición de clientes

## Stack de UI

- **Tailwind CSS v4** con tema custom (light + dark mode)
- **Shadcn/Radix UI** para componentes base (Dialog, Select, Tabs, etc.)
- **Lucide React** para iconos
- **next-themes** para dark/light mode toggle
- **react-hook-form + zod** para formularios con validación
- **sonner** para toast notifications
- **date-fns** para formateo de fechas
- **xlsx** para exportación Excel

## Patrones principales

- **State management**: Zustand stores (auth, UI state)
- **Data fetching**: Custom hooks con fetch API + error handling
- **Forms**: react-hook-form + zodResolver + validación server-side
- **Error handling**: Centralizado en `useErrorHandler` + toast notifications
- **Theming**: CSS variables en `globals.css`, switchable via `next-themes`
- **SEO**: Metadata estática por página con OpenGraph
- **Routing**: Next.js App Router, layouts anidados por sección
