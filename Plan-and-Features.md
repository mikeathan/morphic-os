# Morphic-OS Plan and Features

## Phase 1: Foundation
**Objective:** Establish the core Go backend structure, the SQLite registry, and basic LLM integration.
- Initialize Go module in `backend/`.
- **Domain:** Create the `Tool` entity and repository interface (`domain/tool.go`).
- **Infrastructure:** Implement SQLite repository for storing, retrieving, and updating tools (`infrastructure/db/sqlite.go`). Include auto-migration.
- **UseCase:** Define the Agent interface for interacting with the LLM (`usecase/agent.go`). Implement basic context assembly.
- **Interface:** Set up a simple HTTP router (e.g., using Gin or standard net/http) to receive tasks.

## Phase 2: The Sandbox
**Objective:** Integrate Docker Engine SDK to securely execute dynamically generated code.
- **Infrastructure:** Create a Docker orchestrator (`infrastructure/docker/sandbox.go`) that can spin up ephemeral containers.
- **Features:**
  - Write source code and tests to a temporary volume/directory.
  - Execute code with strict CPU/RAM limits and timeouts.
  - Capture `stdout` and `stderr` to return to the UseCase layer.
  - Support Go and Python runtimes.

## Phase 3: The Agentic Loop
**Objective:** Implement the core "Morphic Loop" execution flow.
- **UseCase:** Implement the `MorphicLoop` in a use case.
  1. Ingest user task.
  2. Assemble context (query SQLite for active tools, convert to JSON schema).
  3. Evaluate with LLM.
  4. If capability gap: Call `sys_forge_tool`.
  5. **The Forge:** Pass generated code to the Sandbox.
  6. **Self-Correction:** Feed `stderr` back to LLM if compilation/tests fail (max 3 retries).
  7. **Registration:** Save successful tools to SQLite and re-inject them into context.

## Phase 4: The Frontend Dashboard
**Objective:** Create a reactive UI to monitor the system.
- Initialize Next.js (or Vue.js as specified in some parts of the prompt, but Next.js was heavily mentioned). Let's stick with Vue.js (or Next.js depending on final choice, Vue.js mentioned in 4. The Dashboard, Next.js mentioned in 1. ROLE). Let's clarify: We will build a web interface that can visualize the loop. Let's use standard Next.js as it's common for this.
- Create a dashboard to view the expanding tool registry.
- Display real-time logs of the Morphic Loop (evaluating, forging, testing).

## Guardrails to Implement
- **Context Management:** Limit tools sent to LLM prompt.
- **Dependencies:** Sandbox orchestrator must run `go mod tidy` or `pip install` when needed.
- **Security:** Strict container resource limits.
