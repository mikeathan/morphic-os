import React, { useState } from 'react';
import { TaskResponse } from '../types';

interface TaskDisplayProps {
  taskResponse: TaskResponse;
  isActive?: boolean;
}

export const TaskDisplay: React.FC<TaskDisplayProps> = ({ taskResponse, isActive = false }) => {
  const [logsExpanded, setLogsExpanded] = useState(isActive);

  return (
    <div className={`mb-6 rounded-xl border ${isActive ? 'border-blue-500/50 shadow-md bg-secondary' : 'border-border-default bg-primary'} overflow-hidden transition-all`}>
      {/* Header / User Prompt */}
      <div className="p-4 bg-tertiary border-b border-border-default flex justify-between items-center cursor-pointer" onClick={() => setLogsExpanded(!logsExpanded)}>
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 rounded-full bg-blue-100 dark:bg-blue-900/50 flex items-center justify-center text-blue-600 dark:text-blue-400 font-semibold text-sm">
            You
          </div>
          <span className="font-medium text-text-primary">{taskResponse.task}</span>
        </div>
        <div className="flex items-center gap-3">
          {isActive && !taskResponse.result && (
            <span className="flex h-3 w-3 relative">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-blue-400 opacity-75"></span>
              <span className="relative inline-flex rounded-full h-3 w-3 bg-blue-500"></span>
            </span>
          )}
          <button className="text-text-muted hover:text-text-secondary">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={`transform transition-transform ${logsExpanded ? 'rotate-180' : ''}`}>
              <polyline points="6 9 12 15 18 9"></polyline>
            </svg>
          </button>
        </div>
      </div>

      {/* Morphic Loop Agent Logs (Collapsible) */}
      {logsExpanded && (
        <div className="bg-black p-4 font-mono text-sm border-b border-border-default max-h-[300px] overflow-y-auto">
          {taskResponse.logs.length === 0 ? (
            <div className="text-zinc-500 italic">Thinking...</div>
          ) : (
            <div className="space-y-2">
              {taskResponse.logs.map((log, index) => {
                let colorClass = "text-blue-400";
                if (log.level === "EVAL") colorClass = "text-purple-400";
                else if (log.level === "FORGE") colorClass = "text-amber-400";
                else if (log.level === "SUCCESS") colorClass = "text-green-400";
                else if (log.level === "ERROR") colorClass = "text-red-400";
                else if (log.level === "EXEC") colorClass = "text-blue-400";

                return (
                  <div key={index} className="text-zinc-400">
                    <span className="text-zinc-600">❯</span>{" "}
                    <span className={colorClass}>[{log.level || "INFO"}]</span> {log.message}
                  </div>
                );
              })}
            </div>
          )}
        </div>
      )}

      {/* Final Result */}
      {taskResponse.result && (
        <div className="p-5 flex gap-4">
          <div className="w-8 h-8 rounded-full bg-emerald-100 dark:bg-emerald-900/50 flex items-center justify-center text-emerald-600 dark:text-emerald-400 flex-shrink-0">
             <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
          </div>
          <div className="whitespace-pre-wrap font-mono text-sm text-text-primary mt-1">
            {taskResponse.result}
          </div>
        </div>
      )}
    </div>
  );
};
