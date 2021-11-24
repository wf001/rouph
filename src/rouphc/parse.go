package main

import (
	"strconv"
)

var Locals *VarList
var Globals *VarList
var cnt int = -1

/*
* Node
 */
type NodeKind int

const (
	ND_KIND_ADD       NodeKind = iota + 1 // +
	ND_KIND_SUB                           // -
	ND_KIND_MUL                           // *
	ND_KIND_DIV                           // /
	ND_KIND_REM                           // %
	ND_KIND_EQ                            // ==
	ND_KIND_NE                            // !=
	ND_KIND_LT                            // <
	ND_KIND_LE                            // =<
	ND_KIND_GT                            // >
	ND_KIND_GE                            // >=
	ND_KIND_ASSIGN                        // =
	ND_KIND_ADDR                          // unary &
	ND_KIND_DEREF                         // unary *
	ND_KIND_RETURN                        // return
	ND_KIND_IF                            // if
	ND_KIND_FOR                           // for
	ND_KIND_SIZEOF                        // size_of
	ND_KIND_BLOCK                         // {...}
	ND_KIND_FUNCALL                       // function call
	ND_KIND_EXPR_STMT                     // Expression Statement
	ND_KIND_VAR                           // Local Variables
	ND_KIND_NUM                           // Integer
	ND_KIND_NULL                          // null
	ND_KIND_STDLIB                        // standard library
)

type Node struct {
	Kind NodeKind
	Next *Node
	Ty   *Type
	Lhs  *Node
	Rhs  *Node
	Cond *Node
	Then *Node
	Else *Node
	Init *Node
	Inc  *Node
	Body *Node
	Func string
	Args *Node
	Var  *Var
	Val  int
}

type Var struct {
	Name     string
	Ty       *Type
	isLocal  bool
	Offset   int
	Contents string
	ContLen  int
}
type VarList struct {
	Next *VarList
	V    *Var
}

type Function struct {
	Next      *Function
	Name      string
	N         *Node
	Params    *VarList
	Locals    *VarList
	StackSize int
}
type Prog struct {
	Globals *VarList
	Fns     *Function
}

func Program(tok *Token) *Prog {
	printTokenAndeNode("program", tok)
	head := new(Function)
	head.Next = nil
	cur := head
	Globals = nil

	for {
		if tok.Kind == TK_KIND_EOF {
			break
		}
		if token, isFunc := isFunction(tok); isFunc {
			tok = token
			tok, cur.Next = function(tok)
			cur = cur.Next
		} else {
			tok = globalVar(tok)
		}
		Info("%s\n", "prg")
	}
	prg := new(Prog)
	prg.Globals = Globals
	prg.Fns = head.Next
	return prg
}

func function(tok *Token) (*Token, *Function) {
	printTokenAndeNode("function", tok)

	Locals = nil

	fn := new(Function)
	tok, _ = baseType(tok)
	if tok.Kind != TK_IDENT {
		panic("identifier not found")
	}
	fn.Name = tok.Str
	tok = tok.Next

	if tok.Val != "(" {
		panic("( not found")
	}
	tok = tok.Next

	tok, fn.Params = readFuncParams(tok)

	if tok.Val != "{" {
		panic("{ not found")
	}
	tok = tok.Next

	head := new(Node)
	head.Next = nil
	cur := head
	for {
		if tok.Val == "}" {
			tok = tok.Next
			break
		}
		tok, cur.Next = stmt(tok)
		cur = cur.Next
	}

	fn.N = head.Next
	fn.Locals = Locals

	return tok, fn
}

func stmt(tok *Token) (*Token, *Node) {
	printTokenAndeNode("stmt", tok)
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
		tok = tok.Next
		node := newNode(ND_KIND_FOR, nil, nil)

		if tok.Val != "(" {
			panic("( not found")
		}
		tok = tok.Next

		if tok.Val != ";" {
			tok, node.Init = readExprStmt(tok)
			if tok.Val != ";" {
				panic("; not found")
			}
		}
		tok = tok.Next

		if tok.Val != ";" {
			tok, node.Cond = expr(tok)
			if tok.Val != ";" {
				panic("; not found")
			}
		}
		tok = tok.Next

		if tok.Val != ")" {
			tok, node.Inc = readExprStmt(tok)
			if tok.Val != ")" {
				panic(") not found")
			}
		}
		tok = tok.Next
		tok, node.Then = stmt(tok)
		return tok, node
	}

	if tok.Val == "{" {
		tok = tok.Next
		head := new(Node)
		head.Next = nil
		cur := head

		for {
			if tok.Val == "}" {
				tok = tok.Next
				break
			}
			tok, cur.Next = stmt(tok)
			cur = cur.Next
		}
		node := newNode(ND_KIND_BLOCK, nil, nil)
		node.Body = head.Next
		return tok, node
	}

	if isTypeName(tok) {
		return declaration(tok)
	}

	tok, node := readExprStmt(tok)

	if tok.Val != ";" {
		panic("; not found")
	}
	tok = tok.Next
	return tok, node
}

func declaration(tok *Token) (*Token, *Node) {
	printTokenAndeNode("declaration", tok)
	tok, ty := baseType(tok)
	if tok.Kind != TK_IDENT {
		panic("identifier not found.")
	}
	name := tok.Str
	tok = tok.Next

	tok, ty = readTypeSuffix(tok, ty)
	v := pushVar(name, ty, true)

	if tok.Val == ";" {
		tok = tok.Next
		return tok, newNode(ND_KIND_NULL, nil, nil)
	}

	if tok.Val != "=" {
		panic("= not found.")
	}
	tok = tok.Next
	Lhs := newVar(v)
	tok, Rhs := expr(tok)
	if tok.Val != ";" {
		panic("; not found.")
	}
	tok = tok.Next
	node := newNode(ND_KIND_ASSIGN, Lhs, Rhs)

	return tok, newNode(ND_KIND_EXPR_STMT, node, nil)

}

func expr(tok *Token) (*Token, *Node) {
	printTokenAndeNode("expr", tok)
	return assign(tok)
}
func assign(tok *Token) (*Token, *Node) {
	printTokenAndeNode("assign", tok)
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
	printTokenAndeNode("equality", tok)
	var m_node *Node
	tok, node := relational(tok)
	for {
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
	printTokenAndeNode("relational", tok)
	var m_node *Node
	tok, node := add(tok)
	for {
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
	printTokenAndeNode("add", tok)
	var m_node *Node
	tok, node := mul(tok)
	for {
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
	printTokenAndeNode("mul", tok)
	var p_node *Node
	tok, node := unary(tok)
	for {
		if tok.Val == "*" {
			tok = tok.Next
			tok, p_node = unary(tok)
			node = newNode(ND_KIND_MUL, node, p_node)
		} else if tok.Val == "/" {
			tok = tok.Next
			tok, p_node = unary(tok)
			node = newNode(ND_KIND_DIV, node, p_node)
		} else if tok.Val == "%" {
			tok = tok.Next
			tok, p_node = unary(tok)
			node = newNode(ND_KIND_REM, node, p_node)
		} else {
			return tok, node
		}
	}
}
func unary(tok *Token) (*Token, *Node) {
	printTokenAndeNode("unary", tok)
	var u_node *Node
	if tok.Val == "+" {
		tok = tok.Next
		return unary(tok)
	}
	if tok.Val == "-" {
		tok = tok.Next
		tok, u_node = unary(tok)
		return tok, newNode(ND_KIND_SUB, newNodeNum(0), u_node)
	}
	if tok.Val == "&" {
		tok = tok.Next
		tok, u_node = unary(tok)
		return tok, newNode(ND_KIND_ADDR, u_node, nil)
	}
	if tok.Val == "*" {
		tok = tok.Next
		tok, u_node = unary(tok)
		return tok, newNode(ND_KIND_DEREF, u_node, nil)
	}
	return postFix(tok)
}
func postFix(tok *Token) (*Token, *Node) {
	tok, node := primary(tok)
	var p_node *Node
	for {
		if tok.Val != "[" {
			break
		}
		tok = tok.Next
		tok, p_node = expr(tok)
		exp := newNode(ND_KIND_ADD, p_node, node)
		if tok.Val != "]" {
			panic("not found ]")
		}
		tok = tok.Next
		node = newNode(ND_KIND_DEREF, exp, nil)
	}
	return tok, node
}
func primary(tok *Token) (*Token, *Node) {
	printTokenAndeNode("primary", tok)
	if tok.Val == "(" {
		tok = tok.Next
		tok, node := expr(tok)
		if tok.Val != ")" {
			panic("invalid closing.")
		}
		tok = tok.Next
		return tok, node
	}

	if tok.Str == "sizeof" {
		tok = tok.Next
		tok, node := unary(tok)
		return tok, newNode(ND_KIND_SIZEOF, node, nil)
	}
	if tok.Kind == TK_STDLIB {
		lib_name := tok.Val

		tok = tok.Next
		tok, node := unary(tok)
		node.Func = lib_name
		return tok, newNode(ND_KIND_STDLIB, node, nil)
	}

	if tok.Kind == TK_IDENT {
		i_tok := tok
		tok = tok.Next

		// function
		if tok.Val == "(" {
			tok = tok.Next
			node := newNode(ND_KIND_FUNCALL, nil, nil)
			node.Func = i_tok.Str
			tok, node.Args = funcArgs(tok)
			return tok, node
		}

		// variable
		v := findVar(i_tok)
		if v == nil {
			panic("undefined variable.")
		}
		return tok, newVar(v)
	}
	// type 'string'
	if tok.Kind == TK_STR {
		ty := arrayOf(charType(), tok.ContLen)
		v := pushVar(newLabel(), ty, false)
		v.Contents = tok.Contents
		Info("%d\n", v.ContLen)
		Info("%d\n", tok.ContLen)
		v.ContLen = tok.ContLen
		tok = tok.Next
		return tok, newVar(v)
	}
	i := isNumber(tok.Val)
	tok = tok.Next
	return tok, newNodeNum(i)
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
func pushVar(name string, ty *Type, isLocal bool) *Var {
	v := new(Var)
	v.Name = name
	v.Ty = ty
	v.isLocal = isLocal

	vl := new(VarList)
	vl.V = v

	if isLocal {
		vl.Next = Locals
		Locals = vl
	} else {
		vl.Next = Globals
		Globals = vl
	}

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
		Info("##\x1b[32m node %p\x1b[0m\n", node)
		switch node.Kind {

		case ND_KIND_NULL:
			return
		case ND_KIND_NUM:
			Info("## %+v\n", node)
			Info("## type -> %+v\n", node.Ty)
			return
		case ND_KIND_STDLIB:
			Info("## %+v\n", node)
			return
		case ND_KIND_EXPR_STMT:
			Info("## %+v\n", node)
			printNode(node.Lhs)
			return
		case ND_KIND_VAR:
			Info("## %+v\n", node)
			Info("->### var %p\n", node.Var)
			Info("->### var %+v\n", node.Var)
			return
		case ND_KIND_ADDR:
			printNode(node.Lhs)
			return
		case ND_KIND_DEREF:
			printNode(node.Lhs)
			return
		case ND_KIND_FUNCALL:
			for arg := node.Args; arg != nil; arg = arg.Next {
				gen(arg)
				printNode(arg)
			}
			return
		case ND_KIND_FOR:
			Info("## %+v\n", node)
			Info("->### var %p\n", node.Var)
			Info("->### var %+v\n", node.Var)
			if node.Init != nil {
				printNode(node.Init)
			}
			if node.Cond != nil {
				printNode(node.Cond)
			}
			gen(node.Then)
			if node.Inc != nil {
				printNode(node.Inc)
			}
			return
		case ND_KIND_ASSIGN:
			Info("## %+v\n", node)
			printNode(node.Lhs)
			printNode(node.Rhs)
			return
		case ND_KIND_IF:
			if node.Else != nil {
				gen(node.Cond)
				gen(node.Then)
				gen(node.Else)
			} else {
				gen(node.Cond)
				gen(node.Then)
			}
			return
		case ND_KIND_BLOCK:
			for n := node.Body; n != nil; n = n.Next {
				printNode(n)
			}
			return
		case ND_KIND_RETURN:
			Info("## %+v\n", node)
			printNode(node.Lhs)
			return
		}
		Info("## %+v\n", node)
		printNode(node.Lhs)
		printNode(node.Rhs)
	}
}

func findVar(tok *Token) *Var {
	for vl := Locals; vl != nil; vl = vl.Next {
		if vl.V.Name == tok.Str {
			return vl.V
		}
	}
	for vl := Globals; vl != nil; vl = vl.Next {
		if vl.V.Name == tok.Str {
			return vl.V
		}
	}
	return nil
}
func readExprStmt(tok *Token) (*Token, *Node) {
	printTokenAndeNode("readExprStmt", tok)
	tok, r_node := expr(tok)
	node := newNode(ND_KIND_EXPR_STMT, r_node, nil)
	return tok, node

}
func funcArgs(tok *Token) (*Token, *Node) {
	if tok.Val == ")" {
		tok = tok.Next
		return tok, nil
	}
	tok, head := assign(tok)
	cur := head
	for {
		if tok.Val != "," {
			break
		}
		tok = tok.Next

		tok, cur.Next = assign(tok)
		cur = cur.Next
	}
	if tok.Val != ")" {
		panic("Invalid closing.")
	}
	tok = tok.Next
	return tok, head
}
func printTokenAndeNode(name string, tok *Token) {
	if DEBUG {
		sep()
		Info("##\x1b[33m %s\x1b[0m\n", name)
		Info("## tok %+v\n", tok)
	}
}

func baseType(tok *Token) (*Token, *Type) {
	var ty *Type

	if tok.Str == "char" {
		tok = tok.Next
		ty = charType()
	} else {
		if tok.Str != "int" {
			panic("tok.Str must be char or int")
		}
		tok = tok.Next
		ty = intType()
	}
	for {
		if tok.Val != "*" {
			break
		}
		tok = tok.Next
		ty = pointerTo(ty)
	}
	return tok, ty
}

func readFuncParam(tok *Token) (*Token, *VarList) {
	tok, ty := baseType(tok)
	if tok.Kind != TK_IDENT {
		panic("not found identifier.")
	}
	name := tok.Str
	tok = tok.Next
	tok, ty = readTypeSuffix(tok, ty)

	vl := new(VarList)
	vl.V = pushVar(name, ty, true)
	return tok, vl
}

func readFuncParams(tok *Token) (*Token, *VarList) {
	if tok.Val == ")" {
		tok = tok.Next
		return tok, nil
	}

	tok, head := readFuncParam(tok)
	cur := head

	for {
		if tok.Val == ")" {
			break
		}
		if tok.Val != "," {
			panic("not found ','")
		}
		tok = tok.Next
		tok, cur.Next = readFuncParam(tok)
		cur = cur.Next
	}
	tok = tok.Next
	return tok, head

}
func readTypeSuffix(tok *Token, base *Type) (*Token, *Type) {
	if tok.Val != "[" {
		return tok, base
	}
	tok = tok.Next
	if tok.Kind != TK_KIND_NUM {
		panic("not number")
	}
	sz, _ := strconv.Atoi(tok.Val)
	tok = tok.Next
	if tok.Val != "]" {
		panic("not found ]")
	}
	tok = tok.Next
	tok, base = readTypeSuffix(tok, base)
	return tok, arrayOf(base, sz)
}
func isFunction(tok *Token) (*Token, bool) {
	token := tok
	tok, _ = baseType(tok)
	isFunc := false
	if tok.Kind == TK_IDENT {
		tok = tok.Next
		if tok.Val == "(" {
			tok = tok.Next
			isFunc = true
		}
	}
	tok = token
	return tok, isFunc

}
func globalVar(tok *Token) *Token {
	tok, ty := baseType(tok)
	if tok.Kind != TK_IDENT {
		panic("not identifier")
	}
	name := tok.Str
	tok = tok.Next
	tok, ty = readTypeSuffix(tok, ty)
	if tok.Val != ";" {
		panic("not found ;")
	}
	tok = tok.Next
	pushVar(name, ty, false)
	return tok
}
func isTypeName(tok *Token) bool {
	return tok.Str == "char" || tok.Str == "int"
}
func newLabel() string {
	cnt += 1
	return ".L.data." + strconv.Itoa(cnt)
}
