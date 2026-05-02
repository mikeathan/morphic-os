"use client";

import React, { useState } from 'react';
import { ToolRegistry } from './components/ToolRegistry';
import { MorphicMap } from './components/MorphicMap';
import { VFSExplorer } from './components/VFSExplorer';
import { Terminal } from './components/Terminal';
import { HardwareMetrics } from './components/HardwareMetrics';
import { TaskSubmissionForm } from './components/TaskSubmissionForm';
import { TaskDisplay } from './components/TaskDisplay';
import { useMorphicLoop } from './hooks/useMorphicLoop';

export default function Home() {
  const { tools, currentTask, history, isSubmitting, runTask } = useMorphicLoop();
  const [taskInput, setTaskInput] = useState("");

  const handleTaskSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!taskInput.trim()) return;

    await runTask(taskInput);
    setTaskInput("");
  };

  return (
    <div className="min-h-screen bg-primary p-8 font-sans text-text-primary">
      <header className="mb-8">
        <h1 className="text-3xl font-bold tracking-tight">Morphic-OS Dashboard</h1>
        <p className="text-text-secondary mt-2">Real-time monitoring and tool registry</p>
      </header>

      <main className="grid grid-cols-1 lg:grid-cols-3 gap-8">

        {/* Left Column: Tools */}
        <div className="lg:col-span-1">
          <ToolRegistry tools={tools} />
          <HardwareMetrics />
          <VFSExplorer />
          <MorphicMap tools={tools} />
          <Terminal />
        </div>

        {/* Right Column: Interaction & History */}
        <div className="lg:col-span-2 flex flex-col h-full">
          <TaskSubmissionForm
            taskInput={taskInput}
            setTaskInput={setTaskInput}
            isSubmitting={isSubmitting}
            handleSubmitTask={handleTaskSubmit}
            taskResult={null} // Task result is now handled by TaskDisplay
          />

          <div className="mt-4 flex-1">
            {currentTask && (
              <TaskDisplay taskResponse={currentTask} isActive={true} />
            )}

            {history.map((task) => (
              <TaskDisplay key={task.id} taskResponse={task} isActive={false} />
            ))}

            {!currentTask && history.length === 0 && (
              <div className="text-center py-12 text-text-secondary border border-dashed border-border-default rounded-xl">
                Submit a task to see Morphic-OS in action.
              </div>
            )}
          </div>
        </div>

      </main>
    </div>
  );
}
