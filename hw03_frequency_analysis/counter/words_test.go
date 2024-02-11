package counter

import (
	"sort"
	"testing"
)

func TestWordsLess(t *testing.T) {
	ww := words{
		word{name: "а", count: 2},
		word{name: "и", count: 2},
	}
	sort.Sort(sort.Reverse(ww))
	if ww[0].name != "а" {
		t.Error()
	}
}
