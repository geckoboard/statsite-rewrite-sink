package main

import (
	"regexp"

	"github.com/geckoboard/statsite-rewrite-sink/dsl"
)

var rules = []dsl.Rule{
	{
		MatchMetric: "consul.http.{method}.{path}",
		ReplaceWith: "consul.http",
		CustomPatterns: map[string]*regexp.Regexp{
			"path": regexp.MustCompile(`.+`),
		},
	},
}
