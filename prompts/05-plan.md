Read docs/constitution.md and specs/001-calculator/spec.md. DO NOT write
any code. Generate specs/001-calculator/plan.md with:

- Module structure for both backend and frontend (backend: calculator/,
  handlers/, main.go; frontend: components, reducer, API client — propose
  the exact layout)
- REST API contract: every endpoint, request body example, success
  response example, and error response example, with HTTP status codes
- Domain error → HTTP status mapping table (e.g. division by zero,
  negative square root)
- Justified technical decisions (and their discarded alternative) for
  anything not already fixed in docs/decisions.md — reference
  docs/decisions.md instead of repeating a decision already made there
- Test strategy for both layers, consistent with the TDD principle in the
  constitution

No data model section: the calculator is stateless per constitution
principle 7, there is nothing to persist.

Everything must respect the constitution and cover all RF. Mark which RF
each part covers.
