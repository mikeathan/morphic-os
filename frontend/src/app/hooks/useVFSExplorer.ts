import { useState, useEffect, useCallback } from 'react';
import { VirtualFile } from '../types';
import { fetchVirtualFiles, fetchVirtualFileContent } from '../services/api';
import { decodeBase64 } from '../utils/vfsUtils';

export const useVFSExplorer = () => {
  const [files, setFiles] = useState<VirtualFile[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedFile, setSelectedFile] = useState<VirtualFile | null>(null);
  const [fileContent, setFileContent] = useState<string | null>(null);
  const [contentLoading, setContentLoading] = useState(false);

  const loadFiles = useCallback(async () => {
    setLoading(true);
    try {
      const data = await fetchVirtualFiles();
      setFiles(data);
    } catch (err) {
      console.error("Failed to load VFS files:", err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    let isMounted = true;
    const fetchInitialFiles = async () => {
      try {
        const data = await fetchVirtualFiles();
        if (isMounted) {
          setFiles(data);
        }
      } catch (err) {
        console.error("Failed to load VFS files:", err);
      } finally {
        if (isMounted) {
          setLoading(false);
        }
      }
    };
    fetchInitialFiles();
    return () => { isMounted = false; };
  }, []);

  const handleFileClick = useCallback(async (file: VirtualFile) => {
    if (file.is_dir) return; // For now, we just display contents of files

    setSelectedFile(file);
    setContentLoading(true);
    try {
      const fullFile = await fetchVirtualFileContent(file.id);
      if (fullFile.content) {
        setFileContent(decodeBase64(fullFile.content));
      } else {
        setFileContent("");
      }
    } catch (err) {
      console.error("Failed to load file content:", err);
      setFileContent("Error loading content.");
    } finally {
      setContentLoading(false);
    }
  }, []);

  return {
    files,
    loading,
    selectedFile,
    fileContent,
    contentLoading,
    loadFiles,
    handleFileClick,
  };
};
