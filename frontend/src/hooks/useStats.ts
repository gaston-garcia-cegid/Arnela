import { useState, useEffect } from 'react';
import { api, type DashboardStats } from '@/lib/api';
import { useAuthStore } from '@/stores/useAuthStore';

interface UseStatsReturn {
  stats: DashboardStats | null;
  loading: boolean;
  error: string | null;
  refetch: () => Promise<void>;
}

/**
 * Hook to fetch and manage dashboard statistics
 * Automatically fetches stats on mount if token is available
 */
export function useStats(): UseStatsReturn {
  const token = useAuthStore((state) => state.token);
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchStats = async () => {
    if (!token) {
      setError('No authentication token available');
      setLoading(false);
      return;
    }

    try {
      setLoading(true);
      setError(null);
      const data = await api.stats.getDashboardStats(token);
      setStats(data);
    } catch (err) {
      console.error('Error fetching stats:', err);
      setError(err instanceof Error ? err.message : 'Error al cargar estadísticas');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchStats();
  }, [token]);

  return {
    stats,
    loading,
    error,
    refetch: fetchStats,
  };
}
