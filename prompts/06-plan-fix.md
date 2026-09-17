Before I approve, one technical gap in "Justified technical decisions" #2
(generic decodeRequest[T]): as described, it can't actually produce the
field-specific messages the API contract examples promise (e.g. "b is
required"), because Go's encoding/json leaves a missing numeric field at
its zero value (0) — indistinguishable from a field that was explicitly
sent as 0 — when decoding into a plain float64 field. It CAN produce a
field-specific message for a wrong-type value (e.g. a string where a
number is expected), since that does error out with the field name
included.

Please revise decision #2 to use pointer fields (*float64) for every
required operand in each per-operation request struct (e.g.
AddRequest{A, B *float64}), so a nil pointer after decoding
unambiguously means "field missing," distinct from an explicitly-sent
zero. Update decodeRequest[T]'s described behavior (or the per-handler
logic that follows it) to check for nil on each required field and
produce the specific "<field> is required" message, and to surface the
existing field-name-inclusive error Go's decoder already gives for
wrong-JSON-type values. Keep this change scoped to decision #2 and the
backend module structure's request struct descriptions — no other
section should need to change.
