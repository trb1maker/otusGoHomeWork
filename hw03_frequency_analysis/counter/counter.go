package counter

import "sort"

type Counter map[string]uint

func (c Counter) GetTop() []string {
	ww := make(words, len(c))
	for w, c := range c {
		ww = append(ww, word{name: w, count: c})
	}
	sort.Sort(sort.Reverse(ww))
	result := make([]string, 0)
	for i, w := range ww {
		if i > 9 {
			break
		}
		result = append(result, w.name)
	}
	return result
}
