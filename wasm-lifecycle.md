# Morphic-OS: WebAssembly Lifecycle & Execution Flow

## 1. Objective
This document defines the exact step-by-step lifecycle of a "Forged Tool" within `morphic-os`. Agents building this system must strictly adhere to this pipeline when implementing the `internal/usecase` and `internal/infrastructure/wasm` packages.

## 2. Phase 1: The Forge (Code Generation & Compilation)
When the LLM identifies a capability gap, it outputs source code (AssemblyScript/TinyGo) designed to solve the problem.

**Execution Steps:**
1. **Receive Source:** The Go Kernel receives the raw `source_code` string and a `JSON Schema` defining the tool's inputs/outputs.
2. **Write Temp File:** The Kernel writes the `source_code` to a temporary directory (e.g., `/tmp/morphic/build/tool_name.ts`).
3. **Invoke Compiler:** The Kernel executes a system shell command to invoke the local compiler.
   * *Example:* `asc /tmp/morphic/build/tool_name.ts -b /tmp/morphic/build/tool_name.wasm`
4. **Read Binary:** If compilation succeeds, the Kernel reads the output `.wasm` file into memory as a `[]byte`.
5. **Cleanup:** The temporary files are deleted.

## 3. Phase 2: Registry & Persistence
The system must persist the tool so it can be used instantly in the future without recompilation.

**Execution Steps:**
1. **Database Transaction:** The Kernel initiates a SQLite transaction.
2. **Insert Record:** It inserts the Tool metadata (Name, Description, JSON Schema, Source Code) AND the `[]byte` Wasm binary into the `WasmBinary` (BLOB) column.
3. **Context Update:** The Kernel reloads the active tool schemas into the LLM's system prompt context.

## 4. Phase 3: The Wasm Bridge (Execution & Memory Management)
When the LLM calls an existing tool, the Kernel must load the BLOB and execute it via `wazero`. This is the most technically complex phase.

**Execution Steps:**
1. **Fetch Binary:** The Kernel queries SQLite for the `WasmBinary` BLOB associated with the tool name.
2. **Instantiate Module:** `wazero` instantiates the Wasm module from the `[]byte`.
3. **Memory Allocation (Host to Guest):**
   * WebAssembly only understands numbers, not strings.
   * The Go Kernel must call an exported Wasm function (e.g., `allocate(size)`) to reserve memory inside the Wasm guest for the JSON arguments.
   * The Kernel writes the JSON argument string into that specific memory offset.
4. **Execution:** The Kernel calls the target Wasm function, passing the memory pointer and length of the input payload.
5. **Memory Retrieval (Guest to Host):**
   * The Wasm function processes the data, writes the JSON result to its own memory, and returns a new pointer/length to the Go host.
   * The Go Kernel reads the bytes at that pointer, parses them as a JSON string, and returns the result to the LLM context.
6. **Module Closure:** The `wazero` instance is closed to free resources.
