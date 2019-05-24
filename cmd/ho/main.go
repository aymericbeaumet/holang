package main

import (
	"errors"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"

	"holang/pkg/lexer"
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

	tokens, err := lexer.Tokenize(reader, filepath)
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(tokens)
}
