package lexer

import (
	"fmt"
	gotoken "go/token"
	"io"
	"log"
	"text/scanner"

	"holang/pkg/token"
)

type Token struct {
	Type     token.Token
	Position scanner.Position
}

func Tokenize(reader io.Reader, filepath string) []Token {
	tokens := []Token{}

	s := scanner.Scanner{}
	s.Init(reader)
	s.Filename = filepath

	for {
		lexeme := s.Scan()
		if lexeme == scanner.EOF {
			break
		}
		var t Token
		tt := s.TokenText()
		switch tt {

		// Holang

		case "**":
			t = Token{Type: token.POW, Position: s.Position}
		case "**=":
			t = Token{Type: token.POW_ASSIGN, Position: s.Position}

		case "<":
			t = Token{Type: token.LCHEVR, Position: s.Position}
		case ">":
			t = Token{Type: token.RCHEVR, Position: s.Position}

		case "enum":
			t = Token{Type: token.ENUM, Position: s.Position}
		case "match":
			t = Token{Type: token.MATCH, Position: s.Position}

		// Golang

		case "+":
			t = Token{Type: gotoken.ADD, Position: s.Position}
		case "-":
			t = Token{Type: gotoken.SUB, Position: s.Position}
		case "*":
			t = Token{Type: gotoken.MUL, Position: s.Position}
		case "/":
			t = Token{Type: gotoken.QUO, Position: s.Position}
		case "%":
			t = Token{Type: gotoken.REM, Position: s.Position}

		case "&":
			t = Token{Type: gotoken.AND, Position: s.Position}
		case "|":
			t = Token{Type: gotoken.OR, Position: s.Position}
		case "^":
			t = Token{Type: gotoken.XOR, Position: s.Position}
		case "<<":
			t = Token{Type: gotoken.SHL, Position: s.Position}
		case ">>":
			t = Token{Type: gotoken.SHR, Position: s.Position}
		case "&^":
			t = Token{Type: gotoken.AND_NOT, Position: s.Position}

		case "+=":
			t = Token{Type: gotoken.ADD_ASSIGN, Position: s.Position}
		case "-=":
			t = Token{Type: gotoken.SUB_ASSIGN, Position: s.Position}
		case "*=":
			t = Token{Type: gotoken.MUL_ASSIGN, Position: s.Position}
		case "/=":
			t = Token{Type: gotoken.QUO_ASSIGN, Position: s.Position}
		case "%=":
			t = Token{Type: gotoken.REM_ASSIGN, Position: s.Position}

		case "&=":
			t = Token{Type: gotoken.AND_ASSIGN, Position: s.Position}
		case "|=":
			t = Token{Type: gotoken.OR_ASSIGN, Position: s.Position}
		case "^=":
			t = Token{Type: gotoken.XOR_ASSIGN, Position: s.Position}
		case "<<=":
			t = Token{Type: gotoken.SHL_ASSIGN, Position: s.Position}
		case ">>=":
			t = Token{Type: gotoken.SHR_ASSIGN, Position: s.Position}
		case "&^=":
			t = Token{Type: gotoken.AND_NOT_ASSIGN, Position: s.Position}

		case "&&":
			t = Token{Type: gotoken.LAND, Position: s.Position}
		case "||":
			t = Token{Type: gotoken.LOR, Position: s.Position}
		case "<-":
			t = Token{Type: gotoken.ARROW, Position: s.Position}
		case "++":
			t = Token{Type: gotoken.INC, Position: s.Position}
		case "--":
			t = Token{Type: gotoken.DEC, Position: s.Position}

		case "==":
			t = Token{Type: gotoken.EQL, Position: s.Position}
		case "=":
			t = Token{Type: gotoken.ASSIGN, Position: s.Position}
		case "!":
			t = Token{Type: gotoken.NOT, Position: s.Position}

		case "!=":
			t = Token{Type: gotoken.NEQ, Position: s.Position}
		case "<=":
			t = Token{Type: gotoken.LEQ, Position: s.Position}
		case ">=":
			t = Token{Type: gotoken.GEQ, Position: s.Position}
		case ":=":
			t = Token{Type: gotoken.DEFINE, Position: s.Position}
		case "...":
			t = Token{Type: gotoken.ELLIPSIS, Position: s.Position}

		case "(":
			t = Token{Type: gotoken.LPAREN, Position: s.Position}
		case "[":
			t = Token{Type: gotoken.LBRACK, Position: s.Position}
		case "{":
			t = Token{Type: gotoken.LBRACE, Position: s.Position}
		case ",":
			t = Token{Type: gotoken.COMMA, Position: s.Position}
		case ".":
			t = Token{Type: gotoken.PERIOD, Position: s.Position}

		case ")":
			t = Token{Type: gotoken.RPAREN, Position: s.Position}
		case "]":
			t = Token{Type: gotoken.RBRACK, Position: s.Position}
		case "}":
			t = Token{Type: gotoken.RBRACE, Position: s.Position}
		case ";":
			t = Token{Type: gotoken.SEMICOLON, Position: s.Position}
		case ":":
			t = Token{Type: gotoken.COLON, Position: s.Position}

		case "break":
			t = Token{Type: gotoken.BREAK, Position: s.Position}
		case "case":
			t = Token{Type: gotoken.CASE, Position: s.Position}
		case "chan":
			t = Token{Type: gotoken.CHAN, Position: s.Position}
		case "const":
			t = Token{Type: gotoken.CONST, Position: s.Position}
		case "continue":
			t = Token{Type: gotoken.CONTINUE, Position: s.Position}
		case "default":
			t = Token{Type: gotoken.DEFAULT, Position: s.Position}
		case "defer":
			t = Token{Type: gotoken.DEFER, Position: s.Position}
		case "else":
			t = Token{Type: gotoken.ELSE, Position: s.Position}
		case "fallthrough":
			t = Token{Type: gotoken.FALLTHROUGH, Position: s.Position}
		case "for":
			t = Token{Type: gotoken.FOR, Position: s.Position}
		case "func":
			t = Token{Type: gotoken.FUNC, Position: s.Position}
		case "go":
			t = Token{Type: gotoken.GO, Position: s.Position}
		case "goto":
			t = Token{Type: gotoken.GOTO, Position: s.Position}
		case "if":
			t = Token{Type: gotoken.IF, Position: s.Position}
		case "import":
			t = Token{Type: gotoken.IMPORT, Position: s.Position}
		case "interface":
			t = Token{Type: gotoken.INTERFACE, Position: s.Position}
		case "map":
			t = Token{Type: gotoken.MAP, Position: s.Position}
		case "package":
			t = Token{Type: gotoken.PACKAGE, Position: s.Position}
		case "range":
			t = Token{Type: gotoken.RANGE, Position: s.Position}
		case "return":
			t = Token{Type: gotoken.RETURN, Position: s.Position}
		case "select":
			t = Token{Type: gotoken.SELECT, Position: s.Position}
		case "struct":
			t = Token{Type: gotoken.STRUCT, Position: s.Position}
		case "switch":
			t = Token{Type: gotoken.SWITCH, Position: s.Position}
		case "type":
			t = Token{Type: gotoken.TYPE, Position: s.Position}
		case "var":
			t = Token{Type: gotoken.VAR, Position: s.Position}

		default:
			log.Print(fmt.Errorf("Unsupported at %s: %s\n", s.Position, s.TokenText()))
		}

		tokens = append(tokens, t)
	}

	return tokens
}
