import { renderHook, act } from '@testing-library/react';
import { useAppointments } from '../useAppointments';
import { api } from '@/lib/api';
import { useAuthStore } from '@/stores/useAuthStore';
import { vi, describe, it, expect, beforeEach } from 'vitest';

// Mocks
vi.mock('@/lib/api', () => ({
    api: {
        appointments: {
            getMyAppointments: vi.fn(),
            list: vi.fn(),
            getById: vi.fn(),
            create: vi.fn(),
            update: vi.fn(),
            cancel: vi.fn(),
            confirm: vi.fn(),
            getEmployees: vi.fn(),
            getAvailableSlots: vi.fn(),
        },
        clients: {
            search: vi.fn(),
        }
    }
}));

vi.mock('@/stores/useAuthStore', () => ({
    useAuthStore: vi.fn(),
}));

describe('useAppointments Hook', () => {
    // Default mock setup
    const mockToken = 'valid-token';

    beforeEach(() => {
        vi.clearAllMocks();
        (useAuthStore as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
            token: mockToken,
            user: { role: 'client' },
        });
    });

    describe('getMyAppointments', () => {
        it('fetches appointments successfully', async () => {
            const mockData = { appointments: [{ id: 1 }], total: 1 };
            (api.appointments.getMyAppointments as any).mockResolvedValue(mockData);

            const { result } = renderHook(() => useAppointments());

            let response;
            await act(async () => {
                response = await result.current.getMyAppointments(1, 10);
            });

            expect(response).toEqual({
                appointments: [{ id: 1 }],
                total: 1,
                page: 1,
                pageSize: 10
            });
            expect(api.appointments.getMyAppointments).toHaveBeenCalledWith(mockToken, 1, 10);
            expect(result.current.loading).toBe(false);
            expect(result.current.error).toBeNull();
        });

        it('handles error state', async () => {
            (api.appointments.getMyAppointments as any).mockRejectedValue(new Error('API Error'));

            const { result } = renderHook(() => useAppointments());

            await act(async () => {
                await result.current.getMyAppointments();
            });

            expect(result.current.error).toBe('API Error');
            expect(result.current.loading).toBe(false);
        });

        it('handles missing token (unauthenticated)', async () => {
            (useAuthStore as unknown as ReturnType<typeof vi.fn>).mockReturnValue({ token: null });
            const { result } = renderHook(() => useAppointments());

            let response;
            await act(async () => {
                response = await result.current.getMyAppointments();
            });

            expect(response).toBeNull();
            expect(result.current.error).toBe('No autenticado');
            expect(api.appointments.getMyAppointments).not.toHaveBeenCalled();
        });
    });

    describe('createAppointment', () => {
        it('creates an appointment successfully', async () => {
            const newAppt = { date: '2023-01-01' };
            const createdAppt = { id: 1, ...newAppt };
            (api.appointments.create as any).mockResolvedValue(createdAppt);

            const { result } = renderHook(() => useAppointments());

            let response;
            await act(async () => {
                response = await result.current.createAppointment(newAppt as any);
            });

            expect(response).toEqual(createdAppt);
            expect(result.current.loading).toBe(false);
        });

        it('handles create error', async () => {
            (api.appointments.create as any).mockRejectedValue(new Error('Create Failed'));
            const { result } = renderHook(() => useAppointments());

            await act(async () => {
                await result.current.createAppointment({} as any);
            });

            expect(result.current.error).toBe('Create Failed');
        });
    });

    describe('cancelAppointment', () => {
        it('cancels appointment successfully', async () => {
            (api.appointments.cancel as any).mockResolvedValue({});
            const { result } = renderHook(() => useAppointments());

            let success;
            await act(async () => {
                success = await result.current.cancelAppointment('1', 'reason');
            });

            expect(success).toBe(true);
        });

        it('handles cancellation error', async () => {
            (api.appointments.cancel as any).mockRejectedValue(new Error('Fail'));
            const { result } = renderHook(() => useAppointments());

            let success;
            await act(async () => {
                success = await result.current.cancelAppointment('1', 'reason');
            });

            expect(success).toBe(false);
            expect(result.current.error).toBe('Fail');
        });
    });

    describe('getAvailableSlots', () => {
        it('returns slots on success', async () => {
            (api.appointments.getAvailableSlots as any).mockResolvedValue({ slots: ['10:00'] });
            const { result } = renderHook(() => useAppointments());

            let slots;
            await act(async () => {
                slots = await result.current.getAvailableSlots('emp1', 'date', 60);
            });

            expect(slots).toEqual(['10:00']);
        });

        it('returns empty array on error', async () => {
            (api.appointments.getAvailableSlots as any).mockRejectedValue(new Error('Fail'));
            const { result } = renderHook(() => useAppointments());

            let slots;
            await act(async () => {
                slots = await result.current.getAvailableSlots('emp1', 'date', 60);
            });

            expect(slots).toEqual([]);
            expect(result.current.error).toBe('Fail');
        });
    });

    // Additional Coverage for listAllAppointments
    describe('listAllAppointments', () => {
        it('lists appointments with filters', async () => {
            const mockData = { appointments: [], total: 0 };
            (api.appointments.list as any).mockResolvedValue(mockData);

            const { result } = renderHook(() => useAppointments());

            let response;
            await act(async () => {
                response = await result.current.listAllAppointments({ page: 2 });
            });

            expect(response?.page).toBe(2);
            expect(api.appointments.list).toHaveBeenCalledWith(mockToken, { page: 2 });
        });
    });
});
