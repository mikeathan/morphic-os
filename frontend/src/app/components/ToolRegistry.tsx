import React from 'react';
import { Tool } from '../types';

interface ToolRegistryProps {
  tools: Tool[];
}

export const ToolRegistry: React.FC<ToolRegistryProps> = ({ tools }) => {
  return (
    <section className="bg-panel rounded-xl shadow-sm border border-border-default p-6">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-xl font-semibold text-text-primary">Tool Registry</h2>
        <span className="bg-tertiary text-text-primary text-xs font-medium px-2.5 py-0.5 rounded-full">Active</span>
      </div>

      <div className="space-y-4">
        {tools.length === 0 ? (
          <p className="text-text-secondary text-sm">No active tools found.</p>
        ) : (
          tools.map((tool) => (
            <div key={tool.id} className={`border border-border-default rounded-lg p-4 transition-colors hover:bg-primary ${!tool.active && 'opacity-60'}`}>
              <div className="flex justify-between items-start">
                <div>
                  <h3 className="font-medium text-lg text-text-primary">{tool.name}</h3>
                  <p className="text-sm text-text-secondary mt-1">{tool.description}</p>
                </div>
                {tool.active ? (
                  <span className="bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400 text-xs font-medium px-2 py-1 rounded">Active</span>
                ) : (
                  <span className="bg-tertiary text-text-secondary text-xs font-medium px-2 py-1 rounded">Inactive</span>
                )}
              </div>
              <div className="mt-3 flex items-center text-xs text-text-muted">
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
