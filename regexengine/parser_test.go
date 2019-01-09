package regexengine

import (
	"strings"
	"testing"
)

func TestMeasurementParser(t *testing.T) {
	type expectedMeasurement struct {
		Name       string
		Prefix     string
		Suffix     string
		MetaSuffix string
	}
	type example struct {
		Input    string
		Expected []expectedMeasurement
	}

	examples := []example{
		{
			Input: "counts.consul.health.service.query.pegasus|1.000000|1545077605",
			Expected: []expectedMeasurement{
				{
					Name:       "consul.health.service.query.pegasus",
					Prefix:     "counts.",
					MetaSuffix: "|1.000000|1545077605",
				},
			},
		},
		{
			Input: "timers.envoy.cluster.tuner-grpc.upstream_rq_time.count|182|154695525",
			Expected: []expectedMeasurement{
				{
					Name:       "envoy.cluster.tuner-grpc.upstream_rq_time",
					Prefix:     "timers.",
					Suffix:     ".count",
					MetaSuffix: "|182|154695525",
				},
			},
		},
	}

	for _, ex := range examples {
		in := strings.NewReader(ex.Input)
		p := NewParser(in)

		measurements := []Measurement{}
		for p.Parse() {
			measurements = append(measurements, p.Measurement())
		}

		for i, expectation := range ex.Expected {
			m := measurements[i]

			if m.Name() != expectation.Name {
				t.Errorf("wanted %q, got %q", expectation.Name, m.Name())
			}

			if m.NamePrefix() != expectation.Prefix {
				t.Errorf("wanted %q, got %q", expectation.Prefix, m.NamePrefix())
			}

			if m.NameSuffix() != expectation.Suffix {
				t.Errorf("wanted %q, got %q", expectation.Suffix, m.NameSuffix())
			}

			if m.MetaSuffix() != expectation.MetaSuffix {
				t.Errorf("wanted %q, got %q", expectation.MetaSuffix, m.MetaSuffix())
			}
		}

	}
}
