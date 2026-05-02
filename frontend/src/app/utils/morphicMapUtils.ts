import { Tool } from '../types';

// Use CSS variables for proper dark/light mode support
export const getToolNodeStyle = (isActive: boolean) => ({
  background: isActive ? 'var(--tool-active-bg, #ecfdf5)' : 'var(--tool-inactive-bg, #f4f4f5)',
  border: `1px solid ${isActive ? 'var(--tool-active-border, #10b981)' : 'var(--tool-inactive-border, #a1a1aa)'}`,
  borderRadius: '8px',
  padding: '10px',
  color: isActive ? 'var(--tool-active-text, #065f46)' : 'var(--tool-inactive-text, #52525b)',
  fontWeight: 'bold',
  opacity: isActive ? 1 : 0.6,
});

export const buildToolNodes = (tools: Tool[]) => {
  return tools.map((tool, index) => ({
    id: tool.id,
    position: { x: 250 + (index % 3) * 200, y: 100 + Math.floor(index / 3) * 150 },
    data: { label: tool.name },
    style: getToolNodeStyle(tool.active),
  }));
};
