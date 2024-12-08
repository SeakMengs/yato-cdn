import { CDN_URL } from "@/util";
import { useState } from "react";

export function useFileDelete() {
  const [deleting, setDeleting] = useState(false);
  const [deleteError, setDeleteError] = useState<string | null>(null);
  const [deleteSuccess, setDeleteSuccess] = useState(false);

  const deleteFile = async (filename: string) => {
    setDeleting(true);
    setDeleteError(null);
    setDeleteSuccess(false);

    try {
      const response = await fetch(`${CDN_URL}/api/v1/cdn/${filename}`, {
        method: "DELETE",
        cache: "no-cache",
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      setDeleteSuccess(true);
    } catch (error) {
      setDeleteError(
        error instanceof Error
          ? error.message
          : "An error occurred during delete",
      );
    } finally {
      setDeleting(false);
    }
  };

  return { deleteFile, deleting, deleteError, deleteSuccess };
}
