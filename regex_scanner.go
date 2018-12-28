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

type convertMetric struct {
	match regexp.Regexp
}

var consulHTTPRequest = regexp.MustCompile(`^consul.http.(?P<method>[^\.]+)\.(?P<path>.+)$`)

func extractTagsFromMetricName(namespace string) *match {
	if m := consulHTTPRequest.FindStringSubmatch(namespace); m != nil {
		return &match{
			name: "consul.http",
			tags: map[string]string{
				"method": m[1],
				"path":   m[2],
			},
		}
	}

	return nil
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

func RegexScanner(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	namesToIndex := map[string]int{}

	for index, name := range statsitePattern.SubexpNames() {
		namesToIndex[name] = index
	}

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
		if taggedName := extractTagsFromMetricName(metricName); taggedName != nil {
			output = formatter(encodeTags(*taggedName), metric, namesToIndex)
		}

		fmt.Fprintln(out, output)
	}
}
