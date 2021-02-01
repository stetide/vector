package vector

import (
	"errors"
	"strconv"
)

// Parser is type Parser
type Parser struct {
	tokens []Token
	pos    int
	curTok Token
}

// NewParser returns new Parser
func NewParser(tokens []Token) *Parser {
	p := &Parser{tokens: tokens, pos: -1}
	p.advance()
	return p
}

func (p *Parser) advance() {
	p.pos++
	if p.pos >= len(p.tokens) {
		// log.Println("advance err", p.curTok.ttype)
		p.curTok = Token{}
		return
	}
	p.curTok = p.tokens[p.pos]
}

func (p *Parser) previous() Token {
	if p.pos-1 == -1 {
		return Token{}
	}
	return p.tokens[p.pos-1]
}

func (p *Parser) peek() Token {
	if p.pos+1 == len(p.tokens) {
		return Token{}
	}
	return p.tokens[p.pos+1]
}

func (p *Parser) expr() (Node, error) {
	left, err := p.term()
	if err != nil {
		return nil, err
	}
	for p.curTok.ttype == tPLUS || p.curTok.ttype == tMINUS {
		var node OperationNode
		node.left = left
		node.op = p.curTok
		p.advance()
		node.right, err = p.term()
		if err != nil {
			return nil, err
		}
		left = node
	}
	return left, nil
}

func (p *Parser) term() (Node, error) {
	left, err := p.atom()
	if err != nil {
		return nil, err
	}
	for p.curTok.ttype == tMUL || p.curTok.ttype == tDIV {
		var node OperationNode
		node.left = left
		node.op = p.curTok
		p.advance()
		node.right, err = p.atom()
		if err != nil {
			return nil, err
		}
		left = node
	}
	return left, nil
}

func (p *Parser) atom() (Node, error) {
	left, err := p.factor()
	if err != nil {
		return nil, err
	}
	for p.curTok.ttype == tPOW || p.curTok.ttype == tROOT {
		var node OperationNode
		node.left = left
		node.op = p.curTok
		p.advance()
		node.right, err = p.factor()
		if err != nil {
			return nil, err
		}
		left = node
	}
	return left, nil
}

func (p *Parser) makeUnaryNode() (UnaryNode, error) {
	var node UnaryNode
	var err error

	if p.curTok.ttype == tABS {
		node.op = Token{ttype: tABSQ, val: "?"}
		p.advance()
		node.node, err = p.factor()
		if err != nil {
			return node, err
		}
		if p.curTok.ttype != tABS {
			return node, errors.New("Expected |")
		}
		p.advance()
		return node, nil
	}

	node.op = p.curTok
	p.advance()
	node.node, err = p.factor()
	return node, err
}

func (p *Parser) makeParens() (Node, error) {
	var node Node
	var err error
	p.advance()
	// log.Println(p.pos)
	node, err = p.expr()
	if p.curTok.ttype != tRPAREN {
		err = errors.New("Expected )")
		// log.Println(p.curTok, p.pos)
	}
	p.advance()
	return node, err
}

func (p *Parser) makeVarNode() (VarNode, error) {
	var node VarNode
	var err error
	node.ident = p.curTok
	p.advance()
	if p.curTok.ttype != tEQ {
		// return VarNode{}, errors.New("Expected =")
		return node, nil
	}
	p.advance()
	node.val, err = p.expr()
	return node, err
}

func (p *Parser) makeNumNode() (NumberNode, error) {
	f, err := strconv.ParseFloat(p.curTok.val, 64)
	p.advance()
	return NumberNode(f), err
}

func (p *Parser) makeAns() Node {
	p.advance()
	return VarNode{Token{ttype: tIDENT, val: "ans"}, nil}
}

func (p *Parser) makeKeywNode() (Node, error) {
	var node Node
	var err error
	switch p.curTok.val {
	case kwVEC.name:
		node, err = p.makeVecNode()
	case kwQUIT.name, kwQUIT.getNameByAlias(p.curTok.val):
		err = ExitErr{}
	case kwHELP.name:
		err = HelpErr{}
	case kwANS.name:
		node = p.makeAns()
	case kwCLEAR.name, kwCLEAR.getNameByAlias(p.curTok.val):
		err = ClearErr{}
	default:
		err = errors.New("Keyword not implemented")
	}
	return node, err
}

func (p *Parser) makeVecNode() (VecNode, error) {
	var node VecNode

	var startTok = Token{tLVECPAR, "["}
	var endTok = Token{tRVECPAR, "]"}
	p.advance()
	if p.previous().ttype != tLVECPAR {
		if p.curTok.ttype != tLPAREN {
			return node, errors.New("Expected (")
		}
		p.advance()
		startTok = Token{tLPAREN, "("}
		endTok = Token{tRPAREN, ")"}
	}

	for p.curTok.ttype != endTok.ttype {
		switch p.curTok.ttype {
		case tEMPTY:
			return node, errors.New("Expected " + endTok.val)
		case tSPACE:
			p.advance()
			continue
		case tDLM:
			var i int
			if p.previous().ttype == startTok.ttype {
				// log.Println("prev")
				i++
			}
			for p.curTok.ttype == tDLM {
				// log.Println("curr")
				i++
				p.advance()
			}

			if p.curTok.ttype != endTok.ttype {
				// log.Println("next")
				i--
			}

			for i > 0 {
				i--
				node.fields = append(node.fields, NumberNode(0))
			}
			continue
		}

		n, err := p.expr()
		if err != nil {
			return node, err
		}
		switch n.(type) {
		case VecNode:
			return node, errors.New("Vec in vec not allowed")
		case VarNode:
			if n.(VarNode).val != nil {
				return node, errors.New("Cannot assign var in vec")
			}
		}
		node.fields = append(node.fields, n)
	}
	p.advance()

	return node, nil
}

func (p *Parser) factor() (Node, error) {
	var node Node
	var err error

	switch p.curTok.ttype {
	case tNUM:
		node, err = p.makeNumNode()
	case tMINUS, tPLUS, tABSQ, tABS:
		node, err = p.makeUnaryNode()
	case tLPAREN:
		node, err = p.makeParens()
	case tLVECPAR:
		node, err = p.makeVecNode()
	case tIDENT:
		node, err = p.makeVarNode()
	case tKEYW:
		node, err = p.makeKeywNode()
	default:
		err = errors.New("Expected expression")
	}
	return node, err
}

// Parse creates AST from tokens
func (p *Parser) Parse() (Node, error) {
	node, err := p.expr()
	if err != nil {
		return nil, err
	}
	if p.pos != len(p.tokens) {
		return nil, errors.New("Expected expression")
	}
	return node, nil
}
