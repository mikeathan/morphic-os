export const FILE_ICONS = {
  dir: '📁',
  file: '📄',
};

export const decodeBase64 = (content: string): string => {
  try {
    // Robust base64 decoding with TextDecoder for UTF-8 support
    const binaryString = atob(content);
    const bytes = new Uint8Array(binaryString.length);
    for (let i = 0; i < binaryString.length; i++) {
        bytes[i] = binaryString.charCodeAt(i);
    }
    return new TextDecoder('utf-8').decode(bytes);
  } catch {
     return content; // If it's not base64, it might just be text
  }
};
