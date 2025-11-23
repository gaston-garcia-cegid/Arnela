// Zustand store for appointments state management
import { Appointment } from '@/types/appointment';
import { create } from 'zustand';

interface PaginationState {
  page: number;
  pageSize: number;
  total: number;
}

interface AppointmentStore {
  // Current appointments list
  appointments: Appointment[];
  
  // Selected appointment for viewing/editing
  selectedAppointment: Appointment | null;
  
  // Pagination
  pagination: PaginationState;

  // Actions
  setAppointments: (appointments: Appointment[]) => void;
  setSelectedAppointment: (appointment: Appointment | null) => void;
  setPagination: (pagination: Partial<PaginationState>) => void;
  addAppointment: (appointment: Appointment) => void;
  updateAppointment: (id: string, updates: Partial<Appointment>) => void;
  removeAppointment: (id: string) => void;
  clearAppointments: () => void;
}

export const useAppointmentStore = create<AppointmentStore>((set) => ({
  // ✅ Inicializar como array vacío en lugar de null
  appointments: [],
  selectedAppointment: null,
  pagination: {
    page: 1,
    pageSize: 10,
    total: 0,
  },

  setAppointments: (appointments) => set({ appointments }),

  setSelectedAppointment: (appointment) => set({ selectedAppointment: appointment }),

  setPagination: (pagination) =>
    set((state) => ({
      pagination: { ...state.pagination, ...pagination },
    })),

  addAppointment: (appointment) =>
    set((state) => ({
      appointments: [...state.appointments, appointment],
      pagination: { ...state.pagination, total: state.pagination.total + 1 },
    })),

  updateAppointment: (id, updates) =>
    set((state) => ({
      appointments: state.appointments.map((apt) =>
        apt.id === id ? { ...apt, ...updates } : apt
      ),
      selectedAppointment:
        state.selectedAppointment?.id === id
          ? { ...state.selectedAppointment, ...updates }
          : state.selectedAppointment,
    })),

  removeAppointment: (id) =>
    set((state) => ({
      appointments: state.appointments.filter((apt) => apt.id !== id),
      pagination: { ...state.pagination, total: state.pagination.total - 1 },
    })),

  clearAppointments: () => set({ appointments: [], pagination: { page: 1, pageSize: 10, total: 0 } }),
}));
