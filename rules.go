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
		PartialMatch:   "listener.{envoy_listener}.",
		ReplaceWith:    "listener.",
		CustomPatterns: map[string]*regexp.Regexp{
			// Inspired by https://github.com/envoyproxy/envoy/blob/87553968ec2258919d986e9a76512b0009d01575/source/common/config/well_known_names.cc#L102
			// My understanding is that this matches an ip/port combo,
			// or a UUID?
			// https://www.envoyproxy.io/docs/envoy/latest/api-v2/api/v2/lds.proto#envoy-api-msg-listener
			"envoy_listener": regexp.MustCompile(`(?:[\d\._]+|[_\[\]aAbBcCdDeEfF\d]+)`),
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
	{
		RequiredPrefix: "envoy.",
		PartialMatch:   "ssl.ciphers.{envoy_ssl_cipher}",
		ReplaceWith:    "ssl.ciphers",
	},
	{
		RequiredPrefix:         "envoy.",
		DropMeasurementsOfZero: true,
	},
}
