package dsl

import "regexp"

type Rule struct {
	// CompleteMatch indicates that the full metric name must match this exact pattern
	CompleteMatch string
	// PartialMatch allows you to extract and modify part of a metric's name
	PartialMatch string
	// Whatever is matched will be replaced with this
	ReplaceWith    string
	RequiredPrefix string
	// By default placeholders match any text that does not contain dots
	// Some systems (e.g. consul) emit metrics that use dots in significant
	// values
	CustomPatterns map[string]*regexp.Regexp
}
