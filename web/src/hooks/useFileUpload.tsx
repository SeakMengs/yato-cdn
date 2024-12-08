import { CDN_URL } from "@/util";
import { useState } from "react";

export function useFileUpload() {
  const [uploading, setUploading] = useState(false);
  const [uploadError, setUploadError] = useState<string | null>(null);
  const [uploadSuccess, setUploadSuccess] = useState(false);

  const uploadFile = async (file: File) => {
    setUploading(true);
    setUploadError(null);
    setUploadSuccess(false);

    const formData = new FormData();
    formData.append("file", file);

    try {
      const response = await fetch(`${CDN_URL}/api/v1/cdn/upload`, {
        method: "POST",
        body: formData,
        cache: "no-cache",
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      setUploadSuccess(true);
    } catch (error) {
      setUploadError(
        error instanceof Error
          ? error.message
          : "An error occurred during upload",
      );
    } finally {
      setUploading(false);
    }
  };

  return { uploadFile, uploading, uploadError, uploadSuccess };
}
