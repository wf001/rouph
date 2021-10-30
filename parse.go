package main

type NodeKind int

const (
	ND_KIND_ADD NodeKind = iota + 1
	ND_KIND_SUB
	ND_KIND_MUL
	ND_KIND_DIV
	ND_KIND_EQ
	ND_KIND_NE
	ND_KIND_LT     // <
	ND_KIND_LE     // =<
	ND_KIND_GT     // >
	ND_KIND_GE     // >=
	ND_KIND_RETURN // return
	ND_KIND_NUM
)

type Node struct {
	Kind NodeKind
	Next *Node
	Lhs  *Node
	Rhs  *Node
	Val  int
}

/*
Node Func
*/
func newNode(kind NodeKind, lhs *Node, rhs *Node) *Node {
	node := new(Node)
	node.Kind = kind
	node.Lhs = lhs
	node.Rhs = rhs
	return node
}
func newNodeNum(val int) *Node {
	node := new(Node)
	node.Kind = ND_KIND_NUM
	node.Val = val
	return node
}

func printNode(node *Node) {
	if DEBUG {
		if node.Kind == ND_KIND_NUM {
			Info("## node %p\n", node)
			Info("## %+v\n", node)
			return
		}
		if node.Kind == ND_KIND_RETURN {
			Info("## node %p\n", node)
			Info("## %+v\n", node)
			printNode(node.Lhs)
			return
		}
		Info("## node %p\n", node)
		Info("## %+v\n", node)
		printNode(node.Lhs)
		printNode(node.Rhs)
	}
}
func Program(tok *Token) (*Token, *Node) {
	head := new(Node)
	head.Next = nil
	cur := head
	for {
		if tok.Kind == TK_KIND_EOF {
			break
		}
		tok, cur.Next = stmt(tok)
		cur = cur.Next
	}
	return tok, head.Next
}
func stmt(tok *Token) (*Token, *Node) {
	if tok.Val == "return" {
		tok = tok.Next
		tok, e_node := expr(tok)
		node := newNode(ND_KIND_RETURN, e_node, nil)

		if tok.Val != ";" {
			panic("; not found")
		}
		tok = tok.Next
		return tok, node
	}

	tok, node := expr(tok)
	if tok.Val != ";" {
		panic("; not found")
	}
	tok = tok.Next
	return tok, node
}
func expr(tok *Token) (*Token, *Node) {
	return equality(tok)
}
func equality(tok *Token) (*Token, *Node) {
	var m_node *Node
	Info("%s\n", "eq")
	Info("%p\n", tok)
	tok, node := relational(tok)
	for {
		Info("%s\n", "eq for")
		Info("%s\n", tok.Val)
		if tok.Val == "==" {
			tok = tok.Next
			tok, m_node = relational(tok)
			node = newNode(ND_KIND_EQ, node, m_node)
		} else if tok.Val == "!=" {
			tok = tok.Next
			tok, m_node = relational(tok)
			node = newNode(ND_KIND_NE, node, m_node)
		} else {
			return tok, node
		}
	}
}
func relational(tok *Token) (*Token, *Node) {
	var m_node *Node
	Info("%s\n", "rel")
	Info("%p\n", tok)
	tok, node := add(tok)
	for {
		Info("%s\n", "rel for")
		Info("%s\n", tok.Val)
		if tok.Val == "<" {
			tok = tok.Next
			tok, m_node = add(tok)
			node = newNode(ND_KIND_LT, node, m_node)
		} else if tok.Val == "=<" {
			tok = tok.Next
			tok, m_node = add(tok)
			node = newNode(ND_KIND_LE, node, m_node)
		} else if tok.Val == ">" {
			tok = tok.Next
			tok, m_node = add(tok)
			node = newNode(ND_KIND_GT, node, m_node)
		} else if tok.Val == ">=" {
			tok = tok.Next
			tok, m_node = add(tok)
			node = newNode(ND_KIND_GE, node, m_node)
		} else {
			return tok, node
		}
	}
}

func add(tok *Token) (*Token, *Node) {
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
			return tok, node
		}
	}
}
func mul(tok *Token) (*Token, *Node) {
	var p_node *Node
	Info("%s\n", "mul")
	Info("%p\n", tok)
	tok, node := unary(tok)
	for {
		Info("%s\n", "mul for")
		Info("%s\n", tok.Val)
		if tok.Val == "*" {
			tok = tok.Next
			tok, p_node = unary(tok)
			node = newNode(ND_KIND_MUL, node, p_node)
		} else if tok.Val == "/" {
			tok = tok.Next
			tok, p_node = unary(tok)
			node = newNode(ND_KIND_DIV, node, p_node)
		} else {
			return tok, node
		}
	}
}
func unary(tok *Token) (*Token, *Node) {
	var u_node *Node
	if tok.Val == "+" {
		tok = tok.Next
		return primary(tok)
	} else if tok.Val == "-" {
		tok = tok.Next
		tok, u_node = primary(tok)
		return tok, newNode(ND_KIND_SUB, newNodeNum(0), u_node)
	}
	return primary(tok)
}
func primary(tok *Token) (*Token, *Node) {
	Info("%s\n", "pri")
	Info("%p\n", tok)
	Info("%s\n", tok.Val)
	if tok.Val == "(" {
		tok = tok.Next
		tok, node := expr(tok)
		Info("'%s'", tok.Val)
		if tok.Val != ")" {
			panic("invalid closing.")
		}
		tok = tok.Next
		return tok, node
	}
	i := isNumber(tok.Val)
	tok = tok.Next
	return tok, newNodeNum(i)
}
