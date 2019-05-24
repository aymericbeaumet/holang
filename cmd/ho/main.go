package main

import (
	"errors"
	"log"
	"os"

	"holang/pkg/lexer"
	"holang/pkg/parser"
	"holang/pkg/printer"
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

	file, _, err := parser.ParseFile(tokens, 0)
	if err != nil {
		log.Fatal(err)
	}

	printer.FprintFile(os.Stdout, file)
}
