package vector

import (
	"fmt"
	"runtime"
)

// VERSION is version
const VERSION = "2.8.13"

// Memory stores Ident, Node values
type Memory map[string]Node

var memory = Memory{}

// Run runs txt
func Run(txt string) (Node, error) {
	lexer := NewLexer(txt)
	tokens, err := lexer.GenerateTokens()
	if err != nil {
		return nil, err
	}

	parser := NewParser(tokens)
	ast, err := parser.Parse()
	if err != nil {
		switch err.(type) {
		case ExitErr:
			return nil, err
		}
		return nil, err
	}

	res, err := ast.resolve()
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Execute executes syntax tree
func Execute(ast Node) (Node, error) {
	res, err := ast.resolve()
	if err != nil {
		return nil, err
	}

	switch res.(type) {
	case VecNode, NumberNode:
		memory["ans"] = res
	}

	switch runtime.GOOS {
	case "js":
		fmt.Println(memory)
		// case "windows", "linux":
		// 	fmt.Println(memory)
	}

	return res, nil
}
