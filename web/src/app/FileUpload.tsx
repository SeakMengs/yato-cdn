import { useState, useCallback, useEffect } from "react";
import { useDropzone } from "react-dropzone";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Upload } from "lucide-react";

interface FileUploadProps {
  uploadFile: (file: File) => Promise<void>;
  refetchFiles: () => Promise<void>;
  uploading: boolean;
  error: string | null;
  success: boolean;
}

export function FileUpload({
  uploadFile,
  refetchFiles,
  uploading,
  error,
  success,
}: FileUploadProps) {
  const [file, setFile] = useState<File | null>(null);

  const onDrop = useCallback((acceptedFiles: File[]) => {
    setFile(acceptedFiles[0]);
  }, []);

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: { "image/*": [] },
  });

  const handleUpload = async () => {
    if (file) {
      await uploadFile(file);
      setFile(null);
    }
  };

  useEffect(() => {
    if (success) {
      refetchFiles();
    }
  }, [success]);

  return (
    <Card>
      <CardHeader>
        <CardTitle>Upload File</CardTitle>
      </CardHeader>
      <CardContent>
        <div
          {...getRootProps()}
          className={`border-2 border-dashed rounded-md p-8 text-center cursor-pointer ${
            isDragActive ? "border-primary" : "border-gray-300"
          }`}
        >
          <input {...getInputProps()} />
          {file ? (
            <p>{file.name}</p>
          ) : isDragActive ? (
            <p>Drop the file here ...</p>
          ) : (
            <p>Drag and drop a file here, or click to select a file</p>
          )}
        </div>
        <Button
          onClick={handleUpload}
          disabled={!file || uploading}
          className="mt-4 w-full"
        >
          {uploading ? "Uploading..." : "Upload"}
          <Upload className="ml-2 h-4 w-4" />
        </Button>
        {error && (
          <Alert variant="destructive" className="mt-4">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}
        {success && (
          <Alert className="mt-4">
            <AlertDescription>File uploaded successfully!</AlertDescription>
          </Alert>
        )}
      </CardContent>
    </Card>
  );
}
