import { useState, useEffect, useCallback } from 'react';

export interface UseQueryOptions {
  refetchOnMount?: boolean;
}

export interface UseQueryResult<T> {
  data: T | null;
  loading: boolean;
  error: string | null;
  refetch: () => Promise<void>;
}

// I would typically use a library like react-query for this, but as requested I'm avoiding 3rd party libraries
export function useQuery<T>(
  queryFn: () => Promise<T>,
  options: UseQueryOptions = {}
): UseQueryResult<T> {
  const { refetchOnMount = true } = options;

  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const executeQuery = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const result = await queryFn();
      setData(result);
    } catch (error) {
      console.error('Error fetching data:', error);
      const errorMessage = 'An error occurred while fetching data';
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  }, [queryFn]);

  const refetch = useCallback(async () => {
    await executeQuery();
  }, [executeQuery]);

  useEffect(() => {
    if (refetchOnMount) {
      executeQuery();
    }
  }, [executeQuery, refetchOnMount]);

  return {
    data,
    loading,
    error,
    refetch,
  };
}
