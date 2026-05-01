import React from 'react';
import { LogEvent } from '../types';

interface LiveLogsProps {
  logs: LogEvent[];
}

export const LiveLogs: React.FC<LiveLogsProps> = ({ logs }) => {
  return (
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
  );
};
