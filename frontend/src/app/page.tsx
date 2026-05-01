"use client";

import React, { useState, useEffect } from 'react';

export default function Home() {
  const [tools, setTools] = useState<any[]>([]);
  const [logs, setLogs] = useState<any[]>([]);
  const [taskInput, setTaskInput] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [taskResult, setTaskResult] = useState<string | null>(null);

  useEffect(() => {
    fetch('/api/tools')
      .then((res) => res.json())
      .then((data) => {
        if (Array.isArray(data)) {
          setTools(data);
        }
      })
      .catch((err) => console.error("Failed to fetch tools:", err));

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

  const handleSubmitTask = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!taskInput.trim()) return;

    setIsSubmitting(true);
    setTaskResult(null);
    setLogs([]); // Clear logs for new task

    try {
      const res = await fetch('/api/task', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ task: taskInput }),
      });

      if (!res.ok) {
        throw new Error(`Error: ${res.status} ${res.statusText}`);
      }

      const data = await res.json();
      setTaskResult(data.result || "No result returned.");

      // Refresh tool registry
      const toolsRes = await fetch('/api/tools');
      const toolsData = await toolsRes.json();
      if (Array.isArray(toolsData)) {
        setTools(toolsData);
      }
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

      <div className="mb-8 bg-white dark:bg-zinc-900 rounded-xl shadow-sm border border-zinc-200 dark:border-zinc-800 p-6">
        <h2 className="text-xl font-semibold mb-4">Submit a Task</h2>
        <form onSubmit={handleSubmitTask} className="flex gap-4">
          <input
            type="text"
            value={taskInput}
            onChange={(e) => setTaskInput(e.target.value)}
            placeholder="E.g., Calculate the sum of primes up to 100"
            className="flex-1 rounded-md border border-zinc-300 dark:border-zinc-700 bg-transparent px-4 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
            disabled={isSubmitting}
          />
          <button
            type="submit"
            disabled={isSubmitting || !taskInput.trim()}
            className="rounded-md bg-blue-600 px-6 py-2 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isSubmitting ? "Running..." : "Run Task"}
          </button>
        </form>
        {taskResult && (
          <div className="mt-4 p-4 rounded-md bg-zinc-100 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700">
            <h3 className="text-sm font-semibold text-zinc-500 dark:text-zinc-400 mb-2">Result:</h3>
            <div className="whitespace-pre-wrap font-mono text-sm">{taskResult}</div>
          </div>
        )}
      </div>

      <main className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <section className="bg-white dark:bg-zinc-900 rounded-xl shadow-sm border border-zinc-200 dark:border-zinc-800 p-6">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-xl font-semibold">Tool Registry</h2>
            <span className="bg-zinc-100 dark:bg-zinc-800 text-xs font-medium px-2.5 py-0.5 rounded-full">Active</span>
          </div>

          <div className="space-y-4">
            {tools.length === 0 ? (
              <p className="text-zinc-500 text-sm">No active tools found.</p>
            ) : (
              tools.map((tool) => (
                <div key={tool.id} className={`border border-zinc-200 dark:border-zinc-800 rounded-lg p-4 transition-colors hover:bg-zinc-50 dark:hover:bg-zinc-800/50 ${!tool.active && 'opacity-60'}`}>
                  <div className="flex justify-between items-start">
                    <div>
                      <h3 className="font-medium text-lg">{tool.name}</h3>
                      <p className="text-sm text-zinc-600 dark:text-zinc-400 mt-1">{tool.description}</p>
                    </div>
                    {tool.active ? (
                      <span className="bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400 text-xs font-medium px-2 py-1 rounded">Active</span>
                    ) : (
                      <span className="bg-zinc-100 text-zinc-800 dark:bg-zinc-800 dark:text-zinc-300 text-xs font-medium px-2 py-1 rounded">Inactive</span>
                    )}
                  </div>
                  <div className="mt-3 flex items-center text-xs text-zinc-500">
                    <span className="mr-3">ID: {tool.id.substring(0, 8)}...</span>
                    <span>Language: {tool.language === "go" ? "Go (WASM)" : tool.language}</span>
                  </div>
                </div>
              ))
            )}
          </div>
        </section>

        <section className="bg-white dark:bg-zinc-900 rounded-xl shadow-sm border border-zinc-200 dark:border-zinc-800 p-6 flex flex-col h-[600px]">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-xl font-semibold">Real-time Logs (Morphic Loop)</h2>
            <div className="flex items-center gap-2">
              <span className="relative flex h-3 w-3">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-blue-400 opacity-75"></span>
                <span className="relative inline-flex rounded-full h-3 w-3 bg-blue-500"></span>
              </span>
              <span className="text-sm text-zinc-600 dark:text-zinc-400">Live</span>
            </div>
          </div>

          <div className="bg-zinc-950 rounded-lg p-4 font-mono text-sm overflow-y-auto flex-1 border border-zinc-800">
            <div className="space-y-3">
              {logs.length === 0 ? (
                <div className="text-zinc-500 italic">Waiting for events...</div>
              ) : (
                logs.map((log, index) => {
                  let colorClass = "text-blue-400";
                  if (log.level === "EVAL") colorClass = "text-purple-400";
                  else if (log.level === "FORGE") colorClass = "text-amber-400";
                  else if (log.level === "SUCCESS") colorClass = "text-green-400";
                  else if (log.level === "ERROR") colorClass = "text-red-400";
                  else if (log.level === "EXEC") colorClass = "text-blue-400";

                  return (
                    <div key={index} className="text-zinc-400">
                      <span className="text-zinc-500">[{new Date().toLocaleTimeString()}]</span>{" "}
                      <span className={colorClass}>[{log.level || "INFO"}]</span> {log.message}
                    </div>
                  );
                })
              )}
            </div>
          </div>
        </section>
      </main>
    </div>
  );
}
