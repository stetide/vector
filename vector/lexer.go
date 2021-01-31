package vector

import (
	"errors"
	"strings"
)

const (
	sLOWER   = "abcdefghijklmnopqrstuvwxyz"
	sUPPER   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	sLETTERS = sLOWER + sUPPER
	sDIGITS  = "0123456789"
)

// Lexer lexes given input
type Lexer struct {
	text   string
	pos    int
	char   rune
	tokens []Token
}

// NewLexer returns new Lexer
func NewLexer(input string) *Lexer {
	l := &Lexer{text: input, pos: -1}
	l.advance()
	return l
}

func (l *Lexer) advance() {
	l.pos++
	if l.pos >= len(l.text) {
		l.char = 0
		return
	}
	l.char = rune(l.text[l.pos])
}

func (l *Lexer) makeNum() error {
	var numStr string
	var dotCount int
	for strings.ContainsRune(sDIGITS+".", l.char) {
		if l.char == '.' {
			dotCount++
		}
		if dotCount > 1 {
			return errors.New("Invalid number")
		}
		numStr += string(l.char)
		l.advance()
	}
	l.tokens = append(l.tokens, Token{tNUM, numStr})
	return nil
}

func (l *Lexer) makeIdentKw() {
	var identStr string
	for strings.ContainsRune(sLETTERS+"_"+sDIGITS, l.char) {
		identStr += string(l.char)
		l.advance()
	}

	if isKeyword(identStr) {
		l.tokens = append(l.tokens, Token{tKEYW, identStr})
		return
	}

	l.tokens = append(l.tokens, Token{tIDENT, identStr})
}

// GenerateTokens generates token slice from text
func (l *Lexer) GenerateTokens() ([]Token, error) {
	for l.pos < len(l.text) {
		if strings.ContainsRune(sDIGITS, l.char) {
			if err := l.makeNum(); err != nil {
				return nil, err
			}
			continue
		}
		if strings.ContainsRune(sLETTERS+"_", l.char) {
			l.makeIdentKw()
			continue
		}
		switch l.char {
		case ' ':
		case '+':
			l.tokens = append(l.tokens, Token{tPLUS, string(l.char)})
		case '-':
			l.tokens = append(l.tokens, Token{tMINUS, string(l.char)})
		case '*':
			l.tokens = append(l.tokens, Token{tMUL, string(l.char)})
		case '/', ':':
			l.tokens = append(l.tokens, Token{tDIV, string(l.char)})
		case '^':
			l.tokens = append(l.tokens, Token{tPOW, string(l.char)})
		case '\\':
			l.tokens = append(l.tokens, Token{tROOT, string(l.char)})
		case '=':
			l.tokens = append(l.tokens, Token{tEQ, string(l.char)})
		case '(':
			l.tokens = append(l.tokens, Token{tLPAREN, string(l.char)})
		case ')':
			l.tokens = append(l.tokens, Token{tRPAREN, string(l.char)})
		case '[':
			l.tokens = append(l.tokens, Token{tLVECPAR, string(l.char)})
		case ']':
			l.tokens = append(l.tokens, Token{tRVECPAR, string(l.char)})
		case '?':
			l.tokens = append(l.tokens, Token{tABSQ, string(l.char)})
		case '|':
			l.tokens = append(l.tokens, Token{tABS, string(l.char)})
		case ';':
			l.tokens = append(l.tokens, Token{tDLM, string(l.char)})
		default:
			return nil, errors.New("Invalid Character")
		}
		l.advance()
	}
	return l.tokens, nil
}
