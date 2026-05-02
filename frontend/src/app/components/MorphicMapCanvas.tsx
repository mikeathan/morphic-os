import React, { useMemo } from 'react';
import { ReactFlow, MiniMap, Controls, Background, useNodesState, useEdgesState } from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import { Tool } from '../types';
import { buildToolNodes } from '../utils/morphicMapUtils';

interface MorphicMapCanvasProps {
  tools: Tool[];
}

export const MorphicMapCanvas: React.FC<MorphicMapCanvasProps> = ({ tools }) => {
  const initialNodes = useMemo(() => buildToolNodes(tools), [tools]);

  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, , onEdgesChange] = useEdgesState([]);

  // Update nodes when tools change
  React.useEffect(() => {
    setNodes(initialNodes);
  }, [initialNodes, setNodes]);

  return (
    <div style={{ height: 400, width: '100%' }} className="border border-border-default rounded-lg overflow-hidden">
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
  );
};
