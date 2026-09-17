# sezzle-calculator

A full-stack calculator: a Go REST API backend and a React + TypeScript
frontend. Built following a spec-driven, TDD workflow — see
[Design process & decisions](#design-process--decisions) below for the
full rationale behind every architectural choice.

**Live demo:** https://sezzle-calculator-production.up.railway.app/ · **API:** https://sezzle-calculator-api-production.up.railway.app/

## Features

- Seven operations: addition, subtraction, multiplication, division,
  power (`^`, plus an `x²` shortcut), square root, percentage
- Full keyboard support (digits, operators, Enter, Escape, Backspace)
- Screen-reader accessible (`aria-live` result announcements,
  `aria-label`led controls, tested with Windows Narrator)
- Touch-accessible (44px+ tap targets, 56px on mobile)
- Domain-aware error handling (division by zero, negative square root,
  non-finite results) distinct from request validation errors

## Tech stack

| | |
|---|---|
| Backend | Go 1.27+, standard library `net/http` only (no framework) |
| Frontend | React 19 + TypeScript, Vite, native `fetch`, `useReducer` |
| Backend tests | Go's built-in `testing` package |
| Frontend tests | Vitest + React Testing Library |
| Containerization | Docker (multi-stage builds), Docker Compose |
| CI | GitHub Actions (`.github/workflows/ci.yml`) |
| Deploy | Railway (two independent services) |

## Setup

### Prerequisites
- Go 1.27+
- Node 24+
- Docker (optional, for containerized run)

### Backend

```bash
cd backend
go run main.go
```

Runs on `http://localhost:8080`. Configure the allowed CORS origin via
the `CORS_ALLOWED_ORIGIN` environment variable (defaults to
`http://localhost:5173` if unset).

### Frontend

```bash
cd frontend
npm install
npm run dev
```

Runs on `http://localhost:5173`. Configure the backend URL via the
`VITE_API_BASE_URL` build-time environment variable (defaults to
`http://localhost:8080`).

### Running both with Docker

```bash
docker-compose up --build
```

Frontend: `http://localhost:3000` · Backend: `http://localhost:8080`

## API reference

All endpoints accept `POST` with a JSON body and return JSON.

**Success shape:** `{"result": number}`
**Error shape:** `{"error": {"code": string, "message": string}}`

| Endpoint | Body | Example |
|---|---|---|
| `POST /add` | `{a, b}` | `{"a": 2, "b": 3}` → `{"result": 5}` |
| `POST /subtract` | `{a, b}` | `{"a": 5, "b": 3}` → `{"result": 2}` |
| `POST /multiply` | `{a, b}` | `{"a": 4, "b": 2.5}` → `{"result": 10}` |
| `POST /divide` | `{a, b}` | `{"a": 10, "b": 4}` → `{"result": 2.5}` |
| `POST /power` | `{base, exponent}` | `{"base": 2, "exponent": 10}` → `{"result": 1024}` |
| `POST /sqrt` | `{value}` | `{"value": 16}` → `{"result": 4}` |
| `POST /percentage` | `{value, percentage}` | `{"value": 200, "percentage": 15}` → `{"result": 30}` |

**Example with curl:**
```bash
curl -X POST http://localhost:8080/divide \
  -H "Content-Type: application/json" \
  -d '{"a": 10, "b": 0}'
# → 422 {"error":{"code":"DIVISION_BY_ZERO","message":"cannot divide by zero"}}
```

### Error codes

| Code | HTTP status | Meaning |
|---|---|---|
| `INVALID_INPUT` | 400 | Malformed JSON, missing or non-numeric field |
| `DIVISION_BY_ZERO` | 422 | Divisor is zero |
| `NEGATIVE_SQRT` | 422 | Square root of a negative number |
| `NON_FINITE_RESULT` | 422 | Result is not a finite real number (e.g. `(-8)^0.5`, overflow) |

Full field-by-field contract, edge cases, and acceptance criteria:
[`specs/001-calculator/spec.md`](specs/001-calculator/spec.md).

## Running tests

```bash
# Backend
cd backend
go test ./... -cover

# Frontend
cd frontend
npm run test -- --coverage
```

### Coverage

| Package | Coverage | Notes |
|---|---|---|
| `backend/calculator/` | 83.3% | Pure business logic |
| `backend/handlers/` | 75.0% | Critical paths (thin `main.go` wiring excluded) |
| `frontend/api/` | 100% | |
| `frontend/components/` (Display, Keypad, OperationButtons) | 100% | |
| `frontend/state/` (reducer) | 88.9% | |
| `frontend/components/Calculator.tsx` | 69.6% | Uncovered lines mostly in network-failure/keyboard-focus-guard branches |

Coverage is intentionally uneven by design, not an oversight — see
decision #13 in `docs/decisions.md` for the rationale (near-complete on
pure logic, focused on critical paths elsewhere).

## Design process & decisions

This project was built with a spec-driven, test-first workflow:

**Constitution → spec → QA review → plan → tasks → TDD implementation**

- [`docs/constitution.md`](docs/constitution.md) — non-negotiable
  project principles
- [`specs/001-calculator/spec.md`](specs/001-calculator/spec.md) — 27
  functional requirements + 5 non-functional requirements, in EARS
  notation, covering both API and UI behavior
- [`specs/001-calculator/plan.md`](specs/001-calculator/plan.md) —
  module structure, API contract, test strategy
- [`specs/001-calculator/tasks.md`](specs/001-calculator/tasks.md) —
  the implementation broken into small, ordered, TDD tasks
- [`docs/decisions.md`](docs/decisions.md) — **21 architecture
  decisions**, each with context, alternatives considered, rationale,
  trade-offs, and how it would scale. This is the best starting point
  for understanding *why* the codebase looks the way it does.

### AI tooling

This project was built using Claude (both for architectural discussion
and as the coding agent, via Claude Code) throughout the entire
process — from the constitution through implementation. Every prompt
used is saved in [`prompts/`](prompts/), in the order they were used.
Decisions were made by me; Claude was used to explore trade-offs,
generate code against an approved spec under a strict TDD discipline
(test first, confirmed failing, then implemented), and catch issues
(a QA pass on the spec surfaced 20 findings before implementation
began; manual testing afterward surfaced two more bugs, documented in
decisions #18 and #19).

## Deployment

Deployed on Railway as two independent services (backend and frontend),
each built from its own `Dockerfile`. See decision #8 in
`docs/decisions.md` for why this is a single-domain, layered service
rather than a distributed microservice architecture.

CI runs on every push/PR to `main` via GitHub Actions
(`.github/workflows/ci.yml`): backend and frontend test suites run in
parallel and must pass before a PR is merged.