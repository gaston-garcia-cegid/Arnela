// Custom hook for managing appointments
'use client';

import { useState, useCallback } from 'react';
import { api } from '@/lib/api';
import { useAuthStore } from '@/stores/useAuthStore';
import type {
  Appointment,
  CreateAppointmentRequest,
  UpdateAppointmentRequest,
  CancelAppointmentRequest,
  ConfirmAppointmentRequest,
  Therapist,
} from '@/types/appointment';

interface AppointmentFilters {
  clientId?: string;
  therapistId?: string;
  status?: string;
  startDate?: string;
  endDate?: string;
  page?: number;
  pageSize?: number;
}

export function useAppointments() {
  const { token } = useAuthStore();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Get my appointments (client only - uses /appointments/me)
  const getMyAppointments = useCallback(
    async (page: number = 1, pageSize: number = 10) => {
      if (!token) {
        setError('No autenticado');
        return null;
      }

      setLoading(true);
      setError(null);

      try {
        const response = await api.appointments.getMyAppointments(token, page, pageSize);
        return {
          appointments: response.appointments || [],
          total: response.total || 0,
          page: response.page || 1,
          pageSize: response.pageSize || 10,
        };
      } catch (err: any) {
        setError(err.message || 'Error al cargar las citas');
        return null;
      } finally {
        setLoading(false);
      }
    },
    [token]
  );

  // List all appointments with filters (admin/employee - uses /appointments)
  const listAllAppointments = useCallback(
    async (filters?: AppointmentFilters) => {
      if (!token) {
        setError('No autenticado');
        return null;
      }

      setLoading(true);
      setError(null);

      try {
        const response = await api.appointments.list(token, filters);
        return {
          appointments: response.appointments || [],
          total: response.total || 0,
          page: filters?.page || 1,
          pageSize: filters?.pageSize || 10,
        };
      } catch (err: any) {
        setError(err.message || 'Error al cargar las citas');
        return null;
      } finally {
        setLoading(false);
      }
    },
    [token]
  );

  // Get appointment by ID
  const getAppointment = useCallback(
    async (id: string) => {
      if (!token) {
        setError('No autenticado');
        return null;
      }

      setLoading(true);
      setError(null);

      try {
        const appointment = await api.appointments.getById(id, token);
        return appointment;
      } catch (err: any) {
        setError(err.message || 'Error al cargar la cita');
        return null;
      } finally {
        setLoading(false);
      }
    },
    [token]
  );

  // Create appointment
  const createAppointment = useCallback(
    async (data: CreateAppointmentRequest) => {
      if (!token) {
        setError('No autenticado');
        return null;
      }

      setLoading(true);
      setError(null);

      try {
        const appointment = await api.appointments.create(data, token);
        return appointment;
      } catch (err: any) {
        setError(err.message || 'Error al crear la cita');
        return null;
      } finally {
        setLoading(false);
      }
    },
    [token]
  );

  // Update appointment
  const updateAppointment = useCallback(
    async (id: string, data: UpdateAppointmentRequest) => {
      if (!token) {
        setError('No autenticado');
        return null;
      }

      setLoading(true);
      setError(null);

      try {
        const appointment = await api.appointments.update(id, data, token);
        return appointment;
      } catch (err: any) {
        setError(err.message || 'Error al actualizar la cita');
        return null;
      } finally {
        setLoading(false);
      }
    },
    [token]
  );

  // Cancel appointment
  const cancelAppointment = useCallback(
    async (id: string, reason: string) => {
      if (!token) {
        setError('No autenticado');
        return false;
      }

      setLoading(true);
      setError(null);

      try {
        await api.appointments.cancel(id, { reason }, token);
        return true;
      } catch (err: any) {
        setError(err.message || 'Error al cancelar la cita');
        return false;
      } finally {
        setLoading(false);
      }
    },
    [token]
  );

  // Confirm appointment (admin/employee only)
  const confirmAppointment = useCallback(
    async (id: string, notes?: string) => {
      if (!token) {
        setError('No autenticado');
        return null;
      }

      setLoading(true);
      setError(null);

      try {
        await api.appointments.confirm(id, { notes }, token);
        return true;
      } catch (err: any) {
        setError(err.message || 'Error al confirmar la cita');
        return false;
      } finally {
        setLoading(false);
      }
    },
    [token]
  );

  // Get therapists
  const getTherapists = useCallback(async (): Promise<Therapist[]> => {
    if (!token) {
      setError('No autenticado');
      return [];
    }

    setLoading(true);
    setError(null);

    try {
      const response = await api.appointments.getTherapists(token);
      return response.therapists || [];
    } catch (err: any) {
      setError(err.message || 'Error al cargar los terapeutas');
      return [];
    } finally {
      setLoading(false);
    }
  }, [token]);

  // Get available slots
  const getAvailableSlots = useCallback(
    async (therapistId: string, date: string, duration: 45 | 60): Promise<string[]> => {
      if (!token) {
        setError('No autenticado');
        return [];
      }

      setLoading(true);
      setError(null);

      try {
        const response = await api.appointments.getAvailableSlots(token, therapistId, date, duration);
        return response.slots || [];
      } catch (err: any) {
        setError(err.message || 'Error al cargar horarios disponibles');
        return [];
      } finally {
        setLoading(false);
      }
    },
    [token]
  );

  return {
    loading,
    error,
    getMyAppointments,       // ✅ For clients: /appointments/me?page=1&pageSize=10
    listAllAppointments,     // ✅ For admin/employee: /appointments?filters...
    getAppointment,
    createAppointment,
    updateAppointment,
    cancelAppointment,
    confirmAppointment,
    getTherapists,
    getAvailableSlots,
  };
}
