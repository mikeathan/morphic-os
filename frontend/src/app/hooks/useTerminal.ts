import { useState, useEffect } from 'react';
import { LogEvent } from '../types';

export const useTerminal = () => {
  const [logs, setLogs] = useState<LogEvent[]>([]);

  useEffect(() => {
    const eventSource = new EventSource('/api/logs');

    eventSource.onmessage = (event) => {
      try {
        const parsedLog = JSON.parse(event.data);
        setLogs((prev) => [...prev, parsedLog]);
      } catch (err) {
        console.error("Failed to parse log event in terminal:", err);
      }
    };

    eventSource.onerror = (err) => {
      console.error("EventSource failed:", err);
    };

    return () => {
      eventSource.close();
    };
  }, []);

  return { logs };
};
