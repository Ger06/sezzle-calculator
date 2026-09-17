# Tasks 001 — Calculator API and UI

> **If anything below is ambiguous or blocks progress while you're
> executing it — a missing decision, an unclear requirement, a spec/plan
> mismatch — stop and ask. Do not guess or work around it silently.**

Each task follows the constitution's TDD principle (#5): write the test,
confirm it fails, implement the minimum to pass, refactor, then mark the
task done here and stop — don't start the next task unprompted
(`AGENTS.md`). Backend tasks (Phase 1) must all be green before any
Frontend task (Phase 2) starts.

## Phase 1 — Backend

- [x] **T1 — Domain errors, non-finite guard, rounding helper.** In
  `calculator/`, write failing tests for `ErrDivisionByZero`,
  `ErrNegativeSqrt`, `ErrNonFiniteResult` (`errors.go`), `checkFinite`
  (`finite.go` — errors iff `NaN`/`±Inf`), and `Round` (`round.go` —
  10-decimal rounding); implement all three. **RF:** RF-4, RF-6, RF-9,
  RF-10. **Done when:** `go test ./calculator/...` passes, including a
  `NaN`/`+Inf`/`-Inf` case for `checkFinite` and a >10-decimal case for
  `Round`.

- [x] **T2 — Add, Subtract.** Write failing table-driven tests (including
  a negative-operand case each), then implement `Add`/`Subtract` in
  `calculator.go`, each calling `checkFinite` then `Round`. **RF:** RF-1,
  RF-2. **Done when:** `go test ./calculator/...` green for both.

- [x] **T3 — Multiply, Divide.** Write failing tests including the
  `divide(0, 0)` edge case (must be `ErrDivisionByZero`, not `NaN`);
  implement, checking `b == 0` before computing. **RF:** RF-3, RF-4.
  **Done when:** all table cases pass, including `divide(0,0)`.

- [x] **T4 — Sqrt.** Write failing tests including `sqrt(-0)` (must
  succeed, result `0`) and `sqrt(-1e-15)` (must error — accepted
  floating-point-boundary behavior per spec.md); implement, checking
  `value < 0` before computing. **RF:** RF-6. **Done when:** both edge
  cases behave exactly as specified.

- [x] **T5 — Power.** Write failing tests including `power(0, 0) == 1`
  (hard requirement), `power(-8, 0.5)`, and `power(0, -1)` (both of the
  latter two → `ErrNonFiniteResult` via `checkFinite`, no special-casing
  per input shape); implement. **RF:** RF-5, RF-9. **Done when:** all
  three cases pass.

- [x] **T6 — Percentage.** Write failing tests including a negative
  `value` and a negative `percentage`; implement. **RF:** RF-7. **Done
  when:** `percentage(-50, 10)` and `percentage(50, -10)` both return
  `-5`.

- [x] **T7 — Shared handler helpers.** In `handlers/handlers.go`, write
  failing tests for `decodeRequest[T]`, `requireField(name string, p
  *float64) error` (nil → `"<name> is required"`), `writeResult`, and
  `writeError` (exact `{"result"}` / `{"error":{"code","message"}}`
  shapes per plan.md's data contract); implement, using `*float64`
  request fields per plan.md decision #2. **RF:** RF-8, RF-12. **Done
  when:** a wrong-JSON-type case and a missing-field case each produce
  the correct, distinct message.

- [x] **T8 — Add/Subtract/Multiply handlers.** Write failing `httptest`
  tests (200 + correct body; a missing operand → 400 `INVALID_INPUT`),
  then implement the three handlers + request structs and wire them in
  `main.go`. **RF:** RF-1, RF-2, RF-3, RF-8. **Done when:** all three
  endpoints return correct 200 bodies and a 400 for a missing operand.

- [x] **T9 — Divide/Sqrt handlers.** Write failing tests including each
  operation's 422 domain-error path; implement and wire in `main.go`.
  **RF:** RF-4, RF-6, RF-8, RF-9, RF-12. **Done when:** `POST /divide
  {"a":1,"b":0}` → 422 `DIVISION_BY_ZERO`; `POST /sqrt {"value":-4}` →
  422 `NEGATIVE_SQRT`.

- [x] **T10 — Power/Percentage handlers.** Write failing tests including
  power's `NON_FINITE_RESULT` 422 path; implement and wire in `main.go`.
  **RF:** RF-5, RF-7, RF-8, RF-9, RF-12. **Done when:** `POST /power
  {"base":-8,"exponent":0.5}` → 422 `NON_FINITE_RESULT`; `POST
  /percentage {"value":200,"percentage":15}` → 200 `{"result":30}`.

- [x] **T11 — CORS middleware + full route wiring.** Write a failing
  `httptest` test asserting CORS headers on a request from
  `http://localhost:5173` and `204` on `OPTIONS`; implement
  `handlers/cors.go` as a decorator (plan.md) and confirm all seven
  routes are wired in `main.go` through it. **RF:** RF-11. **Done when:**
  the CORS test passes and `go run main.go` starts without error.

- [x] **T12 — Backend integration pass + green gate.** Write
  `handlers/handlers_integration_test.go` exercising all seven endpoints
  end-to-end (one success + every documented error code), per
  `AGENTS.md`'s naming convention. **RF:** RF-1–RF-12. **Done when:**
  `go test ./... -cover` passes with zero failures. **This is the phase
  gate — do not start any Phase 2 task until this is green.**

## Phase 2 — Frontend

- [x] **T13 — Test tooling + template cleanup.** Add Vitest + React
  Testing Library as dev dependencies (already mandated by `AGENTS.md`,
  not a new undecided dependency) with minimal config; remove the
  default Vite template content (`App.tsx`, `App.css`,
  `src/assets/react.svg` etc.). **RF:** none (setup). **Done when:**
  `npm run test` runs (even with zero tests yet) and the default Vite
  counter demo is gone.

- [x] **T14 — Reducer: entry & editing actions.** Write failing tests,
  then implement digit entry, decimal point (no-op if the operand
  already has one), sign toggle, backspace (including down to
  empty/zero), and clear, in `state/calculatorReducer.ts`. **RF:** RF-13,
  RF-19, RF-20, RF-22, RF-23. **Done when:** all five actions pass their
  tests in isolation (no rendering).

- [x] **T15 — Reducer: operator selection & chaining.** Write failing
  tests, then implement operator selection (defaulting to `0` as the
  first operand per RF-14), the chaining-with-a-second-operand vs.
  replace-the-pending-operator-without-one branch, and operator-after-a-
  displayed-result. **RF:** RF-14, RF-18, RF-27. **Done when:** both the
  chain/replace branch and the RF-27 case have a passing test.

- [x] **T16 — Reducer: submission lifecycle & error states.** Write
  failing tests, then implement the in-flight guard, result display,
  error display, error-recovery-on-new-input, and network-failure state.
  **RF:** RF-16, RF-17, RF-24, RF-25, RF-26. **Done when:** a simulated
  "start → success" and a "start → error" sequence each land in the
  correct terminal state.

- [x] **T17 — calculatorApi.** Write failing tests (mocking `fetch`,
  including the `response.ok` gotcha per `docs/decisions.md` #10) for
  one function per operation; implement, returning the `{ok:true,result}
  | {ok:false,kind:'domain',...} | {ok:false,kind:'network'}` union.
  **RF:** RF-8 (field names), RF-12, RF-25. **Done when:** all seven
  functions send the correct request field names (spec.md's Data
  contract) and correctly discriminate all three outcomes.

- [x] **T18 — Display + Keypad.** Write failing RTL tests, then
  implement `Display` (value / error / connection-error rendering) and
  `Keypad` (digits, `.`, `±`, `⌫`, Clear — wired to the reducer). Add basic responsive support to Display, Touch targets (all Keypad and OperationButtons buttons) must be at least 44×44px, regardless of viewport size. Do not rely on :hover as the only feedback for button state — use :active (which fires on touch too) for pressed-state styling, so mobile users get visual feedback without a hover capability. **RF:**
  RF-13, RF-16, RF-17, RF-19, RF-20, RF-22, RF-23, RF-25. **Done when:**
  a digit press updates the rendered display in a test.

- [x] **T19 — OperationButtons + disabled-state wiring.** Write failing
  tests, then implement the operator/unary/equals buttons and their
  `disabled` wiring to RF-18/RF-26. **RF:** RF-14, RF-15, RF-18, RF-26.
  **Done when:** equals is disabled with no second operand entered and
  re-enabled after a response, per test.

- [x] **T20 — Calculator container: wiring + keyboard.** Write failing
  integration tests, then implement `Calculator.tsx`: owns `useReducer`,
  is the only place calling `calculatorApi` (dispatching start/success/
  error around it), and the inline `keydown` listener. **RF:** RF-21,
  plus integrates RF-13–RF-27. **Done when:**
  `Calculator.integration.test.tsx` covers one full binary-operation flow
  (including `^` reaching general power), one unary shortcut, one error
  path, and one keyboard-equivalence case — all green.

- [x] **T21 — Manual CORS verification + full green check.** Run the
  backend (`go run main.go`) and frontend (`npm run dev`) together;
  perform one calculation for each of the seven operations by hand in
  the browser. **RF:** RF-11 (manual — cross-process behavior isn't
  expressible as an automated test within this repo). **Done when:**
  `npm run test` is green and all seven operations succeed in the
  browser with no CORS errors in the console.
