# Plan 001 — Calculator API and UI

Technical plan implementing `spec.md`'s RF-1–RF-27, governed by
`docs/constitution.md` and `docs/decisions.md`. This is the HOW: module
structure, API contract, error mapping, remaining technical decisions,
and test strategy. No data model section — the calculator is stateless
(constitution #6, RNF-2).

## Module structure — Backend

Folder layout fixed by `docs/decisions.md` #5; file-level layout below.

```
backend/
├── go.mod
├── main.go                        — wires all seven routes on Go 1.22's
│                                     method-aware net/http mux, applies
│                                     the CORS middleware, starts the
│                                     server. [RF-1–RF-7, RF-11]
├── calculator/
│   ├── calculator.go               — Add, Subtract, Multiply, Divide,
│   │                                  Power, Sqrt, Percentage: pure
│   │                                  functions, each calling
│   │                                  checkFinite then Round before
│   │                                  returning. [RF-1–RF-7, RF-10]
│   ├── calculator_test.go          — table-driven tests per operation,
│   │                                  including spec.md's edge cases.
│   │                                  [RF-1–RF-7, edge cases]
│   ├── errors.go                   — typed domain errors:
│   │                                  ErrDivisionByZero, ErrNegativeSqrt,
│   │                                  ErrNonFiniteResult (constitution
│   │                                  #10 — no net/http import here).
│   │                                  [RF-4, RF-6, RF-9]
│   ├── errors_test.go
│   ├── finite.go                   — checkFinite(result float64) error:
│   │                                  the single generic non-finite
│   │                                  guard from decisions.md #6, run
│   │                                  before rounding. [RF-9]
│   ├── finite_test.go
│   ├── round.go                    — Round(value float64) float64,
│   │                                  10-decimal rounding (decisions.md
│   │                                  #14). [RF-10]
│   └── round_test.go
└── handlers/
    ├── handlers.go                  — shared helpers: a generic
    │                                   decodeRequest[T], a
    │                                   requireField(name, *float64)
    │                                   nil-check helper, writeResult,
    │                                   writeError, and the error→status
    │                                   mapping table below. [RF-8, RF-12]
    ├── add.go                       — AddRequest{A, B *float64} + handler.
    │                                   [RF-1, RF-8]
    ├── subtract.go                  — SubtractRequest{A, B *float64}.
    │                                   [RF-2, RF-8]
    ├── multiply.go                  — MultiplyRequest{A, B *float64}.
    │                                   [RF-3, RF-8]
    ├── divide.go                    — DivideRequest{A, B *float64}.
    │                                   [RF-4, RF-8, RF-9 via calculator/]
    ├── power.go                     — PowerRequest{Base, Exponent
    │                                   *float64}. [RF-5, RF-8, RF-9 via
    │                                   calculator/]
    ├── sqrt.go                      — SqrtRequest{Value *float64}.
    │                                   [RF-6, RF-8]
    ├── percentage.go                — PercentageRequest{Value,
    │                                   Percentage *float64}.
    │                                   [RF-7, RF-8]
    ├── cors.go                      — CORS middleware, Decorator +
    │                                   short-circuiting Chain of
    │                                   Responsibility per decisions.md
    │                                   #14. [RF-11]
    ├── handlers_test.go             — unit-level: decoding, error
    │                                   mapping.
    └── handlers_integration_test.go — httptest full round-trip, named
                                        per AGENTS.md convention.
                                        [RF-1–RF-12]
```

## Module structure — Frontend

Folder layout fixed by `docs/decisions.md` #15 (type-based: `components/`,
`state/`, `api/`, `types/` — no fifth top-level folder, see "Justified
technical decisions" below for why keyboard wiring doesn't get one).

```
frontend/src/
├── types/
│   └── calculator.ts        — State, Action, OperationName, per-
│                               operation request types, ErrorCode union,
│                               ApiResult discriminated union. [supports
│                               all RFs]
├── api/
│   ├── calculatorApi.ts     — one function per operation; checks
│   │                          response.ok (decisions.md #10's fetch
│   │                          gotcha) and returns a discriminated union:
│   │                          {ok:true, result} |
│   │                          {ok:false, kind:'domain', code, message} |
│   │                          {ok:false, kind:'network'}.
│   │                          [RF-8, RF-9 (indirectly), RF-12, RF-25]
│   └── calculatorApi.test.ts
├── state/
│   ├── calculatorReducer.ts  — pure reducer, zero DOM/React imports
│   │                           (constitution #4): digit entry (RF-13),
│   │                           operator select/chain (RF-14), unary
│   │                           shortcuts (RF-15), result (RF-16), error
│   │                           (RF-17), submit-guard (RF-18), clear
│   │                           (RF-19), backspace (RF-20), decimal
│   │                           (RF-22), sign toggle (RF-23), error
│   │                           recovery (RF-24), in-flight guard (RF-26),
│   │                           operator-after-result (RF-27).
│   └── calculatorReducer.test.ts
└── components/
    ├── Calculator.tsx             — container: owns useReducer, is the
    │                                only place that calls calculatorApi
    │                                (dispatches a synchronous "submit
    │                                start" action, awaits the call,
    │                                dispatches success/domain-error/
    │                                network-error — RF-25, RF-26), and
    │                                wires the keydown listener inline
    │                                (RF-21, decisions.md #17).
    ├── Calculator.test.tsx
    ├── Calculator.integration.test.tsx — full-flow tests, AGENTS.md
    │                                      naming convention.
    ├── Display.tsx                 — shows the current value, an error
    │                                  message, or a connection-error
    │                                  message. [RF-16, RF-17, RF-25]
    ├── Display.test.tsx
    ├── Keypad.tsx                  — digits (RF-13), decimal point
    │                                 (RF-22), sign toggle (RF-23),
    │                                 delete-last-digit icon button
    │                                 (RF-20), Clear button (RF-19).
    ├── Keypad.test.tsx
    ├── OperationButtons.tsx        — `+ - × ÷ % ^`, `√x`, `x²`, `=`;
    │                                 disabled-state wiring. [RF-14,
    │                                 RF-15, RF-18, RF-26]
    └── OperationButtons.test.tsx
```

## REST API contract

All seven endpoints share the Data contract from `spec.md` (success
`{"result": number}`; error `{"error": {"code": string, "message":
string}}`). `OPTIONS` preflight requests (RF-11) return `204 No Content`
with CORS headers, no body.

### `POST /add` — RF-1
Request: `{"a": 2, "b": 3}` → 200: `{"result": 5}`
Error (400, RF-8): `{"error": {"code": "INVALID_INPUT", "message": "b is required"}}`

### `POST /subtract` — RF-2
Request: `{"a": 5, "b": 3}` → 200: `{"result": 2}`

### `POST /multiply` — RF-3
Request: `{"a": 4, "b": 2.5}` → 200: `{"result": 10}`

### `POST /divide` — RF-4
Request: `{"a": 10, "b": 4}` → 200: `{"result": 2.5}`
Domain error (422): `{"a": 10, "b": 0}` →
`{"error": {"code": "DIVISION_BY_ZERO", "message": "cannot divide by zero"}}`

### `POST /power` — RF-5
Request: `{"base": 2, "exponent": 10}` → 200: `{"result": 1024}`
Domain error (422): `{"base": -8, "exponent": 0.5}` →
`{"error": {"code": "NON_FINITE_RESULT", "message": "result is not a finite real number"}}`

### `POST /sqrt` — RF-6
Request: `{"value": 16}` → 200: `{"result": 4}`
Domain error (422): `{"value": -4}` →
`{"error": {"code": "NEGATIVE_SQRT", "message": "cannot take the square root of a negative number"}}`

### `POST /percentage` — RF-7
Request: `{"value": 200, "percentage": 15}` → 200: `{"result": 30}`

Malformed-request example shared by all seven endpoints (RF-8):
non-JSON body or a numeric-looking string → 400:
`{"error": {"code": "INVALID_INPUT", "message": "a must be a JSON number, got string"}}`

## Domain error → HTTP status mapping

| Code | HTTP status | Triggered by | RF |
|---|---|---|---|
| `INVALID_INPUT` | 400 | Malformed JSON, missing field, or a field not of JSON type `number` | RF-8 |
| `DIVISION_BY_ZERO` | 422 | `/divide` with `b == 0` | RF-4 |
| `NEGATIVE_SQRT` | 422 | `/sqrt` with `value < 0` | RF-6 |
| `NON_FINITE_RESULT` | 422 | Any operation whose raw computed result is `NaN`/`±Inf`, not already covered above (power's non-real cases, `0^negative`, overflow) | RF-9 |

## Justified technical decisions
(Only decisions `docs/decisions.md` doesn't already make. Folder-level
layout, the 400/422 split, `fetch`/`useReducer`/CORS/rounding-in-backend
choices, and test-coverage *shape* are all already decided there and
simply applied above, not repeated.)

1. **`calculator/` is one file (`calculator.go`), not one file per
   operation.** Unlike `handlers/` (decisions.md #2 — each operation gets
   its own request struct, justifying separate files), each operation
   function here is a one-to-three-line pure computation sharing the same
   imports and helpers. Discarded: per-operation files mirroring
   `handlers/` — adds navigation overhead across seven near-trivial
   functions with no corresponding benefit.
2. **A generic `decodeRequest[T any](r *http.Request) (T, error)` helper
   in `handlers.go`, with every required operand declared as a pointer
   (`*float64`) in its per-operation request struct** (e.g.
   `AddRequest{A, B *float64}`), not a plain `float64`. Go's
   `encoding/json` leaves a missing numeric field at its zero value
   (`0`) when decoding into a plain `float64` — indistinguishable from a
   client explicitly sending `0` — so a plain-`float64` struct cannot
   produce RF-8's field-specific "`b` is required" message; it can only
   catch a *wrong-typed* value (e.g. a string), which `encoding/json`
   already reports with the field name included, and `decodeRequest`
   surfaces as-is. With `*float64` fields, decoding leaves a genuinely
   missing field as `nil`, distinguishable from an explicit `0` (which
   decodes to a non-nil pointer to `0`). Each handler then calls a small
   `requireField(name string, p *float64) error` helper per required
   field — checked after `decodeRequest` succeeds — to produce the
   missing-field message; the wrong-type case never reaches this check,
   since `encoding/json` already fails decoding before any field is
   inspected. Discarded: plain `float64` fields with manual
   `json.NewDecoder(r.Body).Decode(&req)` repeated per handler — both
   duplicates decode boilerplate seven times *and* cannot distinguish
   "missing" from "explicitly zero," silently producing wrong
   `INVALID_INPUT` messages (or none at all) for a legitimate missing
   field.
3. **`checkFinite` is a single shared function in `calculator/finite.go`**,
   called at the end of every operation, implementing decisions.md #6's
   "one generic post-computation check" concretely. Discarded: inlining
   `math.IsNaN(result) || math.IsInf(result, 0)` separately in each of the
   seven functions — duplicates the exact same two-line check seven times
   and risks drift if one copy is edited without the others.
4. **Rounding (RF-10) happens inside `calculator/`**, as the last step of
   each operation function (after `checkFinite` passes), not in
   `handlers/`. `docs/decisions.md` #14 says rounding happens "in the
   backend" but doesn't pin the layer. Rounding is a business rule about
   acceptable precision, not an HTTP/serialization concern, so it belongs
   with the pure functions per constitution #3 — this also means any
   future non-HTTP caller of `calculator/` gets correctly-rounded results
   for free. Discarded: rounding in `handlers/` after calling
   `calculator/` — would leave `calculator/`'s own return values
   un-rounded and inconsistent with what callers actually receive.
5. **Async orchestration lives in `Calculator.tsx`'s event handlers, not
   the reducer.** A handler dispatches a synchronous "submit start"
   action, `await`s `calculatorApi`, then dispatches a synchronous
   success/domain-error/network-error action with the outcome — the
   reducer itself stays 100% synchronous and pure (constitution #4).
   Discarded: a `useEffect` keyed on a "pending" state field that
   triggers the fetch as a side effect — adds an extra render cycle and a
   layer of indirection for what is always exactly one request per user
   action here.
6. **The keydown listener (RF-21) is inlined in `Calculator.tsx` via
   `useEffect`**, not extracted to a custom hook. Discarded: a new
   `hooks/useKeyboardShortcuts.ts` — would add a fifth top-level folder
   beyond the four decisions.md #15 already fixed, for what
   decisions.md #17 already characterizes as a handful of lines mapping
   keys to the same `dispatch` calls the buttons use.
7. **`OPTIONS` preflight responses are `204 No Content`** with the CORS
   headers set, no body. Discarded: `200 OK` with an empty body — `204`
   is the more precise status for "request understood, nothing to
   return," and pairs naturally with the CORS middleware's early-return
   short-circuit branch (decisions.md #14).

## Test strategy

Coverage *shape* (near-total on pure logic, critical-path elsewhere) is
fixed by `docs/decisions.md` #13; TDD process (test → confirm red →
implement → refactor, one task at a time) is fixed by the constitution
(#5) and `AGENTS.md`. Applied concretely here:

### Backend
- **`calculator/*_test.go`** — near-complete, table-driven: valid inputs
  for all seven operations (RF-1–RF-7); every domain-error trigger
  (`b == 0` for RF-4, `value < 0` for RF-6, and RF-9's non-finite cases:
  `power(-8, 0.5)`, `power(0, -1)`, plus `power(0, 0) == 1` as a hard
  pass-through case per RF-5); 10-decimal rounding (RF-10); and spec.md's
  edge cases verbatim (`divide(0,0)`, `sqrt(-0)`, `1e20 + 1`,
  `sqrt(-1e-15)` — asserted as accepted `NEGATIVE_SQRT` behavior, not a
  bug).
- **`handlers/handlers_test.go` + `handlers_integration_test.go`** —
  critical paths per decisions.md #13: valid request → 200 with correct
  body shape; malformed JSON / missing field / wrong JSON type (e.g.
  `"a": "5"`) → 400 `INVALID_INPUT` with a field-specific message
  (RF-8); each domain error → 422 with the correct code (RF-4, RF-6,
  RF-9); CORS headers present on a cross-origin request and on an
  `OPTIONS` preflight (RF-11); full error-contract shape assertion
  (RF-12).

### Frontend
- **`state/calculatorReducer.test.ts`** — near-complete: every action
  from RF-13–RF-27 in isolation — digit append and marking an operand
  "entered" (RF-13, RF-18), decimal point no-op when already present
  (RF-22), sign toggle (RF-23), backspace down to empty/zero (RF-20),
  clear (RF-19), operator selection defaulting to `0` as first operand
  (RF-14), chaining-with-a-second-operand vs. replacing-the-operator-
  without-one (RF-14), operator-after-a-displayed-result (RF-27),
  error-recovery on new input while an error is shown (RF-24), and
  in-flight guard toggling (RF-26).
- **`api/calculatorApi.test.ts`** — one test per operation asserting the
  correct request field names (Data contract) and correct discrimination
  of success / RF-12 error body / RF-25 network failure, mocking
  `fetch`'s `response.ok` gotcha (decisions.md #10).
- **`components/*.test.tsx` + `Calculator.integration.test.tsx`** — key
  flows, not exhaustive per-button coverage (decisions.md #13): one
  parametrized digit-press test; one full flow per operation category
  (a binary op via `+`, `^` reaching general power, and `%`; a unary
  shortcut via `√x`/`x²`; an error path RF-17; a network-failure path
  RF-25); disabled-state assertions for RF-18/RF-26; a keyboard-
  equivalence smoke test for RF-21.

All of RF-1–RF-27 and RNF-1–RNF-4 are covered above: RNF-1 (trust
boundary) by RF-8 backend tests existing independently of any frontend
test; RNF-2 (statelessness) by `calculator/` and `handlers/` never
importing a storage package; RNF-3 (consistent contract) by the shared
`handlers.go` helpers and the Data-contract-driven API contract above;
RNF-4 (no silent data loss) by the domain-error mapping table.
