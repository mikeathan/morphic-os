import React from 'react';
import { useMetrics } from '../hooks/useMetrics';

export const HardwareMetrics: React.FC = () => {
  const { metrics, error } = useMetrics();

  const formatMemory = (bytes: number) => {
    return (bytes / 1024 / 1024).toFixed(2) + ' MB';
  };

  return (
    <section className="bg-panel rounded-xl shadow-sm border border-border-default p-6 mt-6">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-xl font-semibold text-text-primary">Hardware & Metrics</h2>
      </div>

      {error ? (
        <div className="text-sm text-red-500">Error loading metrics: {error.message}</div>
      ) : !metrics ? (
        <div className="text-sm text-text-secondary">Loading metrics...</div>
      ) : (
        <div className="grid grid-cols-2 gap-4">
          <div className="bg-primary p-4 rounded-lg border border-border-default">
            <div className="text-xs text-text-secondary uppercase tracking-wider mb-1">Goroutines</div>
            <div className="text-2xl font-semibold text-text-primary">{metrics.goroutines}</div>
          </div>
          <div className="bg-primary p-4 rounded-lg border border-border-default">
            <div className="text-xs text-text-secondary uppercase tracking-wider mb-1">Allocated Mem</div>
            <div className="text-2xl font-semibold text-text-primary">{formatMemory(metrics.allocated_mem)}</div>
          </div>
          <div className="bg-primary p-4 rounded-lg border border-border-default">
            <div className="text-xs text-text-secondary uppercase tracking-wider mb-1">Sys Mem</div>
            <div className="text-2xl font-semibold text-text-primary">{formatMemory(metrics.sys_mem)}</div>
          </div>
          <div className="bg-primary p-4 rounded-lg border border-border-default">
            <div className="text-xs text-text-secondary uppercase tracking-wider mb-1">Pruned Count</div>
            <div className="text-2xl font-semibold text-text-primary">{metrics.pruned_count}</div>
          </div>
        </div>
      )}
    </section>
  );
};
