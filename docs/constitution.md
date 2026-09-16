# Constitution — sezzle-calculator

1. **Minimal stack.** No dependency is added to `go.mod` or `package.json`
   without a spec section justifying why the standard library / already
   approved packages can't do it.
2. **Spec governs.** No endpoint, operation, or UI behavior is implemented
   unless it is described in an active spec under `specs/`.
3. **Backend separation.** `calculator/` contains pure functions with zero
   `net/http` imports; `handlers/` only parses requests, calls
   `calculator/`, and writes responses.
4. **Frontend separation.** All state transitions live in the reducer,
   with zero DOM/React imports; components only render state and dispatch
   actions.
5. **Strict TDD.** Every business-logic change starts with a test that is
   run and confirmed failing before any implementation line is written.
6. **Stateless by design.** The backend stores nothing between requests —
   no database, file, or in-memory cache holds calculation history.
7. **Trust boundary.** Every validation rule enforced in the frontend has
   an independent, equivalent check in the backend; the backend never
   assumes the frontend ran first.
8. **Cross-cutting concerns are spec'd upfront.** CORS allowed
   origins/methods and floating-point rounding rules (precision, rounding
   mode) are defined in the spec before implementation begins.
9. **English only.** All code, identifiers, comments, commit messages, and
   documentation are written in English.
10. **Domain errors, mapped at the boundary.** Business-rule failures
   (e.g. division by zero, negative square root) are typed errors
   defined in `calculator/`. Only `handlers/` maps them to HTTP status
   codes — `calculator/` never references HTTP.