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

// CompileRulesIntoMatchers takes user-friendly, user-defined rules for
// extracting tags from metric names and compiles them into regex patterns that
// can be applied to metric names.
func CompileRulesIntoMatchers(rules []dsl.Rule) matchers {
	matches := make(matchers, 0, len(rules))

	for _, r := range rules {
		var patternPrefix, patternSuffix string

		matcher := matcher{RequiredPrefix: r.RequiredPrefix, ReplaceWith: r.ReplaceWith}

		rule := r.PartialMatch
		if r.CompleteMatch != "" {
			rule = r.CompleteMatch
			patternPrefix = "^"
			patternSuffix = "$"
		}

		pattern := compileRuleIntoPattern(rule, r.CustomPatterns)

		matcher.Pattern = regexp.MustCompile(
			fmt.Sprintf("%s%s%s", patternPrefix, pattern, patternSuffix),
		)

		matches = append(matches, matcher)
	}

	return matches
}

func compileRuleIntoPattern(rule string, customPatterns map[string]*regexp.Regexp) string {
	var regexPattern strings.Builder

	// A rule looks like
	// `some.literal.string.{a_tag_placeholder}.more.literal.strings`.
	// We need to find the position of all placeholders in the string so
	// that we can replace them with a regex pattern.
	positionsOfPlaceholders := tagPlaceholderPattern.FindAllStringSubmatchIndex(rule, -1)

	// No placeholders detected, so this is a literal string.
	// Not really sure why you'd have a string that only matched a literal,
	// but let's handle it to be on the safe side
	if positionsOfPlaceholders == nil {
		return regexp.QuoteMeta(rule)
	}

	// Assume that the rule begins with a literal string
	startOfNextLiteral := 0

	for _, indexes := range positionsOfPlaceholders {
		begin := indexes[0]
		end := indexes[1]

		literalBeforePlaceholder := rule[startOfNextLiteral:begin]
		regexPattern.WriteString(regexp.QuoteMeta(literalBeforePlaceholder))

		// Slice the `tag_name` part out of `{tag_name}`
		name := rule[begin+1 : end-1]
		pattern := `[^\.]+`

		// Allow custom patterns for the tag value. Useful if
		// the tag value contains dots (e.g. consul's HTTP
		// metrics contain a dot-separated path)
		if custom, ok := customPatterns[name]; ok {
			pattern = custom.String()
		}

		regexPattern.WriteString(fmt.Sprintf(`(?P<%s>%s)`, name, pattern))
		startOfNextLiteral = end
	}

	// Append any literal parts of the string that come after the final placeholder
	regexPattern.WriteString(regexp.QuoteMeta(rule[startOfNextLiteral:]))

	return regexPattern.String()
}
