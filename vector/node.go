package vector

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

// Node is node type
type Node interface {
	resolve() (Node, error)
	String() string
}

// NumberNode represents number
type NumberNode float64

func (n NumberNode) add(a NumberNode) NumberNode {
	return n + a
}

func (n NumberNode) min(a NumberNode) NumberNode {
	return n - a
}

func (n NumberNode) mul(a NumberNode) NumberNode {
	return n * a
}

func (n NumberNode) div(a NumberNode) (NumberNode, error) {
	if a == 0 {
		return 0, errors.New("Division by zero")
	}
	return n / a, nil
}

func (n NumberNode) pow(a NumberNode) NumberNode {
	return NumberNode(math.Pow(float64(n), float64(a)))
}

func (n NumberNode) rot(a NumberNode) (NumberNode, error) {
	if n < 0 {
		return 0, errors.New("Negative root")
	} else if n == 0 {
		return 0, nil
	}
	return NumberNode(math.Pow(float64(n), 1/float64(a))), nil
}

func (n NumberNode) resolve() (Node, error) {
	return n, nil
}

func (n NumberNode) String() string {
	/*if float64(int(n)) == float64(n) {
		return fmt.Sprintf("%d", int(n))
	}
	return fmt.Sprintf("%f", float64(n))*/
	return strconv.FormatFloat(float64(n), 'f', -1, 64)
}

// UnaryNode is node with one Token
type UnaryNode struct {
	op   Token
	node Node
}

func (n UnaryNode) resolve() (Node, error) {
	var err error
	n.node, err = n.node.resolve()
	if err != nil {
		return nil, err
	}
	if n.node == nil {
		return nil, errors.New("Invalid syntax")
	}

	switch n.op.ttype {
	case tPLUS:
		break
	case tMINUS:
		switch n.node.(type) {
		case NumberNode:
			return -n.node.(NumberNode), nil
		case VecNode:
			n.node, err = n.node.resolve()
			if err != nil {
				return nil, err
			}
			return n.node.(VecNode).scalarMul(-1), nil
		}
	case tABSQ:
		switch n.node.(type) {
		case NumberNode:
			return NumberNode(math.Abs(float64(n.node.(NumberNode)))), nil
		case VecNode:
			return n.node.(VecNode).abs(), nil
		}
	default:
		return nil, errors.New("Unary operator not implemented")
	}
	return n, nil
}

func (n UnaryNode) String() string {
	return fmt.Sprintf("(%s%s)", n.op.val, n.node)
}

// OperationNode is Node with two nodes
type OperationNode struct {
	left  Node
	op    Token
	right Node
}

func (n OperationNode) conflicts() bool {
	switch n.left.(type) {
	case VecNode:
		switch n.right.(type) {
		case NumberNode:
			return true
		}
	case NumberNode:
		switch n.right.(type) {
		case VecNode:
			return true
		}
	}
	return false
}

func (n OperationNode) resolve() (Node, error) {
	var node Node
	var err error
	n.left, err = n.left.resolve()
	if err != nil {
		return nil, err
	}

	n.right, err = n.right.resolve()
	if err != nil {
		return nil, err
	}

	switch n.op.ttype {
	case tPLUS:
		if n.conflicts() {
			return nil, errors.New("Conflicting types")
		}
		switch n.left.(type) {
		case VecNode:
			node = n.left.(VecNode).add(n.right.(VecNode))
		case NumberNode:
			node = n.left.(NumberNode).add(n.right.(NumberNode))
		default:
			return nil, errors.New("Unexpected type")
		}
	case tMINUS:
		if n.conflicts() {
			return nil, errors.New("Conflicting types")
		}
		switch n.left.(type) {
		case VecNode:
			node = n.left.(VecNode).min(n.right.(VecNode))
		case NumberNode:
			node = n.left.(NumberNode).min(n.right.(NumberNode))
		default:
			return nil, errors.New("Unexpected type")
		}
	case tMUL:
		switch n.left.(type) {
		case VecNode:
			switch n.right.(type) {
			case NumberNode:
				node = n.left.(VecNode).scalarMul(n.right.(NumberNode))
			default:
				return nil, errors.New("Not implemented")
			}
		case NumberNode:
			switch n.right.(type) {
			case NumberNode:
				node = n.left.(NumberNode).mul(n.right.(NumberNode))
			case VecNode:
				node = n.right.(VecNode).scalarMul(n.left.(NumberNode))
			}
		default:
			return nil, errors.New("Unexpected type")
		}
	case tDIV:
		switch n.left.(type) {
		case VecNode:
			switch n.right.(type) {
			case NumberNode:
				if node, err = n.left.(VecNode).scalarDiv(n.right.(NumberNode)); err != nil {
					return nil, err
				}
			default:
				return nil, errors.New("Not implemented")
			}
		case NumberNode:
			switch n.right.(type) {
			case NumberNode:
				if node, err = n.left.(NumberNode).div(n.right.(NumberNode)); err != nil {
					return nil, err
				}
			case VecNode:
				return nil, errors.New("Cannot divide by Vec")
			}
		default:
			return nil, errors.New("Unexpected type")
		}
	case tPOW:
		switch n.left.(type) {
		case VecNode:
			return nil, errors.New("Pow for vec not implemented")
		case NumberNode:
			switch n.right.(type) {
			case VecNode:
				return nil, errors.New("Pow for vec not implemented")
			case NumberNode:
				node = n.left.(NumberNode).pow(n.right.(NumberNode))
			}
		}
	case tROOT:
		switch n.left.(type) {
		case VecNode:
			return nil, errors.New("Pow for vec not implemented")
		case NumberNode:
			switch n.right.(type) {
			case VecNode:
				return nil, errors.New("Pow for vec not implemented")
			case NumberNode:
				if node, err = n.right.(NumberNode).rot(n.left.(NumberNode)); err != nil {
					return nil, err
				}
			}
		}
	}
	return node, nil
}

func (n OperationNode) String() string {
	return fmt.Sprintf("(%s %s %s)", n.left, n.op.val, n.right)
}

// VarNode holds Ident and value
type VarNode struct {
	ident Token
	val   Node
}

func (n VarNode) resolve() (Node, error) {
	// var called
	if n.val == nil {
		if v, ok := memory[n.ident.val]; ok {
			return v.resolve()
		}
		return nil, errors.New(n.ident.val + " is not defined")
	}

	// test value
	switch n.val.(type) {
	case VarNode:
		if _, err := n.val.resolve(); err != nil {
			return nil, err
		}
		tmp, ok := memory[n.val.(VarNode).ident.val]
		if !ok {
			return nil, errors.New(n.val.(VarNode).ident.val + "is not defined")
		}
		n.val = tmp
	}

	memory[n.ident.val] = n.val
	return nil, nil
}

func (n VarNode) String() string {
	if n.val == nil {
		return fmt.Sprintf("(%s::called)", n.ident)
	}
	return fmt.Sprintf("(%s::%s)", n.ident, n.val)
}

// VecNode represents Vector
type VecNode struct {
	fields []Node
}

func (n VecNode) add(a VecNode) VecNode {
	var node VecNode

	nlen := len(n.fields)
	alen := len(a.fields)
	if alen != nlen {
		if alen < nlen {
			for i := 0; i < nlen-alen; i++ {
				a.fields = append(a.fields, NumberNode(0))
			}
		} else {
			for i := 0; i < alen-nlen; i++ {
				n.fields = append(n.fields, NumberNode(0))
			}
		}
	}

	for i, f := range n.fields {
		f = f.(NumberNode).add(a.fields[i].(NumberNode))
		node.fields = append(node.fields, f)
	}
	return node
}

func (n VecNode) min(a VecNode) VecNode {
	var node VecNode

	nlen := len(n.fields)
	alen := len(a.fields)
	if alen != nlen {
		if alen < nlen {
			for i := 0; i < nlen-alen; i++ {
				a.fields = append(a.fields, NumberNode(0))
			}
		} else {
			for i := 0; i < alen-nlen; i++ {
				n.fields = append(n.fields, NumberNode(0))
			}
		}
	}

	for i, f := range n.fields {
		f = f.(NumberNode).min(a.fields[i].(NumberNode))
		node.fields = append(node.fields, f)
	}
	return node
}

func (n VecNode) mul(a VecNode) {}

func (n VecNode) div(a VecNode) {}

func (n VecNode) scalarMul(a NumberNode) VecNode {
	var node VecNode
	for _, f := range n.fields {
		f = f.(NumberNode).mul(a)
		node.fields = append(node.fields, f)
	}
	return node
}

func (n VecNode) scalarDiv(a NumberNode) (VecNode, error) {
	var err error
	var node VecNode
	if a == 0 {
		return VecNode{}, errors.New("Divisino by zero")
	}
	for _, f := range n.fields {
		f, err = f.(NumberNode).div(a)
		node.fields = append(node.fields, f)
		if err != nil {
			return n, err
		}
	}
	return node, nil
}

func (n VecNode) abs() NumberNode {
	var res NumberNode
	for _, f := range n.fields {
		res += f.(NumberNode).pow(2)
	}
	res, _ = res.rot(2)
	return res
}

func (n VecNode) resolve() (Node, error) {
	var node VecNode
	for _, f := range n.fields {
		var err error
		f, err = f.resolve()
		if err != nil {
			return nil, err
		}
		switch f.(type) {
		case VecNode:
			return nil, errors.New("Vec in vec not allowed")
		}
		node.fields = append(node.fields, f)
	}
	return node, nil
}

func (n VecNode) String() string {
	var str = "vec("
	for i, nd := range n.fields {
		if i == len(n.fields)-1 {
			str += fmt.Sprintf("%s", nd)
			continue
		}
		str += fmt.Sprintf("%s ", nd)
	}
	return str + ")"
}
