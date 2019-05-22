package main

import (
	"errors"
	"log"
	"os"

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

	tokens := lexer.Tokenize(reader, filepath)
	log.Printf("%+v", tokens)
}
