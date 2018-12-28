package regexengine

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

var matchTimerMetric = `(?:(?P<timertype>timers\.)(?P<timermetric>[^|]+)(?P<timerattr>.(?:sum|sum_sq|mean|lower|upper|count|stdev|p\d\d?))`
var matchGenericMetric = `(?P<wildtype>[^\.]+\.)(?P<wildmetric>[^|]+))`
var matchMeasurement = `(?P<measurements>\|.+)`

var statsitePattern = regexp.MustCompile(fmt.Sprintf(`(?:%s|%s)%s`, matchTimerMetric, matchGenericMetric, matchMeasurement))

func NewParser(in io.Reader) *parser {
	namesToIndex := map[string]int{}

	for index, name := range statsitePattern.SubexpNames() {
		namesToIndex[name] = index
	}

	return &parser{
		scanner:      bufio.NewScanner(in),
		namesToIndex: namesToIndex,
	}
}

type parser struct {
	scanner      *bufio.Scanner
	namesToIndex map[string]int
}

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

func (p *parser) Measurement() Measurement {
	line := p.scanner.Text()

	metric := statsitePattern.FindStringSubmatch(line)

	m := measurement{metric, p}

	if metric[p.namesToIndex["typertype"]] != "" {
		return timerMeasurement{m}
	}

	return m
}

func (p *parser) Parse() bool {
	return p.scanner.Scan()
}
