import { useState, useCallback, useRef } from 'react';
import { toast } from 'sonner';
import { logger } from '@/lib/logger';

/**
 * Hook for optimistic UI updates
 * 
 * Provides a pattern for updating UI immediately while an async operation
 * is in progress, with automatic rollback on error.
 * 
 * @example
 * ```tsx
 * const { execute, isLoading } = useOptimisticUpdate();
 * 
 * const handleToggleActive = async (clientId: string, currentState: boolean) => {
 *   await execute({
 *     optimisticFn: () => {
 *       // Update UI immediately
 *       setClients(prev => prev.map(c => 
 *         c.id === clientId ? { ...c, isActive: !currentState } : c
 *       ));
 *     },
 *     asyncFn: () => api.clients.update(clientId, { isActive: !currentState }),
 *     rollbackFn: () => {
 *       // Revert UI on error
 *       setClients(prev => prev.map(c => 
 *         c.id === clientId ? { ...c, isActive: currentState } : c
 *       ));
 *     },
 *     successMessage: 'Cliente actualizado',
 *     errorMessage: 'Error al actualizar cliente',
 *   });
 * };
 * ```
 */

interface UseOptimisticUpdateOptions<T> {
  /**
   * Function that updates the UI optimistically (runs immediately)
   */
  optimisticFn: () => void;
  
  /**
   * Async function that performs the actual server operation
   */
  asyncFn: () => Promise<T>;
  
  /**
   * Function to rollback the optimistic update on error
   */
  rollbackFn: () => void;
  
  /**
   * Success message to show in toast (optional)
   */
  successMessage?: string;
  
  /**
   * Error message to show in toast (optional, uses default if not provided)
   */
  errorMessage?: string;
  
  /**
   * Callback to execute after successful operation (optional)
   */
  onSuccess?: (data: T) => void;
  
  /**
   * Callback to execute on error (optional)
   */
  onError?: (error: Error) => void;
  
  /**
   * Whether to show loading indicator while async operation is in progress
   * @default true
   */
  showLoading?: boolean;
}

interface UseOptimisticUpdateReturn {
  /**
   * Execute an optimistic update
   */
  execute: <T>(options: UseOptimisticUpdateOptions<T>) => Promise<T | null>;
  
  /**
   * Whether an optimistic update is currently in progress
   */
  isLoading: boolean;
  
  /**
   * Last error that occurred (null if no error)
   */
  error: Error | null;
}

export function useOptimisticUpdate(): UseOptimisticUpdateReturn {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const loadingToastId = useRef<string | number | null>(null);

  const execute = useCallback(async <T,>(
    options: UseOptimisticUpdateOptions<T>
  ): Promise<T | null> => {
    const {
      optimisticFn,
      asyncFn,
      rollbackFn,
      successMessage,
      errorMessage,
      onSuccess,
      onError,
      showLoading = true,
    } = options;

    try {
      setIsLoading(true);
      setError(null);

      // 1. Apply optimistic update immediately
      optimisticFn();

      // 2. Show loading toast if enabled
      if (showLoading) {
        loadingToastId.current = toast.loading('Guardando...');
      }

      // 3. Execute async operation
      const result = await asyncFn();

      // 4. Dismiss loading toast
      if (loadingToastId.current !== null) {
        toast.dismiss(loadingToastId.current);
        loadingToastId.current = null;
      }

      // 5. Show success message
      if (successMessage) {
        toast.success(successMessage);
      }

      // 6. Execute success callback
      if (onSuccess) {
        onSuccess(result);
      }

      logger.info('[useOptimisticUpdate] Operation succeeded', { successMessage });
      return result;

    } catch (err) {
      // 7. Rollback optimistic update
      rollbackFn();

      // 8. Dismiss loading toast
      if (loadingToastId.current !== null) {
        toast.dismiss(loadingToastId.current);
        loadingToastId.current = null;
      }

      const error = err instanceof Error ? err : new Error('Unknown error');
      setError(error);

      // 9. Show error message
      const finalErrorMessage = errorMessage || 'Ocurrió un error. Por favor, intenta de nuevo.';
      toast.error(finalErrorMessage);

      // 10. Execute error callback
      if (onError) {
        onError(error);
      }

      logger.error('[useOptimisticUpdate] Operation failed', error, { errorMessage });
      return null;

    } finally {
      setIsLoading(false);
    }
  }, []);

  return {
    execute,
    isLoading,
    error,
  };
}
