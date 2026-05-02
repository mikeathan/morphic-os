import React from 'react';
import { useVFSExplorer } from '../hooks/useVFSExplorer';
import { FILE_ICONS } from '../utils/vfsUtils';

export const VFSExplorer: React.FC = () => {
  const {
    files,
    loading,
    selectedFile,
    fileContent,
    contentLoading,
    loadFiles,
    handleFileClick,
  } = useVFSExplorer();

  return (
    <section className="bg-white dark:bg-zinc-900 rounded-xl shadow-sm border border-zinc-200 dark:border-zinc-800 p-6 mt-6 flex flex-col h-[500px]">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-xl font-semibold">VFS Explorer</h2>
        <button onClick={loadFiles} className="text-sm text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300">
          Refresh
        </button>
      </div>

      <div className="flex flex-1 overflow-hidden gap-4">
        {/* Left pane: File tree/list */}
        <div className="w-1/3 border-r border-zinc-200 dark:border-zinc-800 overflow-y-auto pr-4">
          {loading ? (
            <div className="text-sm text-zinc-500">Loading files...</div>
          ) : files.length === 0 ? (
            <div className="text-sm text-zinc-500">No files found in VFS.</div>
          ) : (
            <ul className="space-y-1">
              {files.map(file => (
                <li key={file.id}>
                  <button
                    onClick={() => handleFileClick(file)}
                    className={`flex items-center w-full text-left px-2 py-1.5 rounded-md text-sm transition-colors ${selectedFile?.id === file.id ? 'bg-blue-50 dark:bg-blue-900/20 text-blue-700 dark:text-blue-300' : 'hover:bg-zinc-100 dark:hover:bg-zinc-800'}`}
                  >
                    <span className="mr-2 text-zinc-400">
                      {file.is_dir ? FILE_ICONS.dir : FILE_ICONS.file}
                    </span>
                    <span className="truncate">{file.path}</span>
                  </button>
                </li>
              ))}
            </ul>
          )}
        </div>

        {/* Right pane: File content preview */}
        <div className="w-2/3 flex flex-col overflow-hidden bg-zinc-50 dark:bg-zinc-950 rounded-lg border border-zinc-200 dark:border-zinc-800">
          {selectedFile ? (
            <>
              <div className="px-4 py-2 border-b border-zinc-200 dark:border-zinc-800 bg-zinc-100 dark:bg-zinc-900 flex justify-between items-center text-xs text-zinc-500">
                <span className="font-mono truncate">{selectedFile.path}</span>
                <span>{selectedFile.size} bytes</span>
              </div>
              <div className="flex-1 overflow-y-auto p-4">
                {contentLoading ? (
                  <div className="text-sm text-zinc-500">Loading content...</div>
                ) : (
                  <pre className="text-xs font-mono whitespace-pre-wrap text-zinc-800 dark:text-zinc-300">
                    {fileContent}
                  </pre>
                )}
              </div>
            </>
          ) : (
            <div className="flex-1 flex items-center justify-center text-sm text-zinc-400">
              Select a file to view its contents.
            </div>
          )}
        </div>
      </div>
    </section>
  );
};
