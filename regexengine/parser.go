package regexengine

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

// timer are metrics are special because statsite emits ~10 aggregates
// for each timer. Aggregates have a special suffix
var matchTimerMetric = `(?:(?P<timertype>timers\.)(?P<timermetric>[^|]+)(?P<timerattr>.(?:sum|sum_sq|mean|lower|upper|count|stdev|p\d\d?))`
var matchGenericMetric = `(?P<wildtype>[^\.]+\.)(?P<wildmetric>[^|]+))`
var matchMeasurement = `(?P<measurements>\|.+)`

var statsitePattern = regexp.MustCompile(fmt.Sprintf(`(?:%s|%s)%s`, matchTimerMetric, matchGenericMetric, matchMeasurement))

func NewParser(in io.Reader) *parser {
	namesToIndex := map[string]int{}

	for index, name := range statsitePattern.SubexpNames() {
		namesToIndex[name] = index
	}

	timerTypeIndex, ok := namesToIndex["timertype"]
	if !ok {
		panic("could not find timer type index")
	}

	return &parser{
		scanner:        bufio.NewScanner(in),
		namesToIndex:   namesToIndex,
		timerTypeIndex: timerTypeIndex,
	}
}

type parser struct {
	scanner        *bufio.Scanner
	namesToIndex   map[string]int
	timerTypeIndex int
}

func (p *parser) Measurement() Measurement {
	line := p.scanner.Text()

	metric := statsitePattern.FindStringSubmatch(line)

	m := measurement{metric, p}

	if metric[p.timerTypeIndex] != "" {
		return timerMeasurement{m}
	}

	return m
}

func (p *parser) Parse() bool {
	return p.scanner.Scan()
}
