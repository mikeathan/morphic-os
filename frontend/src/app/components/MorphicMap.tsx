import React from 'react';
import { Tool } from '../types';
import { MorphicMapCanvas } from './MorphicMapCanvas';

interface MorphicMapProps {
  tools: Tool[];
}

export const MorphicMap: React.FC<MorphicMapProps> = ({ tools }) => {
  return (
    <section className="bg-white dark:bg-zinc-900 rounded-xl shadow-sm border border-zinc-200 dark:border-zinc-800 p-6 mt-6">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-xl font-semibold">Morphic Map</h2>
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
