package token

import "go/token"

type Token = token.Token

const (
	holang_beg Token = token.VAR<<1 + iota

	// Operators
	POW
	POW_ASSIGN

	// Keywords
	ENUM
	MATCH

	holang_end
)
