import React from 'react';
import { Tool } from '../types';
import { MorphicMapCanvas } from './MorphicMapCanvas';

interface MorphicMapProps {
  tools: Tool[];
}

export const MorphicMap: React.FC<MorphicMapProps> = ({ tools }) => {
  return (
    <section className="bg-panel rounded-xl shadow-sm border border-border-default p-6 mt-6">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-xl font-semibold text-text-primary">Morphic Map</h2>
        <span className="bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-400 text-xs font-medium px-2.5 py-0.5 rounded-full">Visualizer</span>
      </div>

      <MorphicMapCanvas tools={tools} />

      {tools.length === 0 && (
        <div className="text-center text-zinc-500 mt-4 text-sm">
          No tools available to map.
        </div>
      )}
    </section>
  );
};
