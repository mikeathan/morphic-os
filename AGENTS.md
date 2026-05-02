# Morphic-OS Agent Guidelines

Welcome to the Morphic-OS codebase. When working on this project, you must adhere to the following global rules:

## 1. Documentation & Constitution Updates
- **Every structural, architectural, or significant feature step must update the `SDD-Constitution.md` and/or `Plan-and-Features.md` documents.**
- The `SDD-Constitution.md` is the absolute source of truth for architectural guidelines and implementation details. Keep it up to date with the latest multi-LLM approaches, backend configurations, and frontend structures.

## 2. SOLID and Clean Architecture
- **Backend (Go):** Strictly follow Clean Architecture (Domain, UseCase, Interface, Infrastructure). Utilize Interface Segregation. Configuration must be injected, not globally accessed. Avoid tight coupling to specific LLM providers.
- **Frontend (Next.js):** Follow senior engineering practices. Abstract components, utilize feature-based or standard domain-driven folder structures, manage state cleanly, and prefer strict TypeScript typing over `any`.

## 3. The Morphic Loop
- The core loop is sacrosanct. Any changes to how tools are generated, compiled, or executed must maintain the capability for self-correction (feeding `stderr` and errors back to the LLM) and secure sandboxing.
