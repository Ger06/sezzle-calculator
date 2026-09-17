Update specs/001-calculator/spec.md with two more corrections:

1. RF-20 (Delete last digit) gets an on-screen control in addition to the
   Backspace key — remove the "keyboard-only; there is no on-screen
   control" wording. The on-screen control is an icon button (⌫-style,
   matching the standard Android calculator convention), not a text
   label — same delete-last-character behavior as before, just now also
   triggerable by tapping this button, not only by pressing Backspace.
   Update RF-21 (keyboard input) accordingly if it referenced RF-20 as
   keyboard-only.

2. Add a new requirement: when a result is currently displayed (RF-16)
   and the user presses a binary operator, the system shall use that
   displayed result as the first operand of the new operation — the same
   flow as if the user had just entered it as a fresh first operand
   (RF-14). This is distinct from RF-14's existing chaining clause, which
   covers pressing a new operator before equals; this covers pressing an
   operator AFTER equals has already produced and displayed a result
   (e.g. 2+2=4, then +3=7).

Update the completion criteria to include both of these. Show me the
updated sections (RF-19-21 area, the new requirement, and completion
criteria) — no need to repaste the entire document if only these parts
changed.
