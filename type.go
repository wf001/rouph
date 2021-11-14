package main

func intType() *Type {
	t := new(Type)
	t.Kind = TY_INT
	return t
}
func pointerTo(base *Type) *Type {
	t := new(Type)
	t.Kind = TY_PTR
	t.Base = base
	return t
}
func arrayOf(base *Type, size int) *Type {
	ty := new(Type)
	ty.Kind = TY_ARRAY
	ty.Base = base
	ty.ArraySize = size
	return ty
}
func sizeOf(ty *Type) int {
	if ty.Kind == TY_INT || ty.Kind == TY_PTR {
		return 8
	}
	if ty.Kind != TY_ARRAY {
		panic("type is allowed only array.")
	}
	return sizeOf(ty.Base) * ty.ArraySize
}

func visit(node *Node) {
	if node == nil {
		return
	}
	visit(node.Lhs)
	visit(node.Rhs)
	visit(node.Cond)
	visit(node.Then)
	visit(node.Else)
	visit(node.Init)
	visit(node.Inc)

	for b := node.Body; b != nil; b = b.Next {
		visit(b)
	}
	for a := node.Args; a != nil; a = a.Next {
		visit(a)
	}

	switch node.Kind {
	case ND_KIND_MUL,
		ND_KIND_DIV,
		ND_KIND_EQ,
		ND_KIND_NE,
		ND_KIND_LE,
		ND_KIND_LT,
		ND_KIND_GE,
		ND_KIND_GT,
		ND_KIND_FUNCALL,
		ND_KIND_NUM:
		node.Ty = intType()
		return
	case ND_KIND_VAR:
		node.Ty = node.Var.Ty
		return
	case ND_KIND_ADD:
		if node.Rhs.Ty.Base != nil {
			tmp := node.Lhs
			node.Lhs = node.Rhs
			node.Rhs = tmp
		}
		if node.Rhs.Ty.Base != nil {
			panic("invalid token")
		}
		node.Ty = node.Lhs.Ty
		return
	case ND_KIND_SUB:
		if node.Rhs.Ty.Base != nil {
			panic("invalid token")
		}
		node.Ty = node.Lhs.Ty
		return
	case ND_KIND_ASSIGN:
		node.Ty = node.Lhs.Ty
		return
	case ND_KIND_ADDR:
		node.Ty = pointerTo(node.Lhs.Ty)
		if node.Lhs.Ty.Kind == TY_ARRAY {
			node.Ty = pointerTo(node.Lhs.Ty.Base)
		} else {
			node.Ty = pointerTo(node.Lhs.Ty)
		}
		return
	case ND_KIND_DEREF:
		if node.Lhs.Ty.Base == nil {
			panic("Invalid pointer dereference")
		}
		node.Ty = node.Lhs.Ty.Base
		return
	case ND_KIND_SIZEOF:
		node.Kind = ND_KIND_NUM
		node.Ty = intType()
		node.Val = sizeOf(node.Lhs.Ty)
		node.Lhs = nil
	}
}
func addType(prg *Function) {
	for fn := prg; fn != nil; fn = fn.Next {
		for n := fn.N; n != nil; n = n.Next {
			visit(n)
		}
	}
}
