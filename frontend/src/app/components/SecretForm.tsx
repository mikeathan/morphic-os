import React, { useState } from 'react';

interface SecretFormProps {
  onSubmit: (key: string, value: string) => Promise<void>;
  isSubmitting: boolean;
}

export function SecretForm({ onSubmit, isSubmitting }: SecretFormProps) {
  const [newKey, setNewKey] = useState('');
  const [newValue, setNewValue] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newKey.trim() || !newValue.trim()) return;

    await onSubmit(newKey, newValue);
    setNewKey('');
    setNewValue('');
  };

  return (
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
  );
}
