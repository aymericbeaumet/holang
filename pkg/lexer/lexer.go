package lexer

import (
	"fmt"
	gotoken "go/token"
	"io"
	"text/scanner"
	"unicode"

	"holang/pkg/token"
)

type Token struct {
	Type     token.Token
	Position scanner.Position
	Value    string
}

func Tokenize(reader io.Reader, filepath string) ([]Token, error) {
	tokens := []Token{}

	s := scanner.Scanner{}
	s.Init(reader)
	s.Filename = filepath

	for {
		t := readToken(&s)
		switch t.Type {
		default:
			tokens = append(tokens, t)
		case gotoken.ILLEGAL:
			return tokens, fmt.Errorf("%+v", t)
		case gotoken.EOF:
			return tokens, nil
		}
	}
}

func readToken(s *scanner.Scanner) Token {
	v := ""

	for isWhitespace(s.Peek()) || isNewline(s.Peek()) {
		s.Next()
	}

	if s.Peek() == scanner.EOF {
		return Token{Type: gotoken.EOF}
	}

	// TODO: decimal, floating, imaginary

	if s.Peek() == '0' {
		v += string(s.Next())
		if s.Peek() == 'b' {
			v += string(s.Next())
			if isBinaryDigit(s.Peek()) {
				v += string(s.Next())
				for isBinaryDigit(s.Peek()) || s.Peek() == '_' {
					v += string(s.Next())
				}
				return Token{Type: gotoken.INT, Value: v}
			}
			return Token{Type: gotoken.ILLEGAL, Value: v}
		}
		if isOctalDigit(s.Peek()) {
			v += string(s.Next())
			for isOctalDigit(s.Peek()) || s.Peek() == '_' {
				v += string(s.Next())
			}
			return Token{Type: gotoken.INT, Value: v}
		}
		if s.Peek() == 'o' {
			v += string(s.Next())
			if isOctalDigit(s.Peek()) {
				v += string(s.Next())
				for isOctalDigit(s.Peek()) || s.Peek() == '_' {
					v += string(s.Next())
				}
				return Token{Type: gotoken.INT, Value: v}
			}
			return Token{Type: gotoken.ILLEGAL, Value: v}
		}
		if s.Peek() == 'x' {
			v += string(s.Next())
			if isHexDigit(s.Peek()) {
				v += string(s.Next())
				for isHexDigit(s.Peek()) || s.Peek() == '_' {
					v += string(s.Next())
				}
				return Token{Type: gotoken.INT, Value: v}
			}
			return Token{Type: gotoken.ILLEGAL, Value: v}
		}
	}

	if isLetter(s.Peek()) {
		v += string(s.Next())
		for isLetter(s.Peek()) || isUnicodeDigit(s.Peek()) {
			v += string(s.Next())
		}
		switch v {
		// holang
		case "enum":
			return Token{Type: token.ENUM, Value: v}
		case "match":
			return Token{Type: token.MATCH, Value: v}
		// golang
		case "break":
			return Token{Type: gotoken.BREAK, Value: v}
		case "case":
			return Token{Type: gotoken.CASE, Value: v}
		case "chan":
			return Token{Type: gotoken.CHAN, Value: v}
		case "const":
			return Token{Type: gotoken.CONST, Value: v}
		case "continue":
			return Token{Type: gotoken.CONTINUE, Value: v}
		case "default":
			return Token{Type: gotoken.DEFAULT, Value: v}
		case "defer":
			return Token{Type: gotoken.DEFER, Value: v}
		case "else":
			return Token{Type: gotoken.ELSE, Value: v}
		case "fallthrough":
			return Token{Type: gotoken.FALLTHROUGH, Value: v}
		case "for":
			return Token{Type: gotoken.FOR, Value: v}
		case "func":
			return Token{Type: gotoken.FUNC, Value: v}
		case "go":
			return Token{Type: gotoken.GO, Value: v}
		case "goto":
			return Token{Type: gotoken.GOTO, Value: v}
		case "if":
			return Token{Type: gotoken.IF, Value: v}
		case "import":
			return Token{Type: gotoken.IMPORT, Value: v}
		case "interface":
			return Token{Type: gotoken.INTERFACE, Value: v}
		case "map":
			return Token{Type: gotoken.MAP, Value: v}
		case "package":
			return Token{Type: gotoken.PACKAGE, Value: v}
		case "range":
			return Token{Type: gotoken.RANGE, Value: v}
		case "return":
			return Token{Type: gotoken.RETURN, Value: v}
		case "select":
			return Token{Type: gotoken.SELECT, Value: v}
		case "struct":
			return Token{Type: gotoken.STRUCT, Value: v}
		case "switch":
			return Token{Type: gotoken.SWITCH, Value: v}
		case "type":
			return Token{Type: gotoken.TYPE, Value: v}
		case "var":
			return Token{Type: gotoken.VAR, Value: v}
		default:
			return Token{Type: gotoken.IDENT, Value: v}
		}
	}

	if s.Peek() == '\'' {
		v += string(s.Next())
		if s.Peek() == '\\' {
			v += string(s.Next())
		}
		v += string(s.Next())
		if s.Peek() == '\'' {
			v += string(s.Next())
			return Token{Type: gotoken.CHAR, Value: v}
		}
		return Token{Type: gotoken.ILLEGAL, Value: v}
	}

	if s.Peek() == '"' {
		v += string(s.Next())
		escaped := false
		for !isNewline(s.Peek()) {
			cur := s.Next()
			v += string(cur)
			if escaped {
				escaped = false
			} else if cur == '\\' {
				escaped = true
			} else if cur == '"' {
				return Token{Type: gotoken.STRING, Value: v}
			}
		}
		return Token{Type: gotoken.ILLEGAL, Value: v}
	}

	if s.Peek() == '`' {
		v += string(s.Next())
		escaped := false
		for s.Peek() != scanner.EOF {
			cur := s.Next()
			v += string(cur)
			if escaped {
				escaped = false
			} else if cur == '\\' {
				escaped = true
			} else if cur == '`' {
				return Token{Type: gotoken.STRING, Value: v}
			}
		}
		return Token{Type: gotoken.ILLEGAL, Value: v}
	}

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
		if s.Peek() == '/' {
			v += string(s.Next())
			for !isNewline(s.Peek()) {
				v += string(s.Next())
			}
			return Token{Type: gotoken.COMMENT, Value: v}
		}
		if s.Peek() == '*' {
			v += string(s.Next())
			for s.Peek() != scanner.EOF {
				cur := s.Next()
				v += string(cur)
				if cur == '*' && s.Peek() == '/' {
					v += string(s.Next())
					return Token{Type: gotoken.COMMENT, Value: v}
				}
			}
			return Token{Type: gotoken.ILLEGAL, Value: v}
		}
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

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t'
}
