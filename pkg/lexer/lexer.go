package lexer

import (
	"fmt"
	gotoken "go/token"
	"io"
	"log"
	"text/scanner"
	"unicode"

	"holang/pkg/token"
)

type Token struct {
	Type     token.Token
	Position scanner.Position
	Value    string
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
		case "<":
			t = Token{Type: gotoken.LSS, Position: s.Position}
		case ">":
			t = Token{Type: gotoken.GTR, Position: s.Position}
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
			switch {
			case isFloat(tt):
				t = Token{Type: gotoken.FLOAT, Position: s.Position}
			case isIdent(tt):
				t = Token{Type: gotoken.IDENT, Position: s.Position}
			case isImag(tt):
				t = Token{Type: gotoken.IMAG, Position: s.Position}
			case isIntLit(tt):
				t = Token{Type: gotoken.INT, Position: s.Position}
			case isRuneLit(tt):
				t = Token{Type: gotoken.CHAR, Position: s.Position}
			case isStringLit(tt):
				t = Token{Type: gotoken.STRING, Position: s.Position}
			default:
				log.Print(fmt.Errorf("%s: %s\n", s.Position, s.TokenText()))
			}
		}

		t.Value = tt
		tokens = append(tokens, t)
	}

	return tokens
}

func isFloat(s string) bool {
	return false
}

func isImag(s string) bool {
	return false
}

func isIdent(s string) bool {
	for i, r := range s {
		switch i {
		case 0:
			if !isLetter(r) {
				return false
			}
		default:
			if !(isLetter(r) || isUnicodeDigit(r)) {
				return false
			}
		}
	}
	return true
}

func isIntLit(s string) bool {
	return isBinaryLit(s) || isOctalLit(s) || isDecimalLit(s) || isHexLit(s)
}

func isBinaryLit(s string) bool {
	valid := false
	for i, r := range s {
		switch i {
		case 0:
			if !(r == '0') {
				return false
			}
		case 1:
			if !(r == 'b' || r == 'B') {
				return false
			}
		default:
			if !isBinaryDigit(r) {
				return false
			}
			valid = true
		}
	}
	return valid
}

func isOctalLit(s string) bool {
	valid := false
	for i, r := range s {
		switch i {
		case 0:
			if !(r == '0') {
				return false
			}
			valid = true
		default:
			if !isOctalDigit(r) {
				return false
			}
		}
	}
	return valid
}

func isDecimalLit(s string) bool {
	valid := false
	for i, r := range s {
		switch i {
		case 0:
			if !(r >= '1' && r <= '9') {
				return false
			}
			valid = true
		default:
			if !isDecimalDigit(r) {
				return false
			}
		}
	}
	return valid
}

func isHexLit(s string) bool {
	valid := false
	for i, r := range s {
		switch i {
		case 0:
			if !(r == '0') {
				return false
			}
		case 1:
			if !(r == 'x' || r == 'X') {
				return false
			}
		default:
			if !isOctalDigit(r) {
				return false
			}
			valid = true
		}
	}
	return valid
}

func isStringLit(s string) bool {
	return isRawStringLit(s) || isInterpretedStringLit(s)
}

func isRawStringLit(s string) bool {
	terminated := false
	for i, r := range s {
		switch i {
		case 0:
			if !(r == '`') {
				return false
			}
		default:
			if r == '`' {
				if terminated {
					return false
				}
				terminated = true
			} else if !(isUnicodeChar(r) || isNewline(r)) {
				return false
			}
		}
	}
	return terminated
}

func isInterpretedStringLit(s string) bool {
	terminated := false
	for i, r := range s {
		switch i {
		case 0:
			if !(r == '"') {
				return false
			}
		default:
			if r == '"' {
				if terminated {
					return false
				}
				terminated = true
			} else if !(isUnicodeValue(r) || isByteValue(r)) {
				return false
			}
		}
	}
	return terminated
}

func isRuneLit(s string) bool {
	terminated := false
	for i, r := range s {
		switch i {
		case 0:
			if !(r == '\'') {
				return false
			}
		default:
			if r == '\'' {
				if terminated {
					return false
				}
				terminated = true
			} else if !(isUnicodeValue(r) || isByteValue(r)) {
				return false
			}
		}
	}
	return terminated
}

// TODO
func isUnicodeValue(r rune) bool {
	return true
}

// TODO
func isByteValue(r rune) bool {
	return true
}

func isLetter(r rune) bool {
	return isUnicodeLetter(r) || r == '_'
}

func isBinaryDigit(r rune) bool {
	return r >= '0' && r <= '1'
}

func isOctalDigit(r rune) bool {
	return r >= '0' && r <= '7'
}

func isDecimalDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isHexDigit(r rune) bool {
	return isDecimalDigit(r) || (r >= 'A' && r <= 'F') || (r >= 'a' && r <= 'f')
}

func isUnicodeChar(r rune) bool {
	return !isNewline(r)
}

func isUnicodeDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func isUnicodeLetter(r rune) bool {
	return unicode.IsLetter(r)
}

func isNewline(r rune) bool {
	return r == '\n'
}
