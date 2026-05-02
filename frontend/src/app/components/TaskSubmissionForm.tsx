import React from 'react';

interface TaskSubmissionFormProps {
  taskInput: string;
  setTaskInput: (val: string) => void;
  isSubmitting: boolean;
  handleSubmitTask: (e: React.FormEvent) => void;
  taskResult: string | null;
}

export const TaskSubmissionForm: React.FC<TaskSubmissionFormProps> = ({
  taskInput,
  setTaskInput,
  isSubmitting,
  handleSubmitTask,
  taskResult
}) => {
  return (
    <div className="mb-8 bg-panel rounded-xl shadow-sm border border-border-default p-6">
      <h2 className="text-xl font-semibold mb-4 text-text-primary">Submit a Task</h2>
      <form onSubmit={handleSubmitTask} className="flex gap-4">
        <input
          type="text"
          value={taskInput}
          onChange={(e) => setTaskInput(e.target.value)}
          placeholder="E.g., Calculate the sum of primes up to 100"
          className="flex-1 rounded-md border border-border-default bg-transparent px-4 py-2 text-sm text-text-primary focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
          disabled={isSubmitting}
        />
        <button
          type="submit"
          disabled={isSubmitting || !taskInput.trim()}
          className="rounded-md bg-blue-600 px-6 py-2 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {isSubmitting ? "Running..." : "Run Task"}
        </button>
      </form>
      {taskResult && (
        <div className="mt-4 p-4 rounded-md bg-tertiary border border-border-default">
          <h3 className="text-sm font-semibold text-text-secondary mb-2">Result:</h3>
          <div className="whitespace-pre-wrap font-mono text-sm text-text-primary">{taskResult}</div>
        </div>
      )}
    </div>
  );
};
