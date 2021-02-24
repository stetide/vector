package vector

import (
	"fmt"
	"strings"
)

const (
	tEMPTY = iota
	tNUM
	tIDENT
	tKEYW
	tFUNC
	tDLM
	tSPACE
	tPLUS
	tMINUS
	tMUL
	tDIV
	tPOW
	tROOT
	tEQ
	tLPAREN
	tRPAREN
	tLVECPAR
	tRVECPAR
	tABSQ
	tABS
)

var sTypes = []string{
	"EMPTY",
	"NUM",
	"IDENT",
	"KEYW",
	"FUNC",
	"DLM",
	"SPACE",
	"PLUS",
	"MINUS",
	"MUL",
	"DIV",
	"POW",
	"ROOT",
	"EQ",
	"LPAREN",
	"RPAREN",
	"LVECPAR",
	"RVECPAR",
	"ABSQ",
	"ABS",
}

// TokenType is Token typ
type TokenType int

func (t TokenType) String() string {
	return sTypes[t]
}

// Token is Token
type Token struct {
	ttype TokenType
	val   string
}

func (t Token) String() string {
	if t.ttype == tKEYW || t.ttype == tFUNC {
		return strings.ToUpper(fmt.Sprintf("%s", t.val))
	} else if t.ttype > tFUNC {
		return fmt.Sprintf("%s", t.ttype)
	}
	return fmt.Sprintf("%s:%s", t.ttype, t.val)
}
