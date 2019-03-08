package regexengine

import "strings"

type dropper struct {
	RequiredPrefix string
}

func (d dropper) ShouldDrop(m Measurement) bool {
	if !strings.HasPrefix(m.Name(), d.RequiredPrefix) {
		return false
	}

	v, err := m.Value()

	return err == nil && v == 0
}
