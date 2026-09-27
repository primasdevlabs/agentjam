package policy

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/primasdevlabs/agentjam/pkg/core"
)

// Hygiene and insecure-sink evaluators — SAST-style checks gated on loaded
// policy categories.

var (
	evalSinkPattern     = regexp.MustCompile(`(?i)\b(eval\s*\(|new\s+Function\s*\(|innerHTML\s*=|outerHTML\s*=|dangerouslySetInnerHTML|document\.write\s*\()`)
	tlsBypassPattern    = regexp.MustCompile(`(?i)(InsecureSkipVerify|verify\s*=\s*False|rejectUnauthorized\s*:\s*false|CURLOPT_SSL_VERIFYPEER.*0)`)
	shellExecPattern    = regexp.MustCompile(`(?i)(os\.system|child_process\.(exec|execSync)|Runtime\.getRuntime\(\)\.exec|subprocess\.(call|run|Popen)\s*\(\s*['"]?[a-z]|exec\.Command\()\s*[^)]*\+`)
	weakCryptoPattern   = regexp.MustCompile(`(?i)(md5|sha1)\s*\([^)]*(password|passwd|secret|token)|(password|passwd|secret|token)[^;]{0,40}(md5|sha1)\s*\(`)
	plainHTTPPattern    = regexp.MustCompile(`(?i)["']http://[^"'\s]+`)
	debugPatterns       = regexp.MustCompile(`(?i)(console\.(log|debug|warn)\s*\(|\bdebugger\b|pdb\.set_trace|binding\.pry|fmt\.Print(ln|f)?\s*\(|println!|System\.out\.print)`)
	todoPattern         = regexp.MustCompile(`\b(TODO|FIXME|HACK|XXX)\b`)
	localhostURLPattern = regexp.MustCompile(`(?i)^["']?http://(localhost|127\.0\.0\.1|0\.0\.0\.0|\[::1\])`)
)

// EvaluateSinkRules flags insecure API sinks. Requires a loaded "security"
// category policy.
func EvaluateSinkRules(policies []core.PolicyManifest, content string) []EvaluationViolation {
	violations := make([]EvaluationViolation, 0)
	if !policyActive(policies, "security") {
		return violations
	}
	type sinkRule struct {
		re      *regexp.Regexp
		rule    string
		name    string
		enforce core.PolicyEnforcement
		message string
	}
	for _, sr := range []sinkRule{
		{evalSinkPattern, "no-eval-or-dom-injection", "Insecure API Sinks", core.EnforceStrictBlock,
			"Unsafe evaluation or raw HTML injection (eval, innerHTML, dangerouslySetInnerHTML) — use safe rendering/parsing APIs."},
		{tlsBypassPattern, "no-tls-verify-bypass", "Insecure API Sinks", core.EnforceStrictBlock,
			"TLS certificate verification disabled. Never ship InsecureSkipVerify/verify=False/rejectUnauthorized:false."},
		{shellExecPattern, "no-shell-string-exec", "Insecure API Sinks", core.EnforceStrictBlock,
			"Shell command built via string concatenation — command injection risk. Use argument arrays, never joined strings."},
		{weakCryptoPattern, "no-weak-crypto-for-secrets", "Insecure API Sinks", core.EnforceStrictBlock,
			"Weak hash (MD5/SHA1) applied to credentials. Use bcrypt, argon2, or scrypt."},
	} {
		if sr.re.MatchString(content) {
			violations = append(violations, EvaluationViolation{
				PolicyID:    "05-insecure-api-sinks",
				PolicyName:  sr.name,
				Enforcement: sr.enforce,
				RuleName:    sr.rule,
				Message:     sr.message,
				Line:        firstMatchLine(sr.re, content),
			})
		}
	}

	// Plain http:// — RE2 lacks lookaheads, so localhost URLs are filtered
	// out after matching.
	for _, m := range plainHTTPPattern.FindAllString(content, -1) {
		if localhostURLPattern.MatchString(m) {
			continue
		}
		violations = append(violations, EvaluationViolation{
			PolicyID:    "05-insecure-api-sinks",
			PolicyName:  "Insecure API Sinks",
			Enforcement: core.EnforceWarning,
			RuleName:    "prefer-https",
			Message:     "Plain http:// URL — remote endpoints must use TLS.",
			Line:        firstMatchLine(plainHTTPPattern, content),
		})
		break
	}
	return violations
}

// EvaluateHygieneRules flags debug leftovers and deferred-work accumulation.
// Requires a loaded "coding-standards" category policy.
func EvaluateHygieneRules(policies []core.PolicyManifest, content string) []EvaluationViolation {
	violations := make([]EvaluationViolation, 0)
	if !policyActive(policies, "coding-standards") {
		return violations
	}

	if debugPatterns.MatchString(content) {
		violations = append(violations, EvaluationViolation{
			PolicyID:    "08-debug-leftovers",
			PolicyName:  "Debug Statement Hygiene",
			Enforcement: core.EnforceWarning,
			RuleName:    "no-debug-statements",
			Message:     "Debug statements left in source (console.log, debugger, pdb, fmt.Println). Remove before committing.",
			Line:        firstMatchLine(debugPatterns, content),
		})
	}

	maxMarkers := 5
	if p := policyByID(policies, "todo-hygiene"); p != nil {
		maxMarkers = ruleInt(p, "max_markers_per_file", maxMarkers)
	}
	markers := len(todoPattern.FindAllString(content, -1))
	if markers > maxMarkers {
		violations = append(violations, EvaluationViolation{
			PolicyID:    "09-todo-hygiene",
			PolicyName:  "Deferred Work Hygiene",
			Enforcement: core.EnforceInfo,
			RuleName:    "max_markers_per_file",
			Message:     messagef("%d TODO/FIXME/HACK markers in one file (limit %d) — resolve or ticket them.", markers, maxMarkers),
			Line:        firstMatchLine(todoPattern, content),
		})
	}
	return violations
}

// deprecatedAPIPatterns are source-level usages of deprecated or removed
// language/runtime APIs. Kept to high-signal patterns to avoid noise.
var deprecatedAPIRules = []struct {
	re      *regexp.Regexp
	rule    string
	enforce core.PolicyEnforcement
	message string
}{
	{regexp.MustCompile(`\bio/ioutil\b`), "no-deprecated-stdlib", core.EnforceWarning,
		"io/ioutil is deprecated since Go 1.16 — use io/os equivalents."},
	{regexp.MustCompile(`\brand\.Seed\s*\(`), "no-deprecated-stdlib", core.EnforceWarning,
		"rand.Seed is deprecated since Go 1.20 — use rand.New(rand.NewSource(seed)) or the global generator."},
	{regexp.MustCompile(`\bnew\s+Buffer\s*\(`), "no-deprecated-node-apis", core.EnforceWarning,
		"new Buffer() is deprecated and unsafe (uninitialized memory) — use Buffer.alloc/Buffer.from."},
	{regexp.MustCompile(`\bmysql_(query|connect|select_db|fetch_)\w*\s*\(`), "no-removed-php-apis", core.EnforceStrictBlock,
		"mysql_* functions were removed in PHP 7 — use PDO or mysqli."},
	{regexp.MustCompile(`\bcreate_function\s*\(`), "no-removed-php-apis", core.EnforceStrictBlock,
		"create_function() was removed in PHP 8 and is an injection risk — use closures."},
	{regexp.MustCompile(`(?m)^\s*(import|from)\s+distutils\b`), "no-removed-py-modules", core.EnforceWarning,
		"distutils was removed in Python 3.12 — migrate to setuptools/packaging."},
	{regexp.MustCompile(`\b(datetime\.)?(utcnow|utcfromtimestamp)\s*\(`), "no-deprecated-py-apis", core.EnforceWarning,
		"utcnow()/utcfromtimestamp() are deprecated in Python 3.12 — use datetime.now(timezone.utc)."},
	{regexp.MustCompile(`\bnp\.(float|int|object|bool)\s*[(\[]`), "no-deprecated-numpy-aliases", core.EnforceWarning,
		"np.float/np.int/np.object/np.bool aliases were removed in NumPy 1.24 — use builtin types."},
	{regexp.MustCompile(`\bnew\s+(Integer|Long|Double|Float|Short|Byte|Character|Boolean)\s*\(`), "no-deprecated-java-boxing", core.EnforceWarning,
		"Primitive wrapper constructors are deprecated for removal — use valueOf() or autoboxing."},
	{regexp.MustCompile(`\b(WebRequest\.Create|HttpWebRequest)\b`), "no-obsolete-dotnet-apis", core.EnforceWarning,
		"WebRequest/HttpWebRequest are obsolete — use HttpClient."},
}

// EvaluateDeprecatedAPIs flags deprecated/removed runtime API usage.
// Requires a loaded "deprecated-packages" security policy.
func EvaluateDeprecatedAPIs(policies []core.PolicyManifest, content string) []EvaluationViolation {
	violations := make([]EvaluationViolation, 0)
	if policyByID(policies, "deprecated-packages") == nil {
		return violations
	}
	for _, dr := range deprecatedAPIRules {
		if dr.re.MatchString(content) {
			violations = append(violations, EvaluationViolation{
				PolicyID:    "07-deprecated-packages",
				PolicyName:  "Deprecated Package Guard",
				Enforcement: dr.enforce,
				RuleName:    dr.rule,
				Message:     dr.message,
				Line:        firstMatchLine(dr.re, content),
			})
		}
	}
	return violations
}

// sensitiveFileBasenames/exns drive the committed-secrets file guard.
var sensitiveFileBasenames = map[string]bool{
	".env": true, ".env.local": true, ".env.production": true,
	"credentials.json": true, "serviceaccount.json": true,
	"id_rsa": true, "id_dsa": true, "id_ed25519": true,
}

var sensitiveFileExtensions = map[string]bool{
	".pem": true, ".key": true, ".p12": true, ".pfx": true, ".keystore": true,
}

// EvaluateSensitivePath flags committed sensitive files by name. Requires a
// loaded "security" category policy.
func (pe *PolicyEngine) EvaluateSensitivePath(relPath string) []EvaluationViolation {
	p := policyByID(pe.policies, "sensitive-files")
	if p == nil {
		return nil
	}
	base := strings.ToLower(filepath.Base(relPath))
	if sensitiveFileBasenames[base] || sensitiveFileExtensions[strings.ToLower(filepath.Ext(base))] {
		return []EvaluationViolation{{
			PolicyID:    p.ID,
			PolicyName:  p.Name,
			Enforcement: p.Enforcement,
			RuleName:    "no-committed-sensitive-files",
			Message:     messagef("%q looks like a committed secret/key material — remove it and rotate its contents.", filepath.Base(relPath)),
		}}
	}
	return nil
}
