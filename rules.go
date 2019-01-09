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
		PartialMatch:   "cluster.{envoy_cluster}",
		ReplaceWith:    "cluster",
	},
	{
		RequiredPrefix: "envoy.",
		PartialMatch:   "vhost.{envoy_vhost}",
		ReplaceWith:    "vhost",
	},
	{
		RequiredPrefix: "envoy.",
		PartialMatch:   "listener.{envoy_listener}.http.",
		ReplaceWith:    "listener.http.",
		CustomPatterns: map[string]*regexp.Regexp{
			"envoy_listener": regexp.MustCompile(`.+`),
		},
	},
	{
		RequiredPrefix: "envoy.",
		PartialMatch:   "_rq_{envoy_http_status_code}",
		ReplaceWith:    "_rq_status_code",
		CustomPatterns: map[string]*regexp.Regexp{
			"envoy_http_status_code": regexp.MustCompile(`\d{3}`),
		},
	},
	{
		RequiredPrefix: "envoy.",
		PartialMatch:   "_rq_{envoy_http_status_class}",
		ReplaceWith:    "_rq_status_class",
		CustomPatterns: map[string]*regexp.Regexp{
			"envoy_http_status_class": regexp.MustCompile(`(?i)\dxx`),
		},
	},
}
