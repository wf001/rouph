package main

type NodeKind int

const (
	ND_KIND_ADD NodeKind = iota + 1
	ND_KIND_SUB
	ND_KIND_MUL
	ND_KIND_DIV
	ND_KIND_NUM
)

type Node struct {
	Kind NodeKind
	Lhs  *Node
	Rhs  *Node
	Val  int
}


/*
Node Func
*/
func newNode(kind NodeKind, lhs *Node, rhs *Node) *Node {
	Info("newNode %d\n", kind)
	node := new(Node)
	node.Kind = kind
	node.Lhs = lhs
	node.Rhs = rhs
	return node
}
func newNodeNum(val int) *Node {
	Info("newNodeNum %d\n", val)
	node := new(Node)
	Info("%p\n", node)
	node.Kind = ND_KIND_NUM
	node.Val = val
	return node
}

func printNode(node *Node) {
	if DEBUG {
		if node.Kind == ND_KIND_NUM {
			Info("node %p\n", node)
			Info("%+v\n", node)
			return
		}
		Info("node %p\n", node)
		Info("%+v\n", node)
		printNode(node.Lhs)
		printNode(node.Rhs)
	}
}
func Expr(tok *Token) *Node {
	var m_node *Node
	Info("%s\n", "expr")
	Info("%p\n", tok)
	tok, node := mul(tok)
	for {
		Info("%s\n", "expr for")
		Info("%s\n", tok.Val)
		if tok.Val == "+" {
			tok = tok.Next
			tok, m_node = mul(tok)
			node = newNode(ND_KIND_ADD, node, m_node)
		} else if tok.Val == "-" {
			tok = tok.Next
			tok, m_node = mul(tok)
			node = newNode(ND_KIND_SUB, node, m_node)
		} else {
			Info("%s\n", "=================")
			printNode(node)
			Info("%s\n", "=================")
			return node
		}
	}
}
func mul(tok *Token) (*Token, *Node) {
    var p_node *Node
	Info("%s\n", "mul")
	Info("%p\n", tok)
	tok, node := primary(tok)
	for {
		Info("%s\n", "mul for")
		Info("%s\n", tok.Val)
		if tok.Val == "*" {
			tok = tok.Next
			tok, p_node = primary(tok)
			node = newNode(ND_KIND_MUL, node, p_node)
		} else if tok.Val == "/" {
			tok = tok.Next
			tok, p_node = primary(tok)
			node = newNode(ND_KIND_DIV, node, p_node)
		} else {
			return tok, node
		}
	}
}
func primary(tok *Token) (*Token, *Node) {
	Info("%s\n", "pri")
	Info("%p\n", tok)
	Info("%s\n", tok.Val)
	if tok.Val == "(" {
		tok = tok.Next
		node := Expr(tok)
		if tok.Val != ")" {
			panic("error")
		}
		tok = tok.Next
		return tok, node
	}
	i := isNumber(tok.Val)
	tok = tok.Next
	return tok, newNodeNum(i)
}
