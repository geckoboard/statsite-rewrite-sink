package dsl

import "regexp"

type Rule struct {
	MatchMetric    string
	ReplaceWith    string
	CustomPatterns map[string]*regexp.Regexp
}
