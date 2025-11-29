// Appointment types matching backend domain models

import { Employee } from './employee';

export type AppointmentStatus = 
  | 'pending' 
  | 'confirmed' 
  | 'cancelled' 
  | 'completed' 
  | 'rescheduled';

// Deprecated: Use Employee instead
export interface Therapist {
  id: string;
  name: string;
  specialties: string[];
  isAvailable: boolean;
  avatarColor: string;
}

export interface Client {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  phone: string;
}

export interface Appointment {
  id: string;
  clientId: string;
  employeeId: string; // Changed from therapistId
  title: string;
  description: string;
  startTime: string; // ISO 8601 date string
  endTime: string;   // ISO 8601 date string
  durationMinutes: number;
  status: AppointmentStatus;
  notes?: string;
  cancellationReason?: string;
  googleCalendarEventId?: string;
  createdBy: string;
  createdAt: string;
  updatedAt: string;
  deletedAt?: string;
  client?: Client;
  employee?: Employee; // Changed from therapist
  
  // Deprecated: For backward compatibility
  therapistId?: string;
  therapist?: Therapist;
}

export interface CreateAppointmentRequest {
  clientId?: string; // Optional: for admin/employee creating appointments for others
  employeeId: string; // Changed from therapistId
  title: string;
  description?: string;
  startTime: string; // ISO 8601 date string
  durationMinutes: 45 | 60;
}

export interface UpdateAppointmentRequest {
  title?: string;
  description?: string;
  startTime?: string;
  durationMinutes?: 45 | 60;
  employeeId?: string; // Changed from therapistId
}

export interface ConfirmAppointmentRequest {
  notes?: string;
}

export interface CancelAppointmentRequest {
  reason: string;
}

export interface AppointmentFilter {
  clientId?: string;
  employeeId?: string; // Changed from therapistId
  status?: AppointmentStatus;
  startDate?: string; // ISO 8601 date string
  endDate?: string;   // ISO 8601 date string
  page?: number;
  pageSize?: number;
  
  // Deprecated: For backward compatibility
  therapistId?: string;
}

export interface ListAppointmentsResponse {
  appointments: Appointment[];
  total: number;
  page: number;
  pageSize: number;
}

export interface GetMyAppointmentsResponse {
  appointments: Appointment[];
  total: number;
  page: number;
  pageSize: number;
}

// Deprecated: Use GetEmployeesResponse from employee.ts
export interface GetTherapistsResponse {
  therapists: Therapist[];
}

export interface GetEmployeesResponse {
  employees: Employee[];
}

export interface GetAvailableSlotsResponse {
  slots: string[]; // Array of ISO 8601 date strings
}
