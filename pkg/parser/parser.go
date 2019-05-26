package parser

import (
	"errors"
	"fmt"
	goast "go/ast"
	gotoken "go/token"

	"holang/pkg/lexer"
)

type Parser struct {
	tokens []lexer.Token
	index  int
}

func NewParser(tokens []lexer.Token) Parser {
	return Parser{
		tokens: tokens,
		index:  0,
	}
}

// https://golang.org/ref/spec#Source_file_organization
func (p *Parser) ParseFile() (goast.File, error) {
	p.eatComments()

	// https://golang.org/ref/spec#Package_clause

	_, err := p.eat(gotoken.PACKAGE)
	if err != nil {
		return goast.File{}, err
	}

	ident, err := p.eat(gotoken.IDENT)
	if err != nil {
		return goast.File{}, err
	}

	if ident.IsBlankIdentifier() {
		return goast.File{}, errors.New("The package name must not be the blank identifier")
	}

	_, err = p.eat(gotoken.SEMICOLON)
	if err != nil {
		return goast.File{}, err
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
	}, nil
}

func (p *Parser) eatComments() ([]lexer.Token, error) {
	comments := []lexer.Token{}
	for {
		comment, err := p.eat(gotoken.COMMENT)
		if err != nil {
			return comments, nil
		}
		comments = append(comments, comment)
	}
}

func (p *Parser) eat(_type gotoken.Token) (lexer.Token, error) {
	token, err := p.currentToken()
	if err != nil {
		return lexer.Token{}, err
	}
	if token.Type != _type {
		return lexer.Token{}, fmt.Errorf("Expected %s token, but got %s at: %+v", _type, token.Type, token)
	}
	p.index++
	return token, nil
}

func (p *Parser) currentToken() (lexer.Token, error) {
	if p.index >= len(p.tokens) {
		return lexer.Token{}, errors.New("Reached end of file")
	}
	return p.tokens[p.index], nil
}
