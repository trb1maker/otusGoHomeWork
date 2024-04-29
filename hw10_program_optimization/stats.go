package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	"github.com/mailru/easyjson"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	// Файл читаю построчно, память на анмаршалинг выделяю только 1 раз
	s := bufio.NewScanner(r)
	u := &User{}
	domainStat := make(DomainStat)
	for s.Scan() {
		// Заменил библиотеку, отвечающую за анмаршалинг
		if err := easyjson.Unmarshal(s.Bytes(), u); err != nil {
			return nil, err
		}
		// Функции strings.ToLower и strings.SplitN вызываю только 1 раз
		email := strings.ToLower(u.Email)
		if strings.Contains(email, domain) {
			domainStat[strings.SplitN(email, "@", 2)[1]]++
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return domainStat, nil
}
