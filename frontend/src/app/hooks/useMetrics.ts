import { useState, useEffect } from 'react';
import { MetricsData } from '../types';

export const useMetrics = (pollingIntervalMs = 5000) => {
  const [metrics, setMetrics] = useState<MetricsData | null>(null);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    let isMounted = true;

    const fetchMetrics = async () => {
      try {
        const response = await fetch('/api/metrics');
        if (!response.ok) {
          throw new Error('Failed to fetch metrics');
        }
        const data = await response.json();
        if (isMounted) {
          setMetrics(data);
          setError(null);
        }
      } catch (err) {
        if (isMounted) {
          setError(err as Error);
        }
      }
    };

    fetchMetrics();
    const interval = setInterval(fetchMetrics, pollingIntervalMs);

    return () => {
      isMounted = false;
      clearInterval(interval);
    };
  }, [pollingIntervalMs]);

  return { metrics, error };
};
