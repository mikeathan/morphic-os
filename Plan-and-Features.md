# Morphic-OS Plan and Features

## Phase 1: Foundation (Completed)
**Objective:** Establish the core Go backend structure, the SQLite registry, and basic LLM integration.
- Initialize Go module in `backend/`.
- **Domain:** Create the `Tool` entity and repository interface (`domain/tool.go`).
- **Infrastructure:** Implement SQLite repository for storing, retrieving, and updating tools (`infrastructure/db/sqlite.go`). Include auto-migration.
- **UseCase:** Define the Agent interface for interacting with the LLM (`usecase/agent.go`). Implement basic context assembly.
- **Interface:** Set up a simple HTTP router (e.g., using standard net/http) to receive tasks.

## Phase 2: The Sandbox (Completed)
**Objective:** Integrate `wazero` to securely execute dynamically generated WebAssembly code.
- **Infrastructure:** Create a WASM orchestrator (`infrastructure/wasm/sandbox.go`) that can instantiate and execute WebAssembly modules using `github.com/tetratelabs/wazero`.
- **Features:**
  - Compile generated Go code to WebAssembly (`GOOS=wasip1 GOARCH=wasm go build`).
  - Execute the compiled WASM module in an isolated Wazero runtime with strict memory limits and timeouts. Support parsing Wasm execution arguments using `wasip1` to read from `os.Args`.
  - Capture `stdout` and `stderr` to return to the UseCase layer.

## Phase 3: The Agentic Loop (Completed)
**Objective:** Implement the core "Morphic Loop" execution flow.
- **UseCase:** Implement the `MorphicLoop` in a use case.
  1. Ingest user task.
  2. Assemble context (query SQLite for active tools, convert to JSON schema).
  3. Evaluate with LLM.
  4. If capability gap: Call `sys_forge_tool`.
  5. **The Forge:** Pass generated code to the Sandbox.
  6. **Self-Correction:** Feed `stderr` back to LLM if compilation/tests fail (max 3 retries).
  7. **Registration:** Save successful tools to SQLite and re-inject them into context.

## Phase 4: The Frontend Dashboard (Completed)
**Objective:** Create a reactive UI to monitor the system.
- Initialize Next.js app.
- Create a dashboard to view the expanding tool registry (`ToolRegistry.tsx`).
- Display real-time logs of the Morphic Loop (evaluating, forging, testing) using Server-Sent Events (`TaskDisplay.tsx` and `TaskSubmissionForm.tsx`).

## Guardrails to Implement (Completed)
- **Context Management:** Limit tools sent to LLM prompt (e.g., limited to 10 tools max).
- **Dependencies:** Sandbox orchestrator runs `go mod init mytool` and `go mod tidy` when needed.
- **Security:** Strict container resource limits implemented using Wazero's `WithMemoryLimitPages(1024)`.

## Phase 5: Workspaces and Advanced Sandboxing (In Progress)
**Objective:** Add robust isolation and configuration for workspaces, supporting tools that require databases and external configurations.
- **Workspaces:** Implement `Workspace` domain models. Link tools to a specific `WorkspaceID`.
- **Sandbox Isolation:** Mount workspace-specific directories into the Wazero runtime (`WithFS`) and inject workspace environment variables (`WithEnv`), enabling the tools to access local files safely and use external API/Database credentials.
- **Configurable Timeouts:** Make execution timeouts configurable instead of hardcoded to 30 seconds.
- **Real-Time Log Streaming:** Stream the WASM tools' standard output via the `/api/logs` endpoint, so the user can see what a long-running tool is doing live.
