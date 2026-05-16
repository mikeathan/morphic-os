import React, { useEffect, useState } from 'react';

interface MemoryVector {
  id: string;
  workspace_id: string;
  content: string;
  access_frequency: number;
  last_recall: string;
  core_memory: boolean;
  created_at: string;
  updated_at: string;
}

export const MemoryManager: React.FC = () => {
  const [memories, setMemories] = useState<MemoryVector[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchMemories = async () => {
      try {
        const res = await fetch('/api/memory?workspace_id=default');
        if (!res.ok) {
          throw new Error('Failed to fetch memories');
        }
        const data = await res.json();
        setMemories(data || []);
      } catch (err) {
        setError(err instanceof Error ? err.message : String(err));
      } finally {
        setLoading(false);
      }
    };

    fetchMemories();
  }, []);

  return (
    <section className="bg-panel rounded-xl shadow-sm border border-border-default p-6 mt-6">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-xl font-semibold text-text-primary">Cognitive Memory</h2>
      </div>

      {loading ? (
        <div className="text-sm text-text-secondary">Loading memories...</div>
      ) : error ? (
        <div className="text-sm text-red-500">Error: {error}</div>
      ) : (
        <div className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div className="bg-primary p-4 rounded-lg border border-border-default">
               <div className="text-xs text-text-secondary uppercase tracking-wider mb-1">Total Memories</div>
               <div className="text-2xl font-semibold text-text-primary">{memories.length}</div>
            </div>
            {/* Pruning stats are intentionally left to HardwareMetrics as it already tracks 'Pruned Count' */}
          </div>

          <div>
             <h3 className="text-sm font-semibold text-text-primary mb-2">Extracted Facts & Preferences</h3>
             {memories.length === 0 ? (
               <div className="text-sm text-text-secondary italic border border-dashed border-border-default rounded-lg p-4 text-center">
                 No core memories extracted yet.
               </div>
             ) : (
               <ul className="space-y-2 max-h-48 overflow-y-auto">
                 {memories.map((mem) => (
                   <li key={mem.id} className="text-sm text-text-primary bg-primary border border-border-default rounded-md p-2">
                     <div className="flex items-start justify-between">
                        <span>{mem.content}</span>
                        {mem.core_memory && (
                          <span className="ml-2 bg-purple-100 text-purple-800 text-[10px] font-medium px-2 py-0.5 rounded-full dark:bg-purple-900 dark:text-purple-300">CORE</span>
                        )}
                     </div>
                     <div className="text-[10px] text-text-secondary mt-1">
                       Recalled {mem.access_frequency} times • Last: {new Date(mem.last_recall).toLocaleDateString()}
                     </div>
                   </li>
                 ))}
               </ul>
             )}
          </div>
        </div>
      )}
    </section>
  );
};
