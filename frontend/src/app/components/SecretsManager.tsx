import React, { useState, useEffect } from 'react';
import { Secret } from '../types';
import { fetchSecrets, addSecret, deleteSecret } from '../services/api';
import { SecretsList } from './SecretsList';
import { SecretForm } from './SecretForm';

export function SecretsManager() {
  const [secrets, setSecrets] = useState<Secret[]>([]);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadSecrets = async () => {
    try {
      const data = await fetchSecrets();
      setSecrets(data);
    } catch (err) {
      console.error(err);
      setError('Error connecting to backend');
    }
  };

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    loadSecrets();
  }, []);

  const handleAddSecret = async (key: string, value: string) => {
    setIsSubmitting(true);
    setError(null);

    try {
      await addSecret(key, value);
      await loadSecrets();
    } catch (err) {
      console.error(err);
      setError('Failed to create secret');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDeleteSecret = async (id: string) => {
    if (!window.confirm('Are you sure you want to delete this secret?')) return;

    try {
      await deleteSecret(id);
      await loadSecrets();
    } catch (err) {
      console.error(err);
      setError('Failed to delete secret');
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

      <SecretForm onSubmit={handleAddSecret} isSubmitting={isSubmitting} />

      <SecretsList secrets={secrets} onDelete={handleDeleteSecret} />
    </div>
  );
}
