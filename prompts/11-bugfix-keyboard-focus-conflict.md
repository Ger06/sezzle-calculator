Bug found: when navigating the calculator with Tab (keyboard/screen-
reader accessibility) and pressing Enter or Space while focus is on an
operator/digit button, the button doesn't activate as expected — the
global keydown listener (RF-21) appears to intercept Enter (mapped to
equals) regardless of which element currently has focus, competing with
the browser's native Enter/Space button-activation behavior.

Fix: the global keydown listener in Calculator.tsx should skip its own
handling when the event's target is already one of the calculator's own
interactive elements (a <button>) — let the browser's native click-on-
Enter/Space behavior handle it in that case, instead of the global
listener also acting on the same keypress. A simple check: if
(event.target instanceof HTMLButtonElement) return early before doing
any of the digit/operator/equals/clear/backspace dispatch logic.

Add a test (or manual verification note, if this specific interaction is
hard to express in RTL) confirming that Tab-focusing an operator button
and pressing Enter activates that button's own action, not the global
equals shortcut, and that Tab-focusing has no keydown side effects when
the event originates from a button.
