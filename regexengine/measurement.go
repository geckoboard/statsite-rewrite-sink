package regexengine

type Measurement interface {
	Name() string
	NameSuffix() string
	NamePrefix() string
	WholeLine() string
	MetaSuffix() string
}

type measurement struct {
	match []string
	p     *parser
}

func (m measurement) WholeLine() string {
	return m.match[0]
}

func (m measurement) Name() string {
	return m.match[m.p.namesToIndex["wildmetric"]]
}

func (m measurement) NamePrefix() string {
	return m.match[m.p.namesToIndex["wildtype"]]
}

func (m measurement) NameSuffix() string { return "" }

func (m measurement) MetaSuffix() string {
	return m.match[m.p.namesToIndex["measurements"]]
}

type timerMeasurement struct {
	measurement
}

func (t timerMeasurement) Name() string {
	return t.measurement.match[t.measurement.p.namesToIndex["timermetric"]]
}
func (t timerMeasurement) NamePrefix() string { return "timers." }

func (t timerMeasurement) NameSuffix() string {
	return t.measurement.match[t.measurement.p.namesToIndex["timerattr"]]
}
