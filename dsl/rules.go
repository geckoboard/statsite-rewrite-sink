package dsl

import "regexp"

type Rule struct {
	// CompleteMatch indicates that the full metric name must match this exact pattern
	// You can use `{tag_name}` placeholders to extract a value for the tag `tag_name`.
	// By default the placeholder will match anything that is not a period `.`, but this
	// can be overriden using `CustomPatterns`
	CompleteMatch string

	// PartialMatch allows you to extract and modify part of a metric's name, rather
	// than matching/replacing the whole metric name.
	PartialMatch string

	// Whatever is matched will be replaced with this. You can use `{tag_name}` placeholders
	// to interpolate tag values into the replacement.
	ReplaceWith string

	// RequiredPrefix is useful for restricting which metrics partial matches are applied to.
	// The full metric name must have this prefix.
	RequiredPrefix string

	// By default placeholders match any text that does not contain dots
	// Some systems (e.g. consul) emit metrics that use dots in significant
	// values. The key should match the name of the placeholder (e.g.
	// `tag_name` in `{tag_name}`). The default value is `[^\.]+`
	CustomPatterns map[string]*regexp.Regexp
}
