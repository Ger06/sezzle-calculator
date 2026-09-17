Go through specs/001-calculator/spec.md requirement by requirement (all
RF-x, backend and frontend). For each one, state: which test(s) cover it
(unit and/or integration), and the result of running them. If any RF is
not covered or fails, say so clearly — do not soften it.

Then run the full test suite for both layers (`go test ./...` in backend/,
`npm run test` in frontend/) and report the results.

Finally, check the completion criteria from the spec and give me a
verdict: is the spec fully satisfied? List anything left unmet.
