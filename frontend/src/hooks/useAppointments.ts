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

export function useAppointments() {
  const { token } = useAuthStore();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Get my appointments (client)
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
          appointments: response.appointments || [], // ✅ Fallback a array vacío
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
      return response.therapists || []; // ✅ Fallback a array vacío
    } catch (err: any) {
      setError(err.message || 'Error al cargar los terapeutas');
      return []; // ✅ Devolver array vacío en caso de error
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
        return response.slots || []; // ✅ Fallback a array vacío
      } catch (err: any) {
        setError(err.message || 'Error al cargar horarios disponibles');
        return []; // ✅ Devolver array vacío en caso de error
      } finally {
        setLoading(false);
      }
    },
    [token]
  );

  return {
    loading,
    error,
    getMyAppointments,
    getAppointment,
    createAppointment,
    updateAppointment,
    cancelAppointment,
    confirmAppointment,
    getTherapists,
    getAvailableSlots,
  };
}
