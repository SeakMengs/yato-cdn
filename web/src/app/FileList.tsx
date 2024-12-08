"use client";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Skeleton } from "@/components/ui/skeleton";
import { File, Region } from "./page";
import Image from "next/image";
import Link from "next/link";
import { useFetch } from "@/hooks/useFetch";
import { CDN_URL } from "@/util";
import { useFileDelete } from "@/hooks/useFileDelete";
import { Button } from "@/components/ui/button";
import { TrashIcon } from "lucide-react";
import { useEffect } from "react";

interface FileListProps {
  regions: Region[] | null;
  files: File[] | null;
  loading: boolean;
  error: Error | null;
  refetchFiles: () => Promise<void>;
}

interface CDNFile {
  id: string;
  name: string;
  region: string;
  url: string;
}

export function FileList({
  files,
  loading,
  error,
  regions,
  refetchFiles,
}: FileListProps) {
  if (loading) {
    return (
      <div>
        <Skeleton className="h-8 w-full mb-4" />
        <Skeleton className="h-128 w-full" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-red-500">Error loading files: {error.message}</div>
    );
  }

  return (
    <div>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="cursor-pointer">File Name</TableHead>
            <TableHead>Image</TableHead>
            <TableHead className="cursor-pointer">
              Region Served From{" "}
            </TableHead>
            <TableHead>Available in Regions</TableHead>
            <TableHead>Action</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {files &&
            files.map((file) => (
              <CDNFileRow
                key={file.id}
                file={file}
                regions={regions}
                refetchFiles={refetchFiles}
              />
            ))}
        </TableBody>
      </Table>
    </div>
  );
}

function CDNFileRow({
  file,
  regions,
  refetchFiles,
}: {
  refetchFiles: () => Promise<void>;
  file: File;
  regions: Region[] | null;
}) {
  const CDN_SERVEFILE = "/api/v1/cdn/";
  const REGION_SERVEFILE = "/api/v1/files/";
  const {
    data: cdnFile,
    loading: filesLoading,
    error: filesError,
  } = useFetch<CDNFile>(`${CDN_SERVEFILE}${file.name}`);
  const { deleteFile, deleteSuccess, deleting } = useFileDelete();

  useEffect(() => {
    if (deleteSuccess) {
      refetchFiles();
    }
  }, [deleteSuccess]);

  function getCDNUrl(name: string) {
    return `${CDN_URL}${CDN_SERVEFILE}${name}`;
  }

  function getRegionUrl(name: string, regionUrl: string) {
    return `${regionUrl}${REGION_SERVEFILE}${name}`;
  }

  if (filesLoading) {
    return (
      <TableRow>
        <TableCell>
          <Skeleton className="h-8 w-128" />
        </TableCell>
        <TableCell>
          <Skeleton className="h-8 w-128" />
        </TableCell>
        <TableCell>
          <Skeleton className="h-8 w-128" />
        </TableCell>
      </TableRow>
    );
  }

  if (filesError) {
    return (
      <TableRow>
        <TableCell colSpan={5} className="text-red-500">
          Error loading files
        </TableCell>
      </TableRow>
    );
  }

  if (!cdnFile) {
    return (
      <TableRow>
        <TableCell colSpan={5}>No CDN file found</TableCell>
      </TableRow>
    );
  }

  return (
    <TableRow key={cdnFile.id}>
      <TableCell>{file.name}</TableCell>
      <TableCell>
        <Link href={getCDNUrl(cdnFile.name)}>
          <Image
            width={128}
            height={128}
            src={cdnFile.url}
            alt={cdnFile.name}
            className="w-16 h-16 object-cover"
          />
        </Link>
      </TableCell>
      <TableCell>{cdnFile.region}</TableCell>
      <TableCell>
        <div className="flex flex-col gap-2">
          {regions &&
            regions.map((region) => (
              <div className="flex flex-row gap-2 items-center" key={region.id}>
                {" "}
                <Link href={getRegionUrl(file.name, region.domain)}>
                  <Image
                    width={128}
                    height={128}
                    key={region.id}
                    src={getRegionUrl(file.name, region.domain)}
                    alt={region.name}
                    className="w-16 h-16 object-cover"
                  />
                </Link>
                <span>{region.name}</span>
              </div>
            ))}
        </div>
      </TableCell>
      <TableCell>
        <Button
          size={"icon"}
          className="bg-red-500"
          onClick={() => {
            deleteFile(cdnFile.name);
          }}
        >
          <TrashIcon className="h4 w-4" />
        </Button>
        {deleting ? "Deleting..." : ""}
      </TableCell>
    </TableRow>
  );
}
