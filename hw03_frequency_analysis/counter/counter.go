package counter

import (
	"sort"
	"strings"
	"unicode"
)

type Counter map[string]uint

// Модифицирует слово, приводя его к нижнему регистру и удаляя знаки препинания
func modifyWord(r rune) rune {
	if unicode.IsPunct(r) {
		return -1
	}
	return unicode.ToLower(r)
}

// Добавляет слово в счетчик
func (c Counter) Add(s string) {
	s = strings.Map(modifyWord, s)
	if s == "" {
		return
	}
	c[s]++
}

// Возвращает топ 10 слов в нужном порядке
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
