---
trigger: always_on
---

# ⚛️ Frontend Architecture (Next.js 16)
## Stack Tecnológico
- **Framework**: Next.js 16 (App Router)
- **Lenguaje**: TypeScript
- **Estado Global**: Zustand
- **UI Components**: Shadcn UI
- **Styling**: Tailwind CSS / React CSS
- **HTTP Client**: Fetch API / Axios

## Estructura de Carpetas
frontend/
├── app/                          # App Router (Next.js 16)
│   ├── layout.tsx               # Root layout
│   ├── page.tsx                 # Landing page
│   ├── (auth)/                  # Auth routes group
│   │   ├── login/
│   │   │   └── page.tsx
│   │   └── register/
│   │       └── page.tsx
│   ├── (client)/                # Client area group
│   │   ├── layout.tsx
│   │   ├── dashboard/
│   │   │   └── page.tsx
│   │   └── appointments/
│   │       ├── page.tsx
│   │       └── [id]/
│   │           └── page.tsx
│   └── (backoffice)/            # Admin/Employee group
│       ├── layout.tsx
│       ├── dashboard/
│       │   └── page.tsx
│       ├── clients/
│       │   ├── page.tsx
│       │   ├── new/
│       │   │   └── page.tsx
│       │   └── [id]/
│       │       └── page.tsx
│       ├── employees/
│       │   └── page.tsx
│       └── appointments/
│           └── page.tsx
├── components/
│   ├── ui/                      # Shadcn components
│   │   ├── button.tsx
│   │   ├── input.tsx
│   │   ├── dialog.tsx
│   │   └── ...
│   ├── common/                  # Shared components
│   │   ├── Navbar.tsx
│   │   ├── Sidebar.tsx
│   │   ├── LoadingSpinner.tsx
│   │   └── ErrorBoundary.tsx
│   ├── client/                  # Client-specific
│   │   ├── AppointmentCard.tsx
│   │   └── ClientProfile.tsx
│   └── backoffice/              # Admin/Employee
│       ├── ClientTable.tsx
│       ├── ClientForm.tsx
│       ├── EmployeeTable.tsx
│       └── Calendar.tsx
├── stores/                       # Zustand stores
│   ├── useAuthStore.ts
│   ├── useClientStore.ts
│   ├── useEmployeeStore.ts
│   └── useAppointmentStore.ts
├── hooks/                        # Custom hooks
│   ├── useDebounce.ts
│   ├── useAuth.ts
│   ├── useClients.ts
│   └── useAppointments.ts
├── lib/                          # Utilities
│   ├── api.ts                   # API client
│   ├── constants.ts
│   ├── validators.ts
│   └── formatters.ts
└── types/                        # TypeScript types
    ├── user.ts
    ├── client.ts
    ├── employee.ts
    └── appointment.ts

## Zustand State Management
### Auth Store
typescript
// stores/useAuthStore.ts
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface User {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  role: 'admin' | 'employee' | 'client';
}

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  
  login: (user: User, token: string) => void;
  logout: () => void;
  setUser: (user: User) => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,

      login: (user, token) => 
        set({ user, token, isAuthenticated: true }),

      logout: () => 
        set({ user: null, token: null, isAuthenticated: false }),

      setUser: (user) => 
        set({ user }),
    }),
    {
      name: 'auth-storage', // LocalStorage key
    }
  )
);

### Client Store
typescript
// stores/useClientStore.ts
import { create } from 'zustand';

interface Client {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  dni: string;
  city?: string;
  isActive: boolean;
}

interface ClientState {
  clients: Client[];
  selectedClient: Client | null;
  isLoading: boolean;
  error: string | null;

  setClients: (clients: Client[]) => void;
  addClient: (client: Client) => void;
  updateClient: (id: string, client: Partial<Client>) => void;
  deleteClient: (id: string) => void;
  selectClient: (client: Client | null) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
}

export const useClientStore = create<ClientState>((set) => ({
  clients: [],
  selectedClient: null,
  isLoading: false,
  error: null,

  setClients: (clients) => set({ clients }),
  
  addClient: (client) => 
    set((state) => ({ clients: [...state.clients, client] })),
  
  updateClient: (id, updatedClient) =>
    set((state) => ({
      clients: state.clients.map((c) =>
        c.id === id ? { ...c, ...updatedClient } : c
      ),
    })),
  
  deleteClient: (id) =>
    set((state) => ({
      clients: state.clients.filter((c) => c.id !== id),
    })),
  
  selectClient: (client) => set({ selectedClient: client }),
  
  setLoading: (loading) => set({ isLoading: loading }),
  
  setError: (error) => set({ error }),
}));

## API Client
typescript
// lib/api.ts
import { useAuthStore } from '@/stores/useAuthStore';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

class ApiClient {
  private baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const token = useAuthStore.getState().token;

    const config: RequestInit = {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...(token && { Authorization: `Bearer ${token}` }),
        ...options.headers,
      },
    };

    const response = await fetch(`${this.baseURL}${endpoint}`, config);

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || 'Request failed');
    }

    return response.json();
  }

  // Auth endpoints
  async login(email: string, password: string) {
    return this.request('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
  }

  async register(data: RegisterData) {
    return this.request('/auth/register', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async getMe() {
    return this.request('/auth/me');
  }

  // Client endpoints
  async getClients(filters?: ClientFilters) {
    const params = new URLSearchParams(filters as any);
    return this.request(`/clients?${params}`);
  }

  async getClient(id: string) {
    return this.request(`/clients/${id}`);
  }

  async createClient(data: CreateClientData) {
    return this.request('/clients', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async updateClient(id: string, data: Partial<CreateClientData>) {
    return this.request(`/clients/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async deleteClient(id: string) {
    return this.request(`/clients/${id}`, {
      method: 'DELETE',
    });
  }

  // Employee endpoints (future)
  async getEmployees() {
    return this.request('/employees');
  }

  // Appointment endpoints (future)
  async getAppointments(filters?: AppointmentFilters) {
    const params = new URLSearchParams(filters as any);
    return this.request(`/appointments?${params}`);
  }
}

export const api = new ApiClient(API_BASE_URL);
## Custom Hooks
### useAuth Hook

typescript
// hooks/useAuth.ts
import { useAuthStore } from '@/stores/useAuthStore';
import { api } from '@/lib/api';
import { useRouter } from 'next/navigation';

export function useAuth() {
  const router = useRouter();
  const { user, token, login, logout } = useAuthStore();

  const handleLogin = async (email: string, password: string) => {
    try {
      const response = await api.login(email, password);
      login(response.user, response.token);
      
      // Redirect basado en rol
      if (response.user.role === 'client') {
        router.push('/client/dashboard');
      } else {
        router.push('/backoffice/dashboard');
      }
    } catch (error) {
      throw error;
    }
  };

  const handleLogout = () => {
    logout();
    router.push('/login');
  };

  const isAdmin = user?.role === 'admin';
  const isEmployee = user?.role === 'employee';
  const isClient = user?.role === 'client';
  const canAccessBackoffice = isAdmin || isEmployee;

  return {
    user,
    token,
    isAuthenticated: !!token,
    isAdmin,
    isEmployee,
    isClient,
    canAccessBackoffice,
    login: handleLogin,
    logout: handleLogout,
  };
}

### useClients Hook

```typescript
// hooks/useClients.ts
import { useEffect } from 'react';
import { useClientStore } from '@/stores/useClientStore';
import { api } from '@/lib/api';

export function useClients() {
  const {
    clients,
    isLoading,
    error,
    setClients,
    setLoading,
    setError,
    addClient,
    updateClient,
    deleteClient,
  } = useClientStore();

  const fetchClients = async (filters?: ClientFilters) => {
    setLoading(true);
    setError(null);
    try {
      const data = await api.getClients(filters);
      setClients(data.data || data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const createClient = async (clientData: CreateClientData) => {
    try {
      const newClient = await api.createClient(clientData);
      addClient(newClient);
      return newClient;
    } catch (err) {
      throw err;
    }
  };

  const handleUpdateClient = async (id: string, data: Partial<CreateClientData>) => {
    try {
      const updated = await api.updateClient(id, data);
      updateClient(id, updated);
      return updated;
    } catch (err) {
      throw err;
    }
  };

  const handleDeleteClient = async (id: string) => {
    try {
      await api.deleteClient(id);
      deleteClient(id);
    } catch (err) {
      throw err;
    }
  };

  return {
    clients,
    isLoading,
    error,
    fetchClients,
    createClient,
    updateClient: handleUpdateClient,
    deleteClient: handleDeleteClient,
  };
}

## Middleware de Autenticación

```typescript
// middleware.ts (raíz del proyecto)
import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

export function middleware(request: NextRequest) {
  const token = request.cookies.get('auth-storage')?.value;
  // Rutas públicas
  const publicPaths = ['/login', '/register', '/'];
  const isPublicPath = publicPaths.includes(request.nextUrl.pathname);
  // Si no hay token y ruta es privada
  if (!token && !isPublicPath) {
    return NextResponse.redirect(new URL('/login', request.url));
  }
  // Si hay token en rutas auth, redirigir a dashboard
  if (token && (request.nextUrl.pathname === '/login' || request.nextUrl.pathname === '/register')) {
    return NextResponse.redirect(new URL('/backoffice/dashboard', request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    '/((?!api|_next/static|_next/image|favicon.ico).*)',
  ],
};
## Consideraciones

- **App Router**: Usar Server Components por defecto, Client Components solo cuando necesario (`'use client'`)
- **Zustand**: Preferir sobre Context API para estado global
- **Shadcn UI**: Componentes accesibles y customizables
- **TypeScript**: Tipos estrictos, interfaces compartidas con backend (camelCase)
- **Error Handling**: Usar Error Boundaries para componentes
- **Loading States**: Mostrar spinners durante fetch
- **Optimistic Updates**: Actualizar UI antes de confirmar con backend (UX)