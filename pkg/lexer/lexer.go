package lexer

import (
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
		t := readToken(&s)
		switch t.Type {
		case gotoken.EOF:
			return tokens
		case gotoken.ILLEGAL:
			log.Println("warning: ", t)
			continue
		default:
			tokens = append(tokens, t)
		}
	}
}

func readToken(s *scanner.Scanner) Token {
	v := ""

	if s.Peek() == scanner.EOF {
		return Token{Type: gotoken.EOF, Value: v}
	}

	// TODO: literals

	if s.Peek() == '+' {
		v += string(s.Next())
		if s.Peek() == '=' {
			v += string(s.Next())
			return Token{Type: gotoken.ADD_ASSIGN, Value: v}
		}
		if s.Peek() == '+' {
			v += string(s.Next())
			return Token{Type: gotoken.INC, Value: v}
		}
		return Token{Type: gotoken.ADD, Value: v}
	}

	if s.Peek() == '-' {
		v += string(s.Next())
		if s.Peek() == '=' {
			v += string(s.Next())
			return Token{Type: gotoken.SUB_ASSIGN, Value: v}
		}
		if s.Peek() == '-' {
			v += string(s.Next())
			return Token{Type: gotoken.DEC, Value: v}
		}
		return Token{Type: gotoken.SUB, Value: v}
	}

	if s.Peek() == '*' {
		v += string(s.Next())
		if s.Peek() == '=' {
			v += string(s.Next())
			return Token{Type: gotoken.MUL_ASSIGN, Value: v}
		}
		if s.Peek() == '*' {
			v += string(s.Next())
			if s.Peek() == '=' {
				v += string(s.Next())
				return Token{Type: token.POW_ASSIGN, Value: v}
			}
			return Token{Type: token.POW, Value: v}
		}
		return Token{Type: gotoken.MUL, Value: v}
	}

	if s.Peek() == '/' {
		v += string(s.Next())
		if s.Peek() == '=' {
			v += string(s.Next())
			return Token{Type: gotoken.QUO_ASSIGN, Value: v}
		}
		return Token{Type: gotoken.QUO, Value: v}
	}

	if s.Peek() == '%' {
		v += string(s.Next())
		if s.Peek() == '=' {
			v += string(s.Next())
			return Token{Type: gotoken.REM_ASSIGN, Value: v}
		}
		return Token{Type: gotoken.REM, Value: v}
	}

	if s.Peek() == '&' {
		v += string(s.Next())
		if s.Peek() == '^' {
			v += string(s.Next())
			if s.Peek() == '=' {
				v += string(s.Next())
				return Token{Type: gotoken.AND_NOT_ASSIGN, Value: v}
			}
			return Token{Type: gotoken.AND_NOT, Value: v}
		}
		if s.Peek() == '=' {
			v += string(s.Next())
			return Token{Type: gotoken.AND_ASSIGN, Value: v}
		}
		if s.Peek() == '&' {
			v += string(s.Next())
			return Token{Type: gotoken.LAND, Value: v}
		}
		return Token{Type: gotoken.AND, Value: v}
	}

	if s.Peek() == '|' {
		v += string(s.Next())
		if s.Peek() == '=' {
			v += string(s.Next())
			return Token{Type: gotoken.OR_ASSIGN, Value: v}
		}
		if s.Peek() == '|' {
			v += string(s.Next())
			return Token{Type: gotoken.LOR, Value: v}
		}
		return Token{Type: gotoken.OR, Value: v}
	}

	if s.Peek() == '^' {
		v += string(s.Next())
		if s.Peek() == '=' {
			v += string(s.Next())
			return Token{Type: gotoken.XOR_ASSIGN, Value: v}
		}
		return Token{Type: gotoken.XOR, Value: v}
	}

	if s.Peek() == '<' {
		v += string(s.Next())
		if s.Peek() == '<' {
			v += string(s.Next())
			if s.Peek() == '=' {
				v += string(s.Next())
				return Token{Type: gotoken.SHL_ASSIGN, Value: v}
			}
			return Token{Type: gotoken.SHL, Value: v}
		}
		if s.Peek() == '-' {
			v += string(s.Next())
			return Token{Type: gotoken.ARROW, Value: v}
		}
		if s.Peek() == '=' {
			v += string(s.Next())
			return Token{Type: gotoken.LEQ, Value: v}
		}
		return Token{Type: gotoken.LSS, Value: v}
	}

	if s.Peek() == '>' {
		v += string(s.Next())
		if s.Peek() == '>' {
			v += string(s.Next())
			if s.Peek() == '=' {
				v += string(s.Next())
				return Token{Type: gotoken.SHR_ASSIGN, Value: v}
			}
			return Token{Type: gotoken.SHR, Value: v}
		}
		if s.Peek() == '=' {
			v += string(s.Next())
			return Token{Type: gotoken.GEQ, Value: v}
		}
		return Token{Type: gotoken.GTR, Value: v}
	}

	if s.Peek() == '=' {
		v += string(s.Next())
		if s.Peek() == '=' {
			v += string(s.Next())
			return Token{Type: gotoken.EQL, Value: v}
		}
		return Token{Type: gotoken.ASSIGN, Value: v}
	}

	if s.Peek() == '!' {
		v += string(s.Next())
		if s.Peek() == '=' {
			v += string(s.Next())
			return Token{Type: gotoken.NEQ, Value: v}
		}
		return Token{Type: gotoken.NOT, Value: v}
	}

	if s.Peek() == ':' {
		v += string(s.Next())
		if s.Peek() == '=' {
			v += string(s.Next())
			return Token{Type: gotoken.DEFINE, Value: v}
		}
		return Token{Type: gotoken.COLON, Value: v}
	}

	if s.Peek() == '.' {
		v += string(s.Next())
		if s.Peek() == '.' {
			v += string(s.Next())
			if s.Peek() == '.' {
				v += string(s.Next())
				return Token{Type: gotoken.ELLIPSIS, Value: v}
			}
			return Token{Type: gotoken.ILLEGAL, Value: v}
		}
		return Token{Type: gotoken.PERIOD, Value: v}
	}

	if s.Peek() == '(' {
		v += string(s.Next())
		return Token{Type: gotoken.LPAREN, Value: v}
	}

	if s.Peek() == ')' {
		v += string(s.Next())
		return Token{Type: gotoken.RPAREN, Value: v}
	}

	if s.Peek() == '[' {
		v += string(s.Next())
		return Token{Type: gotoken.LBRACK, Value: v}
	}

	if s.Peek() == ']' {
		v += string(s.Next())
		return Token{Type: gotoken.RBRACK, Value: v}
	}

	if s.Peek() == '{' {
		v += string(s.Next())
		return Token{Type: gotoken.LBRACE, Value: v}
	}

	if s.Peek() == '}' {
		v += string(s.Next())
		return Token{Type: gotoken.RBRACE, Value: v}
	}

	if s.Peek() == ',' {
		v += string(s.Next())
		return Token{Type: gotoken.COMMA, Value: v}
	}

	if s.Peek() == ';' {
		v += string(s.Next())
		return Token{Type: gotoken.SEMICOLON, Value: v}
	}

	// TODO: keywords

	v += string(s.Next())
	return Token{Type: gotoken.ILLEGAL, Value: v}
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
