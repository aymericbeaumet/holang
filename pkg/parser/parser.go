package parser

import (
	"errors"
	"fmt"
	goast "go/ast"
	gotokens "go/token"

	"holang/pkg/lexer"
)

// https://golang.org/ref/spec#Source_file_organization
func ParseFile(tokens []lexer.Token, i int) (goast.File, int, error) {
	i = consumeComments(tokens, i)

	// https://golang.org/ref/spec#Package_clause

	i, _, err := consume(tokens, i, gotokens.PACKAGE)
	if err != nil {
		return goast.File{}, 0, err
	}

	i, ident, err := consume(tokens, i, gotokens.IDENT)
	if err != nil {
		return goast.File{}, 0, err
	}

	i, _, err = consume(tokens, i, gotokens.SEMICOLON)
	if err != nil {
		return goast.File{}, 0, err
	}

	// https://golang.org/ref/spec#Import_declarations

	imports := []*goast.ImportSpec{}

	// https://golang.org/ref/spec#TopLevelDecl

	decls := []goast.Decl{}

	return goast.File{
		Name: &goast.Ident{
			Name: ident.Value,
		},
		Decls:   decls,
		Imports: imports,
	}, i, nil
}

func consume(tokens []lexer.Token, i int, _type gotokens.Token) (int, lexer.Token, error) {
	if i >= len(tokens) {
		return 0, lexer.Token{}, errors.New("Reached end of file")
	}
	token := tokens[i]
	if token.Type != _type {
		return 0, lexer.Token{}, fmt.Errorf("Expected %s token, but got %s at: %+v", _type, token.Type, token)
	}
	return i + 1, token, nil
}

func consumeComments(tokens []lexer.Token, i int) int {
	for i < len(tokens) && tokens[i].Type == gotokens.COMMENT {
		i++
	}
	return i
}
