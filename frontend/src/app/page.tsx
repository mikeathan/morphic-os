"use client";

import React, { useState, useEffect } from 'react';
import { fetchTools, submitTask } from './services/api';
import { Tool, LogEvent } from './types';
import { ToolRegistry } from './components/ToolRegistry';
import { LiveLogs } from './components/LiveLogs';
import { TaskSubmissionForm } from './components/TaskSubmissionForm';

export default function Home() {
  const [tools, setTools] = useState<Tool[]>([]);
  const [logs, setLogs] = useState<LogEvent[]>([]);
  const [taskInput, setTaskInput] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [taskResult, setTaskResult] = useState<string | null>(null);

  useEffect(() => {
    // Initial fetch of tools
    fetchTools()
      .then(setTools)
      .catch((err) => console.error("Failed to fetch tools:", err));

    // Setup SSE for real-time logs
    const eventSource = new EventSource('/api/logs');
    eventSource.onmessage = (event) => {
      try {
        const parsedLog = JSON.parse(event.data);
        setLogs((prev) => [...prev, parsedLog]);
      } catch (err) {
        console.error("Failed to parse log event:", err);
      }
    };
    return () => {
      eventSource.close();
    };
  }, []);

  const handleTaskSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!taskInput.trim()) return;

    setIsSubmitting(true);
    setTaskResult(null);
    setLogs([]); // Clear logs for new task

    try {
      const result = await submitTask(taskInput);
      setTaskResult(result);

      // Refresh tool registry after task completes
      const updatedTools = await fetchTools();
      setTools(updatedTools);
    } catch (err: any) {
      setTaskResult(`Error: ${err.message}`);
    } finally {
      setIsSubmitting(false);
      setTaskInput("");
    }
  };

  return (
    <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950 p-8 font-sans text-zinc-900 dark:text-zinc-50">
      <header className="mb-8">
        <h1 className="text-3xl font-bold tracking-tight">Morphic-OS Dashboard</h1>
        <p className="text-zinc-600 dark:text-zinc-400 mt-2">Real-time monitoring and tool registry</p>
      </header>

      <TaskSubmissionForm
        taskInput={taskInput}
        setTaskInput={setTaskInput}
        isSubmitting={isSubmitting}
        handleSubmitTask={handleTaskSubmit}
        taskResult={taskResult}
      />

      <main className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <ToolRegistry tools={tools} />
        <LiveLogs logs={logs} />
      </main>
    </div>
  );
}
