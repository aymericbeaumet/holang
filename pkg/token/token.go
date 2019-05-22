package token

import "go/token"

type Token = token.Token

const (
	_ Token = token.VAR<<1 + iota

	// Operators
	POW
	POW_ASSIGN

	// Delimiters
	LCHEVR
	RCHEVR

	// Keywords
	ENUM
	MATCH
)
