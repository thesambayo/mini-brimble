import { useEffect, useRef, useState } from "react";
import { getLogsUrl } from "@/lib/api";

interface LogViewerProps {
  deploymentId: string;
}

export function LogViewer({ deploymentId }: LogViewerProps) {
  const [lines, setLines] = useState<string[]>([]);
  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    setLines([]);
    const es = new EventSource(getLogsUrl(deploymentId));

    es.onmessage = (e) => {
      setLines((prev) => [...prev, e.data]);
    };

    es.onerror = () => {
      es.close();
    };

    return () => es.close();
  }, [deploymentId]);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [lines]);

  return (
    <div className="h-full min-h-64 bg-zinc-950 text-green-400 font-mono text-xs p-3 overflow-y-auto">
      {lines.length === 0 ? (
        <span className="text-zinc-500">Waiting for logs…</span>
      ) : (
        lines.map((line, i) => (
          <div key={i} className="whitespace-pre-wrap leading-5">
            {line}
          </div>
        ))
      )}
      <div ref={bottomRef} />
    </div>
  );
}
