# Demo Log — Mini Brimble

> Inspired by the Demo-First Project System.
> Each demo proves one piece works in isolation before connecting it to anything else.

---

## Demo Log

```
[x] Demo 1:  Railpack builds a sample Go HTTP server locally
[x] Demo 2:  Railpack builds a sample Vite React app locally
[x] Demo 3:  Go backend runs locally — POST /deployments + GET /deployments works
[x] Demo 4:  Backend clones a Git URL into a local workspace folder
[x] Demo 5:  Backend invokes Railpack as a subprocess and streams output to terminal
[x] Demo 6:  SSE endpoint works — curl receives a stream of lines
[x] Demo 7:  Docker socket mounted — backend starts the built container programmatically
[x] Demo 8:  Caddy Admin API — backend registers a route to the running container
[ ] Demo 9:  Frontend up — form submits a Git URL, list shows deployments and status
[ ] Demo 10: Frontend streams live build logs over SSE
[ ] Demo 11: Full pipeline end to end
[ ] Demo 12: docker compose up brings the whole stack up on a clean machine
[ ] Demo 13: Deployed on Brimble + feedback written
[ ] Demo 14: README + optional Loom walkthrough
```

---

## Demo Details

### Demo 1: Railpack builds a sample Go HTTP server locally
Railpack installed locally. A simple Go HTTP server sits in a folder. Run `railpack build` on it, see it produce a Docker image, and run that image successfully.
**Goal:** Understand what Railpack does before touching any backend code.

---

### Demo 2: Railpack builds a sample Vite React app locally
Same as Demo 1 but with a Vite React app.
**Goal:** Confirm Railpack handles both app types you'll be testing against throughout the project.

---

### Demo 3: Go backend runs locally — POST /deployments + GET /deployments works
No pipeline yet. Just the Go API running, SQLite connected, and basic deployment CRUD working. Test with curl.
**Goal:** API skeleton is alive and storing state.

---

### Demo 4: Backend clones a Git URL into a local workspace folder
When you POST a Git URL, the backend clones it into a temp folder. Log the path to terminal.
**Goal:** The first real pipeline step works — source code is on disk.

---

### Demo 5: Backend invokes Railpack as a subprocess and streams output to terminal
Backend runs `railpack build` on the cloned repo as a subprocess. Railpack's output prints live in the terminal. No SSE yet — just confirm the subprocess works and produces an image.
**Goal:** Backend can drive Railpack programmatically.

---

### Demo 6: SSE endpoint works — curl receives a stream of lines
Add a `GET /deployments/:id/logs` endpoint that streams lines over SSE. Wire it to the Railpack subprocess output. Test with curl — log lines arrive in real time while the build runs.
**Goal:** Streaming works before the frontend touches it.

---

### Demo 7: Docker socket mounted — backend starts the built container programmatically
Backend calls the Docker socket to run the image Railpack just built. Container starts, you can curl it on its port.
**Goal:** Backend owns the full build → run pipeline.

---

### Demo 8: Caddy Admin API — backend registers a route to the running container
Caddy running locally via Docker. After the container starts, backend POSTs to Caddy's Admin API to add a route. Hitting that route in the browser reaches the running container.
**Goal:** Dynamic routing works.

---

### Demo 9: Frontend up — form submits a Git URL, list shows deployments and status
Vite + TanStack Router + Query frontend running. Submit a Git URL, see it appear in the deployments list with a status that updates. No live logs yet.
**Goal:** Frontend and backend are connected.

---

### Demo 10: Frontend streams live build logs over SSE
Frontend connects to the SSE endpoint and renders log lines in real time while a build is running. Logs persist so you can scroll back after the build finishes.
**Goal:** The most important UI feature works.

---

### Demo 11: Full pipeline end to end
Submit a Git URL → backend clones it → Railpack builds it → container starts → Caddy route registered → live URL appears in the UI → logs streamed throughout.
**Goal:** Everything works together for the first time.

---

### Demo 12: docker compose up brings the whole stack up on a clean machine
Frontend, backend, Caddy, and BuildKit all come up with one command. The pipeline works exactly as in Demo 11 but fully containerized.
**Goal:** Hard requirement met.

---

### Demo 13: Deployed on Brimble + feedback written
Deploy frontend (or any hello world) on Brimble. Write honest feedback — bugs, friction, confusing UI, missing features.
**Goal:** Mandatory submission requirement done.

---

### Demo 14: README + optional Loom walkthrough
README covers setup instructions, architecture decisions, what you'd do with more time. Loom is a 5-10 minute walkthrough of the running system.
**Goal:** Project is closed and submittable.

---

## The Critical Path

| Stage | Demos | Focus |
|-------|-------|-------|
| Exploration | 1 – 2 | No backend code — just understand Railpack |
| Backend only | 3 – 8 | Don't touch the frontend yet |
| Frontend | 9 – 10 | Wire in once backend is solid |
| Integration | 11 | The moment it all clicks |
| Wrap up | 12 – 14 | Containerize, deploy, document |

---

## The Rule
> If it's been more than a day since the last demo, something is wrong.
> Ask: *"What's the next smallest thing I can do to earn a new entry?"*
