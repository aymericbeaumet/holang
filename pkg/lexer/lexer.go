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

func (t *Token) IsBlankIdentifier() bool {
	return t.Type == gotoken.IDENT && t.Value == "_"
}

func Tokenize(reader io.Reader, filepath string) ([]Token, error) {
	tokens := []Token{}

	s := scanner.Scanner{}
	s.Init(reader)
	s.Filename = filepath

	asi := false // automatic semicolon insertion on upcoming newline?
	for {
		t := readToken(&s, asi)
		asi = false
		switch t.Type {
		case gotoken.IDENT:
			fallthrough
		case gotoken.INT:
			fallthrough
		case gotoken.FLOAT:
			fallthrough
		case gotoken.IMAG:
			fallthrough
		case gotoken.CHAR:
			fallthrough
		case gotoken.STRING:
			fallthrough
		case gotoken.BREAK:
			fallthrough
		case gotoken.CONTINUE:
			fallthrough
		case gotoken.FALLTHROUGH:
			fallthrough
		case gotoken.RETURN:
			fallthrough
		case gotoken.INC:
			fallthrough
		case gotoken.DEC:
			fallthrough
		case gotoken.RPAREN:
			fallthrough
		case gotoken.RBRACK:
			fallthrough
		case gotoken.RBRACE:
			asi = true
			fallthrough
		default:
			tokens = append(tokens, t)
		case gotoken.ILLEGAL:
			return tokens, fmt.Errorf("%+v", t)
		case gotoken.EOF:
			return tokens, nil
		}
	}
}

func readToken(s *scanner.Scanner, asi bool) Token {
	v := []rune{}

	for isWhitespace(s.Peek()) || isNewline(s.Peek()) {
		cur := s.Next()
		if asi && isNewline(cur) {
			return Token{Type: gotoken.SEMICOLON}
		}
	}

	if s.Peek() == scanner.EOF {
		return Token{Type: gotoken.EOF}
	}

	if isDecimalDigit(s.Peek()) || s.Peek() == '.' {
		cur := s.Next()
		v = append(v, cur)
		// Binary notation 0b...
		if cur == '0' && (s.Peek() == 'B' || s.Peek() == 'b') {
			v = append(v, s.Next())
			if isBinaryDigit(s.Peek()) {
				v = append(v, s.Next())
				for isBinaryDigit(s.Peek()) || isDigitSeparator(s.Peek()) {
					v = append(v, s.Next())
				}
				return Token{Type: gotoken.INT, Value: string(v)}
			}
			return Token{Type: gotoken.ILLEGAL, Value: string(v)}
		}
		// Octal notation 0o...
		if cur == '0' && (s.Peek() == 'O' || s.Peek() == 'o') {
			v = append(v, s.Next())
			if isOctalDigit(s.Peek()) {
				v = append(v, s.Next())
				for isOctalDigit(s.Peek()) || isDigitSeparator(s.Peek()) {
					v = append(v, s.Next())
				}
				return Token{Type: gotoken.INT, Value: string(v)}
			}
			return Token{Type: gotoken.ILLEGAL, Value: string(v)}
		}
		// Hexadecimal notation 0x...
		if cur == '0' && (s.Peek() == 'X' || s.Peek() == 'x') {
			v = append(v, s.Next())
			if isHexDigit(s.Peek()) {
				v = append(v, s.Next())
				for isHexDigit(s.Peek()) || isDigitSeparator(s.Peek()) {
					v = append(v, s.Next())
				}
				return Token{Type: gotoken.INT, Value: string(v)}
			}
			return Token{Type: gotoken.ILLEGAL, Value: string(v)}
		}
		// Alternative octal notation 0...
		if cur == '0' {
			if isOctalDigit(s.Peek()) {
				v = append(v, s.Next())
				for isOctalDigit(s.Peek()) || isDigitSeparator(s.Peek()) {
					v = append(v, s.Next())
				}
			}
			if !(isDecimalDigit(s.Peek()) || isDigitSeparator(s.Peek()) || s.Peek() == '.' || s.Peek() == 'i' || s.Peek() == 'E' || s.Peek() == 'e') {
				return Token{Type: gotoken.INT, Value: string(v)}
			}
			// fallthrough
		}
		// Decimals, Floating-point, Imaginary
		float := false
		if isDecimalDigit(v[0]) {
			for isDecimalDigit(s.Peek()) || isDigitSeparator(s.Peek()) {
				v = append(v, s.Next())
			}
		}
		if (v[0] == '.' || s.Peek() == '.') && v[0] != s.Peek() {
			float = true
			if s.Peek() == '.' {
				v = append(v, s.Next())
			}
			if isDecimalDigit(s.Peek()) {
				v = append(v, s.Next())
				for isDecimalDigit(s.Peek()) || isDigitSeparator(s.Peek()) {
					v = append(v, s.Next())
				}
			}
		}
		if s.Peek() == 'E' || s.Peek() == 'e' {
			float = true
			v = append(v, s.Next())
			if s.Peek() == '+' || s.Peek() == '-' {
				v = append(v, s.Next())
			}
			if isDecimalDigit(s.Peek()) {
				v = append(v, s.Next())
				for isDecimalDigit(s.Peek()) || isDigitSeparator(s.Peek()) {
					v = append(v, s.Next())
				}
			}
		}
		if s.Peek() == 'i' {
			v = append(v, s.Next())
			return Token{Type: gotoken.IMAG, Value: string(v)}
		}
		if float {
			return Token{Type: gotoken.FLOAT, Value: string(v)}
		}
		if v[0] == '0' {
			return Token{Type: gotoken.ILLEGAL, Value: string(v)}
		}
		return Token{Type: gotoken.INT, Value: string(v)}
	}

	if isLetter(s.Peek()) {
		v = append(v, s.Next())
		for isLetter(s.Peek()) || isUnicodeDigit(s.Peek()) {
			v = append(v, s.Next())
		}
		switch string(v) {
		case "break":
			return Token{Type: gotoken.BREAK, Value: string(v)}
		case "case":
			return Token{Type: gotoken.CASE, Value: string(v)}
		case "chan":
			return Token{Type: gotoken.CHAN, Value: string(v)}
		case "const":
			return Token{Type: gotoken.CONST, Value: string(v)}
		case "continue":
			return Token{Type: gotoken.CONTINUE, Value: string(v)}
		case "default":
			return Token{Type: gotoken.DEFAULT, Value: string(v)}
		case "defer":
			return Token{Type: gotoken.DEFER, Value: string(v)}
		case "else":
			return Token{Type: gotoken.ELSE, Value: string(v)}
		case "enum":
			return Token{Type: token.ENUM, Value: string(v)}
		case "fallthrough":
			return Token{Type: gotoken.FALLTHROUGH, Value: string(v)}
		case "for":
			return Token{Type: gotoken.FOR, Value: string(v)}
		case "func":
			return Token{Type: gotoken.FUNC, Value: string(v)}
		case "go":
			return Token{Type: gotoken.GO, Value: string(v)}
		case "goto":
			return Token{Type: gotoken.GOTO, Value: string(v)}
		case "if":
			return Token{Type: gotoken.IF, Value: string(v)}
		case "import":
			return Token{Type: gotoken.IMPORT, Value: string(v)}
		case "interface":
			return Token{Type: gotoken.INTERFACE, Value: string(v)}
		case "map":
			return Token{Type: gotoken.MAP, Value: string(v)}
		case "match":
			return Token{Type: token.MATCH, Value: string(v)}
		case "package":
			return Token{Type: gotoken.PACKAGE, Value: string(v)}
		case "range":
			return Token{Type: gotoken.RANGE, Value: string(v)}
		case "return":
			return Token{Type: gotoken.RETURN, Value: string(v)}
		case "select":
			return Token{Type: gotoken.SELECT, Value: string(v)}
		case "struct":
			return Token{Type: gotoken.STRUCT, Value: string(v)}
		case "switch":
			return Token{Type: gotoken.SWITCH, Value: string(v)}
		case "type":
			return Token{Type: gotoken.TYPE, Value: string(v)}
		case "var":
			return Token{Type: gotoken.VAR, Value: string(v)}
		default:
			return Token{Type: gotoken.IDENT, Value: string(v)}
		}
	}

	if s.Peek() == '\'' {
		v = append(v, s.Next())
		if s.Peek() == '\\' {
			v = append(v, s.Next())
		}
		v = append(v, s.Next())
		if s.Peek() == '\'' {
			v = append(v, s.Next())
			return Token{Type: gotoken.CHAR, Value: string(v)}
		}
		return Token{Type: gotoken.ILLEGAL, Value: string(v)}
	}

	if s.Peek() == '"' {
		v = append(v, s.Next())
		escaped := false
		for !isNewline(s.Peek()) {
			cur := s.Next()
			v = append(v, cur)
			if escaped {
				escaped = false
			} else if cur == '\\' {
				escaped = true
			} else if cur == '"' {
				return Token{Type: gotoken.STRING, Value: string(v)}
			}
		}
		return Token{Type: gotoken.ILLEGAL, Value: string(v)}
	}

	if s.Peek() == '`' {
		v = append(v, s.Next())
		escaped := false
		for s.Peek() != scanner.EOF {
			cur := s.Next()
			v = append(v, cur)
			if escaped {
				escaped = false
			} else if cur == '\\' {
				escaped = true
			} else if cur == '`' {
				return Token{Type: gotoken.STRING, Value: string(v)}
			}
		}
		return Token{Type: gotoken.ILLEGAL, Value: string(v)}
	}

	if s.Peek() == '+' {
		v = append(v, s.Next())
		if s.Peek() == '=' {
			v = append(v, s.Next())
			return Token{Type: gotoken.ADD_ASSIGN, Value: string(v)}
		}
		if s.Peek() == '+' {
			v = append(v, s.Next())
			return Token{Type: gotoken.INC, Value: string(v)}
		}
		return Token{Type: gotoken.ADD, Value: string(v)}
	}

	if s.Peek() == '-' {
		v = append(v, s.Next())
		if s.Peek() == '=' {
			v = append(v, s.Next())
			return Token{Type: gotoken.SUB_ASSIGN, Value: string(v)}
		}
		if s.Peek() == '-' {
			v = append(v, s.Next())
			return Token{Type: gotoken.DEC, Value: string(v)}
		}
		return Token{Type: gotoken.SUB, Value: string(v)}
	}

	if s.Peek() == '*' {
		v = append(v, s.Next())
		if s.Peek() == '=' {
			v = append(v, s.Next())
			return Token{Type: gotoken.MUL_ASSIGN, Value: string(v)}
		}
		if s.Peek() == '*' {
			v = append(v, s.Next())
			if s.Peek() == '=' {
				v = append(v, s.Next())
				return Token{Type: token.POW_ASSIGN, Value: string(v)}
			}
			return Token{Type: token.POW, Value: string(v)}
		}
		return Token{Type: gotoken.MUL, Value: string(v)}
	}

	if s.Peek() == '/' {
		v = append(v, s.Next())
		if s.Peek() == '/' {
			v = append(v, s.Next())
			for !isNewline(s.Peek()) {
				v = append(v, s.Next())
			}
			return Token{Type: gotoken.COMMENT, Value: string(v)}
		}
		if s.Peek() == '*' {
			v = append(v, s.Next())
			for s.Peek() != scanner.EOF {
				cur := s.Next()
				v = append(v, cur)
				if cur == '*' && s.Peek() == '/' {
					v = append(v, s.Next())
					return Token{Type: gotoken.COMMENT, Value: string(v)}
				}
			}
			return Token{Type: gotoken.ILLEGAL, Value: string(v)}
		}
		if s.Peek() == '=' {
			v = append(v, s.Next())
			return Token{Type: gotoken.QUO_ASSIGN, Value: string(v)}
		}
		return Token{Type: gotoken.QUO, Value: string(v)}
	}

	if s.Peek() == '%' {
		v = append(v, s.Next())
		if s.Peek() == '=' {
			v = append(v, s.Next())
			return Token{Type: gotoken.REM_ASSIGN, Value: string(v)}
		}
		return Token{Type: gotoken.REM, Value: string(v)}
	}

	if s.Peek() == '&' {
		v = append(v, s.Next())
		if s.Peek() == '^' {
			v = append(v, s.Next())
			if s.Peek() == '=' {
				v = append(v, s.Next())
				return Token{Type: gotoken.AND_NOT_ASSIGN, Value: string(v)}
			}
			return Token{Type: gotoken.AND_NOT, Value: string(v)}
		}
		if s.Peek() == '=' {
			v = append(v, s.Next())
			return Token{Type: gotoken.AND_ASSIGN, Value: string(v)}
		}
		if s.Peek() == '&' {
			v = append(v, s.Next())
			return Token{Type: gotoken.LAND, Value: string(v)}
		}
		return Token{Type: gotoken.AND, Value: string(v)}
	}

	if s.Peek() == '|' {
		v = append(v, s.Next())
		if s.Peek() == '=' {
			v = append(v, s.Next())
			return Token{Type: gotoken.OR_ASSIGN, Value: string(v)}
		}
		if s.Peek() == '|' {
			v = append(v, s.Next())
			return Token{Type: gotoken.LOR, Value: string(v)}
		}
		return Token{Type: gotoken.OR, Value: string(v)}
	}

	if s.Peek() == '^' {
		v = append(v, s.Next())
		if s.Peek() == '=' {
			v = append(v, s.Next())
			return Token{Type: gotoken.XOR_ASSIGN, Value: string(v)}
		}
		return Token{Type: gotoken.XOR, Value: string(v)}
	}

	if s.Peek() == '<' {
		v = append(v, s.Next())
		if s.Peek() == '<' {
			v = append(v, s.Next())
			if s.Peek() == '=' {
				v = append(v, s.Next())
				return Token{Type: gotoken.SHL_ASSIGN, Value: string(v)}
			}
			return Token{Type: gotoken.SHL, Value: string(v)}
		}
		if s.Peek() == '-' {
			v = append(v, s.Next())
			return Token{Type: gotoken.ARROW, Value: string(v)}
		}
		if s.Peek() == '=' {
			v = append(v, s.Next())
			return Token{Type: gotoken.LEQ, Value: string(v)}
		}
		return Token{Type: gotoken.LSS, Value: string(v)}
	}

	if s.Peek() == '>' {
		v = append(v, s.Next())
		if s.Peek() == '>' {
			v = append(v, s.Next())
			if s.Peek() == '=' {
				v = append(v, s.Next())
				return Token{Type: gotoken.SHR_ASSIGN, Value: string(v)}
			}
			return Token{Type: gotoken.SHR, Value: string(v)}
		}
		if s.Peek() == '=' {
			v = append(v, s.Next())
			return Token{Type: gotoken.GEQ, Value: string(v)}
		}
		return Token{Type: gotoken.GTR, Value: string(v)}
	}

	if s.Peek() == '=' {
		v = append(v, s.Next())
		if s.Peek() == '=' {
			v = append(v, s.Next())
			return Token{Type: gotoken.EQL, Value: string(v)}
		}
		return Token{Type: gotoken.ASSIGN, Value: string(v)}
	}

	if s.Peek() == '!' {
		v = append(v, s.Next())
		if s.Peek() == '=' {
			v = append(v, s.Next())
			return Token{Type: gotoken.NEQ, Value: string(v)}
		}
		return Token{Type: gotoken.NOT, Value: string(v)}
	}

	if s.Peek() == ':' {
		v = append(v, s.Next())
		if s.Peek() == '=' {
			v = append(v, s.Next())
			return Token{Type: gotoken.DEFINE, Value: string(v)}
		}
		return Token{Type: gotoken.COLON, Value: string(v)}
	}

	if s.Peek() == '.' {
		v = append(v, s.Next())
		if s.Peek() == '.' {
			v = append(v, s.Next())
			if s.Peek() == '.' {
				v = append(v, s.Next())
				return Token{Type: gotoken.ELLIPSIS, Value: string(v)}
			}
			return Token{Type: gotoken.ILLEGAL, Value: string(v)}
		}
		return Token{Type: gotoken.PERIOD, Value: string(v)}
	}

	if s.Peek() == '(' {
		v = append(v, s.Next())
		return Token{Type: gotoken.LPAREN, Value: string(v)}
	}

	if s.Peek() == ')' {
		v = append(v, s.Next())
		return Token{Type: gotoken.RPAREN, Value: string(v)}
	}

	if s.Peek() == '[' {
		v = append(v, s.Next())
		return Token{Type: gotoken.LBRACK, Value: string(v)}
	}

	if s.Peek() == ']' {
		v = append(v, s.Next())
		return Token{Type: gotoken.RBRACK, Value: string(v)}
	}

	if s.Peek() == '{' {
		v = append(v, s.Next())
		return Token{Type: gotoken.LBRACE, Value: string(v)}
	}

	if s.Peek() == '}' {
		v = append(v, s.Next())
		return Token{Type: gotoken.RBRACE, Value: string(v)}
	}

	if s.Peek() == ',' {
		v = append(v, s.Next())
		return Token{Type: gotoken.COMMA, Value: string(v)}
	}

	if s.Peek() == ';' {
		v = append(v, s.Next())
		return Token{Type: gotoken.SEMICOLON, Value: string(v)}
	}

	v = append(v, s.Next())
	return Token{Type: gotoken.ILLEGAL, Value: string(v)}
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

func isDigitSeparator(r rune) bool {
	return r == '_'
}
