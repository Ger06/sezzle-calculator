# AGENTS.md — sezzle-calculator

## Project
Full-stack calculator application. Supports the four basic operations
(addition, subtraction, multiplication, division) plus power, square root,
and percentage. Go backend exposing a REST API (one endpoint per
operation), React + TypeScript frontend consuming it. Monorepo hosting two
independent projects (`backend/`, `frontend/`) — no shared package manager
or build tooling between them.

## Stack
- Backend: Go 1.22+ (uses method-aware routing in `net/http`,
  e.g. `mux.HandleFunc("POST /add", ...)`, available since 1.22 — no
  external router needed). Exact version pinned in `go.mod`.
- Frontend: React 18+ + TypeScript, native `fetch`, `useReducer` for state.
  Exact version pinned in `package.json`.
- Backend tests: Go's built-in `testing` package
- Frontend tests: Vitest + React Testing Library

## Commands
Backend (run from `backend/`):
- Run: `go run main.go`
- Tests (unit + integration, same command — no external infra to isolate):
  `go test ./...`
- Tests with coverage: `go test ./... -cover`

Frontend (run from `frontend/`):
- Run: `npm run dev`
- Tests (unit + integration, same command): `npm run test`

## Test file naming
- Integration tests (full HTTP round-trip via `httptest`, or full component
  flow via React Testing Library) are named explicitly so they're
  distinguishable from isolated unit tests, even though they run in the
  same command: `handlers_integration_test.go` on the backend,
  `Calculator.integration.test.tsx` on the frontend.

## Style and conventions
- Go: idiomatic Go (`gofmt`-formatted), explicit error handling
  (`value, err := fn()` / `if err != nil`), no third-party HTTP framework.
- TypeScript: strict mode, function components only, no class components.
- Identifiers and code comments in English; commit messages in English.
- No new dependency (backend or frontend) without an approved spec that
  justifies why it's not achievable with the current stack.

## Naming conventions
- JSON contract (request/response bodies): `camelCase`, lowercase first
  letter — e.g. `{"a": 2, "b": 3}`, `{"result": 5}`. Go structs MUST use
  explicit `json:"..."` tags to enforce this; do not rely on Go's default
  (capitalized) field serialization.
- Go identifiers: `PascalCase` for exported (public) names, `camelCase`
  for unexported. File names: lowercase, `snake_case` if multi-word
  (`operations.go`, `operations_test.go`).
- TypeScript/React: `camelCase` for variables, functions, and hooks;
  `PascalCase` only for components and types. File names: `PascalCase`
  for component files (`Display.tsx`), `camelCase` for everything else
  (`calculatorReducer.ts`).
- REST endpoint paths: lowercase, no trailing slash (`/add`, not `/Add`
  or `/add/`).

## Architecture rules
- Business logic lives in pure functions with no HTTP/React dependency:
  `backend/calculator/` and the frontend reducer. Handlers and components
  stay thin — they parse/render and delegate, they don't compute.
- Domain errors are typed in `calculator/` and mapped to HTTP status codes
  only in `handlers/`. Business logic must not import `net/http`.
- Never trust the client: every rule enforced in the frontend (validation,
  disabled buttons) must also be enforced in the backend. Frontend-only
  validation is not acceptable as the sole guard.
- Round floating-point results in the backend before building the JSON
  response (10 decimal places). Don't round in the frontend.

## TDD workflow (mandatory)
- No business-logic code is written without a failing test first. Flow:
  write test → confirm it fails → implement the minimum to pass → refactor.
- When given a task, implement ONLY that task, following `plan.md` and
  `docs/constitution.md`. Write the test first, show it failing, then
  implement, then show it passing.
- After finishing a task: mark it done in `tasks.md`, state which
  requirement it covers, and STOP. Do not start the next task unprompted.

## Rules
- Read `docs/constitution.md` and the active spec in `specs/` before
  touching any code.
- Do not add dependencies or change the API contract (request/response
  shape) without updating the spec first.
- Do not modify files inside `specs/` unless explicitly asked to.
- CORS and floating-point rounding are explicit requirements — see the
  active spec. Don't treat them as implementation details to improvise.

## On finishing any task
- Backend: run `go test ./...` and confirm in your response that everything
  passes.
- Frontend: run `npm run test` and confirm in your response that everything
  passes.