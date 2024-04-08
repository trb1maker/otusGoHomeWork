package main

import (
	"flag"
	"log"
)

var useMyProgressBar bool

func main() {
	var (
		from, to      string
		offset, limit int64
	)

	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.BoolVar(&useMyProgressBar, "experimental", false, "use my progress bar")
	flag.Parse()

	if err := Copy(from, to, offset, limit); err != nil {
		log.Println(err)
	}
}
