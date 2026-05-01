import { useState, useEffect, useCallback, useRef } from 'react';
import { fetchTools, submitTask } from '../services/api';
import { Tool, TaskResponse } from '../types';

export const useMorphicLoop = () => {
  const [tools, setTools] = useState<Tool[]>([]);
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Keep history of previous tasks
  const [history, setHistory] = useState<TaskResponse[]>([]);

  // Current active task state
  const [currentTask, setCurrentTask] = useState<TaskResponse | null>(null);

  // Use refs to keep event listener callbacks up-to-date with current state without re-binding
  const currentTaskRef = useRef<TaskResponse | null>(null);

  // Sync ref with state
  useEffect(() => {
    currentTaskRef.current = currentTask;
  }, [currentTask]);

  const loadTools = useCallback(async () => {
    try {
      const data = await fetchTools();
      setTools(data);
    } catch (err) {
      console.error("Failed to fetch tools:", err);
    }
  }, []);

  useEffect(() => {
    let isMounted = true;
    fetchTools().then(data => {
      if (isMounted) setTools(data);
    }).catch(console.error);

    const eventSource = new EventSource('/api/logs');
    eventSource.onmessage = (event) => {
      try {
        const parsedLog = JSON.parse(event.data);
        // Only append logs if we have a current task running
        if (currentTaskRef.current) {
          setCurrentTask(prev => {
            if (!prev) return prev;
            return {
              ...prev,
              logs: [...prev.logs, parsedLog]
            };
          });
        }
      } catch (err) {
        console.error("Failed to parse log event:", err);
      }
    };

    return () => {
      isMounted = false;
      eventSource.close();
    };
  }, []);

  const runTask = async (taskInput: string) => {
    if (!taskInput.trim() || isSubmitting) return;

    // If there is an existing currentTask, move it to history
    if (currentTask) {
      setHistory(prev => [currentTask, ...prev]);
    }

    setIsSubmitting(true);

    const newTask: TaskResponse = {
      id: Math.random().toString(36).substring(7),
      task: taskInput,
      result: null,
      logs: [],
      timestamp: Date.now()
    };

    setCurrentTask(newTask);

    try {
      const result = await submitTask(taskInput);

      setCurrentTask(prev => {
        if (!prev) return prev;
        return {
          ...prev,
          result: result
        };
      });

      await loadTools();
    } catch (err) {
      setCurrentTask(prev => {
        if (!prev) return prev;
        return {
          ...prev,
          result: err instanceof Error ? `Error: ${err.message}` : `Error: ${String(err)}`
        };
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  return {
    tools,
    currentTask,
    history,
    isSubmitting,
    runTask
  };
};
