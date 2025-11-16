// API Client for Arnela Backend
// Base URL: http://localhost:8080/api/v1

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

// Types matching backend responses
export interface ApiError {
  error: string;
  details?: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
  dni?: string;
  phoneNumber?: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface AuthResponse {
  token: string;
  user: {
    id: string;
    email: string;
    firstName: string;
    lastName: string;
    role: 'admin' | 'employee' | 'client';
    isActive: boolean;
  };
}

export interface Client {
  id: string;
  userId?: string;
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  nif: string;
  dni?: string;
  address?: string;
  city?: string;
  postalCode?: string;
  province?: string;
  dateOfBirth?: string;
  isActive: boolean;
  notes?: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateClientRequest {
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  nif: string;
  dni?: string;
  address?: string;
  city?: string;
  postalCode?: string;
  province?: string;
  dateOfBirth?: string;
  notes?: string;
}

export interface ListClientsResponse {
  clients: Client[];
  total: number;
}

// Helper function to make authenticated requests
async function fetchWithAuth(url: string, token?: string, options: RequestInit = {}) {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_BASE_URL}${url}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    const error: ApiError = await response.json().catch(() => ({
      error: 'Unknown error occurred',
    }));
    throw new Error(error.error || `HTTP ${response.status}: ${response.statusText}`);
  }

  return response.json();
}

// Auth endpoints
export const api = {
  auth: {
    register: async (data: RegisterRequest): Promise<AuthResponse> => {
      return fetchWithAuth('/auth/register', undefined, {
        method: 'POST',
        body: JSON.stringify(data),
      });
    },

    login: async (data: LoginRequest): Promise<AuthResponse> => {
      return fetchWithAuth('/auth/login', undefined, {
        method: 'POST',
        body: JSON.stringify(data),
      });
    },

    getMe: async (token: string): Promise<AuthResponse['user']> => {
      return fetchWithAuth('/auth/me', token, {
        method: 'GET',
      });
    },
  },

  clients: {
    list: async (token: string): Promise<ListClientsResponse> => {
      return fetchWithAuth('/clients', token, {
        method: 'GET',
      });
    },

    getById: async (id: string, token: string): Promise<Client> => {
      return fetchWithAuth(`/clients/${id}`, token, {
        method: 'GET',
      });
    },

    create: async (data: CreateClientRequest, token: string): Promise<Client> => {
      return fetchWithAuth('/clients', token, {
        method: 'POST',
        body: JSON.stringify(data),
      });
    },

    update: async (id: string, data: Partial<Client>, token: string): Promise<Client> => {
      return fetchWithAuth(`/clients/${id}`, token, {
        method: 'PUT',
        body: JSON.stringify(data),
      });
    },

    delete: async (id: string, token: string): Promise<void> => {
      return fetchWithAuth(`/clients/${id}`, token, {
        method: 'DELETE',
      });
    },
  },
};
