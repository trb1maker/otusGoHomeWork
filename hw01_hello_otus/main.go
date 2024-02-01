package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	const hello = "Hello, OTUS!"
	fmt.Print(reverse.String(hello))
}
