package counter

type word struct {
	name  string
	count uint
}

type words []word

func (w words) Len() int {
	return len(w)
}

func (w words) Less(i, j int) bool {
	if w[i].count < w[j].count {
		return true
	}
	// Чтобы получить лексикографическую сортировку нарушаю описание метода Less,
	// возможно это не лучший вариант и стоит в дальнейшем вынести логику лексикографической сортировки из этого метода
	if w[i].count == w[j].count {
		return w[i].name > w[j].name
	}
	return false
}

func (w words) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}
