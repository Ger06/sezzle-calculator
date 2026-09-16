# Design Decisions — Sezzle Calculator

Log of architecture decisions made before and during implementation. Format: context, alternatives considered, final decision, accepted trade-off, and how it would scale.

---

## 1. Calculator type: direct operations, not expression evaluation

**Context:** the assignment asks for addition, subtraction, multiplication, division, and optionally power/square root/percentage. One possible interpretation was a "scientific" calculator that evaluates full expressions (e.g. `3 + 4 * 2 - 1`).

**Alternatives considered:**
- Expression calculator: a single `POST /evaluate` endpoint that receives a string and parses it respecting operator precedence.
- Direct calculator: one operation with its operands per request.

**Decision:** direct calculator. Each operation is an isolated call with explicit operands.

**Why:** the assignment lists discrete operations, not arbitrary expression evaluation. An expression parser (tokenizer, precedence, parentheses) is a project on its own that would have consumed the entire timebox (2-4 hours), leaving no room for frontend, tests, or documentation.

**Trade-off:** the calculator doesn't support chaining operations in a single input (`3+4*2`). If that were required later, it's a major design change, not an incremental extension.

**How it would scale:** if the domain genuinely needed complex expressions, the right move would be evaluating an existing parsing library rather than writing a custom parser — but there was no justification for that complexity within the given scope.

---

## 2. Endpoint design: one per operation, not a generic `/calculate`

**Context:** operations have different arities — `add/subtract/multiply/divide` take `{a, b}`, `sqrt` takes `{value}`, `power` takes `{base, exponent}`, `percentage` takes `{value, percentage}`.

**Alternatives considered:**
- Single `POST /calculate` endpoint with `{operation: string, ...operands}` and an internal function map.
- One REST endpoint per operation (`POST /add`, `POST /divide`, etc.), each with its own typed request struct.

**Decision:** one endpoint per operation.

**Why:** Go is a statically typed language. With separate endpoints, each one declares its own request struct (`AddRequest{A, B float64}`, `SqrtRequest{Value float64}`), leveraging the language's type system. A single endpoint would force simulating dynamic typing (generic payload with optional fields or `map[string]interface{}`), fighting the language instead of using it.

**Trade-off:** adding a new operation requires a new route, a new handler, and a new struct — more code surface than adding an entry to a function map.

**How it would scale:** for a finite, known set of operations (this case), this is the right choice. If the domain were "extensible calculator with dynamic operations added frequently," the decision would shift toward `/calculate` + function map, prioritizing extensibility over strong per-operation typing.

---

## 3. Repository structure: monorepo

**Context:** frontend (React/TS) and backend (Go) are two projects with different toolchains.

**Decision:** monorepo, with `backend/` and `frontend/` as sibling folders inside the same Git repository.

**Why:** for a project of this size and timebox, a monorepo simplifies setup, review, and delivery (a single repo link, a single root README). Splitting into two repos adds coordination overhead with no real benefit at this scale.

**Trade-off:** in a monorepo, CI/CD pipelines for each part can trigger unnecessarily if path-specific triggers aren't configured. Not a real issue here since no CI is set up for this exercise.

**How it would scale:** if frontend and backend had independent release cycles or separate teams owning them, splitting into two repos (or a monorepo with tooling like Nx/Turborepo) would be the right call.

---

## 4. Definition of the "percentage" operation

**Context:** "percentage" is ambiguous: it could mean "calculate X% of a value" or "what percentage is A of B."

**Decision:** `percentage(value, percentage) = value * percentage / 100`. Example: `{value: 200, percentage: 15}` → `30`.

**Why:** this is the standard interpretation of the `%` button on any physical or OS calculator (Windows, iOS). It requires no further justification since it's the convention users expect.

**Trade-off:** doesn't cover the "what % is A of B" case — if needed, that would be a separate operation (`percentage_of`), not a reinterpretation of this one.

---

## 5. Layering: pure logic vs. HTTP handlers

**Context:** the JD asks for "clean, readable, idiomatic code" and "testable architecture."

**Decision:**
```
backend/
  calculator/     → pure functions, no HTTP, no external dependencies
  handlers/       → parse JSON, call calculator/, build the response
  main.go         → wires routes, starts the server
```

**Why:** handlers contain no business logic — they only decode the request, delegate to a pure function in `calculator/`, and encode the response. This allows testing all math logic without spinning up an HTTP server (fast tests, no request/response mocks). If the protocol changes tomorrow (REST → gRPC), `calculator/` stays untouched.

**Trade-off:** more files and an extra layer of indirection for such a simple calculator — could look like over-engineering at first glance, but it's exactly the separation of concerns the JD explicitly asks for.

---

## 6. Error handling: typed domain errors mapped to HTTP, split 400/422, uniform error contract

**Context:** `divide` (by zero), `sqrt` (of negatives), and `power` (non-real results) can fail due to business rules, not protocol errors — distinct from requests that are malformed before any calculation is attempted (missing field, wrong JSON type). Originally logged as a single "map everything to 400" rule; refined once the spec's clarifying Q&A distinguished the two failure categories.

**Decision:**
- **400 (Bad Request) — `INVALID_INPUT`:** malformed request shape, detected in `handlers/` before `calculator/` is ever called (missing operand, non-numeric JSON value, invalid JSON body). Error message identifies the specific field that failed (e.g. "b is required"), not a generic string.
- **422 (Unprocessable Entity) — domain errors:** the request was well-formed, but the calculation itself is mathematically undefined. Three codes, decided via one generic post-computation check rather than special-casing every input combination: `DIVISION_BY_ZERO` (divisor is zero), `NEGATIVE_SQRT` (operand is negative), `NON_FINITE_RESULT` (any other case where the result is `NaN` or `±Inf` — covers power with a negative base and non-integer exponent, `0^negative`, and numeric overflow, all with a single check: `math.IsNaN(result) || math.IsInf(result, 0)`).
- **Uniform error contract**, identical across all seven endpoints and both status codes: `{"error": {"code": "SCREAMING_SNAKE_CASE", "message": "human-readable string"}}`. Success responses: `{"result": number}`.

**Why:** separates "the request itself was malformed" (400, caught before touching business logic) from "the request was valid but the domain rejects it" (422, caught inside business logic) — two different failure origins get two different status codes, instead of collapsing both into a generic 400. The single non-finite-result check (rather than pre-validating every risky input combination per operation) keeps `calculator/` simple: it doesn't need to know in advance every combination that produces an undefined result, only how to recognize one after computing. The stable `code` field lets any client (the frontend or a future one) branch on error type without parsing the human-readable message, which is free to reword without breaking anything downstream.

**Trade-off:** one more layer of indirection (define typed error → map it to a code + status) compared to returning an HTTP error directly from the math function; two status codes to keep straight instead of one blanket 400.

---

## 7. No database or persistence

**Context:** older versions of this exercise from other candidates were found online that included a database for operation history.

**Decision:** no persistence is implemented. The calculator is stateless: each request is independent.

**Why:** the assignment doesn't mention history or persistence anywhere. Adding it without being asked is scope creep, consuming timebox (2-4 hours) on DB setup, schema, and migrations — time taken away from what's explicitly evaluated (React, Go, tests, clean architecture).

**Trade-off:** no record of past operations. If requested later, a repository layer (isolated persistence layer) would be added without touching `calculator/` or `handlers/`, thanks to the layering already in place.

---

## 8. Internal architecture: layered monolith, single service

**Context:** the assignment calls for a "backend microservice." It's worth clarifying what that means structurally versus at the deployment/scope level.

**Decision:** the backend is a single deployable service (one binary, one process) structured internally as a layered monolith (`calculator/` → `handlers/` → `main.go`), not a distributed set of services.

**Why:** "microservice" and "monolith" describe two different axes, not opposites. Deployment/scope axis: is it a single service scoped to one domain, decoupled from the frontend? Yes — that satisfies "microservice" as used in the assignment. Internal structure axis: is the code split across services communicating over the network, or living in one process calling functions directly? Here it's the latter — a monolith internally. The vast majority of real-world microservices are internally layered monoliths scoped to a small domain; the true opposite of "monolith" is a distributed architecture, which would be unjustified complexity for a calculator with no cross-domain orchestration.

**Trade-off:** none specific to this decision — it's a clarification of terminology rather than a design trade-off.

**How it would scale:** if the domain grew to include unrelated capabilities (e.g., user accounts, billing, calculation history across users), splitting into separate services communicating over the network would become justified. For a stateless calculator with a handful of math operations, that split would add operational complexity with no corresponding benefit.

---

## 9. Frontend component structure: split into multiple components

**Context:** the calculator UI could be a single component handling everything, or split into smaller pieces.

**Decision:** split into multiple components (`Display`, `Keypad`, `OperationButtons`, etc.), each with a single rendering responsibility.

**Why:** same principle applied on the backend — separating by responsibility enables testing each piece in isolation with React Testing Library, without needing to simulate the whole calculator working end-to-end. Matches the "logic separated from interface" principle from the project constitution.

**Trade-off:** more files, and some prop drilling from parent to children — manageable here because the component tree is shallow.

---

## 10. HTTP client: native `fetch`, no axios

**Context:** the frontend needs to call 7 backend endpoints.

**Decision:** native `fetch`, no external HTTP client library.

**Why:** aligned with the "minimal stack, new dependency requires justification" principle. There's no real need for axios when calling a handful of endpoints on the same backend.

**Trade-off / gotcha to handle explicitly:** `fetch` does **not** reject its promise on HTTP error statuses (400, 500) — it only rejects on network failures. A 400 response (e.g., division by zero) comes back as a "successful" fetch with `response.ok === false`. Error handling must explicitly check `response.ok` before reading the body, or errors will fail silently.

---

## 11. State management: `useReducer`, not scattered `useState`

**Context:** calculator state includes several interdependent fields (`display`, `previousValue`, `operation`, `error`) that change together based on discrete user actions (digit press, operation press, equals, clear).

**Decision:** `useReducer` with a single reducer function handling all state transitions, dispatched from UI event handlers.

**Why:** a calculator is inherently a state machine — each user action triggers a well-defined transition between states (entering first operand, operation selected, entering second operand, showing result, error state). With scattered `useState`, that transition logic ends up duplicated across multiple button handlers. With `useReducer`, the transition logic lives in one place (the reducer function), and it's pure and fully testable without rendering anything — the same principle applied to `calculator/` on the backend, now applied to frontend state as well. This keeps the "logic separated from interface" principle consistent on both ends of the project.

**Trade-off:** more upfront boilerplate (action types, reducer function with a switch) than a handful of direct `setState` calls — justified here by the interdependent, state-machine nature of the domain, not by the size of the state alone.

---

## 12. Error handling on both backend and frontend, for different reasons

**Context:** with typed domain errors already defined on the backend (decision #6), it's worth being explicit about why the frontend also needs its own error handling, rather than relying solely on the backend.

**Decision:** both layers handle errors, but with distinct responsibilities:
- **Backend:** the only source of truth for validation. Never optional, regardless of what the frontend does.
- **Frontend:** two separate responsibilities — (a) pre-submit validation for UX (e.g., disabling "=" on empty input) that reduces unnecessary calls and gives immediate feedback, and (b) handling the error response returned by the backend (checking `response.ok` given the `fetch` gotcha from decision #10, parsing the error message, and displaying it legibly instead of a raw technical error or a silent failure).

**Why:** anyone can call the API directly (curl, Postman), bypassing the React UI entirely — so backend validation is the real security/stability boundary and can't be skipped. Frontend-only validation without backend enforcement leaves the API exposed to malformed requests. Backend-only validation without proper frontend handling produces a poor user experience (raw errors, or silent failures due to the `fetch` gotcha) and doesn't satisfy the assignment's explicit requirement for frontend "input validation and error handling."

**Trade-off:** some validation logic is effectively duplicated in spirit (not in code) across both layers — e.g., "b can't be zero" is checked both as UX in the frontend and as the real guard in the backend. This duplication is intentional, not accidental: never trust the client.

---

## 13. Test coverage strategy: near-total on pure logic, critical-path focus elsewhere

**Context:** the assignment asks for "unit tests covering key functionality for both layers," while the project follows TDD (test-first, red-green-refactor per task). It's worth clarifying that these are two different axes: TDD governs *how* each test is written; "key functionality" governs *what* deserves a dedicated test given the timebox.

**Decision:**
- `calculator/` (pure backend logic) and the frontend `useReducer` reducer: near-complete coverage. Both are pure, cheap to test (no mocks, no network, no rendering), and are the core of the domain — low cost, high value.
- `handlers/` (HTTP layer): critical paths only — valid request → 200 with correct result, invalid input → 400 with correct error, malformed JSON → 400. Not exhaustive coverage of every possible malformed payload shape.
- Frontend components: key user flows, not exhaustive per-button coverage — e.g., one parametrized test covering "any digit press updates the display" instead of one test per digit; full flow tests for entering numbers → choosing operation → equals → seeing result, and for the error-display path.

**Why:** TDD as a process doesn't dictate scope — with a 2-4 hour timebox, exhaustive testing of every line isn't feasible or valuable. The right allocation is coverage proportional to cost and value: pure logic is cheap and central, so it gets near-total coverage; integration layers (HTTP, UI) get focused coverage of the paths that actually matter to a user or API consumer.

**Trade-off:** some malformed-input edge cases in handlers and some rare UI interaction sequences are not explicitly tested. Acceptable given the timebox; would be revisited if this moved toward production.

---

## 14. CORS and floating-point rounding must be explicit in the spec, not assumed

**Context:** two practical issues that don't show up until integration time: (1) frontend and backend run on different ports in development, so the browser blocks cross-origin requests unless the backend sends the right CORS headers; (2) both Go's `float64` and JS suffer from floating-point precision artifacts (`0.1 + 0.2 = 0.30000000000000004`), which a calculator surfaces immediately.

**Decision:**
- CORS support is an explicit requirement in `spec.md` and an explicit task in `tasks.md`, not left implicit. Given the task-by-task workflow ("implement ONLY this task, then stop"), anything not explicitly planned risks being skipped until it breaks at integration time.
- Rounding happens in the backend, before the JSON response is built — not in the frontend. Values are rounded to 10 decimal places (`math.Round(value*1e10) / 1e10`) to eliminate floating-point noise while preserving meaningful precision for operations like `sqrt` or `power`.

**Why:** CORS failures and floating-point artifacts are the kind of issue that costs the most time precisely because they surface late (at integration, or when a user tries a decimal input) rather than during isolated unit testing. Planning them explicitly avoids discovering them under time pressure. Rounding in the backend keeps the API contract consistent — the frontend doesn't have to guess how many decimals to trust or re-implement its own rounding logic.

**Trade-off:** fixed-precision rounding (10 decimals) is a reasonable default, not a rigorous numerical policy — a domain requiring arbitrary precision (e.g., financial calculations) would need a different approach (decimal types instead of `float64`).

**CORS implementation note (design pattern):** Go's standard library has no built-in middleware concept — it's a pattern built by hand through function composition: a function takes an `http.Handler` and returns a new `http.Handler` wrapping it (`corsMiddleware(mux)`). This is most precisely a **Decorator** (same interface preserved, behavior added transparently around the wrapped handler) — with an element of **Chain of Responsibility** where the chain can short-circuit: the CORS middleware returns early on `OPTIONS` preflight requests without calling `next.ServeHTTP`, meaning that particular layer can choose not to pass the request further down the chain. Unlike Express, where `app.use()` manages the middleware chain implicitly, in plain `net/http` the composition is explicit and manual (`loggingMiddleware(corsMiddleware(mux))`) — more verbose, but fully transparent about what runs and in what order, consistent with the minimal-stack principle from decision #10.

---

## Process note: how design patterns were identified

Design patterns in this project were not decided upfront and then forced into the code — they were identified through a dedicated review step after each major module was implemented, by explicitly asking: *"Review this code and identify which recognized design patterns it applies, even if not consciously named while writing it — name each pattern, where it appears, and why it fits (or note if it's a partial approximation)."*

This is deliberately a separate step from writing the code itself. Asking for a pattern to be *used* from the start risks forcing it onto a problem that didn't need it — the opposite of the judgment this project aims to demonstrate. Solving the problem first with straightforward code, then naming what emerged, keeps pattern usage honest: patterns present because they fit the problem, not because they were requested in advance.

---

## 15. Frontend folder structure: type-based, not feature-based

**Context:** frontend code can be organized by technical role (`components/`, `state/`, `api/`, `types/` — type-based) or by business feature (`features/calculator/` containing its own component, state, and API calls together — feature-based).

**Decision:** type-based structure:
```
frontend/src/
  components/   (Calculator, Display, Keypad, OperationButtons + colocated tests)
  state/        (calculatorReducer + test)
  api/          (calculatorApi — the fetch wrapper/Adapter + test)
  types/        (shared State, Action, and API request/response types)
```

**Why:** feature-based organization pays off when an app has multiple business features to keep independently grouped as it grows. This project has exactly one feature — the calculator — so a `features/calculator/` folder would be the only folder under `features/`, adding a layer of organization that organizes nothing. This mirrors the same reasoning already applied to the backend (decision #8): domain-based modularization only makes sense when there's more than one domain. Backend and frontend end up following the same organizing principle, for the same reason, on both sides of the project.

**Trade-off:** if the app grew to include unrelated features later (e.g., a history view, user settings), a type-based structure would need to be revisited in favor of feature-based folders to keep unrelated concerns from mixing under the same `components/`/`state/`/`api/` folders.

---

## 16. Branching strategy: one branch per implementation phase, merged via PR

**Context:** the project has two natural sequential blocks (backend, then frontend, per the established implementation order). Options considered: everything on `main` with incremental commits, one branch per small feature (e.g. per math operation), or one branch per phase.

**Decision:** two branches, `backend` and `frontend`, each merged into `main` via a pull request once its block is complete and green. No branch per individual task or operation.

**Why:** branches exist to isolate in-progress work from a stable `main` — typically valuable for parallel work across multiple people, or to pause/revert one line of work independently of another. This is a single-developer, linear-timeline project with a 2-4 hour timebox — there's no parallel work to isolate. A branch per operation would be the same kind of over-engineering already rejected for microservices (decision #8): applying a team-coordination tool to a context with no team to coordinate. At the same time, everything on `main` forfeits a chance to demonstrate familiarity with the branch + PR workflow the JD explicitly values (release processes, CI/CD practices). One branch per phase is the point where the tool's overhead (two branches, two PRs) is trivial, but the artifact produced (a PR with a description of what it covers, referencing the relevant RFs) is genuinely useful to a reviewer skimming the repo.

**Trade-off:** as the sole reviewer of my own PRs, the review step is self-review rather than a second pair of eyes — the value here is the artifact and the demonstrated workflow, not actual peer review.

**How it would scale:** with more than one developer or genuinely parallel workstreams, a branch per feature (or per RF) would become justified.

---

## 17. Physical keyboard support, as the "unexpected but necessary" feature

**Context:** the Silver.dev take-home guide notes that exceptional submissions often add "an unexpected feature" that demonstrates product sense — something users would expect but that wasn't explicitly requested — as opposed to adding unrequested business capabilities (which would contradict the scope discipline already established in decisions #1, #7, and #14).

**Decision:** support keyboard input in addition to clicking buttons, wired to the same reducer actions the buttons already dispatch. Finalized in `specs/001-calculator/spec.md` (RF-21) after the spec's QA pass, the scope ended up richer than originally sketched here: digits, all six binary operators (`+ - * / % ^`), decimal point, `Enter` for equals, `Escape` for clear (AC), and `Backspace` for delete-last-digit — the last one also exposed as an on-screen icon button (⌫-style, matching the standard Android calculator convention), not keyboard-only as first drafted.

**Why:** anyone using a calculator expects to be able to type, and the assignment doesn't mention it anywhere, so adding it shows initiative beyond the literal spec. Critically, it does not add a new business capability — it's a second input path into behavior already defined by other RFs (same actions the buttons dispatch), so it stays consistent with the project's discipline against unrequested scope. Making delete-last-digit also available as an on-screen control (not just a keyboard shortcut) was a deliberate refinement during the spec's QA pass, once it became clear a mouse/touch-only user needs an equivalent to Backspace, not just keyboard users.

**Trade-off:** none of substance — this is additive UX on top of existing, already-tested logic, not new logic requiring its own test surface beyond confirming each key/icon dispatches the right action.

**Related refinement:** the MVP scope also settled on a single Clear/AC action (RF-19, full reset) rather than a separate "clear current entry" action — the on-screen delete-last-digit icon covers the more granular correction case instead, avoiding a third, redundant clear-type control.