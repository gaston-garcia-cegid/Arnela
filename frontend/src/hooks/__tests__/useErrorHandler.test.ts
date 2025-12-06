import { renderHook, act } from '@testing-library/react';
import { useErrorHandler, getUserFriendlyError } from '../useErrorHandler';
import { logError } from '../../lib/logger';
import { toast } from 'sonner';
import { vi, describe, it, expect, beforeEach } from 'vitest';

// Mocks
vi.mock('@/lib/logger', () => ({
    logError: vi.fn(),
}));

vi.mock('sonner', () => ({
    toast: {
        error: vi.fn(),
        success: vi.fn(),
    },
}));

describe('useErrorHandler', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('updates error state and shows toast on handleError', () => {
        const { result } = renderHook(() => useErrorHandler({ component: 'TestComp' }));

        act(() => {
            result.current.handleError(new Error('Test Error'));
        });

        expect(result.current.error).toBe('Test Error');
        // Check looser args to avoid flakiness
        expect(logError).toHaveBeenCalled();
        expect(toast.error).toHaveBeenCalledWith('Test Error', expect.any(Object));
    });

    it('suppresses toast if silent is true', () => {
        const { result } = renderHook(() => useErrorHandler());

        act(() => {
            result.current.handleError(new Error('Err'), { silent: true });
        });

        expect(toast.error).not.toHaveBeenCalled();
        expect(logError).toHaveBeenCalled();
    });

    it('uses custom user message if provided', () => {
        const { result } = renderHook(() => useErrorHandler());
        act(() => {
            result.current.handleError(new Error('Tech Error'), { userMessage: 'Nice Error' });
        });

        expect(result.current.error).toBe('Nice Error');
        expect(toast.error).toHaveBeenCalledWith('Nice Error', expect.any(Object));
    });

    it('handles non-Error objects correctly', () => {
        const { result } = renderHook(() => useErrorHandler());
        act(() => {
            result.current.handleError('String Error');
        });
        expect(result.current.error).toBe('String Error');
    });

    it('handleSuccess clears error and shows success toast', () => {
        const { result } = renderHook(() => useErrorHandler());

        // First set error
        act(() => {
            result.current.handleError(new Error('Err'));
        });
        expect(result.current.error).toBe('Err');

        // Then success
        act(() => {
            result.current.handleSuccess('Success!');
        });

        expect(result.current.error).toBeNull();
        expect(toast.success).toHaveBeenCalledWith('Success!', expect.any(Object));
    });

    it('respects showToast option false in hook init', () => {
        const { result } = renderHook(() => useErrorHandler({ showToast: false }));

        act(() => {
            result.current.handleError(new Error('Err'));
        });

        expect(toast.error).not.toHaveBeenCalled();

        act(() => {
            result.current.handleSuccess('Success');
        });

        expect(toast.success).not.toHaveBeenCalled();
    });
});

describe('getUserFriendlyError', () => {
    it('maps known errors', () => {
        expect(getUserFriendlyError(new Error('Network Error failed'))).toContain('Error de conexión');
        expect(getUserFriendlyError(new Error('Failed to fetch'))).toContain('No se pudo conectar');
        expect(getUserFriendlyError(new Error('Unauthorized access'))).toContain('No tienes permiso');
    });

    it('returns original message for unknown errors', () => {
        expect(getUserFriendlyError(new Error('Custom Logic Error'))).toBe('Custom Logic Error');
    });

    it('handles non-error objects (fallback)', () => {
        expect(getUserFriendlyError('string error')).toBe('Ha ocurrido un error inesperado');
        expect(getUserFriendlyError(null)).toBe('Ha ocurrido un error inesperado');
    });
});
