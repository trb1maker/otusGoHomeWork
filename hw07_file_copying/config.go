// Конфигурацию перенес в отдельную структуру.
// Парсинг аргументов командной строки перенес в отдельную функцию, чтобы не использовать init и глобальные переменные.
package main

import (
	"flag"
)

type config struct {
	from   string
	to     string
	limit  int64
	offset int64
}

func parseConfig() *config {
	cfg := &config{}

	flag.StringVar(&cfg.from, "from", "", "file to read from")
	flag.StringVar(&cfg.to, "to", "", "file to write to")
	flag.Int64Var(&cfg.limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&cfg.offset, "offset", 0, "offset in input file")
	flag.Parse()

	return cfg
}
