import React from 'react';
import { Tool } from '../types';

interface ToolRegistryProps {
  tools: Tool[];
}

export const ToolRegistry: React.FC<ToolRegistryProps> = ({ tools }) => {
  return (
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
  );
};
