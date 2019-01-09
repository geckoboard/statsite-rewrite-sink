package main

import (
	"regexp"

	"github.com/geckoboard/statsite-rewrite-sink/dsl"
)

var rules = []dsl.Rule{
	{
		CompleteMatch: "consul.http.{method}.{path}",
		ReplaceWith:   "consul.http",
		CustomPatterns: map[string]*regexp.Regexp{
			"path": regexp.MustCompile(`.+`),
		},
	},
	{
		RequiredPrefix: "envoy.",
		PartialMatch:   "cluster.{cluster}",
		ReplaceWith:    "cluster",
	},
	{
		RequiredPrefix: "envoy.",
		PartialMatch:   "vhost.{vhost}",
		ReplaceWith:    "vhost",
	},
	{
		RequiredPrefix: "envoy.",
		PartialMatch:   "listener.{listener}.http.",
		ReplaceWith:    "listener.http.",
		CustomPatterns: map[string]*regexp.Regexp{
			"listener": regexp.MustCompile(`.+`),
		},
	},
}
