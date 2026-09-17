# Task implementation prompts — sezzle-calculator

## Variant A — new pattern (stop after this one task)
Use for the first task of each "type": first pure function + test, first
HTTP handler, typed errors, CORS middleware, the reducer, the first
component, etc.

```
Implement ONLY task T[X] from specs/001-calculator/tasks.md, following
plan.md and the constitution. Write the tests first, then the code. Run
the tests and show me the result. When done: mark T[X] in tasks.md, state
which RF it covers, and STOP. Do not start the next task.
```

## Variant B — repeated pattern (batch 2-3 tasks, one stop)
Use once a pattern has already been reviewed and verified (e.g. after
reviewing `Add`, batch `Subtract`/`Multiply`/`Divide`).

```
Implement tasks T[X], T[X+1], and T[X+2] from specs/001-calculator/tasks.md,
following plan.md, the constitution, and the same pattern used in T[Y].
Write the tests first, then the code, for each task. Run the tests and
show me the result for all three. When done: mark T[X]-T[X+2] in tasks.md,
state which RF each covers, and STOP. Do not start the next task.
```

## Rule of thumb
- First task of a new type (math operation category, handler, component,
  anything touching errors/CORS/reducer/architecture) → Variant A, always.
- Repetition of an already-reviewed pattern → Variant B, batch of 2-3.
- If in doubt, use Variant A. Reviewing one extra time costs less than an
  error repeated across several tasks.
