import { Tool, VirtualFile } from '../types';

export const fetchVirtualFiles = async (workspaceId: string = "default"): Promise<VirtualFile[]> => {
  const res = await fetch(`/api/vfs/files?workspace_id=${workspaceId}`);
  if (!res.ok) {
    throw new Error('Failed to fetch virtual files');
  }
  const data = await res.json();
  return Array.isArray(data) ? data : [];
};

export const fetchVirtualFileContent = async (fileId: string): Promise<VirtualFile> => {
  const res = await fetch(`/api/vfs/files/${fileId}`);
  if (!res.ok) {
    throw new Error('Failed to fetch virtual file content');
  }
  return await res.json();
};

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
