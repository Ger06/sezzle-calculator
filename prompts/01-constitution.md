Let's create the constitution for a new project: a full-stack calculator
with a Go backend (REST API, one endpoint per operation: add, subtract,
multiply, divide, power, sqrt, percentage) and a React + TypeScript
frontend consuming it. This is a take-home assignment for a Software
Engineer II role, meant to demonstrate clean architecture, testable code,
and defensible design decisions — not a toy project.

Propose a docs/constitution.md with 8-9 non-negotiable, short and
verifiable principles, covering:

- Minimal stack (no framework/library added without justification, on
  either backend or frontend)
- The relationship between spec and code (spec governs; no behavior
  without an active spec)
- Separation between business logic and interface, on BOTH sides
  (backend: domain vs. HTTP handlers; frontend: state/reducer vs.
  components)
- Test policy: strict TDD (test first, see it fail, implement, refactor),
  not just "tests before merging"
- Data persistence: none — the calculator is stateless by design
- Trust boundary: every validation enforced in the frontend must also be
  enforced independently in the backend
- Explicit handling of cross-cutting concerns that tend to surface late
  (CORS, floating-point rounding) — they must be spec'd upfront, not
  discovered at integration time
- Language: code, identifiers, comments, and documentation in English

Keep each principle to 1-3 lines, imperative and testable (something a
reviewer could check against the code, not a vague aspiration). Max 20
lines total. Wait for my approval before writing anything else.
