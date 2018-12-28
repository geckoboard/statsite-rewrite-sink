package regexengine

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/geckoboard/statsite-rewrite-sink/dsl"
)

// Heavily inspired by https://github.com/seatgeek/statsd-rewrite-proxy/blob/master/regex.go
func CompileRulesIntoMatchers(rules []dsl.Rule) matchers {
	matches := make(matchers, 0, len(rules))

	for _, r := range rules {
		components := strings.Split(r.MatchMetric, ".")
		regexParts := make([]string, 0, len(components))

		for _, part := range components {
			switch part[:1] {
			case "{": // This is a rule for extracting a tag in `{tag_name}` format
				name := part[1 : len(part)-1]
				pattern := `[^\.]+`

				// Allow custom patterns for the tag value.
				// Useful if the tag value contains dots (e.g.
				// consul's HTTP metrics contain a
				// dot-separated path)
				if custom, ok := r.CustomPatterns[name]; ok {
					pattern = custom.String()
				}

				regexParts = append(regexParts, fmt.Sprintf(`(?P<%s>%s)`, name, pattern))
			case "*": // Allow wildcards for ignoring parts of a metric we don't care about
				regexParts = append(regexParts, ".+?")
			default:
				regexParts = append(regexParts, regexp.QuoteMeta(part))

			}
		}

		pattern := strings.Join(regexParts, `\.+`)
		matches = append(matches, matcher{Pattern: regexp.MustCompile(pattern), ReplaceWith: r.ReplaceWith})
	}

	return matches
}
