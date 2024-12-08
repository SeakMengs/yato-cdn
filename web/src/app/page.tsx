"use client";

import { useFetch } from "../hooks/useFetch";
import { useFileUpload } from "../hooks/useFileUpload";
import { RegionInfo } from "./RegionInfo";
import { FileUpload } from "./FileUpload";
import { FileList } from "./FileList";

export interface Region {
  id: string;
  name: string;
  domain: string;
  ip: string;
}

export interface File {
  id: string;
  name: string;
}

export default function Dashboard() {
  const {
    data: regions,
    loading: regionsLoading,
    error: regionsError,
  } = useFetch<Region[]>("/api/v1/regions");
  const {
    data: files,
    loading: filesLoading,
    error: filesError,
    refetch: refetchFiles,
  } = useFetch<File[]>("/api/v1/files");
  const { uploadFile, uploading, uploadError, uploadSuccess } = useFileUpload();

  return (
    <div className="container mx-auto px-4 py-8">
      <header className="mb-8">
        <h1 className="text-3xl font-bold mb-2">Yato CDN</h1>
        <p className="text-gray-600">
          Manage files and view availability across regions.
        </p>
      </header>
      <div className="grid gap-8">
        <RegionInfo
          regions={regions}
          loading={regionsLoading}
          error={regionsError}
        />
        <FileUpload
          refetchFiles={refetchFiles}
          uploadFile={uploadFile}
          uploading={uploading}
          error={uploadError}
          success={uploadSuccess}
        />
        <FileList
          regions={regions}
          files={files}
          loading={filesLoading}
          error={filesError}
        />
      </div>
    </div>
  );
}
