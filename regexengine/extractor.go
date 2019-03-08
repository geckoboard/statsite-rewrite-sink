package regexengine

import (
	"regexp"
	"strings"
)

type extractor struct {
	Pattern        *regexp.Regexp
	ReplaceWith    string
	RequiredPrefix string
}

func (m extractor) ApplyTo(metricName string) *match {
	// If this matcher is configured to only match part of the metric name,
	// we may need to perform some safety checks to verify we only target metrics
	// we're expecting to
	if m.RequiredPrefix != "" && !strings.HasPrefix(metricName, m.RequiredPrefix) {
		return nil
	}

	// Work out where the match appears in the string
	boundsOfMatch := m.Pattern.FindStringIndex(metricName)
	if boundsOfMatch == nil {
		return nil
	}

	// A "Match" would return the entire string that matched the pattern
	// A "Submatch" returns each of the capture groups found within the match
	submatch := m.Pattern.FindStringSubmatch(metricName)
	if submatch == nil {
		return nil
	}

	replacement := m.ReplaceWith
	tags := map[string]string{}

	// Convert all the things we matched into tags
	for i, groupName := range m.Pattern.SubexpNames() {
		// The 0th group contains the entire string the regexp matched
		if i == 0 {
			continue
		}

		tags[groupName] = submatch[i]
		replacement = strings.Replace(replacement, "{"+groupName+"}", submatch[i], -1)
	}

	// The metric name looks like:
	//
	// `foo.bar.(some.thing.we.care.about).other.stuff`
	//
	// We want to:
	// - cut out the `some.thing.we.care.about` part of the original string
	// - extract tags from it
	// - replace the part of the string we cut out with a suitable replacement
	var out strings.Builder
	out.WriteString(metricName[:boundsOfMatch[0]])
	out.WriteString(replacement)
	out.WriteString(metricName[boundsOfMatch[1]:])

	return &match{
		name: out.String(),
		tags: tags,
	}
}
