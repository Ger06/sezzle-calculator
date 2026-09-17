DO NOT write any code at any point. Let's write the specification for the
first feature of sezzle-calculator: the calculator API and UI. Read
docs/constitution.md.

Initial idea: a Go backend exposing one REST endpoint per operation (add,
subtract, multiply, divide, power, sqrt, percentage), and a React +
TypeScript frontend consuming it. Percentage means "X% of value"
(value * percentage / 100). Division by zero and square root of a negative
number must return domain errors. The backend must allow cross-origin
requests from the frontend during development. Results must be rounded to
10 decimal places before being returned, to avoid floating-point artifacts.

Your job:
1. Ask me ONE question at a time to remove ambiguities (edge cases, error
   behavior, what's out of MVP scope). Max 6 questions.
2. With my answers, generate specs/001-calculator/spec.md with this
   structure: context and goal, users, user stories, numbered functional
   requirements (RF-x) with acceptance criteria in EARS notation, non-
   functional requirements, edge cases, out of scope, completion criteria,
   and open questions marked as [NEEDS CLARIFICATION].
3. The WHAT and the WHY. No stack, architecture, or file names: that goes
   in the plan.
