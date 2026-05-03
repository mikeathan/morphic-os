import React from 'react';
import { Secret } from '../types';

interface SecretsListProps {
  secrets: Secret[];
  onDelete: (id: string) => void;
}

export function SecretsList({ secrets, onDelete }: SecretsListProps) {
  return (
    <div className="space-y-2 max-h-48 overflow-y-auto">
      {secrets.length === 0 ? (
        <p className="text-sm text-text-secondary text-center py-4">No secrets stored.</p>
      ) : (
        secrets.map((secret) => (
          <div key={secret.id} className="flex items-center justify-between p-2 rounded-lg bg-primary border border-border-default">
            <div className="flex flex-col">
              <span className="text-sm font-mono text-text-primary">{secret.key}</span>
              <span className="text-xs text-text-secondary">Added {new Date(secret.created_at).toLocaleDateString()}</span>
            </div>
            <button
              onClick={() => onDelete(secret.id)}
              className="text-text-secondary hover:text-red-500 p-1"
              title="Delete secret"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/><line x1="10" x2="10" y1="11" y2="17"/><line x1="14" x2="14" y1="11" y2="17"/></svg>
            </button>
          </div>
        ))
      )}
    </div>
  );
}
