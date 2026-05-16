import React, { useMemo, useEffect } from 'react';
import { ReactFlow, MiniMap, Controls, Background, useNodesState, useEdgesState, Edge, MarkerType } from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import { Tool } from '../types';
import { buildToolNodes } from '../utils/morphicMapUtils';

interface MorphicMapCanvasProps {
  tools: Tool[];
}

export const MorphicMapCanvas: React.FC<MorphicMapCanvasProps> = ({ tools }) => {
  const initialNodes = useMemo(() => buildToolNodes(tools), [tools]);

  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState<Edge>([]);

  // Update nodes when tools change
  useEffect(() => {
    setNodes(initialNodes);
  }, [initialNodes, setNodes]);

  useEffect(() => {
    let isMounted = true;
    const eventSource = new EventSource('/api/logs');

    eventSource.onmessage = (event) => {
      if (!isMounted) return;
      try {
        const parsedLog = JSON.parse(event.data);
        if (parsedLog.level === "INFO" && parsedLog.message) {
          // Check if it's an IPC event log emitted from the backend
          // E.g., eventBus.Publish might log something like "IPC Event: topic from toolA to toolB"
          // Let's assume the Wasm sandbox or event bus logs it in a specific format to track who published it.
          // Since the requirement says "when an event fires between two tools, dynamically render an animated edge",
          // we need a generic way to represent this. For now, if a log string contains '->' or is a specific IPC format we could parse it.
          // Let's try to parse if the backend sends a structured JSON log for IPC events specifically.

          try {
            // Check if the log message itself is a JSON containing IPC details, e.g., published by Sandbox parse
            const ipcData = JSON.parse(parsedLog.message);
            if (ipcData._ipc_event && ipcData.source_tool && ipcData.target_tool) {
               const edgeId = `e-${ipcData.source_tool}-${ipcData.target_tool}-${Date.now()}`;
               setEdges((eds) => {
                 const newEdge: Edge = {
                   id: edgeId,
                   source: ipcData.source_tool,
                   target: ipcData.target_tool,
                   animated: true,
                   label: ipcData._ipc_event,
                   style: { stroke: '#f59e0b', strokeWidth: 2 },
                   markerEnd: {
                     type: MarkerType.ArrowClosed,
                     color: '#f59e0b',
                   },
                 };
                 return [...eds, newEdge];
               });

               // Optional: remove the edge after a few seconds to make it look like a transient flow
               setTimeout(() => {
                 if (isMounted) {
                   setEdges((eds) => eds.filter(e => e.id !== edgeId));
                 }
               }, 5000);
            }
          } catch {
            // Not a JSON message, ignore
          }
        }
      } catch {
        // ignore parse error
      }
    };

    return () => {
      isMounted = false;
      eventSource.close();
    };
  }, [setEdges]);

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
