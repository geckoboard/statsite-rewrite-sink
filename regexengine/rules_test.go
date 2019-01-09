package regexengine

import (
	"regexp"
	"testing"

	"github.com/geckoboard/statsite-rewrite-sink/dsl"
)

func TestRuleCompilation(t *testing.T) {
	type example struct {
		rule            dsl.Rule
		expectedMatcher matcher
	}

	examples := []example{
		{
			rule: dsl.Rule{
				CompleteMatch: "consul.http.{method}.{path}",
				ReplaceWith:   "consul.http",
				CustomPatterns: map[string]*regexp.Regexp{
					"path": regexp.MustCompile(`.+`),
				},
			},
			expectedMatcher: matcher{
				Pattern:     regexp.MustCompile(`^consul\.http\.(?P<method>[^\.]+)\.(?P<path>.+)$`),
				ReplaceWith: "consul.http",
			},
		},
		{

			rule: dsl.Rule{
				RequiredPrefix: "envoy.",
				PartialMatch:   "cluster.{cluster}",
				ReplaceWith:    "cluster",
			},
			expectedMatcher: matcher{
				Pattern:        regexp.MustCompile(`cluster\.(?P<cluster>[^\.]+)`),
				ReplaceWith:    "cluster",
				RequiredPrefix: "envoy.",
			},
		},
		{

			rule: dsl.Rule{
				RequiredPrefix: "envoy.",
				PartialMatch:   "_rq_{envoy_response_code}",
				ReplaceWith:    "_rq_status_code",
				CustomPatterns: map[string]*regexp.Regexp{
					"envoy_response_code": regexp.MustCompile(`\d{3}`),
				},
			},
			expectedMatcher: matcher{
				Pattern:        regexp.MustCompile(`_rq_(?P<envoy_response_code>\d{3})`),
				ReplaceWith:    "_rq_status_code",
				RequiredPrefix: "envoy.",
			},
		},
		{

			rule: dsl.Rule{
				PartialMatch: "a_prefix_{a_tag_placeholder}_a_suffix",
				ReplaceWith:  "a_thing",
			},
			expectedMatcher: matcher{
				Pattern:     regexp.MustCompile(`a_prefix_(?P<a_tag_placeholder>[^\.]+)_a_suffix`),
				ReplaceWith: "a_thing",
			},
		},
	}

	for _, ex := range examples {
		matchers := CompileRulesIntoMatchers([]dsl.Rule{ex.rule})
		actual := matchers[0]
		expected := ex.expectedMatcher

		if actual.Pattern.String() != expected.Pattern.String() {
			t.Errorf("got %q, expected %q", actual.Pattern.String(), expected.Pattern.String())
		}

		if actual.ReplaceWith != expected.ReplaceWith {
			t.Errorf("got %q, expected %q", actual.ReplaceWith, expected.ReplaceWith)
		}

		if actual.RequiredPrefix != expected.RequiredPrefix {
			t.Errorf("got %q, expected %q", actual.RequiredPrefix, expected.RequiredPrefix)
		}
	}
}
