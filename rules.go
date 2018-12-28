package main

import "regexp"

type rule struct {
	MatchMetric    string
	ReplaceWith    string
	CustomPatterns map[string]*regexp.Regexp
}

var rules = []rule{
	{
		MatchMetric: "consul.http.{method}.{path}",
		ReplaceWith: "consul.http",
		CustomPatterns: map[string]*regexp.Regexp{
			"path": regexp.MustCompile(`.+`),
		},
	},
}
