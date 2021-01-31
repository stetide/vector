package vector

// VERSION is version
const VERSION = "2.3.9"

// Memory stores Ident, Node values
type Memory map[string]Node

var memory = Memory{}

// Execute executes syntax tree
func Execute(ast Node) (Node, error) {
	return ast.resolve()
}
