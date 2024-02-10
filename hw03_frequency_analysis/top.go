package hw03frequencyanalysis

import (
	"strings"

	"github.com/trb1maker/otus_golang_home_work/hw03_frequency_analysis/counter"
)

func Top10(s string) []string {
	c := make(counter.Counter)
	for _, w := range strings.Fields(s) {
		c[w]++
	}
	return c.GetTop()
}
