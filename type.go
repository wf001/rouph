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
		ND_KIND_VAR,
		ND_KIND_FUNCALL,
		ND_KIND_NUM:
		node.Ty = intType()
		return
	case ND_KIND_ADD:
		if node.Rhs.Ty.Kind == TY_PTR {
			tmp := node.Lhs
			node.Lhs = node.Rhs
			node.Rhs = tmp
		}
		if node.Rhs.Ty.Kind == TY_PTR {
			panic("invalid token")
		}
		node.Ty = node.Lhs.Ty
		return
	case ND_KIND_SUB:
		if node.Rhs.Ty.Kind == TY_PTR {
			panic("invalid token")
		}
		node.Ty = node.Lhs.Ty
		return
	case ND_KIND_ASSIGN:
		node.Ty = node.Lhs.Ty
		return
	case ND_KIND_ADDR:
		node.Ty = pointerTo(node.Lhs.Ty)
		return
	case ND_KIND_DEREF:
		if node.Lhs.Ty.Kind == TY_PTR {
			node.Ty = node.Lhs.Ty.Base
		} else {
			node.Ty = intType()
		}
		return
	}
}
func addType(prg *Function) {
	for fn := prg; fn != nil; fn = fn.Next {
		for n := fn.N; n != nil; n = n.Next {
			visit(n)
		}
	}
}
