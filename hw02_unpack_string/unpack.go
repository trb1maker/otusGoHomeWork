package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

const backSlash rune = 92

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var (
		result    = new(strings.Builder)
		items     = []rune(s)
		lastIndex = len(items) - 1
	)
	for i := 0; i < len(items); {
		// Текущий символ не может быть числом
		item := items[i]
		if unicode.IsDigit(item) {
			return "", ErrInvalidString
		}
		// Если текущий символ бэкслэш, нужно сдвинуться на следующий символ
		if item == backSlash {
			i++
			item = items[i]
		}
		// Последний символ просто записываю в результат
		if i == lastIndex {
			result.WriteRune(item)
			break
		}
		// Если ошибок нет, то проверяю, является ли следующий символ числом
		next := items[i+1]
		if unicode.IsDigit(next) {
			count, err := strconv.Atoi(string(next))
			if err != nil {
				panic(err)
			}
			result.WriteString(strings.Repeat(string(item), count))
			// Если следующий символ число, то выполнив распаковку я пропускаю уже обработанный символ
			i += 2
			continue
		}
		// Если следующий символ буква, то сохраняю текущий символ и перехожу к обработке следующего
		result.WriteRune(item)
		i++
	}
	return result.String(), nil
}
