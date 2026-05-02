export interface Tool {
  id: string;
  name: string;
  description: string;
  language: string;
  active: boolean;
}

export interface LogEvent {
  type: string;
  level: string;
  message: string;
}

export interface TaskResponse {
  id: string;
  task: string;
  result: string | null;
  logs: LogEvent[];
  timestamp: number;
}

export interface VirtualFile {
  id: string;
  workspace_id: string;
  path: string;
  name: string;
  is_dir: boolean;
  size: number;
  content?: string; // Content is often passed as base64 or string bytes
  created_at: string;
  updated_at: string;
}
