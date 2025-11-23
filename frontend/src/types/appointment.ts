// Appointment types matching backend domain models

export type AppointmentStatus = 
  | 'pending' 
  | 'confirmed' 
  | 'cancelled' 
  | 'completed' 
  | 'rescheduled';

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
  therapistId: string;
  title: string;
  description: string;
  startTime: string; // ISO 8601 date string
  endTime: string;   // ISO 8601 date string
  durationMinutes: number;
  status: AppointmentStatus;
  notes: string;
  cancellationReason: string;
  googleCalendarEventId: string;
  createdBy: string;
  createdAt: string;
  updatedAt: string;
  deletedAt?: string;
  client?: Client;
  therapist?: Therapist;
}

export interface CreateAppointmentRequest {
  clientId: string;
  therapistId: string;
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
  therapistId?: string;
}

export interface ConfirmAppointmentRequest {
  notes?: string;
}

export interface CancelAppointmentRequest {
  reason: string;
}

export interface AppointmentFilter {
  clientId?: string;
  therapistId?: string;
  status?: AppointmentStatus;
  startDate?: string; // ISO 8601 date string
  endDate?: string;   // ISO 8601 date string
  page?: number;
  pageSize?: number;
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

export interface GetTherapistsResponse {
  therapists: Therapist[];
}

export interface GetAvailableSlotsResponse {
  slots: string[]; // Array of ISO 8601 date strings
}
