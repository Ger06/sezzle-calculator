Bug found: pressing Backspace right after a freshly displayed result
(e.g. 9 × 8 = 72) currently edits that result digit-by-digit (72 → 7),
but pressing a fresh digit in that same state replaces the display
entirely rather than appending — these two actions disagree on whether a
freshly shown result is "editable" or "frozen," which is inconsistent.

Fix: make a freshly displayed result frozen for both actions, matching
standard calculator convention. Backspace pressed while a fresh,
untouched result is displayed shall be a no-op (same rule as "equals
with no operator pending" already being a no-op) — it does nothing until
the user starts a new entry (a digit, which starts a fresh number per
RF-27, or an operator, which uses the result as the first operand per
RF-27). Add a test asserting Backspace is a no-op immediately after a
result, and confirm the existing digit-after-result behavior (replace,
not append) already has a test — add one if it doesn't.
