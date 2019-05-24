package parser

import (
	"fmt"
	goast "go/ast"
	gotokens "go/token"

	"holang/pkg/lexer"
)

func ParseFile(tokens []lexer.Token, i int) (goast.File, error) {
	i = exhaustComments(tokens, i)

	pkg := tokens[i]
	if pkg.Type != gotokens.PACKAGE {
		return goast.File{}, fmt.Errorf("Expected PACKAGE token, got: %+v", pkg)
	}

	i++
	ident := tokens[i]
	if ident.Type != gotokens.IDENT {
		return goast.File{}, fmt.Errorf("Expected IDENT token, got: %+v", ident)
	}

	return goast.File{
		Name: &goast.Ident{
			Name: ident.Value,
		},
	}, nil
}

func exhaustComments(tokens []lexer.Token, i int) int {
	for i < len(tokens) && tokens[i].Type == gotokens.COMMENT {
		i++
	}
	return i
}
