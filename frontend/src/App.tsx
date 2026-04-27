import { useState } from "react";
import { DeploymentForm } from "./components/DeploymentForm";
import { DeploymentList } from "./components/DeploymentList";
import { LogViewer } from "./components/LogViewer";

function App() {
  const [selectedId, setSelectedId] = useState<string | null>(null);

  return (
    <div className="min-h-screen p-6 flex flex-col gap-6">
      <header>
        <h1 className="text-sm font-semibold mb-3">Deploy</h1>
        <DeploymentForm />
      </header>

      <main className="flex gap-4 flex-1">
        <section className="w-80 shrink-0 flex flex-col gap-2">
          <h2 className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
            Deployments
          </h2>
          <DeploymentList onSelect={setSelectedId} />
        </section>

        {selectedId && (
          <section className="flex-1 flex flex-col gap-2 min-w-0">
            <h2 className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
              Logs
            </h2>
            <LogViewer deploymentId={selectedId} />
          </section>
        )}
      </main>
    </div>
  );
}

export default App;
