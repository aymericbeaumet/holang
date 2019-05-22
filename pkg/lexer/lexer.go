package lexer

import (
	"fmt"
	"io"
	"log"
	"text/scanner"
)

type Token struct {
	Type     TokenType
	Position scanner.Position
}

type TokenType int

// Inspiration: https://golang.org/pkg/go/token/#Token
const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF
	COMMENT

	// Identifiers and basic type literals
	// (these tokens stand for classes of literals)
	IDENT  // main
	INT    // 12345
	FLOAT  // 123.45
	IMAG   // 123.45i
	CHAR   // 'a'
	STRING // "abc"

	// Operators and delimiters
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %
	POW // **

	AND     // &
	OR      // |
	XOR     // ^
	SHL     // <<
	SHR     // >>
	AND_NOT // &^

	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	QUO_ASSIGN // /=
	REM_ASSIGN // %=
	POW_ASSIGN // **=

	AND_ASSIGN     // &=
	OR_ASSIGN      // |=
	XOR_ASSIGN     // ^=
	SHL_ASSIGN     // <<=
	SHR_ASSIGN     // >>=
	AND_NOT_ASSIGN // &^=

	LAND  // &&
	LOR   // ||
	ARROW // <-
	INC   // ++
	DEC   // --

	EQL    // ==
	ASSIGN // =
	NOT    // !

	NEQ      // !=
	LEQ      // <=
	GEQ      // >=
	DEFINE   // :=
	ELLIPSIS // ...

	LPAREN // (
	LBRACK // [
	LBRACE // {
	LCHEVR // <
	COMMA  // ,
	PERIOD // .

	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	RCHEVR    // >
	SEMICOLON // ;
	COLON     // :

	// Keywords
	BREAK
	CASE
	CHAN
	CONST
	CONTINUE
	DEFAULT
	DEFER
	ELSE
	ENUM
	FALLTHROUGH
	FOR
	FUNC
	GO
	GOTO
	IF
	IMPORT
	INTERFACE
	MAP
	PACKAGE
	RANGE
	RETURN
	SELECT
	STRUCT
	SWITCH
	TYPE
	VAR
)

func Tokenize(reader io.Reader, filepath string) []Token {
	tokens := []Token{}

	s := scanner.Scanner{}
	s.Init(reader)
	s.Filename = filepath

	for {
		t := s.Scan()
		if t == scanner.EOF {
			break
		}
		var token Token
		tt := s.TokenText()
		switch tt {

		case "+":
			token = Token{Type: ADD, Position: s.Position}
		case "-":
			token = Token{Type: SUB, Position: s.Position}
		case "*":
			token = Token{Type: MUL, Position: s.Position}
		case "/":
			token = Token{Type: QUO, Position: s.Position}
		case "%":
			token = Token{Type: REM, Position: s.Position}
		case "**":
			token = Token{Type: POW, Position: s.Position}

		case "&":
			token = Token{Type: AND, Position: s.Position}
		case "|":
			token = Token{Type: OR, Position: s.Position}
		case "^":
			token = Token{Type: XOR, Position: s.Position}
		case "<<":
			token = Token{Type: SHL, Position: s.Position}
		case ">>":
			token = Token{Type: SHR, Position: s.Position}
		case "&^":
			token = Token{Type: AND_NOT, Position: s.Position}

		case "+=":
			token = Token{Type: ADD_ASSIGN, Position: s.Position}
		case "-=":
			token = Token{Type: SUB_ASSIGN, Position: s.Position}
		case "*=":
			token = Token{Type: MUL_ASSIGN, Position: s.Position}
		case "/=":
			token = Token{Type: QUO_ASSIGN, Position: s.Position}
		case "%=":
			token = Token{Type: REM_ASSIGN, Position: s.Position}
		case "**=":
			token = Token{Type: POW_ASSIGN, Position: s.Position}

		case "&=":
			token = Token{Type: AND_ASSIGN, Position: s.Position}
		case "|=":
			token = Token{Type: OR_ASSIGN, Position: s.Position}
		case "^=":
			token = Token{Type: XOR_ASSIGN, Position: s.Position}
		case "<<=":
			token = Token{Type: SHL_ASSIGN, Position: s.Position}
		case ">>=":
			token = Token{Type: SHR_ASSIGN, Position: s.Position}
		case "&^=":
			token = Token{Type: AND_NOT_ASSIGN, Position: s.Position}

		case "&&":
			token = Token{Type: LAND, Position: s.Position}
		case "||":
			token = Token{Type: LOR, Position: s.Position}
		case "<-":
			token = Token{Type: ARROW, Position: s.Position}
		case "++":
			token = Token{Type: INC, Position: s.Position}
		case "--":
			token = Token{Type: DEC, Position: s.Position}

		case "==":
			token = Token{Type: EQL, Position: s.Position}
		case "=":
			token = Token{Type: ASSIGN, Position: s.Position}
		case "!":
			token = Token{Type: NOT, Position: s.Position}

		case "!=":
			token = Token{Type: NEQ, Position: s.Position}
		case "<=":
			token = Token{Type: LEQ, Position: s.Position}
		case ">=":
			token = Token{Type: GEQ, Position: s.Position}
		case ":=":
			token = Token{Type: DEFINE, Position: s.Position}
		case "...":
			token = Token{Type: ELLIPSIS, Position: s.Position}

		case "(":
			token = Token{Type: LPAREN, Position: s.Position}
		case "[":
			token = Token{Type: LBRACK, Position: s.Position}
		case "{":
			token = Token{Type: LBRACE, Position: s.Position}
		case "<":
			token = Token{Type: LCHEVR, Position: s.Position}
		case ",":
			token = Token{Type: COMMA, Position: s.Position}
		case ".":
			token = Token{Type: PERIOD, Position: s.Position}

		case ")":
			token = Token{Type: RPAREN, Position: s.Position}
		case "]":
			token = Token{Type: RBRACK, Position: s.Position}
		case "}":
			token = Token{Type: RBRACE, Position: s.Position}
		case ">":
			token = Token{Type: RCHEVR, Position: s.Position}
		case ";":
			token = Token{Type: SEMICOLON, Position: s.Position}
		case ":":
			token = Token{Type: COLON, Position: s.Position}

		case "break":
			token = Token{Type: BREAK, Position: s.Position}
		case "case":
			token = Token{Type: CASE, Position: s.Position}
		case "chan":
			token = Token{Type: CHAN, Position: s.Position}
		case "const":
			token = Token{Type: CONST, Position: s.Position}
		case "continue":
			token = Token{Type: CONTINUE, Position: s.Position}
		case "default":
			token = Token{Type: DEFAULT, Position: s.Position}
		case "defer":
			token = Token{Type: DEFER, Position: s.Position}
		case "else":
			token = Token{Type: ELSE, Position: s.Position}
		case "enum":
			token = Token{Type: ENUM, Position: s.Position}
		case "fallthrough":
			token = Token{Type: FALLTHROUGH, Position: s.Position}
		case "for":
			token = Token{Type: FOR, Position: s.Position}
		case "func":
			token = Token{Type: FUNC, Position: s.Position}
		case "go":
			token = Token{Type: GO, Position: s.Position}
		case "goto":
			token = Token{Type: GOTO, Position: s.Position}
		case "if":
			token = Token{Type: IF, Position: s.Position}
		case "import":
			token = Token{Type: IMPORT, Position: s.Position}
		case "interface":
			token = Token{Type: INTERFACE, Position: s.Position}
		case "map":
			token = Token{Type: MAP, Position: s.Position}
		case "package":
			token = Token{Type: PACKAGE, Position: s.Position}
		case "range":
			token = Token{Type: RANGE, Position: s.Position}
		case "return":
			token = Token{Type: RETURN, Position: s.Position}
		case "select":
			token = Token{Type: SELECT, Position: s.Position}
		case "struct":
			token = Token{Type: STRUCT, Position: s.Position}
		case "switch":
			token = Token{Type: SWITCH, Position: s.Position}
		case "type":
			token = Token{Type: TYPE, Position: s.Position}
		case "var":
			token = Token{Type: VAR, Position: s.Position}

		default:
			log.Print(fmt.Errorf("Unsupported at %s: %s\n", s.Position, s.TokenText()))
		}

		tokens = append(tokens, token)
	}

	return tokens
}
