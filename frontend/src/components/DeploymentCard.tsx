import { format } from "date-fns";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import type { Deployment, DeploymentStatus } from "@/lib/api";

const statusColor: Record<DeploymentStatus, string> = {
  pending: "bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400",
  building: "bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400",
  deploying: "bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400",
  running: "bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400",
  failed: "bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400",
};

interface DeploymentCardProps {
  deployment: Deployment;
  onSelect: (id: string) => void;
}

export function DeploymentCard({ deployment, onSelect }: DeploymentCardProps) {
  return (
    <div className="border border-border p-3 flex flex-col gap-2">
      <div className="flex items-center justify-between gap-2">
        <span className="text-xs text-muted-foreground truncate">{deployment.source}</span>
        <Badge className={cn("shrink-0", statusColor[deployment.status])}>
          {deployment.status}
        </Badge>
      </div>

      {deployment.image_tag && (
        <p className="text-xs text-muted-foreground font-mono">{deployment.image_tag}</p>
      )}

      {deployment.deploy_url && (
        <a
          href={deployment.deploy_url}
          target="_blank"
          rel="noreferrer"
          className="text-xs text-primary underline-offset-4 hover:underline truncate"
        >
          {deployment.deploy_url}
        </a>
      )}

      <div className="flex items-center justify-between">
        <span className="text-xs text-muted-foreground">
          {format(new Date(deployment.created_at), "MMM d, yyyy HH:mm")}
        </span>
        <Button size="xs" variant="outline" onClick={() => onSelect(deployment.id)}>
          View Logs
        </Button>
      </div>
    </div>
  );
}
