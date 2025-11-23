// API Client for Arnela Backend
// Base URL: http://localhost:8080/api/v1

import {
  parseApiError,
  handleFetchError,
  ApiError as ApiErrorClass,
} from './errors';

import type {
  Appointment,
  CreateAppointmentRequest,
  UpdateAppointmentRequest,
  ConfirmAppointmentRequest,
  CancelAppointmentRequest,
  ListAppointmentsResponse,
  GetMyAppointmentsResponse,
  GetTherapistsResponse,
  GetAvailableSlotsResponse,
} from '@/types/appointment';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

// Types matching backend responses
export interface ApiErrorResponse {
  error: string;  // Backend sends "error" field with the message
  code?: string;
  details?: Record<string, string[]>;
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

// Helper function to make authenticated requests with retry logic
async function fetchWithAuth(
  url: string,
  token?: string,
  options: RequestInit = {},
  retries = 3
) {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  let lastError: ApiErrorClass | null = null;

  for (let attempt = 0; attempt < retries; attempt++) {
    try {
      const response = await fetch(`${API_BASE_URL}${url}`, {
        ...options,
        headers,
      });

      if (!response.ok) {
        const errorData: ApiErrorResponse = await response.json().catch(() => ({
          error: 'Unknown error occurred',
        }));
        throw parseApiError(response.status, errorData);
      }

      return response.json();
    } catch (error) {
      lastError = handleFetchError(error);

      // Don't retry on client errors (400-499)
      if (lastError.statusCode >= 400 && lastError.statusCode < 500) {
        throw lastError;
      }

      // On last attempt, throw the error
      if (attempt === retries - 1) {
        throw lastError;
      }

      // Exponential backoff: wait 1s, 2s, 4s...
      await new Promise((resolve) => setTimeout(resolve, Math.pow(2, attempt) * 1000));
    }
  }

  throw lastError || new ApiErrorClass('Unknown error', 500);
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

    search: async (query: string, token: string, isActive: boolean = true): Promise<Client[]> => {
      const queryParams = new URLSearchParams();
      queryParams.append('search', query);
      if (isActive) queryParams.append('isActive', 'true');
      queryParams.append('pageSize', '50'); // Get more results for search
      
      const response = await fetchWithAuth(`/clients?${queryParams.toString()}`, token, {
        method: 'GET',
      }) as ListClientsResponse;
      
      return response.clients || [];
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

  appointments: {
    // Create a new appointment
    create: async (data: CreateAppointmentRequest, token: string): Promise<Appointment> => {
      return fetchWithAuth('/appointments', token, {
        method: 'POST',
        body: JSON.stringify(data),
      });
    },

    // Get appointment by ID
    getById: async (id: string, token: string): Promise<Appointment> => {
      return fetchWithAuth(`/appointments/${id}`, token, {
        method: 'GET',
      });
    },

    // Update appointment
    update: async (id: string, data: UpdateAppointmentRequest, token: string): Promise<Appointment> => {
      return fetchWithAuth(`/appointments/${id}`, token, {
        method: 'PUT',
        body: JSON.stringify(data),
      });
    },

    // Cancel appointment
    cancel: async (id: string, data: CancelAppointmentRequest, token: string): Promise<{ message: string }> => {
      return fetchWithAuth(`/appointments/${id}/cancel`, token, {
        method: 'POST',
        body: JSON.stringify(data),
      });
    },

    // Get my appointments (client only)
    getMyAppointments: async (
      token: string, 
      page: number = 1, 
      pageSize: number = 10
    ): Promise<GetMyAppointmentsResponse> => {
      return fetchWithAuth(`/appointments/me?page=${page}&pageSize=${pageSize}`, token, {
        method: 'GET',
      });
    },

    // Confirm appointment (admin/employee only)
    confirm: async (id: string, data: ConfirmAppointmentRequest, token: string): Promise<Appointment> => {
      return fetchWithAuth(`/appointments/${id}/confirm`, token, {
        method: 'POST',
        body: JSON.stringify(data),
      });
    },

    // List all appointments with filters (admin/employee only)
    list: async (
      token: string,
      filters?: {
        clientId?: string;
        therapistId?: string;
        status?: string;
        startDate?: string;
        endDate?: string;
        page?: number;
        pageSize?: number;
      }
    ): Promise<ListAppointmentsResponse> => {
      const queryParams = new URLSearchParams();
      
      if (filters) {
        if (filters.clientId) queryParams.append('clientId', filters.clientId);
        if (filters.therapistId) queryParams.append('therapistId', filters.therapistId);
        if (filters.status) queryParams.append('status', filters.status);
        if (filters.startDate) queryParams.append('startDate', filters.startDate);
        if (filters.endDate) queryParams.append('endDate', filters.endDate);
        if (filters.page) queryParams.append('page', filters.page.toString());
        if (filters.pageSize) queryParams.append('pageSize', filters.pageSize.toString());
      }

      const queryString = queryParams.toString();
      const url = queryString ? `/appointments?${queryString}` : '/appointments';

      return fetchWithAuth(url, token, {
        method: 'GET',
      });
    },

    // Get available therapists
    getTherapists: async (token: string): Promise<GetTherapistsResponse> => {
      return fetchWithAuth('/appointments/therapists', token, {
        method: 'GET',
      });
    },

    // Get available time slots for a therapist on a specific date
    getAvailableSlots: async (
      token: string,
      therapistId: string,
      date: string, // YYYY-MM-DD format
      duration: 45 | 60
    ): Promise<GetAvailableSlotsResponse> => {
      const queryParams = new URLSearchParams({
        therapistId,
        date,
        duration: duration.toString(),
      });

      return fetchWithAuth(`/appointments/available-slots?${queryParams.toString()}`, token, {
        method: 'GET',
      });
    },
  },
};
