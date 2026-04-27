const BASE_URL = "http://localhost:3001";

export type DeploymentStatus =
  | "pending"
  | "building"
  | "deploying"
  | "running"
  | "failed";

export interface Deployment {
  id: string;
  source_type: string;
  source: string;
  status: DeploymentStatus;
  image_tag: string | null;
  deploy_url: string | null;
  created_at: string;
  updated_at: string;
}

export async function createDeployment(
  sourceType: string,
  source: string,
): Promise<Deployment> {
  const res = await fetch(`${BASE_URL}/deployments`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ source_type: sourceType, source }),
  });
  if (!res.ok)
    throw new Error(`Failed to create deployment: ${res.statusText}`);
  return res.json();
}

export async function getDeployments(): Promise<Deployment[]> {
  const res = await fetch(`${BASE_URL}/deployments`);
  if (!res.ok)
    throw new Error(`Failed to fetch deployments: ${res.statusText}`);
  const data = await res.json();
  return data.deployments;
}

export function getLogsUrl(id: string): string {
  return `${BASE_URL}/deployments/${id}/logs`;
}
