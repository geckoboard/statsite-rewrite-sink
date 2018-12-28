package sinkformatter

import (
	"fmt"
	"sort"
	"strings"
)

type Formatter func(namePrefix, name, nameSuffix string, tags map[string]string, measurement string) string

func Librato(namePrefix, name, nameSuffix string, tags map[string]string, measurement string) string {
	return fmt.Sprintf("%s%s#%s%s%s", namePrefix, name, libratoTags(tags), nameSuffix, measurement)
}

func libratoTags(tags map[string]string) string {
	keys := make([]string, 0, len(tags))
	pairs := make([]string, 0, len(tags))

	// go doesn't iterate over maps in a consistent, ordered way
	// We need to sort the keys to ensure stability of our tests
	for k, _ := range tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, tags[k]))
	}

	return strings.Join(pairs, ",")
}
