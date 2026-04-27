import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { createDeployment } from "@/lib/api";

export function DeploymentForm() {
  const [url, setUrl] = useState("");
  const [loading, setLoading] = useState(false);
  const queryClient = useQueryClient();

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!url.trim()) return;
    setLoading(true);
    try {
      await createDeployment("git", url.trim());
      setUrl("");
      await queryClient.invalidateQueries({ queryKey: ["deployments"] });
    } finally {
      setLoading(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} className="flex gap-2">
      <Input
        type="url"
        placeholder="https://github.com/user/repo"
        value={url}
        onChange={(e) => setUrl(e.target.value)}
        className="flex-1"
        required
      />
      <Button type="submit" disabled={loading}>
        {loading ? "Deploying…" : "Deploy"}
      </Button>
    </form>
  );
}
