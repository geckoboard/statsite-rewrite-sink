package regexengine

type match struct {
	name string
	tags map[string]string
}

type matchers struct {
	extractors []extractor
	droppers   []dropper
}

func (m matchers) ExtractTagsFromMetric(metric Measurement) *match {
	result := &match{
		name: metric.Name(),
		tags: map[string]string{},
	}
	matched := false

	for _, matcher := range m.extractors {
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

func (m matchers) ShouldDropMetric(metric Measurement) bool {
	for _, d := range m.droppers {
		if d.ShouldDrop(metric) {
			return true
		}
	}

	return false
}
