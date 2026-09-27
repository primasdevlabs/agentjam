# Debugging Instructions

## Protocol
1. Reproduce the failure reliably before touching code.
2. Form a hypothesis, then isolate variables to confirm it — do not shotgun changes.
3. Add targeted logging at the boundary where behavior diverges from expectation.
4. Fix the root cause, not the symptom; add a regression test that fails without the fix.
