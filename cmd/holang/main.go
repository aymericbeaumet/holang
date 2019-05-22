package main

import (
	"errors"
	"holang/pkg/lexer"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		panic(errors.New("Usage: holang <filepath>"))
	}
	filepath := os.Args[1]
	reader, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	lexer.Tokenize(reader, filepath)
}
