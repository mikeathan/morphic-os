# Morphic-OS SDD Constitution

## 1. Project Overview & Core Philosophy
`morphic-os` is a Self-Synthesizing AI Operating System designed for high-performance local nodes.
Built on SOLID principles and Clean Architecture, it dynamically adapts to arbitrary tasks by writing, compiling, sandboxing, testing, and registering new microservices (tools).
**The Core Paradigm:** When the system encounters a capability gap, the Agentic Loop writes AssemblyScript/TinyGo code, the Go Kernel compiles it into a WebAssembly (`.wasm`) binary, and the binary is stored as a BLOB in SQLite. The system executes these tools in microseconds via the `wazero` runtime.

## 2. Target Use Cases
The architecture is designed to support these autonomous scenarios:
1. **The Zero-Token Auditor:** A scheduled background Wasm daemon scrapes financial/system data and uses IPC to pipe it to an analyzer tool. The LLM is only invoked if a severe anomaly is detected, saving API costs.
2. **The Bulletproof Data Sanitizer:** Sensitive user documents uploaded to the VFS are processed by Ring 3 Wasm tools with strict WASI network blocks. The data is anonymized locally before reaching a cloud LLM.
3. **The Infinite-Context Researcher:** The user feeds the OS hundreds of PDFs. The system pages old context out into the `sqlite-vss` vector DB and dynamically pages it back into the prompt only when relevant, achieving infinite memory.
4. **The Autonomous Local CI/CD:** A daemon watches a local project folder, triggers a Wasm test suite on save, and if tests fail, pipes the error logs directly to the LLM to write a fix.
5. **The System-Safe Web Crawler:** Spawning parallel Wasm workers to scrape unstructured sites. If a worker encounters a memory leak and panics, it dies cleanly within the Wasm runtime without crashing the Go Kernel or exposing the host OS.

## 3. System Constitution & OS Primitives
The architecture strictly separates the Kernel (Go), the Registry (SQLite), and the Userspace (Wasm tools).

### A. Nested Sandboxing (Docker + Wasm)
- **The Outer Sandbox (Docker):** The Go Kernel and frontend are deployed via `docker-compose`. Uses a multi-stage build: compiling the Go binary in a `golang:alpine` container and running it in an ultra-lightweight `alpine:latest` container.
- **The Inner Sandbox (Wazero):** Forged tools operate strictly inside the `wazero` runtime.
- **Development vs. Production:** During active IDE development (`MORPHIC_ENV=dev`), the Outer Sandbox is temporarily disabled (running natively via `go run`) to allow for Delve debugger attachment.

### B. The Unix-Like Virtual File System (VFS)
- WASI filesystem calls must be intercepted by the Go Kernel and routed to a `VirtualFiles` SQLite table.
- **Structure:** Mimics Unix (`/var/logs` for outputs, `/home/agent` for user data, `/usr/bin` for tool metadata).

### C. The Secrets Vault (Encrypted Storage)
- **Encryption:** AES-256-GCM. API keys are stored in a `Secrets` SQLite table. Wasm tools never see the raw keys.
- The Go Kernel intercepts network requests, injects the decrypted headers, and proxies the call.
- **Master Key:** Loaded via a local `.env` file passed into the Docker container environment.

### D. Cognitive Memory & Synaptic Pruning
- **Implementation:** Use `sqlite-vss` for vector embeddings.
- **The Sleep Cycle (Nightly Daemon):**
  1. Consolidates raw VFS chat logs into facts/preferences.
  2. Evaluates existing vectors based on "Access Frequency" and "Time Since Last Recall". Irrelevant data is permanently deleted unless flagged as "Core Memory".
  3. Purges raw daily logs.

### E. Privilege Rings & IPC
- **Ring 0:** Go Kernel (Full access within the Docker container).
- **Ring 1:** Trusted Daemons (Access to specific VFS paths and network).
- **Ring 3:** Untrusted/New Wasm Tools (Compute only).
- **IPC (Event Bus):** Tools communicate via a Go-managed pub/sub channel.

## 4. UI / UX Dashboard Specifications
The Vite/React/Next.js frontend must include:
1. **The Morphic Map (React Flow):** Visualizer showing all active Wasm tools, daemons, and data pipes.
2. **Terminal / stdout:** Live-streaming view of the agent's internal monologue and compiler outputs.
3. **VFS Explorer:** A file-tree UI to browse `/var/logs` or `/home/agent`.
4. **Hardware & Metrics:** Real-time displays of node utilization and memory pruning stats.
5. **Secrets Manager:** Secure UI to input API keys.

## 5. Clean Architecture & Code Guidelines

### 5.1 Clean Architecture Layers
- **Domain Layer:** The core. Contains enterprise-wide logic, entities (e.g., `Tool`, `SystemCall`), and interface definitions. No external dependencies.
- **UseCase Layer:** Contains application-specific business rules. Orchestrates the Morphic Loop, handles capability gaps, and manages the self-correction retry logic.
- **Interface/Adapter Layer:** Adapts data from UseCases to external agencies like HTTP routers (REST/gRPC) and the CLI.
- **Infrastructure Layer:** Frameworks, databases, and external APIs. Includes SQLite registry interactions, WebAssembly (Wazero) sandboxing, and multiple LLM API clients (OpenAI, Gemini, OpenRouter, Mule Router, Nvidia Build). Configuration is injected via a config structure (YAML/JSON) rather than relying solely on global environment variables.

### 5.2 Frontend Architecture
- **Abstraction and Scalability:** The frontend (Next.js) must be abstracted for future expansion. Utilize modular services, custom hooks, and a clean domain-driven or feature-based component folder structure.
- **Features:**
  - Real-time streaming of agent responses with VS Code-style collapsible UI elements.
  - Persistent conversation histories, allowing users to view previous and current interactions without losing context.

### 5.3 Go Idiomatic Patterns
- **Error Handling:** Use `fmt.Errorf("...: %w", err)` for wrapping errors to preserve stack context. Create custom error types for domain-specific failures (e.g., `ErrCompilationFailed`, `ErrSandboxTimeout`).
- **Context Propagation:** All external calls (DB, Sandbox, LLM, HTTP) must accept a `context.Context` to ensure proper cancellation and timeout enforcement.
- **Interface Segregation:** Define small, focused interfaces where they are *used*, not where they are implemented.
- **Concurrency:** Leverage Go routines for parallel test execution or background registry updates, managed strictly with `sync.WaitGroup` and channels.

### 5.4 WebAssembly (Wasm) Execution Strategy
Instead of implementing complex Wasm memory bridge abstractions (e.g., memory allocation/retrieval via exported Wasm functions), Morphic-OS utilizes `wasip1` to leverage standard Go execution patterns natively.
- Generated tools are compiled using `GOOS=wasip1 GOARCH=wasm go build`.
- Inputs to tools are passed as standard command-line arguments and read natively within the generated Wasm using `os.Args`.
- Outputs and errors are natively captured via standard `stdout` and `stderr` streams by configuring the `wazero` sandbox environment `WithArgs`, `WithStdout`, and `WithStderr`.
- Executing a `wasip1` binary that successfully exits will result in a `sys.ExitError` with exit code `0`, which the Sandbox correctly handles as a success state.
- **Workspaces and Isolation:** The Sandbox environment isolates execution per Workspace. It mounts a specific host directory (`WithFS`) providing isolated file system access, and injects workspace-specific environment variables (`WithEnv`), e.g., for database connectivity or API keys.
- **Configurable Timeouts:** Execution timeouts are configurable based on the workspace or tool complexity.
- **Real-Time Streaming:** Tool execution standard output is captured and can be streamed in real-time, allowing long-running tasks to be observed.

## 6. LLM Tool Calling JSON Structure
To standardize interactions, tools will use OpenAI's JSON schema format for function calling:

```json
{
  "type": "function",
  "function": {
    "name": "tool_name",
    "description": "Clear description of what the tool does.",
    "parameters": {
      "type": "object",
      "properties": {
        "arg1": {
          "type": "string",
          "description": "Description of arg1"
        }
      },
      "required": ["arg1"]
    }
  }
}
```

## 7. Directory Structure
```
morphic-os/
├── backend/
│   ├── cmd/
│   │   └── server/          # Main application entrypoint
│   ├── internal/
│   │   ├── domain/          # Entities and core interfaces
│   │   ├── usecase/         # Morphic loop and business logic
│   │   ├── interface/       # HTTP handlers, CLI routers
│   │   └── infrastructure/  # SQLite, Wazero SDK, LLM clients
│   ├── pkg/                 # Reusable utility functions
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── components/      # UI components
│   │   ├── app/             # Next.js App router or views
│   │   └── store/           # State management
│   ├── package.json
│   └── next.config.js       # Or vite.config.js for Vue.js
├── scripts/                 # Build and deployment scripts
└── wasm/                    # Base WebAssembly configurations
```