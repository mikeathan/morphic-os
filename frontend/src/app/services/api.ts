import { Tool, VirtualFile, Secret } from '../types';

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
  const res = await fetch('/api/tasks', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ task }),
  });
  if (!res.ok) {
    throw new Error('Failed to submit task');
  }
  const data = await res.json();
  return data.result;
};

export const fetchSecrets = async (): Promise<Secret[]> => {
  const res = await fetch('/api/secrets');
  if (!res.ok) {
    throw new Error('Failed to fetch secrets');
  }
  const data = await res.json();
  return Array.isArray(data) ? data : [];
};

export const addSecret = async (key: string, value: string): Promise<Secret> => {
  const res = await fetch('/api/secrets', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      workspace_id: 'default',
      key,
      value,
    }),
  });
  if (!res.ok) {
    throw new Error('Failed to create secret');
  }
  return await res.json();
};

export const deleteSecret = async (id: string): Promise<void> => {
  const res = await fetch(`/api/secrets/${id}`, {
    method: 'DELETE',
  });
  if (!res.ok) {
    throw new Error('Failed to delete secret');
  }
};
