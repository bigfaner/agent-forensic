package sanitizer

import "regexp"

// sensitivePattern matches key=value patterns for API keys, secrets, tokens, and passwords.
// Group 1: the key name (api_key, secret, token, password, case-insensitive)
// Group 2: the separator (whitespace, colon, or equals)
// Group 3: optional opening quote
// Group 4: the sensitive value to be masked
var sensitivePattern = regexp.MustCompile(`(?i)(api_key|secret|token|password)([\s:=]+)(["']?)(\S+)`)

// Sanitize replaces sensitive values matching known patterns with ***.
// Returns sanitized content and whether any masking occurred.
func Sanitize(content string) (string, bool) {
	masked := false
	result := sensitivePattern.ReplaceAllStringFunc(content, func(match string) string {
		masked = true
		// ReplaceAllStringFunc only calls this for actual matches, so sub is always non-nil.
		sub := sensitivePattern.FindStringSubmatch(match)
		// Preserve: key + separator + opening_quote + "***"
		return sub[1] + sub[2] + sub[3] + "***"
	})
	return result, masked
}
