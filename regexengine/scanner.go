package regexengine

import (
	"fmt"
	"io"

	"github.com/geckoboard/statsite-rewrite-sink/dsl"
	"github.com/geckoboard/statsite-rewrite-sink/sinkformatter"
)

func Stream(in io.Reader, out io.Writer, rules []dsl.Rule, formatter sinkformatter.Formatter) {
	parser := NewParser(in)
	matchers := CompileRulesIntoMatchers(rules)

	for parser.Parse() {
		measurement := parser.Measurement()

		if matchers.ShouldDropMetric(measurement) {
			continue
		}

		output := measurement.WholeLine()

		// Only change the format if we could extract something from it
		if match := matchers.ExtractTagsFromMetric(measurement); match != nil {
			output = formatter(
				measurement.NamePrefix(),
				match.name,
				measurement.NameSuffix(),
				match.tags,
				measurement.MetaSuffix(),
			)
		}

		fmt.Fprintln(out, output)
	}

}
