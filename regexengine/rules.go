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
		var patternPrefix, patternSuffix string
		var components []string

		matcher := matcher{RequiredPrefix: r.RequiredPrefix, ReplaceWith: r.ReplaceWith}

		if r.CompleteMatch != "" {
			components = strings.Split(r.CompleteMatch, ".")
			patternPrefix = "^"
			patternSuffix = "$"
		} else {
			components = strings.Split(r.PartialMatch, ".")
		}

		pattern := compilePartsIntoPattern(components, r.CustomPatterns)

		regex := regexp.MustCompile(fmt.Sprintf("%s%s%s", patternPrefix, pattern, patternSuffix))

		matcher.Pattern = regex

		matches = append(matches, matcher)
	}

	return matches
}

func compilePartsIntoPattern(parts []string, customPatterns map[string]*regexp.Regexp) string {
	regexParts := make([]string, 0, len(parts))

	for _, part := range parts {
		// The user-supplied pattern is split on `.`, so if a pattern
		// begins or ends with a `.`, then at least one part will be an empty
		// string e.g.
		// `listener.{listener}.http.`
		// is represented as ["listener", "{listener}", "http", ""]
		if part == "" {
			regexParts = append(regexParts, "")
			continue
		}

		switch part[:1] {
		case "{": // This is a rule for extracting a tag in `{tag_name}` format
			name := part[1 : len(part)-1]
			pattern := `[^\.]+`

			// Allow custom patterns for the tag value.
			// Useful if the tag value contains dots (e.g.
			// consul's HTTP metrics contain a
			// dot-separated path)
			if custom, ok := customPatterns[name]; ok {
				pattern = custom.String()
			}

			regexParts = append(regexParts, fmt.Sprintf(`(?P<%s>%s)`, name, pattern))
		case "*": // Allow wildcards for ignoring parts of a metric we don't care about
			regexParts = append(regexParts, ".+?")
		default:
			regexParts = append(regexParts, regexp.QuoteMeta(part))

		}
	}

	return strings.Join(regexParts, `\.+`)
}
