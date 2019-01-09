package regexengine

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/geckoboard/statsite-rewrite-sink/dsl"
)

var (
	tagPlaceholderPattern = regexp.MustCompile(`({[^}]+})`)
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

		// Some metrics do not separate tag names from tag values with
		// periods, and instead interpolate the value as part of the
		// name.
		//
		// An example of this is envoy's "response code"/"response code
		// class" tags which appear as `_rq_200` or `_rq_2XX` in a
		// metric like `envoy.cluster.tuner-grpc.upstream_rq_200`
		positionsOfPlaceholders := tagPlaceholderPattern.FindAllStringSubmatchIndex(part, -1)

		// No placeholders detected, so this is a literal string
		if positionsOfPlaceholders == nil {
			regexParts = append(regexParts, regexp.QuoteMeta(part))
			continue
		}

		var begin, end int
		var regexPart strings.Builder
		startOfLiteral := 0

		for _, indexes := range positionsOfPlaceholders {
			begin = indexes[0]
			end = indexes[1]

			regexPart.WriteString(regexp.QuoteMeta(part[startOfLiteral:begin]))

			// Slice the `tag_name` part out of `{tag_name}`
			name := part[begin+1 : end-1]
			pattern := `[^\.]+`

			// Allow custom patterns for the tag value. Useful if
			// the tag value contains dots (e.g. consul's HTTP
			// metrics contain a dot-separated path)
			if custom, ok := customPatterns[name]; ok {
				pattern = custom.String()
			}

			regexPart.WriteString(fmt.Sprintf(`(?P<%s>%s)`, name, pattern))
		}

		// Append any literal parts of the string after the last placeholder
		regexPart.WriteString(regexp.QuoteMeta(part[end:]))

		regexParts = append(regexParts, regexPart.String())
	}

	return strings.Join(regexParts, `\.+`)
}
