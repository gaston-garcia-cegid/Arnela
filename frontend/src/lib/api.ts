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
  GetEmployeesResponse,
} from '@/types/appointment';

import type {
  Employee,
  CreateEmployeeRequest,
  UpdateEmployeeRequest,
  ListEmployeesResponse,
} from '@/types/employee';

import type {
  Invoice,
  CreateInvoiceRequest,
  UpdateInvoiceRequest,
  Expense,
  CreateExpenseRequest,
  UpdateExpenseRequest,
  ExpenseCategory,
  CreateExpenseCategoryRequest,
  UpdateExpenseCategoryRequest,
  BillingDashboardStats,
  RevenueByMonth,
  ExpensesByCategory,
  InvoiceFilters,
  ExpenseFilters,
  PaginatedResponse,
} from '@/types/billing';

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
  dniCif: string;
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
  dniCif: string;
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

export interface DashboardStats {
  clients: {
    total: number;
    active: number;
    inactive: number;
  };
  employees: {
    total: number;
    active: number;
    inactive: number;
  };
  appointments: {
    total: number;
    pending: number;
    confirmed: number;
    completed: number;
    cancelled: number;
  };
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

  employees: {
    list: async (token: string, page: number = 1, pageSize: number = 50): Promise<ListEmployeesResponse> => {
      const queryParams = new URLSearchParams();
      queryParams.append('page', page.toString());
      queryParams.append('pageSize', pageSize.toString());
      
      const response = await fetchWithAuth(`/employees?${queryParams.toString()}`, token, {
        method: 'GET',
      }) as ListEmployeesResponse;
      
      // Map backend response to include specialty for backward compatibility
      if (response.employees) {
        response.employees = response.employees.map(emp => ({
          ...emp,
          specialty: emp.position || emp.specialties?.[0] || ''
        }));
      }
      
      return response;
    },

    getById: async (id: string, token: string): Promise<Employee> => {
      const employee = await fetchWithAuth(`/employees/${id}`, token, {
        method: 'GET',
      }) as Employee;
      
      // Map backend response to include specialty for backward compatibility
      return {
        ...employee,
        specialty: employee.position || employee.specialties?.[0] || ''
      };
    },

    search: async (query: string, token: string, isActive: boolean = true): Promise<Employee[]> => {
      const queryParams = new URLSearchParams();
      queryParams.append('search', query);
      if (isActive) queryParams.append('isActive', 'true');
      queryParams.append('pageSize', '50');
      
      const response = await fetchWithAuth(`/employees?${queryParams.toString()}`, token, {
        method: 'GET',
      }) as ListEmployeesResponse;
      
      // Map backend response to include specialty for backward compatibility
      const employees = (response.employees || []).map(emp => ({
        ...emp,
        specialty: emp.position || emp.specialties?.[0] || ''
      }));
      
      return employees;
    },

    create: async (data: CreateEmployeeRequest, token: string): Promise<Employee> => {
      const employee = await fetchWithAuth('/employees', token, {
        method: 'POST',
        body: JSON.stringify(data),
      }) as Employee;
      
      // Map backend response to include specialty for backward compatibility
      return {
        ...employee,
        specialty: employee.position || employee.specialties?.[0] || ''
      };
    },

    update: async (id: string, data: UpdateEmployeeRequest, token: string): Promise<Employee> => {
      const employee = await fetchWithAuth(`/employees/${id}`, token, {
        method: 'PUT',
        body: JSON.stringify(data),
      }) as Employee;
      
      // Map backend response to include specialty for backward compatibility
      return {
        ...employee,
        specialty: employee.position || employee.specialties?.[0] || ''
      };
    },

    delete: async (id: string, token: string): Promise<void> => {
      return fetchWithAuth(`/employees/${id}`, token, {
        method: 'DELETE',
      });
    },

    // Get my employee profile (for logged-in employee)
    getMyProfile: async (token: string): Promise<Employee> => {
      const employee = await fetchWithAuth('/employees/me', token, {
        method: 'GET',
      }) as Employee;
      
      // Map backend response to include specialty for backward compatibility
      return {
        ...employee,
        specialty: employee.position || employee.specialties?.[0] || ''
      };
    },

    // Get active employees for appointment scheduling
    getActive: async (token: string): Promise<GetEmployeesResponse> => {
      const queryParams = new URLSearchParams();
      queryParams.append('isActive', 'true');
      queryParams.append('pageSize', '100');
      
      const response = await fetchWithAuth(`/employees?${queryParams.toString()}`, token, {
        method: 'GET',
      }) as ListEmployeesResponse;
      
      // Map backend response to include specialty for backward compatibility
      const employees = (response.employees || []).map(emp => ({
        ...emp,
        specialty: emp.position || emp.specialties?.[0] || ''
      }));
      
      return { employees };
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
        employeeId?: string; // Changed from therapistId
        status?: string;
        startDate?: string;
        endDate?: string;
        page?: number;
        pageSize?: number;
        therapistId?: string; // Deprecated: for backward compatibility
      }
    ): Promise<ListAppointmentsResponse> => {
      const queryParams = new URLSearchParams();
      
      if (filters) {
        if (filters.clientId) queryParams.append('clientId', filters.clientId);
        // Use employeeId if provided, fallback to therapistId for backward compatibility
        const employeeId = filters.employeeId || filters.therapistId;
        if (employeeId) queryParams.append('employeeId', employeeId);
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

    // Get available employees for appointments
    getEmployees: async (token: string): Promise<GetEmployeesResponse> => {
      return api.employees.getActive(token);
    },

    // Deprecated: Use getEmployees instead
    getTherapists: async (token: string): Promise<GetTherapistsResponse> => {
      return fetchWithAuth('/appointments/therapists', token, {
        method: 'GET',
      });
    },

    // Get available time slots for an employee on a specific date
    getAvailableSlots: async (
      token: string,
      employeeId: string, // Changed from therapistId
      date: string, // YYYY-MM-DD format
      duration: 45 | 60
    ): Promise<GetAvailableSlotsResponse> => {
      const queryParams = new URLSearchParams({
        employeeId,
        date,
        duration: duration.toString(),
      });

      return fetchWithAuth(`/appointments/available-slots?${queryParams.toString()}`, token, {
        method: 'GET',
      });
    },
  },

  stats: {
    getDashboardStats: async (token: string): Promise<DashboardStats> => {
      return fetchWithAuth('/stats/dashboard', token, {
        method: 'GET',
      });
    },
  },

  billing: {
    // ============================================
    // Invoice Operations
    // ============================================
    invoices: {
      /**
       * Creates a new invoice with automatic VAT calculation
       * @param data Invoice creation data including client, dates, and amounts
       * @param token Authorization token
       * @returns Created invoice with calculated totals
       */
      create: async (data: CreateInvoiceRequest, token: string): Promise<Invoice> => {
        return fetchWithAuth('/billing/invoices', token, {
          method: 'POST',
          body: JSON.stringify(data),
        });
      },

      /**
       * Retrieves paginated list of invoices with optional filters
       * @param token Authorization token
       * @param filters Optional filters for client, status, date range, pagination
       * @returns Paginated invoice list
       */
      list: async (token: string, filters?: InvoiceFilters): Promise<PaginatedResponse<Invoice>> => {
        if (!filters) {
          return fetchWithAuth('/billing/invoices', token, { method: 'GET' });
        }

        const queryParams = buildQueryParams(filters);
        return fetchWithAuth(`/billing/invoices?${queryParams}`, token, { method: 'GET' });
      },

      /**
       * Gets a single invoice by ID
       * @param id Invoice UUID
       * @param token Authorization token
       * @returns Invoice details
       */
      getById: async (id: string, token: string): Promise<Invoice> => {
        return fetchWithAuth(`/billing/invoices/${id}`, token, { method: 'GET' });
      },

      /**
       * Gets a single invoice by invoice number
       * @param invoiceNumber Invoice number (e.g., F_2025_0001)
       * @param token Authorization token
       * @returns Invoice details
       */
      getByNumber: async (invoiceNumber: string, token: string): Promise<Invoice> => {
        return fetchWithAuth(`/billing/invoices/number/${invoiceNumber}`, token, {
          method: 'GET',
        });
      },

      /**
       * Gets all invoices for a specific client
       * @param clientId Client UUID
       * @param token Authorization token
       * @returns Array of client's invoices
       */
      getByClient: async (clientId: string, token: string): Promise<Invoice[]> => {
        return fetchWithAuth(`/billing/invoices/client/${clientId}`, token, {
          method: 'GET',
        });
      },

      /**
       * Gets all unpaid invoices
       * @param token Authorization token
       * @returns Array of unpaid invoices
       */
      getUnpaid: async (token: string): Promise<Invoice[]> => {
        return fetchWithAuth('/billing/invoices/unpaid', token, { method: 'GET' });
      },

      /**
       * Updates an existing invoice
       * @param id Invoice UUID
       * @param data Updated invoice data
       * @param token Authorization token
       * @returns Updated invoice with recalculated totals
       */
      update: async (
        id: string,
        data: UpdateInvoiceRequest,
        token: string
      ): Promise<Invoice> => {
        return fetchWithAuth(`/billing/invoices/${id}`, token, {
          method: 'PUT',
          body: JSON.stringify(data),
        });
      },

      /**
       * Marks an invoice as paid
       * @param id Invoice UUID
       * @param token Authorization token
       * @returns Updated invoice with paid status
       */
      markAsPaid: async (id: string, token: string): Promise<Invoice> => {
        return fetchWithAuth(`/billing/invoices/${id}/mark-paid`, token, {
          method: 'POST',
        });
      },

      /**
       * Soft deletes an invoice (sets deleted_at timestamp)
       * @param id Invoice UUID
       * @param token Authorization token
       */
      delete: async (id: string, token: string): Promise<void> => {
        return fetchWithAuth(`/billing/invoices/${id}`, token, { method: 'DELETE' });
      },
    },

    // ============================================
    // Expense Operations
    // ============================================
    expenses: {
      /**
       * Creates a new expense record
       * @param data Expense data including supplier, amount, and category
       * @param token Authorization token
       * @returns Created expense
       */
      create: async (data: CreateExpenseRequest, token: string): Promise<Expense> => {
        return fetchWithAuth('/billing/expenses', token, {
          method: 'POST',
          body: JSON.stringify(data),
        });
      },

      /**
       * Retrieves paginated list of expenses with optional filters
       * @param token Authorization token
       * @param filters Optional filters for category, date range, invoice status
       * @returns Paginated expense list
       */
      list: async (token: string, filters?: ExpenseFilters): Promise<PaginatedResponse<Expense>> => {
        if (!filters) {
          return fetchWithAuth('/billing/expenses', token, { method: 'GET' });
        }

        const queryParams = buildQueryParams(filters);
        return fetchWithAuth(`/billing/expenses?${queryParams}`, token, { method: 'GET' });
      },

      /**
       * Gets a single expense by ID
       * @param id Expense UUID
       * @param token Authorization token
       * @returns Expense details with populated category
       */
      getById: async (id: string, token: string): Promise<Expense> => {
        return fetchWithAuth(`/billing/expenses/${id}`, token, { method: 'GET' });
      },

      /**
       * Updates an existing expense
       * @param id Expense UUID
       * @param data Updated expense data
       * @param token Authorization token
       * @returns Updated expense
       */
      update: async (
        id: string,
        data: UpdateExpenseRequest,
        token: string
      ): Promise<Expense> => {
        return fetchWithAuth(`/billing/expenses/${id}`, token, {
          method: 'PUT',
          body: JSON.stringify(data),
        });
      },

      /**
       * Soft deletes an expense (sets deleted_at timestamp)
       * @param id Expense UUID
       * @param token Authorization token
       */
      delete: async (id: string, token: string): Promise<void> => {
        return fetchWithAuth(`/billing/expenses/${id}`, token, { method: 'DELETE' });
      },
    },

    // ============================================
    // Expense Category Operations
    // ============================================
    categories: {
      /**
       * Creates a new expense category (parent or subcategory)
       * @param data Category data with name, code, and optional parentId
       * @param token Authorization token
       * @returns Created category
       */
      create: async (
        data: CreateExpenseCategoryRequest,
        token: string
      ): Promise<ExpenseCategory> => {
        return fetchWithAuth('/billing/expense-categories', token, {
          method: 'POST',
          body: JSON.stringify(data),
        });
      },

      /**
       * Gets all categories (flat list)
       * @param token Authorization token
       * @returns Array of all categories
       */
      list: async (token: string): Promise<ExpenseCategory[]> => {
        return fetchWithAuth('/billing/expense-categories', token, { method: 'GET' });
      },

      /**
       * Gets hierarchical tree of all categories with their subcategories
       * @param token Authorization token
       * @returns Array of parent categories with nested subcategories
       */
      getTree: async (token: string): Promise<ExpenseCategory[]> => {
        return fetchWithAuth('/billing/expense-categories/tree', token, { method: 'GET' });
      },

      /**
       * Gets only parent categories (no subcategories)
       * @param token Authorization token
       * @returns Array of parent categories
       */
      getParents: async (token: string): Promise<ExpenseCategory[]> => {
        return fetchWithAuth('/billing/expense-categories/parents', token, {
          method: 'GET',
        });
      },

      /**
       * Gets a single category by ID
       * @param id Category UUID
       * @param token Authorization token
       * @returns Category details
       */
      getById: async (id: string, token: string): Promise<ExpenseCategory> => {
        return fetchWithAuth(`/billing/expense-categories/${id}`, token, {
          method: 'GET',
        });
      },

      /**
       * Gets all subcategories of a parent category
       * @param parentId Parent category UUID
       * @param token Authorization token
       * @returns Array of subcategories
       */
      getSubcategories: async (parentId: string, token: string): Promise<ExpenseCategory[]> => {
        return fetchWithAuth(
          `/billing/expense-categories/${parentId}/subcategories`,
          token,
          { method: 'GET' }
        );
      },

      /**
       * Updates an existing category
       * @param id Category UUID
       * @param data Updated category data
       * @param token Authorization token
       * @returns Updated category
       */
      update: async (
        id: string,
        data: UpdateExpenseCategoryRequest,
        token: string
      ): Promise<ExpenseCategory> => {
        return fetchWithAuth(`/billing/expense-categories/${id}`, token, {
          method: 'PUT',
          body: JSON.stringify(data),
        });
      },

      /**
       * Deletes a category (cannot delete if it has expenses or subcategories)
       * @param id Category UUID
       * @param token Authorization token
       */
      delete: async (id: string, token: string): Promise<void> => {
        return fetchWithAuth(`/billing/expense-categories/${id}`, token, {
          method: 'DELETE',
        });
      },
    },

    // ============================================
    // Billing Statistics & Reports
    // ============================================
    stats: {
      /**
       * Gets comprehensive billing dashboard statistics
       * @param token Authorization token
       * @returns Dashboard stats with totals and counts
       */
      getDashboard: async (token: string): Promise<BillingDashboardStats> => {
        return fetchWithAuth('/billing/dashboard', token, { method: 'GET' });
      },

      /**
       * Gets revenue breakdown by month for a specific year
       * @param token Authorization token
       * @param year Year to filter (e.g., 2025)
       * @returns Array of monthly revenue data
       */
      getRevenueByMonth: async (token: string, year: number): Promise<RevenueByMonth[]> => {
        const queryParams = new URLSearchParams({ year: year.toString() });
        return fetchWithAuth(`/billing/revenue-by-month?${queryParams}`, token, {
          method: 'GET',
        });
      },

      /**
       * Gets expense breakdown by category
       * @param token Authorization token
       * @param year Optional year filter
       * @returns Array of expense totals per category
       */
      getExpensesByCategory: async (
        token: string,
        year?: number
      ): Promise<ExpensesByCategory[]> => {
        if (!year) {
          return fetchWithAuth('/billing/expenses-by-category', token, { method: 'GET' });
        }

        const queryParams = new URLSearchParams({ year: year.toString() });
        return fetchWithAuth(`/billing/expenses-by-category?${queryParams}`, token, {
          method: 'GET',
        });
      },

      /**
       * Gets current balance (total revenue - total expenses)
       * @param token Authorization token
       * @returns Balance data with revenue and expense totals
       */
      getBalance: async (
        token: string
      ): Promise<{ balance: number; totalRevenue: number; totalExpenses: number }> => {
        return fetchWithAuth('/billing/balance', token, { method: 'GET' });
      },
    },
  },

  // ============================================
  // Global Search
  // ============================================
  search: {
    /**
     * Global search across clients, employees, appointments, and invoices
     * @param query Search query (minimum 2 characters)
     * @param token Authorization token
     * @returns Search results grouped by entity type
     */
    global: async (query: string, token: string): Promise<{
      clients: Array<{
        id: string;
        firstName: string;
        lastName: string;
        email: string;
        phone: string;
        dniCif: string;
      }>;
      employees: Array<{
        id: string;
        name: string;
        email: string;
        phone: string;
        specialties: string[];
        avatarColor: string;
      }>;
      appointments: Array<{
        id: string;
        title: string;
        startTime: string;
        endTime: string;
        status: string;
        clientName: string;
        employeeName: string;
      }>;
      invoices: Array<{
        id: string;
        invoiceNumber: string;
        clientName: string;
        totalAmount: number;
        status: string;
        issueDate: string;
      }>;
      totalResults: number;
    }> => {
      const queryParams = new URLSearchParams({ q: query });
      return fetchWithAuth(`/search?${queryParams}`, token, { method: 'GET' });
    },
  },
};

// ============================================
// Helper Functions
// ============================================

/**
 * Builds URL query parameters from filter object
 * Applies early return pattern for cleaner code
 * @param filters Object with optional filter parameters
 * @returns URL-encoded query string
 */
function buildQueryParams(filters: Record<string, any>): string {
  const params = new URLSearchParams();

  for (const [key, value] of Object.entries(filters)) {
    if (value === undefined || value === null) continue;
    if (value === '') continue;
    
    // Handle arrays and objects properly
    if (typeof value === 'object') {
      params.append(key, JSON.stringify(value));
    } else {
      params.append(key, String(value));
    }
  }

  return params.toString();
}
