# README Notes — Mini Brimble

> Running list of observations, tradeoffs, and decisions to reference when writing the final README.

---

## How It Works

A one-page UI lets you submit a Git URL. The backend clones the repo, builds it into a Docker image using Railpack, runs it as a container, and registers a subdomain route in Caddy — all without any Dockerfiles. Build logs stream to the UI in real time over SSE.

**Stack:** Go backend, Vite + React + TanStack frontend, Caddy ingress, Railpack builder, BuildKit, SQLite.

---

## Architecture

```
User submits Git URL
→ Backend clones repo into /app/workspace/{deploymentId}
→ Backend invokes Railpack as subprocess → BuildKit builds the image
→ Backend runs container via Docker CLI (docker run)
→ Backend registers subdomain route in Caddy via Admin API
→ App live at {deploymentId[:8]}.localhost
→ Logs streamed throughout via SSE
```

---

## Architecture Decisions & Tradeoffs

- **Railpack over hand-written Dockerfiles** — Railpack detects the language and builds the image automatically. Zero Dockerfile maintenance. Tradeoff: larger images (~112MB for Go, ~151MB for React) because Railpack bundles a general-purpose runtime base.
- **Docker CLI over Docker SDK** — backend invokes `docker run` via `exec.Command` rather than the Docker SDK for Go. Simpler and more readable. Tradeoff: backend depends on Docker CLI being installed in the container — handled via the Dockerfile.
- **Subdomain routing over path-based routing** — `{id}.localhost` instead of `localhost/app-{id}`. Path-based routing breaks static asset resolution because the browser requests assets at `/static/styles.css` not `/app-{id}/static/styles.css`. Subdomains solve this cleanly — exactly how Brimble and Vercel work in production.
- **SQLite over Postgres** — simpler setup, zero external dependencies, perfect for a single-node deployment pipeline. Tradeoff: can't scale horizontally.
- **SSE over WebSockets** — SSE is simpler than WebSockets for one-directional streaming (server → client). Build logs only flow one way, so SSE is the right tool.
- **Caddy JSON config over Caddyfile** — using a JSON config file instead of a Caddyfile gives explicit control over the server name (`main`). Caddyfile auto-generates names (`srv0`) which breaks dynamic route registration via the Admin API.
- **In-memory log channels** — each active deployment gets a `chan string` stored in a map on the Application struct. The pipeline writes to it; the SSE handler reads from it. Simple and effective for a single-node setup.

---

## Interesting Observations

- Railpack automatically detects static sites (Vite/React) and bakes Caddy into the image to serve the built `dist` folder — no Node server needed at runtime.
- The Caddy baked into app images has its Admin API disabled — it only serves files. The standalone Caddy ingress service is separate with Admin API enabled.
- Railpack internally uses `mise` to manage language runtimes. When running inside Docker, mise must be pre-installed and symlinked to the exact version path Railpack expects (`/tmp/railpack/mise/mise-{version}`) — otherwise Railpack tries to download it at build time, causing slow or failed builds.
- BuildKit cache is lost on container restart unless a named volume is mounted at `/var/lib/buildkit`. Without it, every restart triggers a full re-download of base images (~300MB).
- `docker-container://buildkit` connection method works by using the mounted Docker socket to exec into the BuildKit container. This requires the Docker socket to be mounted into the backend container.

---

## What I'd Do With More Time

- **Caddy route recovery on restart** — on startup, read all `running` deployments from SQLite and re-register their Caddy routes. Currently routes are lost when Caddy restarts.
- **Container cleanup** — stop and remove containers when a deployment is deleted or redeployed.
- **Build caching per repo** — pass `--cache-key` to Railpack so subsequent deploys of the same repo reuse cached layers.
- **SSL/TLS** — Caddy handles this automatically once you have a real domain.
- **Upload support** — currently only Git URLs are supported. Adding zip upload would complete the source type support.
- **Frontend service in compose** — wire the Vite frontend into docker compose so the whole stack comes up with one command.

---

## What I'd Rip Out

- The random host port assignment in `Run()` — it's unused now that containers communicate via the compose network. It's dead code that should be removed.
- The `containerPort` parameter from `RegisterRoute` — same reason, no longer needed after switching to subdomain routing via container name.
