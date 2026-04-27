import { useQuery } from "@tanstack/react-query";
import { getDeployments } from "@/lib/api";
import { DeploymentCard } from "./DeploymentCard";

interface DeploymentListProps {
  onSelect: (id: string) => void;
}

export function DeploymentList({ onSelect }: DeploymentListProps) {
  const { data, isLoading, isError } = useQuery({
    queryKey: ["deployments"],
    queryFn: getDeployments,
    // refetchInterval: 3000,
  });

  if (isLoading)
    return <p className="text-xs text-muted-foreground">Loading…</p>;
  if (isError)
    return (
      <p className="text-xs text-destructive">Failed to load deployments.</p>
    );
  if (!data?.length)
    return <p className="text-xs text-muted-foreground">No deployments yet.</p>;

  return (
    <div className="flex flex-col gap-2">
      {data.map((deployment) => (
        <DeploymentCard
          key={deployment.id}
          deployment={deployment}
          onSelect={onSelect}
        />
      ))}
    </div>
  );
}
