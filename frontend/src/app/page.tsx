import React from 'react';

export default function Home() {
  return (
    <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950 p-8 font-sans text-zinc-900 dark:text-zinc-50">
      <header className="mb-8">
        <h1 className="text-3xl font-bold tracking-tight">Morphic-OS Dashboard</h1>
        <p className="text-zinc-600 dark:text-zinc-400 mt-2">Real-time monitoring and tool registry</p>
      </header>

      <main className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <section className="bg-white dark:bg-zinc-900 rounded-xl shadow-sm border border-zinc-200 dark:border-zinc-800 p-6">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-xl font-semibold">Tool Registry</h2>
            <span className="bg-zinc-100 dark:bg-zinc-800 text-xs font-medium px-2.5 py-0.5 rounded-full">Active</span>
          </div>

          <div className="space-y-4">
            <div className="border border-zinc-200 dark:border-zinc-800 rounded-lg p-4 transition-colors hover:bg-zinc-50 dark:hover:bg-zinc-800/50">
              <div className="flex justify-between items-start">
                <div>
                  <h3 className="font-medium text-lg">calculator</h3>
                  <p className="text-sm text-zinc-600 dark:text-zinc-400 mt-1">Evaluates mathematical expressions</p>
                </div>
                <span className="bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400 text-xs font-medium px-2 py-1 rounded">Active</span>
              </div>
              <div className="mt-3 flex items-center text-xs text-zinc-500">
                <span className="mr-3">Version: 1.0.0</span>
                <span>Language: Go (WASM)</span>
              </div>
            </div>

            <div className="border border-zinc-200 dark:border-zinc-800 rounded-lg p-4 transition-colors hover:bg-zinc-50 dark:hover:bg-zinc-800/50">
              <div className="flex justify-between items-start">
                <div>
                  <h3 className="font-medium text-lg">file_reader</h3>
                  <p className="text-sm text-zinc-600 dark:text-zinc-400 mt-1">Reads content from local files securely</p>
                </div>
                <span className="bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400 text-xs font-medium px-2 py-1 rounded">Active</span>
              </div>
              <div className="mt-3 flex items-center text-xs text-zinc-500">
                <span className="mr-3">Version: 1.2.0</span>
                <span>Language: Go (WASM)</span>
              </div>
            </div>

            <div className="border border-zinc-200 dark:border-zinc-800 rounded-lg p-4 opacity-60">
              <div className="flex justify-between items-start">
                <div>
                  <h3 className="font-medium text-lg">web_scraper</h3>
                  <p className="text-sm text-zinc-600 dark:text-zinc-400 mt-1">Extracts content from web pages</p>
                </div>
                <span className="bg-zinc-100 text-zinc-800 dark:bg-zinc-800 dark:text-zinc-300 text-xs font-medium px-2 py-1 rounded">Inactive</span>
              </div>
              <div className="mt-3 flex items-center text-xs text-zinc-500">
                <span className="mr-3">Version: 0.9.0</span>
                <span>Language: Go (WASM)</span>
              </div>
            </div>
          </div>
        </section>

        <section className="bg-white dark:bg-zinc-900 rounded-xl shadow-sm border border-zinc-200 dark:border-zinc-800 p-6 flex flex-col h-[600px]">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-xl font-semibold">Real-time Logs (Morphic Loop)</h2>
            <div className="flex items-center gap-2">
              <span className="relative flex h-3 w-3">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-blue-400 opacity-75"></span>
                <span className="relative inline-flex rounded-full h-3 w-3 bg-blue-500"></span>
              </span>
              <span className="text-sm text-zinc-600 dark:text-zinc-400">Live</span>
            </div>
          </div>

          <div className="bg-zinc-950 rounded-lg p-4 font-mono text-sm overflow-y-auto flex-1 border border-zinc-800">
            <div className="space-y-3">
              <div className="text-zinc-400">
                <span className="text-zinc-500">[10:42:01]</span> <span className="text-blue-400">[INFO]</span> Morphic Loop initialized
              </div>
              <div className="text-zinc-400">
                <span className="text-zinc-500">[10:42:05]</span> <span className="text-blue-400">[INFO]</span> Received user task: &quot;Calculate the sum of primes up to 100&quot;
              </div>
              <div className="text-zinc-400">
                <span className="text-zinc-500">[10:42:06]</span> <span className="text-purple-400">[EVAL]</span> Assembling context with 2 active tools
              </div>
              <div className="text-zinc-400">
                <span className="text-zinc-500">[10:42:08]</span> <span className="text-purple-400">[EVAL]</span> Agent response: Action=sys_forge_tool, ToolName=prime_calculator
              </div>
              <div className="text-zinc-400">
                <span className="text-zinc-500">[10:42:08]</span> <span className="text-amber-400">[FORGE]</span> Generating Go code for prime_calculator
              </div>
              <div className="text-zinc-400">
                <span className="text-zinc-500">[10:42:11]</span> <span className="text-amber-400">[FORGE]</span> Compiling tool to WebAssembly (GOOS=wasip1 GOARCH=wasm)
              </div>
              <div className="text-zinc-400">
                <span className="text-zinc-500">[10:42:15]</span> <span className="text-red-400">[ERROR]</span> Compilation failed: undefined variable &apos;n&apos;
              </div>
              <div className="text-zinc-400">
                <span className="text-zinc-500">[10:42:15]</span> <span className="text-amber-400">[FORGE]</span> Self-correction loop: Asking agent to fix code (Attempt 1/3)
              </div>
              <div className="text-zinc-400">
                <span className="text-zinc-500">[10:42:18]</span> <span className="text-amber-400">[FORGE]</span> Re-compiling fixed tool code
              </div>
              <div className="text-zinc-400">
                <span className="text-zinc-500">[10:42:20]</span> <span className="text-green-400">[SUCCESS]</span> Tool compiled successfully
              </div>
              <div className="text-zinc-400">
                <span className="text-zinc-500">[10:42:21]</span> <span className="text-blue-400">[INFO]</span> Saved prime_calculator to registry
              </div>
              <div className="text-zinc-400">
                <span className="text-zinc-500">[10:42:22]</span> <span className="text-blue-400">[EXEC]</span> Executing prime_calculator in Wazero sandbox
              </div>
              <div className="text-zinc-400">
                <span className="text-zinc-500">[10:42:23]</span> <span className="text-green-400">[SUCCESS]</span> Task completed. Result: 1060
              </div>
            </div>
          </div>
        </section>
      </main>
    </div>
  );
}
