import { renderHook, act, waitFor } from '@testing-library/react';
import { useStats } from '../useStats';
import { api } from '@/lib/api';
import { useAuthStore } from '@/stores/useAuthStore';
import { logError } from '@/lib/logger';
import { vi, describe, it, expect, beforeEach } from 'vitest';

// Mocks
vi.mock('@/lib/api', () => ({
    api: {
        stats: {
            getDashboardStats: vi.fn(),
        },
    },
}));

vi.mock('@/stores/useAuthStore', () => ({
    useAuthStore: vi.fn(),
}));

vi.mock('@/lib/logger', () => ({
    logError: vi.fn(),
}));

describe('useStats', () => {
    const mockToken = 'token';

    beforeEach(() => {
        vi.clearAllMocks();
        (useAuthStore as unknown as ReturnType<typeof vi.fn>).mockImplementation((selector) => {
            return selector({ token: mockToken });
        });
    });

    it('fetches stats on mount', async () => {
        const mockStats = { clients: 10, appointments: 5 };
        (api.stats.getDashboardStats as any).mockResolvedValue(mockStats);

        const { result } = renderHook(() => useStats());

        expect(result.current.loading).toBe(true);

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.stats).toEqual(mockStats);
        expect(result.current.error).toBeNull();
    });

    it('handles errors', async () => {
        (api.stats.getDashboardStats as any).mockRejectedValue(new Error('Fetch Error'));

        const { result } = renderHook(() => useStats());

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.error).toBe('Fetch Error');
        expect(logError).toHaveBeenCalled();
    });

    it('handles missing token', async () => {
        (useAuthStore as unknown as ReturnType<typeof vi.fn>).mockImplementation((selector) => {
            return selector({ token: null });
        });

        const { result } = renderHook(() => useStats());

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.error).toContain('No authentication token');
        expect(api.stats.getDashboardStats).not.toHaveBeenCalled();
    });

    it('refetches on demand', async () => {
        const mockStats = { clients: 10 };
        (api.stats.getDashboardStats as any).mockResolvedValue(mockStats);

        const { result } = renderHook(() => useStats());

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        // Update data logic
        const newStats = { clients: 11 };
        (api.stats.getDashboardStats as any).mockResolvedValue(newStats);

        await act(async () => {
            await result.current.refetch();
        });

        expect(result.current.stats).toEqual(newStats);
    });
});
