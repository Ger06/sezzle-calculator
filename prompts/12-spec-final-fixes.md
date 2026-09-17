Update specs/001-calculator/spec.md with three corrections found during
implementation:

1. RF-20 (Delete last digit) needs a clarification: it didn't specify
   behavior when there is no operand currently being entered — e.g.
   immediately after a result is displayed (RF-16), before any new digit
   has been pressed. Resolved as: Backspace is a no-op in that state,
   consistent with how a fresh digit press in the same state already
   replaces the display rather than appending to it (both treat a
   freshly shown, untouched result as frozen, not editable) — not with
   the alternative of editing the result in place, which would have been
   inconsistent with digit-press's existing replace behavior.

2. RF-21 (Keyboard input) needs a clarification: the global keydown
   listener must not override native button activation when focus is
   already on one of the calculator's own interactive controls (e.g.
   Tab-focusing the `+` button, then pressing Enter, must activate that
   focused button via the browser's native Enter/Space handling — not
   be intercepted by the global listener's own Enter-means-equals
   mapping). This was found via manual screen-reader testing (Tab
   navigation + Enter/Space to activate a focused button).

3. Add a new non-functional requirement, RNF-5 (Touch accessibility):
   all interactive controls (Keypad and OperationButtons) shall have a
   minimum touch target size of 44x44px (mobile viewports use 56px per
   the approved reference design), and shall use `:active` state styling
   for pressed feedback rather than relying solely on `:hover`, since
   hover has no equivalent on touch devices.

Show me the updated sections (RF-20, RF-21, and the new RNF-5) — no need
to repaste the entire document.
