# Spec 001 — Calculator API and UI

## Context and goal
This is the first feature of sezzle-calculator, establishing the baseline
the rest of the project builds on. Goal: give a user a reliable,
accurate calculator experience in the browser (four basic operations plus
power, square root, and percentage), backed by a stateless REST API that
handles both everyday arithmetic and mathematically undefined inputs
(division by zero, negative square roots, non-real power results)
predictably — returning a clear domain error instead of a wrong number,
a crash, or a silent `NaN`/`Infinity`.

## Users
- **Calculator user** — interacts with the web UI to perform arithmetic
  (the primary user of this feature).
- **API consumer** — calls the REST endpoints directly (e.g. `curl`,
  Postman, or a future client), bypassing the UI entirely. The API must
  be correct and self-explanatory on its own, independent of the frontend.

## User stories
- **US-1**: As a calculator user, I want to perform addition, subtraction,
  multiplication, and division, so I can get quick, accurate results
  without doing the math by hand.
- **US-2**: As a calculator user, I want to compute powers and square
  roots, so I can handle common non-basic calculations without switching
  tools.
- **US-3**: As a calculator user, I want to compute a percentage of a
  value (e.g. a discount or a tip), so I can answer "what is X% of Y"
  quickly.
- **US-4**: As a calculator user, I want a clear message when I attempt
  something mathematically undefined (e.g. divide by zero), so I
  understand why I didn't get a number, instead of seeing a crash or a
  confusing result like `Infinity`.
- **US-5**: As an API consumer, I want each operation exposed through its
  own endpoint with a consistent, predictable request/response contract,
  so I can integrate without going through the UI.

## Functional requirements

### Data contract
Request fields per operation, and the shared success/error response
shapes referenced throughout this section:

| Operation | Request fields |
|---|---|
| Addition (`+`) | `a`, `b` |
| Subtraction (`-`) | `a`, `b` |
| Multiplication (`×`) | `a`, `b` |
| Division (`÷`) | `a`, `b` |
| Power (`^`) | `base`, `exponent` |
| Square root (`√`) | `value` |
| Percentage (`%`) | `value`, `percentage` |

- Success response: `{"result": number}`.
- Error response: `{"error": {"code": string, "message": string}}`
  (RF-12).

### Operations
- **RF-1 (Addition)**: When a client requests addition with two numeric
  operands `a` and `b`, the system shall return the sum, rounded to 10
  decimal places.
- **RF-2 (Subtraction)**: When a client requests subtraction with two
  numeric operands `a` and `b`, the system shall return `a - b`, rounded
  to 10 decimal places.
- **RF-3 (Multiplication)**: When a client requests multiplication with
  two numeric operands `a` and `b`, the system shall return the product,
  rounded to 10 decimal places.
- **RF-4 (Division)**: When a client requests division with two numeric
  operands `a` and `b` where `b` is not zero, the system shall return
  `a / b`, rounded to 10 decimal places. If `b` is zero, then the system
  shall reject the request with a `DIVISION_BY_ZERO` domain error,
  regardless of the value of `a`.
- **RF-5 (Power)**: When a client requests exponentiation with numeric
  operands `base` and `exponent`, the system shall return `base ^
  exponent`, rounded to 10 decimal places, whenever the mathematically
  correct result is a finite real number — including `base = 0` and
  `exponent = 0`, which shall return `1` (a hard requirement, not merely
  a convention). Non-finite results (e.g. a negative `base` with a
  non-integer `exponent`, or `0` raised to a negative `exponent`) are
  rejected as `NON_FINITE_RESULT` under the generic guard defined in
  RF-9; this requirement does not restate that check.
- **RF-6 (Square root)**: When a client requests the square root of a
  numeric operand `value` that is zero or positive, the system shall
  return the result, rounded to 10 decimal places. If `value` is
  negative, then the system shall reject the request with a
  `NEGATIVE_SQRT` domain error.
- **RF-7 (Percentage)**: When a client requests the percentage operation
  with numeric operands `value` and `percentage`, the system shall return
  `value * percentage / 100`, rounded to 10 decimal places. `value` and
  `percentage` may each be positive, negative, or zero.

### Cross-cutting backend behavior
- **RF-8 (Request validation)**: If a request body is not valid JSON, is
  missing a required operand, or contains a value for an operand that is
  not of JSON type `number` — including a numeric-looking string such as
  `"5"`, which is rejected rather than coerced — then the system shall
  reject it with an `INVALID_INPUT` error whose message identifies the
  specific field that is missing or invalid, distinct from a domain
  error, without attempting the calculation.
- **RF-9 (Generic non-finite guard)**: For any of the seven operations
  (RF-1–RF-7), where the computed result is not a finite real number for
  a reason not already covered by a more specific domain error (RF-4's
  `DIVISION_BY_ZERO`, RF-6's `NEGATIVE_SQRT`), the system shall reject
  the request with a `NON_FINITE_RESULT` domain error rather than
  returning `NaN` or `Infinity` in a successful response. This check
  runs immediately after computing the result and before rounding
  (RF-10) is applied — the check always operates on the raw computed
  value, never on a rounded one.
- **RF-10 (Result rounding)**: The system shall round every numeric
  result to 10 decimal places before returning it in a successful
  response, for all seven operations, applied only after the RF-9
  non-finite check has passed.
- **RF-11 (Cross-origin access)**: While the frontend and backend run on
  different origins during development, the system shall accept `POST`
  requests — and the `OPTIONS` preflight requests browsers send ahead of
  them — from origin `http://localhost:5173`, for all seven operations.
- **RF-12 (Error contract)**: The system shall report every error (both
  `INVALID_INPUT` and domain errors) as a machine-readable `code` plus a
  human-readable `message`, in a consistent shape (see Data contract)
  used identically across all seven endpoints. For `INVALID_INPUT`, the
  message identifies the specific offending field (RF-8); for domain
  errors, the message describes the mathematical rule that was violated.

### Frontend UI
- **RF-13 (Operand entry)**: When the user presses a digit, the system
  shall append it to the operand currently being entered and display it;
  this marks that operand as "entered" for the purposes of RF-18.
- **RF-14 (Binary operation flow)**: When the user selects a binary
  operator (`+`, `-`, `×`, `÷`, `%`, `^`) after entering a first operand,
  the system shall capture that value and await a second operand; when
  the user then presses equals, the system shall submit both operands to
  the corresponding operation and display the result. For `%`, the first
  operand entered is `value` and the second is `percentage`; for `^`,
  the first operand entered is `base` and the second is `exponent` — `^`
  uses the exact same entry flow as `%`. If the user selects a binary
  operator before entering any digit of a first operand, then the system
  shall use the default displayed value (`0`) as the first operand. If
  the user presses equals while no binary operator is pending, then the
  system shall take no action. If the user selects a new binary operator
  while a previously selected operator is still awaiting its second
  operand, and at least one digit of that second operand has already
  been entered, then the system shall immediately compute the pending
  operation using the two operands already entered (the same
  compute-and-display behavior as pressing equals), display the result,
  and use it as the first operand of the newly selected operator. If no
  digit of the second operand has been entered yet, then the system
  shall instead simply replace the pending operator with the newly
  selected one, without triggering any computation. This chaining
  behavior reuses the existing compute-on-equals logic, not expression
  parsing; operator precedence and parentheses remain out of scope.
- **RF-15 (Unary shortcuts)**: When the user presses `√x`, the system
  shall submit the currently displayed value (which defaults to `0` if
  no operand has yet been entered) as the square-root operand (RF-6) and
  display the result. When the user presses `x²`, the system shall
  submit the currently displayed value (same default-to-`0` behavior) as
  `base` with `exponent` fixed at `2` to the power operation (RF-5) and
  display the result. `x²` remains available as a shortcut in addition
  to the general `^` binary operator (RF-14), which supports any
  exponent.
- **RF-16 (Result display)**: When a request succeeds, the system shall
  display the returned result as the new current value.
- **RF-17 (Error display)**: If a request fails with a validation or
  domain error response (RF-12), then the system shall display that
  error's `message` field to the user in place of a result, without
  crashing or showing a raw/technical error. The error's `code` field
  remains present in the response contract for API consumers and future
  use; the frontend does not branch its display logic on `code` in this
  MVP. Network-level failures are distinct from this case — see RF-25.
- **RF-18 (Frontend pre-submit validation)**: While the operand required
  for the pending operation has not had at least one digit explicitly
  entered by the user — the initial or cleared display showing `0` does
  not by itself count as "entered" — the system shall disable the action
  that would submit that operation.
- **RF-19 (Clear)**: When the user presses the on-screen Clear button, or
  the Escape key (RF-21), the system shall reset the calculator to its
  initial state — no operand entered, no pending operator, no displayed
  result or error.
- **RF-20 (Delete last digit)**: When the user presses the Backspace key,
  or taps the on-screen delete-last-digit icon button (⌫-style, matching
  the standard Android calculator convention — an icon, not a text
  label), the system shall remove the last digit from the operand
  currently being entered; deleting the only remaining digit shall reset
  that operand's display to its empty/zero state, not leave it blank in
  an invalid way.
- **RF-21 (Keyboard input)**: When the user presses a keyboard key mapped
  to a digit, a binary operator (`+ - * / ^`), `%`, Enter, Escape, or
  Backspace, the system shall perform the same action as the
  corresponding control — digit entry (RF-13), operator
  selection/chaining (RF-14), submitting the pending operation (RF-14),
  clear (RF-19), or delete-last-digit (RF-20) respectively — producing
  behavior identical to the equivalent mouse/touch input.
- **RF-22 (Decimal point)**: When the user presses the decimal-point
  control, the system shall append a decimal point to the operand
  currently being entered, unless that operand already contains one, in
  which case the press is a no-op.
- **RF-23 (Sign toggle)**: When the user presses the sign-toggle (`±`)
  control, the system shall multiply the currently displayed value by
  `-1` and update the display accordingly, entirely on the frontend,
  without making any backend request.
- **RF-24 (Error recovery on new input)**: While an error message is
  displayed (RF-17), when the user presses any digit, operator, or
  unary shortcut, the system shall clear the displayed error and begin a
  fresh entry with that input, without requiring an explicit Clear
  action (RF-19) first.
- **RF-25 (Network failure handling)**: If a request cannot be completed
  due to a network-level failure — no HTTP response is received (backend
  unreachable, a CORS block, a timeout), or the response must be
  inspected for success/failure before it can be treated as a normal
  result or error body — then the system shall display a distinct
  connection-error message to the user, rather than attempting to parse
  it as a validation/domain error body (RF-17).
- **RF-26 (In-flight submission guard)**: While a request triggered by
  equals or a unary shortcut is in flight, the system shall disable
  equals and all unary-shortcut controls, using the same disabling
  mechanism as RF-18; upon receiving a response — success, error, or
  network failure (RF-25) — the system shall re-enable them.
- **RF-27 (Operator after a result)**: When a result is currently
  displayed (RF-16) and the user selects a binary operator, the system
  shall use that displayed result as the first operand of the new
  operation — the same flow as RF-14 for a freshly entered first operand
  — and await a second operand. This is distinct from RF-14's chaining
  clause, which applies when a new operator is selected before equals
  has produced a result; here, equals has already fired and produced the
  displayed value (e.g. `2 + 2 =` displays `4`; pressing `+ 3 =` next
  then uses `4` as the first operand, yielding `7`).

## Non-functional requirements
- **RNF-1 (Trust boundary)**: Every validation rule enforced by the
  frontend (RF-18) is independently re-enforced by the backend (RF-8);
  the backend never assumes the frontend ran first, since any client can
  call the API directly.
- **RNF-2 (Statelessness)**: No operation result, history, or partial
  calculation state is persisted anywhere in the backend between
  requests.
- **RNF-3 (Consistent contract)**: A successful response always has the
  same shape (`{"result": number}`, see Data contract); an error
  response always has the same shape (`{"error": {"code", "message"}}`,
  RF-12), regardless of which of the seven operations was called.
- **RNF-4 (No silent data loss)**: The system never substitutes a
  mathematically wrong or misleading value (e.g. `0` in place of an
  error) for an operation it cannot complete — it always signals failure
  explicitly (RF-8, RF-9).

## Edge cases
- `divide(0, 0)` → `DIVISION_BY_ZERO`, not `NaN` — the check is on the
  divisor, independent of the dividend.
- `sqrt(-0)` is not negative → succeeds (result `0`), not a domain error.
- `power(0, 0)` → succeeds with result `1` — a hard requirement (RF-5),
  not a convention; the RF-9 finite-result guard treats it as valid.
- `power(0, negative exponent)` → `NON_FINITE_RESULT` (result would be
  infinite).
- `percentage(-50, 10)` and `percentage(50, -10)` → both succeed, both
  returning `-5` — negative `value` and negative `percentage` are valid.
- Operands of very large or mismatched magnitude may lose precision
  silently at the float64 level (e.g. `1e20 + 1 = 1e20`) before rounding
  is even applied; this is an accepted limitation of using `float64`, not
  something the system detects or rejects — most extreme cases already
  surface as `NON_FINITE_RESULT` via overflow instead.
- A value that is mathematically intended to be zero but arrives as a
  small negative float due to floating-point noise from a prior
  computation (e.g. `-1e-15`) is treated as genuinely negative:
  `sqrt(-1e-15)` triggers `NEGATIVE_SQRT` like any other negative input.
  No epsilon tolerance is applied — this is accepted behavior, not a
  gap, consistent with the float64 magnitude-loss limitation above.
- A request body containing a syntactically valid but semantically wrong
  operand (a string, `null`, or a boolean where a number is expected) is
  an `INVALID_INPUT` error, not a domain error.

## Out of scope
- Expression evaluation or chaining multiple operations in one request or
  UI input (e.g. `3 + 4 * 2`); each request is exactly one operation.
  (UI operator-chaining, RF-14, computes eagerly on each new operator
  press — it is not expression parsing with precedence or parentheses.)
- Calculation history or memory (`M+`/`M-`) of any kind.
- Authentication, authorization, or any per-user state.
- Arbitrary-precision or decimal-type math (e.g. for financial-grade
  precision beyond `float64` + 10-decimal rounding).
- DEG/RAD modes, trigonometric functions, or any operation beyond the
  seven listed.
- A precision/rounding selector exposed to the user — 10 decimal places
  is fixed for this feature.
- Display-length or overflow handling for very long typed numbers or
  results — the 10-decimal rounding rule (RF-10) already bounds result
  length to a reasonable range.

## Completion criteria
- All seven operations are reachable — power via both the general `^`
  binary operator (RF-14) and the `x²` shortcut (RF-15) — each producing
  the correct, 10-decimal-rounded result for valid input (RF-1–RF-7,
  RF-10).
- `INVALID_INPUT`, `DIVISION_BY_ZERO`, `NEGATIVE_SQRT`, and
  `NON_FINITE_RESULT` are all reproducible and returned in the contract
  shape defined by RF-12 / Data contract.
- The UI can complete all seven operations end-to-end against the
  running backend, correctly displaying both results and errors
  (RF-13–RF-17), including operator chaining before equals (RF-14),
  starting a new operation from a displayed result after equals (RF-27),
  decimal-point and sign-toggle entry (RF-22, RF-23), error recovery on
  new input (RF-24), distinct network-failure messaging (RF-25),
  in-flight submission guarding (RF-26), and clear/delete/keyboard
  controls (RF-19–RF-21) — delete-last-digit is reachable both via
  Backspace and an on-screen icon button (RF-20).
- Cross-origin `POST`/`OPTIONS` requests from `http://localhost:5173` to
  the backend succeed with no browser CORS failures (RF-11), verified
  manually.
- `go test ./...` and `npm run test` both pass, per `AGENTS.md`.

## Open questions
None. This spec went through a dedicated QA pass (ambiguities,
contradictions between backend and frontend requirements, uncovered edge
cases, and conflicts with `docs/constitution.md`); every finding was
resolved with the user and folded into the requirements above,
including one product decision (`^` added as a fifth binary operator so
the UI can reach general power, not only `x²`).
