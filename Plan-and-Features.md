# Morphic-OS Plan and Features

## Phase 1: Kernel Foundation & Docker Scaffolding
**Objective:** Establish the core Go backend structure, the SQLite registry, and basic LLM integration within Docker.
- Initialize Go module in `backend/`.
- Initialize SQLite/sqlite-vss connection, schemas (`Tools`, `VirtualFiles`, `Secrets`, `KnowledgeBase`).
- **[Completed]** Implement the VFS interface.
- **[Completed]** Set up Docker configurations (multi-stage Debian build with `golang:1.24-bookworm` and `debian:bookworm-slim` for CGO support).
- **Domain:** Create the `Tool` entity and repository interface (`domain/tool.go`).
- **Infrastructure:** Implement SQLite repository for storing, retrieving, and updating tools (`infrastructure/db/sqlite.go`). Include auto-migration.
- **UseCase:** Define the Agent interface for interacting with the LLM (`usecase/agent.go`). Implement basic context assembly.
- **Interface:** Set up a simple HTTP router (e.g., using standard net/http) to receive tasks.

## Phase 2: The Wasm Execution Bridge
**Objective:** Integrate `wazero` to securely execute dynamically generated WebAssembly code and map WASI to the VFS.
- **Infrastructure:** Create a WASM orchestrator (`infrastructure/wasm/sandbox.go`) that can instantiate and execute WebAssembly modules using `github.com/tetratelabs/wazero`.
- **[Completed]** Build the `sys_forge_tool` compilation-to-database pipeline.
- **Features:**
  - **[Completed]** Compile generated Go code to WebAssembly (`GOOS=wasip1 GOARCH=wasm go build`) via OS Compiler Orchestrator.
  - Execute the compiled WASM module in an isolated Wazero runtime with strict memory limits and timeouts. Support parsing Wasm execution arguments using `wasip1` to read from `os.Args`.
  - Capture `stdout` and `stderr` to return to the UseCase layer.
  - **[Completed]** Map WASI filesystem calls to the `VirtualFiles` SQLite table via the Go Kernel.

## Phase 3: Security, Daemons, & Memory
**Objective:** Implement AES Secrets vault, Go Task Scheduler, and the Nightly Sleep Cycle.
- **[Completed]** Implement the AES-256-GCM Secrets vault. Wasm tools never see raw keys; the Go Kernel intercepts and proxies network calls.
- **[Completed]** Implement the Go Task Scheduler for background tasks (e.g., Zero-Token Auditor).
- **[Completed]** Implement the IPC Event Bus for tool-to-tool communication.
- **[Completed]** Implement the Nightly Sleep Cycle (Pruning):
  - **[Completed]** Consolidate raw VFS chat logs into facts/preferences.
  - **[Completed]** Evaluate existing vectors based on access frequency and time since last recall.
  - **[Completed]** Purge raw daily logs.

## Phase 4: The Agentic Loop
**Objective:** Build the dynamic LLM context builder and self-correction compilation loop.
- **UseCase:** Implement the `MorphicLoop` in a use case.
  1. Ingest user task.
  2. Assemble context (query SQLite for active tools, convert to JSON schema).
  3. Evaluate with LLM.
  4. If capability gap: Call `sys_forge_tool`.
  5. **The Forge:** Pass generated code to the Sandbox.
  6. **Self-Correction:** Feed `stderr` back to LLM if compilation/tests fail (max 3 retries).
  7. **Registration:** Save successful tools to SQLite and re-inject them into context.

## Phase 5: The Command Center
**Objective:** Build the React/Vite/Next.js UI and Go SSE APIs.
- Initialize Next.js/Vite app.
- **[Completed]** Build the Go SSE APIs for real-time communication.
- Implement the Dashboard features:
  - **[Completed]** **The Morphic Map (React Flow):** Visualizer showing active tools, daemons, and animated IPC Event Flow pipes.
  - **[Completed]** **Terminal / stdout:** Live-streaming view with Compiler Forge syntax highlighting.
  - **[Completed]** **VFS Explorer:** File-tree UI to browse `/var/logs` or `/home/agent`.
  - **[Completed]** **Hardware & Metrics:** Real-time displays of node utilization and pruning stats.
  - **[Completed]** **Secrets Manager:** Secure UI to input API keys.
  - **[Completed]** **Cognitive Memory Dashboard:** Display extracted Core Memories and user preferences.