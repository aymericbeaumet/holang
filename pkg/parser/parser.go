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

	name, err := p.eatPackageClause()
	if err != nil {
		return goast.File{}, err
	}

	_, err = p.eat(gotoken.SEMICOLON)
	if err != nil {
		return goast.File{}, err
	}

	// https://golang.org/ref/spec#Import_declarations

	imports := []*goast.ImportSpec{}

	for {
		if _, err := p.try(gotoken.IMPORT); err != nil {
			break
		}
		specs, err := p.eatImportDecl()
		if err != nil {
			return goast.File{}, err
		}
		imports = append(imports, specs...)
		if _, err = p.eat(gotoken.SEMICOLON); err != nil {
			return goast.File{}, err
		}
	}

	// https://golang.org/ref/spec#TopLevelDecl

	decls := []goast.Decl{}

	return goast.File{
		Name:    &name,
		Decls:   decls,
		Imports: imports,
	}, nil
}

func (p *Parser) eatComments() ([]lexer.Token, error) {
	comments := []lexer.Token{}
	for {
		comment, err := p.eat(gotoken.COMMENT)
		if err != nil {
			break
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (p *Parser) eatImportDecl() ([]*goast.ImportSpec, error) {
	specs := []*goast.ImportSpec{}
	if _, err := p.eat(gotoken.IMPORT); err != nil {
		return specs, err
	}
	if _, err := p.eat(gotoken.LPAREN); err == nil {
		for {
			if _, err := p.eat(gotoken.RPAREN); err == nil {
				return specs, nil
			}
			spec, err := p.eatImportSpec()
			if err != nil {
				return specs, err
			}
			specs = append(specs, spec)
			if _, err = p.eat(gotoken.SEMICOLON); err != nil {
				return specs, err
			}
		}
	}
	spec, err := p.eatImportSpec()
	if err != nil {
		return specs, err
	}
	specs = append(specs, spec)
	return specs, nil
}

func (p *Parser) eatImportSpec() (*goast.ImportSpec, error) {
	spec := goast.ImportSpec{}
	if name, err := p.try(gotoken.STRING); err == nil && name.Value == "." {
		p.next()
		spec.Name = &goast.Ident{Name: name.Value}
	} else if packageName, err := p.tryPackageName(); err == nil {
		p.next()
		spec.Name = &goast.Ident{Name: packageName.Value}
	}
	path, err := p.eatImportPath()
	if err != nil {
		return nil, err
	}
	spec.Path = &path
	return &spec, nil
}

func (p *Parser) eatImportPath() (goast.BasicLit, error) {
	path, err := p.eat(gotoken.STRING)
	if err != nil {
		return goast.BasicLit{}, err
	}
	return goast.BasicLit{
		Kind:  path.Type,
		Value: path.Value,
	}, nil
}

func (p *Parser) eatPackageClause() (goast.Ident, error) {
	_, err := p.eat(gotoken.PACKAGE)
	if err != nil {
		return goast.Ident{}, err
	}
	ident, err := p.eatPackageName()
	if err != nil {
		return goast.Ident{}, err
	}
	return goast.Ident{Name: ident.Value}, nil
}

func (p *Parser) tryPackageName() (lexer.Token, error) {
	ident, err := p.try(gotoken.IDENT)
	if err != nil {
		return lexer.Token{}, err
	}
	if ident.IsBlankIdentifier() {
		return lexer.Token{}, errors.New("The package name must not be the blank identifier")
	}
	return ident, nil
}

func (p *Parser) eatPackageName() (lexer.Token, error) {
	ident, err := p.tryPackageName()
	if err != nil {
		return lexer.Token{}, err
	}
	p.next()
	return ident, nil
}

func (p *Parser) eat(_type gotoken.Token) (lexer.Token, error) {
	token, err := p.try(_type)
	if err != nil {
		return lexer.Token{}, err
	}
	p.next()
	return token, nil
}

func (p *Parser) try(_type gotoken.Token) (lexer.Token, error) {
	token, err := p.currentToken()
	if err != nil {
		return lexer.Token{}, err
	}
	if token.Type != _type {
		return lexer.Token{}, fmt.Errorf("Expected %s token, but got %s at: %+v", _type, token.Type, token)
	}
	return token, nil
}

func (p *Parser) next() {
	p.index++
}

func (p *Parser) currentToken() (lexer.Token, error) {
	if p.index >= len(p.tokens) {
		return lexer.Token{}, errors.New("Reached end of file")
	}
	return p.tokens[p.index], nil
}
