package sinkformatter

import (
	"fmt"
	"sort"
	"strings"
)

type Formatter func(namePrefix, name, nameSuffix string, tags map[string]string, measurement string) string

func Librato(namePrefix, name, nameSuffix string, tags map[string]string, measurement string) string {
	if len(tags) == 0 {
		return fmt.Sprintf("%s%s%s%s", namePrefix, name, nameSuffix, measurement)
	}

	return fmt.Sprintf("%s%s#%s%s%s", namePrefix, name, libratoTags(tags), nameSuffix, measurement)
}

func libratoTags(tags map[string]string) string {
	pairs := make([]string, 0, len(tags))

	for k, v := range tags {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
	}

	// go doesn't iterate over maps in a consistent, ordered way
	// We need to sort the tags to ensure stability of our tests
	sort.Strings(pairs)

	return strings.Join(pairs, ",")
}
