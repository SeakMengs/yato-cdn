import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Region } from "./page";

interface RegionInfoProps {
  regions: Region[] | null;
  loading: boolean;
  error: Error | null;
}

export function RegionInfo({ regions, loading, error }: RegionInfoProps) {
  if (loading) {
    return (
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {[...Array(3)].map((_, i) => (
          <Card key={i}>
            <CardHeader>
              <Skeleton className="h-4 w-[150px]" />
            </CardHeader>
            <CardContent>
              <Skeleton className="h-4 w-[200px] mb-2" />
              <Skeleton className="h-4 w-[150px]" />
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-red-500">Error loading regions: {error.message}</div>
    );
  }

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
      {regions?.map((region) => (
        <Card key={region.name}>
          <CardHeader>
            <CardTitle>{region.name}</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-gray-600">Domain URL: {region.domain}</p>
            <p className="text-sm text-gray-600">IP Address: {region.ip}</p>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
