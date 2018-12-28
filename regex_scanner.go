package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
)

var oldStatsitePattern = regexp.MustCompile(`\A(?P<type>[^.]+)\.(?P<key>[^|]+)(?P<trailing>.*)\z`)

var matchTimerMetric = `(?:(?P<timertype>timers)\.(?P<timermetric>[^|]+)\.(?P<timerattr>sum|sum_sq|mean|lower|upper|count|stdev|p\d\d?)`
var matchGenericMetric = `(?P<wildtype>[^\.]+)\.(?P<wildmetric>[^|]+))`
var matchMeasurement = `(?P<measurements>\|.+)`

var statsitePattern = regexp.MustCompile(fmt.Sprintf(`(?:%s|%s)%s`, matchTimerMetric, matchGenericMetric, matchMeasurement))

type match struct {
	name string
	tags map[string]string
}

type matchers []matcher

func (m matchers) ExtractTagsFromMetric(metricName string) *match {
	result := &match{
		name: metricName,
		tags: map[string]string{},
	}

	for _, matcher := range m {
		// Merge all results. Some rules will only modify part of a metric name
		if r := matcher.ApplyTo(result.name); r != nil {
			result.name = r.name
			for k, v := range r.tags {
				result.tags[k] = v
			}
		}
	}

	if result.name == metricName && len(result.tags) == 0 {
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

// Heavily inspired by https://github.com/seatgeek/statsd-rewrite-proxy/blob/master/regex.go
func CompileRulesIntoMatchers(rules []rule) matchers {
	matches := make(matchers, 0, len(rules))

	for _, r := range rules {
		components := strings.Split(r.MatchMetric, ".")
		regexParts := make([]string, 0, len(components))

		for _, part := range components {
			switch part[:1] {
			case "{": // This is a rule for extracting a tag in `{tag_name}` format
				name := part[1 : len(part)-1]
				pattern := `[^\.]+`

				// Allow custom patterns for the tag value.
				// Useful if the tag value contains dots (e.g.
				// consul's HTTP metrics contain a
				// dot-separated path)
				if custom, ok := r.CustomPatterns[name]; ok {
					pattern = custom.String()
				}

				regexParts = append(regexParts, fmt.Sprintf(`(?P<%s>%s)`, name, pattern))
			case "*": // Allow wildcards for ignoring parts of a metric we don't care about
				regexParts = append(regexParts, ".+?")
			default:
				regexParts = append(regexParts, regexp.QuoteMeta(part))

			}
		}

		pattern := strings.Join(regexParts, `\.+`)
		matches = append(matches, matcher{Pattern: regexp.MustCompile(pattern), ReplaceWith: r.ReplaceWith})
	}

	return matches
}

func encodeTags(m match) string {
	keys := make([]string, 0, len(m.tags))
	pairs := make([]string, 0, len(m.tags))

	// go doesn't iterate over maps in a consistent, ordered way
	// We need to sort the keys to ensure stability of our tests
	for k, _ := range m.tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, m.tags[k]))
	}

	return fmt.Sprintf("%s#%s", m.name, strings.Join(pairs, ","))
}

func genericFormatter(newMetricName string, originalMatchData []string, matchNameToIndex map[string]int) string {
	typeOfMeasurement := originalMatchData[matchNameToIndex["wildtype"]]
	measurementData := originalMatchData[matchNameToIndex["measurements"]]

	return fmt.Sprintf("%s.%s%s", typeOfMeasurement, newMetricName, measurementData)
}

func timerFormatter(newMetricName string, originalMatchData []string, matchNameToIndex map[string]int) string {
	timerAttribute := originalMatchData[matchNameToIndex["timerattr"]]

	return fmt.Sprintf(
		"timers.%s.%s%s",
		newMetricName,
		timerAttribute,
		originalMatchData[matchNameToIndex["measurements"]],
	)
}

func RegexScanner(in io.Reader, out io.Writer, rules []rule) {
	scanner := bufio.NewScanner(in)

	namesToIndex := map[string]int{}

	for index, name := range statsitePattern.SubexpNames() {
		namesToIndex[name] = index
	}

	matchers := CompileRulesIntoMatchers(rules)

	for scanner.Scan() {
		line := scanner.Text()

		metric := statsitePattern.FindStringSubmatch(line)

		// This is the entire line we found
		output := metric[0]

		// By default assume we're matching a generic metric rather than a timer
		metricName := metric[namesToIndex["wildmetric"]]
		formatter := genericFormatter

		// Timers are unusual - a single timer is represented by ~10 different metrics
		// Each has the same user chosen "name", but a different suffix to indicate the
		// attribute this measurement represents
		if metric[namesToIndex["timertype"]] != "" {
			metricName = metric[namesToIndex["timermetric"]]
			formatter = timerFormatter
		}

		// Only change the format if we could extract something from it
		if taggedName := matchers.ExtractTagsFromMetric(metricName); taggedName != nil {
			output = formatter(encodeTags(*taggedName), metric, namesToIndex)
		}

		fmt.Fprintln(out, output)
	}
}
