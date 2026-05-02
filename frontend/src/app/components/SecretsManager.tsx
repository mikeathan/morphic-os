import React, { useState, useEffect } from 'react';

interface Secret {
  id: string;
  workspace_id: string;
  key: string;
  created_at: string;
}

export function SecretsManager() {
  const [secrets, setSecrets] = useState<Secret[]>([]);
  const [newKey, setNewKey] = useState('');
  const [newValue, setNewValue] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchSecrets = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/secrets');
      if (response.ok) {
        const data = await response.json();
        setSecrets(data || []);
      } else {
        setError('Failed to fetch secrets');
      }
    } catch (err) {
      console.error(err);
      setError('Error connecting to backend');
    }
  };

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    fetchSecrets();
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newKey.trim() || !newValue.trim()) return;

    setIsSubmitting(true);
    setError(null);

    try {
      const response = await fetch('http://localhost:8080/api/secrets', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          workspace_id: 'default',
          key: newKey,
          value: newValue,
        }),
      });

      if (response.ok) {
        setNewKey('');
        setNewValue('');
        fetchSecrets();
      } else {
        setError('Failed to create secret');
      }
    } catch (err) {
      console.error(err);
      setError('Error creating secret');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!window.confirm('Are you sure you want to delete this secret?')) return;

    try {
      const response = await fetch(`http://localhost:8080/api/secrets/${id}`, {
        method: 'DELETE',
      });
      if (response.ok) {
        fetchSecrets();
      } else {
        setError('Failed to delete secret');
      }
    } catch (err) {
      console.error(err);
      setError('Error deleting secret');
    }
  };

  return (
    <div className="bg-surface border border-border-default rounded-xl p-4 shadow-sm mt-8">
      <h2 className="text-lg font-semibold mb-4 text-text-primary flex items-center gap-2">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect width="18" height="11" x="3" y="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
        Secrets Manager
      </h2>

      {error && (
        <div className="mb-4 text-sm text-red-500">
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit} className="mb-6 flex flex-col gap-2">
        <input
          type="text"
          value={newKey}
          onChange={(e) => setNewKey(e.target.value)}
          placeholder="Key (e.g. OPENAI_API_KEY)"
          className="flex-1 bg-primary border border-border-default rounded-lg px-3 py-2 text-sm text-text-primary focus:outline-none focus:ring-1 focus:ring-accent-primary"
          disabled={isSubmitting}
        />
        <input
          type="password"
          value={newValue}
          onChange={(e) => setNewValue(e.target.value)}
          placeholder="Value"
          className="flex-1 bg-primary border border-border-default rounded-lg px-3 py-2 text-sm text-text-primary focus:outline-none focus:ring-1 focus:ring-accent-primary"
          disabled={isSubmitting}
        />
        <button
          type="submit"
          disabled={isSubmitting || !newKey.trim() || !newValue.trim()}
          className="bg-accent-primary hover:bg-accent-hover text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors disabled:opacity-50"
        >
          {isSubmitting ? 'Adding...' : 'Add'}
        </button>
      </form>

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
                onClick={() => handleDelete(secret.id)}
                className="text-text-secondary hover:text-red-500 p-1"
                title="Delete secret"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/><line x1="10" x2="10" y1="11" y2="17"/><line x1="14" x2="14" y1="11" y2="17"/></svg>
              </button>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
