import { Tool } from '../types';

export const fetchTools = async (): Promise<Tool[]> => {
  const res = await fetch('/api/tools');
  if (!res.ok) {
    throw new Error('Failed to fetch tools');
  }
  const data = await res.json();
  return Array.isArray(data) ? data : [];
};

export const submitTask = async (task: string): Promise<string> => {
  const res = await fetch('/api/task', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ task }),
  });

  if (!res.ok) {
    throw new Error(`Error: ${res.status} ${res.statusText}`);
  }

  const data = await res.json();
  return data.result || 'No result returned.';
};
