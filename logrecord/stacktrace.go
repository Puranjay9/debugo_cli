package logrecord

import "regexp"

type StackTrace struct {
	Lines []string
}

var traceStartPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\bpanic\b`),
	regexp.MustCompile(`(?i)\bexception\b`),
	regexp.MustCompile(`(?i)\btraceback\b`),
	regexp.MustCompile(`(?i)\bfatal\b`),
	regexp.MustCompile(`(?i)\berror\b`),
}

var traceContPatterns = []*regexp.Regexp{
	regexp.MustCompile(`^\s+at\s+`),
	regexp.MustCompile(`^\s+File\s+"`),
	regexp.MustCompile(`^\s+/.*:\d+`),
	regexp.MustCompile(`^\s+goroutine`),
	regexp.MustCompile(`^\s+`),
}

func IsTraceStart(line string) bool {
	for _, p := range traceStartPatterns {
		if p.MatchString(line) {
			return true
		}
	}
	return false
}

func IsTraceContinuation(line string) bool {
	for _, p := range traceContPatterns {
		if p.MatchString(line) {
			return true
		}
	}
	return false
}
