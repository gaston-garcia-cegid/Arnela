// Employee types matching backend domain models

export interface Employee {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  dni: string;
  position?: string | null;        // Backend field (nullable)
  specialties: string[];   // Backend field (array)
  specialty?: string;      // Computed/frontend convenience field
  avatarColor: string;
  isActive: boolean;
  hireDate?: string;       // ISO 8601 date string (optional from backend)
  notes?: string | null;
  createdAt: string;
  updatedAt: string;
}

export interface CreateEmployeeRequest {
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  dni: string;
  specialty: string;
  hireDate: string; // ISO 8601 date string
  avatarColor?: string;
  notes?: string;
}

export interface UpdateEmployeeRequest {
  firstName?: string;
  lastName?: string;
  email?: string;
  phone?: string;
  dni?: string;
  specialty?: string;
  hireDate?: string;
  avatarColor?: string;
  notes?: string;
  isActive?: boolean;
}

export interface ListEmployeesResponse {
  employees: Employee[];
  total: number;
}
