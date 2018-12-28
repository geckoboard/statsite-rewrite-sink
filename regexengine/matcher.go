package regexengine

import (
	"regexp"
	"strings"
)

type match struct {
	name string
	tags map[string]string
}

type matchers []matcher

func (m matchers) ExtractTagsFromMetric(metric Measurement) *match {
	result := &match{
		name: metric.Name(),
		tags: map[string]string{},
	}
	matched := false

	for _, matcher := range m {
		// Merge all results. Some rules will only modify part of a metric name
		if r := matcher.ApplyTo(result.name); r != nil {
			matched = true
			result.name = r.name
			for k, v := range r.tags {
				result.tags[k] = v
			}
		}
	}

	if !matched {
		return nil
	}

	return result
}

type matcher struct {
	Pattern     *regexp.Regexp
	ReplaceWith string
}

func (m matcher) ApplyTo(metricName string) *match {
	submatch := m.Pattern.FindStringSubmatch(metricName)
	if submatch == nil {
		return nil
	}

	name := m.ReplaceWith
	tags := map[string]string{}

	// Convert all the things we matched into tags
	for i, groupName := range m.Pattern.SubexpNames() {
		// The 0th group contains the entire string the regexp matched (i.e.
		// `metricName`), which we don't need
		if i == 0 {
			continue
		}

		tags[groupName] = submatch[i]
		name = strings.Replace(name, "{"+groupName+"}", submatch[i], -1)
	}

	return &match{
		name: name,
		tags: tags,
	}
}
