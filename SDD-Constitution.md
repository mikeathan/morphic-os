# Morphic-OS SDD Constitution

## 1. Core Philosophy
Morphic-OS is a "Self-Synthesizing AI Operating System" built on SOLID principles and Clean Architecture. It dynamically adapts to arbitrary tasks by writing, compiling, sandboxing, testing, and registering new microservices (tools).

## 2. Clean Architecture Layers
- **Domain Layer:** The core. Contains enterprise-wide logic, entities (e.g., `Tool`, `SystemCall`), and interface definitions. No external dependencies.
- **UseCase Layer:** Contains application-specific business rules. Orchestrates the Morphic Loop, handles capability gaps, and manages the self-correction retry logic.
- **Interface/Adapter Layer:** Adapts data from UseCases to external agencies like HTTP routers (REST/gRPC) and the CLI.
- **Infrastructure Layer:** Frameworks, databases, and external APIs. Includes SQLite registry interactions, WebAssembly (Wazero) sandboxing, and multiple LLM API clients (OpenAI, Gemini, OpenRouter, Mule Router, Nvidia Build). Configuration is injected via a config structure (YAML/JSON) rather than relying solely on global environment variables.

## 2.1 Frontend Architecture
- **Abstraction and Scalability:** The frontend (Next.js) must be abstracted for future expansion. Utilize modular services, custom hooks, and a clean domain-driven or feature-based component folder structure.
- **Features:**
  - Real-time streaming of agent responses with VS Code-style collapsible UI elements.
  - Persistent conversation histories, allowing users to view previous and current interactions without losing context.

## 3. Go Idiomatic Patterns
- **Error Handling:** Use `fmt.Errorf("...: %w", err)` for wrapping errors to preserve stack context. Create custom error types for domain-specific failures (e.g., `ErrCompilationFailed`, `ErrSandboxTimeout`).
- **Context Propagation:** All external calls (DB, Sandbox, LLM, HTTP) must accept a `context.Context` to ensure proper cancellation and timeout enforcement.
- **Interface Segregation:** Define small, focused interfaces where they are *used*, not where they are implemented.
- **Concurrency:** Leverage Go routines for parallel test execution or background registry updates, managed strictly with `sync.WaitGroup` and channels.

## 4. LLM Tool Calling JSON Structure
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

## 5. Directory Structure
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