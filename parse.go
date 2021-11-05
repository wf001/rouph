package main

type NodeKind int

var Locals *Var

const (
	ND_KIND_ADD       NodeKind = iota + 1 // +
	ND_KIND_SUB                           // -
	ND_KIND_MUL                           // *
	ND_KIND_DIV                           // /
	ND_KIND_EQ                            // ==
	ND_KIND_NE                            // !=
	ND_KIND_LT                            // <
	ND_KIND_LE                            // =<
	ND_KIND_GT                            // >
	ND_KIND_GE                            // >=
	ND_KIND_ASSIGN                        //=
	ND_KIND_RETURN                        // return
	ND_KIND_IF                            //if
	ND_KIND_FOR                           //for
	ND_KIND_EXPR_STMT                     // Expression Statement
	ND_KIND_VAR                           // Local Variables
	ND_KIND_NUM                           // Integer
)

type Node struct {
	Kind NodeKind
	Next *Node
	Lhs  *Node
	Cond *Node
	Then *Node
	Else *Node
	Init *Node
	Inc  *Node
	Rhs  *Node
	Var  *Var
	Val  int
}

type Var struct {
	Next   *Var
	Name   string
	Offset int
}

type Prg struct {
	N         *Node
	Locals    *Var
	StackSize int
}

func Program(tok *Token) *Prg {
	Locals = nil
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

	prg := new(Prg)
	prg.N = head.Next
	prg.Locals = Locals

	return prg
}
func stmt(tok *Token) (*Token, *Node) {
	Info("stmt : %+v\n", tok)
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

	if tok.Str == "if" {
		node := newNode(ND_KIND_IF, nil, nil)
		tok = tok.Next

		if tok.Val != "(" {
			panic("( not found")
		}
		tok = tok.Next
		tok, node.Cond = expr(tok)

		if tok.Val != ")" {
			panic(") not found")
		}
		tok = tok.Next
		tok, node.Then = stmt(tok)

		if tok.Str == "else" {
			tok = tok.Next
			tok, node.Else = stmt(tok)
		}

		return tok, node
	}
	if tok.Str == "for" {
		Info("--%s\n", "for")
		node := newNode(ND_KIND_FOR, nil, nil)
		tok = tok.Next
		if tok.Val != "(" {
			panic("; not found")
		}
		if tok.Val != ";" {
			tok = tok.Next
			tok, node.Init = readExprStmt(tok)
			if tok.Val != ";" {
				panic("; not found")
			}
		}
		if tok.Val != ";" {
			tok = tok.Next
			tok, node.Cond = expr(tok)
			if tok.Val != ";" {
				panic("; not found")
			}
		}
		if tok.Val != ")" {
			tok = tok.Next
			tok, node.Inc = readExprStmt(tok)
			if tok.Val != ";" {
				panic("; not found")
			}
			tok = tok.Next
		}
		Info("%d\n", 5)
		tok, node.Then = stmt(tok)
		return tok, node
	}

	tok, node := readExprStmt(tok)

	Info("tok :: %+v\n", tok)
	if tok.Val != ";" {
		panic("; not found")
	}
	tok = tok.Next
	return tok, node
}
func expr(tok *Token) (*Token, *Node) {
	return assign(tok)
}
func assign(tok *Token) (*Token, *Node) {
	var a_node *Node
	tok, node := equality(tok)
	if tok.Val == "=" {
		tok = tok.Next
		// node is left side value
		// a_node is right side value
		tok, a_node = assign(tok)
		node = newNode(ND_KIND_ASSIGN, node, a_node)
	}
	return tok, node
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

	if tok.Kind == TK_IDENT {
		i_tok := tok
		tok = tok.Next
		v := findVar(i_tok)
		if v == nil {
			v = pushVar(i_tok.Str)
		}
		return tok, newVar(v)
	} else {
		i := isNumber(tok.Val)
		tok = tok.Next
		return tok, newNodeNum(i)
	}
}
func newNode(kind NodeKind, lhs *Node, rhs *Node) *Node {
	node := new(Node)
	node.Kind = kind
	node.Lhs = lhs
	node.Rhs = rhs
	return node
}
func newVar(v *Var) *Node {
	node := newNode(ND_KIND_VAR, nil, nil)
	node.Var = v
	return node
}
func pushVar(name string) *Var {
	v := new(Var)
	v.Next = Locals
	v.Name = name
	Locals = v
	return v
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
		if node.Kind == ND_KIND_EXPR_STMT {
			Info("## node %p\n", node)
			Info("## %+v\n", node)
			printNode(node.Lhs)
			return
		}
		if node.Kind == ND_KIND_VAR {
			Info("## node %p\n", node)
			Info("## %+v\n", node)
			Info("->### var %p\n", node.Var)
			Info("->### var %+v\n", node.Var)
			return
		}
		if node.Kind == ND_KIND_ASSIGN {
			Info("## node %p\n", node)
			Info("## %+v\n", node)
			printNode(node.Lhs)
			printNode(node.Rhs)
			return
		}
		if node.Kind == ND_KIND_IF {
			if node.Else != nil {
				gen(node.Cond)
				gen(node.Then)
				gen(node.Else)
			} else {
				gen(node.Cond)
				gen(node.Then)
			}
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

func findVar(tok *Token) *Var {
	for v := Locals; v != nil; v = v.Next {
		Info("findVar:%s\n", v.Name)
		Info("findVar:%s\n", tok.Str)
		if v.Name == tok.Str {
			return v
		}
	}
	return nil
}
func readExprStmt(tok *Token) (*Token, *Node) {
	tok, r_node := expr(tok)
	node := newNode(ND_KIND_EXPR_STMT, r_node, nil)
	return tok, node

}
