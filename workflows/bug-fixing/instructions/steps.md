# Bug Fixing Protocol

1. Read full error logs and stack trace tracebacks. No guessing root cause without empirical evidence.
2. Fix root cause contract breakage. Never swallow exceptions or introduce dummy fallbacks.
3. Verify fix using automated test runner.
