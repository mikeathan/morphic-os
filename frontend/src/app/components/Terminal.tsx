import React, { useEffect, useRef } from 'react';
import { useTerminal } from '../hooks/useTerminal';

export const Terminal: React.FC = () => {
  const { logs } = useTerminal();
  const terminalEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    // Auto-scroll to bottom on new log
    if (terminalEndRef.current) {
      terminalEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [logs]);

  return (
    <section className="bg-zinc-950 rounded-xl shadow-sm border border-zinc-800 p-6 mt-6 flex flex-col h-[300px]">
      <div className="flex items-center justify-between mb-4 border-b border-zinc-800 pb-2">
        <h2 className="text-xl font-semibold text-zinc-100">Terminal</h2>
        <div className="flex gap-2">
           <div className="w-3 h-3 rounded-full bg-red-500"></div>
           <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
           <div className="w-3 h-3 rounded-full bg-green-500"></div>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto font-mono text-sm">
        {logs.length === 0 ? (
          <div className="text-zinc-500 italic">Waiting for events...</div>
        ) : (
          <div className="space-y-1">
            {logs.map((log, index) => {
              let colorClass = "text-blue-400";
              if (log.level === "EVAL") colorClass = "text-purple-400";
              else if (log.level === "FORGE") colorClass = "text-amber-400";
              else if (log.level === "SUCCESS") colorClass = "text-green-400";
              else if (log.level === "ERROR") colorClass = "text-red-400";
              else if (log.level === "EXEC") colorClass = "text-blue-400";

              return (
                <div key={index} className="text-zinc-300">
                  <span className="text-zinc-600">❯</span>{" "}
                  <span className={colorClass}>[{log.level || "INFO"}]</span> {log.message}
                </div>
              );
            })}
            <div ref={terminalEndRef} />
          </div>
        )}
      </div>
    </section>
  );
};
