/**
 * Custom Hook for Error Handling
 * Provides consistent error handling and user feedback
 */

import { useState, useCallback } from 'react';
import { toast } from 'sonner';
import { logError } from '@/lib/logger';

interface UseErrorHandlerOptions {
    component?: string;
    showToast?: boolean;
    toastDuration?: number;
}

export function useErrorHandler(options: UseErrorHandlerOptions = {}) {
    const {
        component = 'Unknown',
        showToast = true,
        toastDuration = 5000,
    } = options;

    const [error, setError] = useState<string | null>(null);

    /**
     * Handle error with logging and user feedback
     */
    const handleError = useCallback((
        error: Error | unknown,
        context?: {
            action?: string;
            userMessage?: string;
            silent?: boolean;
        }
    ) => {
        const errorMessage = error instanceof Error ? error.message : String(error);
        const userMessage = context?.userMessage || errorMessage;

        // Log error (only in development)
        logError(
            `Error in ${component}${context?.action ? ` - ${context.action}` : ''}`,
            error,
            {
                component,
                action: context?.action,
                timestamp: new Date().toISOString(),
            }
        );

        // Set error state
        setError(userMessage);

        // Show toast notification (unless silent)
        if (showToast && !context?.silent) {
            toast.error(userMessage, {
                duration: toastDuration,
            });
        }

        return userMessage;
    }, [component, showToast, toastDuration]);

    /**
     * Clear error state
     */
    const clearError = useCallback(() => {
        setError(null);
    }, []);

    /**
     * Handle success with user feedback
     */
    const handleSuccess = useCallback((message: string, duration = 3000) => {
        clearError();
        if (showToast) {
            toast.success(message, { duration });
        }
    }, [showToast, clearError]);

    return {
        error,
        handleError,
        clearError,
        handleSuccess,
    };
}

/**
 * Get user-friendly error message from API error
 */
export function getUserFriendlyError(error: Error | unknown): string {
    if (error instanceof Error) {
        // Map common error messages to user-friendly ones
        const errorMappings: Record<string, string> = {
            'Network Error': 'Error de conexión. Por favor, verifica tu internet',
            'Failed to fetch': 'No se pudo conectar con el servidor',
            'Unauthorized': 'No tienes permiso para realizar esta acción',
            'Forbidden': 'Acceso denegado',
            'Not Found': 'Recurso no encontrado',
            'Validation Error': 'Error en los datos ingresados',
            'Internal Server Error': 'Error del servidor. Intenta nuevamente',
        };

        for (const [key, message] of Object.entries(errorMappings)) {
            if (error.message.includes(key)) {
                return message;
            }
        }

        return error.message;
    }

    return 'Ha ocurrido un error inesperado';
}
