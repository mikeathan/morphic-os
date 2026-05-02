import React, { useMemo } from 'react';
import { ReactFlow, MiniMap, Controls, Background, useNodesState, useEdgesState } from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import { Tool } from '../types';

interface MorphicMapProps {
  tools: Tool[];
}

export const MorphicMap: React.FC<MorphicMapProps> = ({ tools }) => {
  const initialNodes = useMemo(() => {
    return tools.map((tool, index) => ({
      id: tool.id,
      position: { x: 250 + (index % 3) * 200, y: 100 + Math.floor(index / 3) * 150 },
      data: { label: tool.name },
      style: {
        background: tool.active ? '#ecfdf5' : '#f4f4f5',
        border: `1px solid ${tool.active ? '#10b981' : '#a1a1aa'}`,
        borderRadius: '8px',
        padding: '10px',
        color: tool.active ? '#065f46' : '#52525b',
        fontWeight: 'bold',
        opacity: tool.active ? 1 : 0.6,
      },
    }));
  }, [tools]);

  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, , onEdgesChange] = useEdgesState([]);

  // Update nodes when tools change
  React.useEffect(() => {
    setNodes(initialNodes);
  }, [initialNodes, setNodes]);

  return (
    <section className="bg-white dark:bg-zinc-900 rounded-xl shadow-sm border border-zinc-200 dark:border-zinc-800 p-6 mt-6">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-xl font-semibold">Morphic Map</h2>
        <span className="bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-400 text-xs font-medium px-2.5 py-0.5 rounded-full">Visualizer</span>
      </div>

      <div style={{ height: 400, width: '100%' }} className="border border-zinc-200 dark:border-zinc-800 rounded-lg overflow-hidden">
        <ReactFlow
          nodes={nodes}
          edges={edges}
          onNodesChange={onNodesChange}
          onEdgesChange={onEdgesChange}
          fitView
          attributionPosition="bottom-right"
        >
          <Controls />
          <MiniMap zoomable pannable />
          <Background color="#aaa" gap={16} />
        </ReactFlow>
      </div>
      {tools.length === 0 && (
        <div className="text-center text-zinc-500 mt-4 text-sm">
          No tools available to map.
        </div>
      )}
    </section>
  );
};
