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
