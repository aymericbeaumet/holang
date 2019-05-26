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

	_, i, err := consume(tokens, i, gotokens.PACKAGE)
	if err != nil {
		return goast.File{}, 0, err
	}

	ident, i, err := consume(tokens, i, gotokens.IDENT)
	if err != nil {
		return goast.File{}, 0, err
	}

	if isBlankIdentifier(ident) {
		return goast.File{}, 0, errors.New("The package name must not be the blank identifier")
	}

	_, i, err = consume(tokens, i, gotokens.SEMICOLON)
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

func consume(tokens []lexer.Token, i int, _type gotokens.Token) (lexer.Token, int, error) {
	if i >= len(tokens) {
		return lexer.Token{}, 0, errors.New("Reached end of file")
	}
	token := tokens[i]
	if token.Type != _type {
		return lexer.Token{}, 0, fmt.Errorf("Expected %s token, but got %s at: %+v", _type, token.Type, token)
	}
	return token, i + 1, nil
}

func consumeComments(tokens []lexer.Token, i int) int {
	for i < len(tokens) && tokens[i].Type == gotokens.COMMENT {
		i++
	}
	return i
}

func isBlankIdentifier(t lexer.Token) bool {
	return t.Value == "_"
}
